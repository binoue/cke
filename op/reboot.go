package op

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cybozu-go/cke"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
)

const defaultEvictionTimeoutSeconds = 600

type rebootOp struct {
	apiservers []*cke.Node
	nodes      []*cke.Node
	index      int64
	config     *cke.Reboot
	step       int

	mu          sync.Mutex
	failedNodes []string
}

func (o *rebootOp) notifyFailedNode(n *cke.Node) {
	o.mu.Lock()
	o.failedNodes = append(o.failedNodes, n.Nodename())
	o.mu.Unlock()
}

// RebootOp returns an Operator to reboot nodes.
func RebootOp(apiservers []*cke.Node, nodes []*cke.Node, index int64, config *cke.Reboot) cke.InfoOperator {
	return &rebootOp{
		apiservers: apiservers,
		nodes:      nodes,
		index:      index,
		config:     config,
	}
}

func (o *rebootOp) Name() string {
	return "reboot"
}

func (o *rebootOp) NextCommand() cke.Commander {
	switch o.step {
	case 0:
		o.step++
		return rebootStartCommand{index: o.index}
	case 1:
		o.step++
		return cordonCommand{
			apiservers:    o.apiservers,
			nodes:         o.nodes,
			unschedulable: true,
		}
	case 2:
		o.step++
		return drainCommand{
			timeoutSeconds:      o.config.EvictionTimeoutSeconds,
			apiserver:           o.apiservers[0],
			nodes:               o.nodes,
			protectedNamespaces: o.config.ProtectedNamespaces,
		}
	case 3:
		o.step++
		return rebootCommand{
			command:          o.config.Command,
			timeoutSeconds:   o.config.CommandTimeoutSeconds,
			nodes:            o.nodes,
			notifyFailedNode: o.notifyFailedNode,
		}
	case 4:
		o.step++
		return cordonCommand{
			apiservers:    o.apiservers,
			nodes:         o.nodes,
			unschedulable: false,
		}
	default:
		return nil
	}
}

func (o *rebootOp) Targets() []string {
	ipAddresses := make([]string, len(o.nodes))
	for i, n := range o.nodes {
		ipAddresses[i] = n.Address
	}
	return ipAddresses
}

func (o *rebootOp) Info() string {
	if len(o.failedNodes) == 0 {
		return ""
	}
	return fmt.Sprintf("failed to reboot some nodes: %v", o.failedNodes)
}

type rebootStartCommand struct {
	index int64
}

func (c rebootStartCommand) Run(ctx context.Context, inf cke.Infrastructure, _ string) error {
	entry, err := inf.Storage().GetRebootsEntry(ctx, c.index)
	if err != nil {
		return err
	}
	entry.Status = cke.RebootStatusRebooting
	return inf.Storage().UpdateRebootsEntry(ctx, entry)
}

func (c rebootStartCommand) Command() cke.Command {
	return cke.Command{
		Name:   "rebootStartCommand",
		Target: strconv.FormatInt(c.index, 10),
	}
}

type cordonCommand struct {
	apiservers    []*cke.Node
	nodes         []*cke.Node
	unschedulable bool
}

func (c cordonCommand) Run(ctx context.Context, inf cke.Infrastructure, _ string) error {
	var cs *kubernetes.Clientset
	var err error
	for _, server := range c.apiservers {
		cs, err = inf.K8sClient(ctx, server)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	nodesAPI := cs.CoreV1().Nodes()
	for _, ckeNode := range c.nodes {
		n, err := nodesAPI.Get(ctx, ckeNode.Nodename(), metav1.GetOptions{})
		if err != nil {
			return err
		}
		if n.Spec.Unschedulable == c.unschedulable {
			continue
		}

		oldData, err := json.Marshal(n)
		if err != nil {
			return err
		}
		n.Spec.Unschedulable = c.unschedulable
		newData, err := json.Marshal(n)
		if err != nil {
			return err
		}
		patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, n)
		if err != nil {
			return fmt.Errorf("failed to create patch for node %s: %v", n.Name, err)
		}
		_, err = nodesAPI.Patch(ctx, n.Name, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("failed to patch node %s: %v", n.Name, err)
		}
	}
	return nil
}

func (c cordonCommand) Command() cke.Command {
	ipAddresses := make([]string, len(c.nodes))
	for i, n := range c.nodes {
		ipAddresses[i] = n.Address
	}
	return cke.Command{
		Name:   "cordonCommand",
		Target: strings.Join(ipAddresses, ","),
	}
}

type drainCommand struct {
	timeoutSeconds      *int
	apiserver           *cke.Node
	nodes               []*cke.Node
	protectedNamespaces *metav1.LabelSelector
}

func listProtectedNamespaces(ctx context.Context, cs *kubernetes.Clientset, ls *metav1.LabelSelector) (map[string]bool, error) {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		// ls should have been validated
		panic(err)
	}
	protected, err := cs.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, err
	}

	nss := make(map[string]bool)
	for _, ns := range protected.Items {
		nss[ns.Name] = true
	}

	return nss, nil
}

func evictOrDeleteNodePod(ctx context.Context, cs *kubernetes.Clientset, n *cke.Node, protected map[string]bool) ([]*corev1.Pod, error) {
	var targets []*corev1.Pod
	podList, err := cs.CoreV1().Pods(corev1.NamespaceAll).List(ctx, metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{"spec.nodeName": n.Nodename()}).String(),
	})
	if err != nil {
		return nil, err
	}
	for _, pod := range podList.Items {
		pod := pod
		owner := metav1.GetControllerOf(&pod)
		if owner != nil && owner.Kind == "DaemonSet" {
			continue
		}
		targets = append(targets, &pod)
		err := cs.CoreV1().Pods(pod.Namespace).Evict(ctx, &policyv1beta1.Eviction{
			ObjectMeta: metav1.ObjectMeta{Name: pod.Name, Namespace: pod.Namespace},
		})
		switch {
		case apierrors.IsNotFound(err):
		case apierrors.IsTooManyRequests(err) && !protected[pod.Namespace]:
			err := cs.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
			if err != nil {
				return nil, err
			}
			log.Warn("deleted non-protected pod", map[string]interface{}{
				"namespace": pod.Namespace,
				"name":      pod.Name,
			})
		case err != nil:
			return nil, fmt.Errorf("failed to evict pod %s/%s: %w", pod.Namespace, pod.Name, err)
		}
	}
	return targets, nil
}

func waitPodDeletion(ctx context.Context, cs *kubernetes.Clientset, pods []*corev1.Pod, ts *int) error {
	timeoutSeconds := defaultEvictionTimeoutSeconds
	if ts != nil {
		timeoutSeconds = *ts
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(timeoutSeconds))
	defer cancel()

OUTER:
	for _, pod := range pods {
		for {
			p, err := cs.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if apierrors.IsNotFound(err) || (p != nil && p.ObjectMeta.UID != pod.ObjectMeta.UID) {
				// pod is deleted, or moved to another node
				continue OUTER
			}
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				msg := "aborted waiting for pod eviction"
				log.Error(msg, map[string]interface{}{
					"namespace": p.Namespace,
					"name":      p.Name,
				})
				return fmt.Errorf(fmt.Sprintf("%s: %s/%s", msg, p.Namespace, p.Name))
			case <-time.After(time.Second * 5):
				log.Info("waiting for pods to be deleted...", nil)
			}
		}
	}
	return nil
}

func (c drainCommand) Run(ctx context.Context, inf cke.Infrastructure, _ string) error {
	cs, err := inf.K8sClient(ctx, c.apiserver)
	if err != nil {
		return err
	}

	protected, err := listProtectedNamespaces(ctx, cs, c.protectedNamespaces)
	if err != nil {
		return err
	}

	var targets []*corev1.Pod
	for _, n := range c.nodes {
		pods, err := evictOrDeleteNodePod(ctx, cs, n, protected)
		if err != nil {
			return err
		}
		targets = append(targets, pods...)
	}

	return waitPodDeletion(ctx, cs, targets, c.timeoutSeconds)
}

func (c drainCommand) Command() cke.Command {
	ipAddresses := make([]string, len(c.nodes))
	for i, n := range c.nodes {
		ipAddresses[i] = n.Address
	}
	return cke.Command{
		Name:   "drainCommand",
		Target: strings.Join(ipAddresses, ","),
	}
}

type rebootCommand struct {
	command          []string
	timeoutSeconds   *int
	nodes            []*cke.Node
	notifyFailedNode func(*cke.Node)
}

func (c rebootCommand) Run(ctx context.Context, inf cke.Infrastructure, _ string) error {
	if c.timeoutSeconds != nil && *c.timeoutSeconds != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(*c.timeoutSeconds))
		defer cancel()
	}

	env := well.NewEnvironment(ctx)
	for _, n := range c.nodes {
		n := n

		env.Go(func(ctx context.Context) error {
			nodeJson, err := json.Marshal(n)
			if err != nil {
				return err
			}

			command := well.CommandContext(ctx, c.command[0], c.command[1:]...)
			command.Stdin = bytes.NewReader(nodeJson)
			err = command.Run()
			if err != nil {
				c.notifyFailedNode(n)
				log.Warn("failed on rebooting node", map[string]interface{}{
					log.FnError: err,
					"node":      n.Nodename(),
				})
			}
			return nil
		})
	}
	env.Stop()
	return env.Wait()
}

func (c rebootCommand) Command() cke.Command {
	ipAddresses := make([]string, len(c.nodes))
	for i, n := range c.nodes {
		ipAddresses[i] = n.Address
	}
	return cke.Command{
		Name:   "rebootCommand",
		Target: strings.Join(ipAddresses, ","),
	}
}
