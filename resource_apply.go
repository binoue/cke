// Code generated by apply_gen.go. DO NOT EDIT.
//go:generate go run ./pkg/apply_gen

package cke

import (
	"fmt"
	"strconv"

	"github.com/cybozu-go/log"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func annotate(meta *metav1.ObjectMeta, rev int64, data []byte) {
	if meta.Annotations == nil {
		meta.Annotations = make(map[string]string)
	}
	meta.Annotations[AnnotationResourceRevision] = strconv.FormatInt(rev, 10)
	meta.Annotations[AnnotationResourceOriginal] = string(data)
}

func applyNamespace(o *corev1.Namespace, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*corev1.Namespace, error), createFunc func(*corev1.Namespace) (*corev1.Namespace, error), patchFunc func(string, types.PatchType, []byte, ...string) (*corev1.Namespace, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(60))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyServiceAccount(o *corev1.ServiceAccount, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*corev1.ServiceAccount, error), createFunc func(*corev1.ServiceAccount) (*corev1.ServiceAccount, error), patchFunc func(string, types.PatchType, []byte, ...string) (*corev1.ServiceAccount, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyConfigMap(o *corev1.ConfigMap, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*corev1.ConfigMap, error), createFunc func(*corev1.ConfigMap) (*corev1.ConfigMap, error), patchFunc func(string, types.PatchType, []byte, ...string) (*corev1.ConfigMap, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyService(o *corev1.Service, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*corev1.Service, error), createFunc func(*corev1.Service) (*corev1.Service, error), patchFunc func(string, types.PatchType, []byte, ...string) (*corev1.Service, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyPodSecurityPolicy(o *policyv1beta1.PodSecurityPolicy, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*policyv1beta1.PodSecurityPolicy, error), createFunc func(*policyv1beta1.PodSecurityPolicy) (*policyv1beta1.PodSecurityPolicy, error), patchFunc func(string, types.PatchType, []byte, ...string) (*policyv1beta1.PodSecurityPolicy, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyNetworkPolicy(o *networkingv1.NetworkPolicy, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*networkingv1.NetworkPolicy, error), createFunc func(*networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error), patchFunc func(string, types.PatchType, []byte, ...string) (*networkingv1.NetworkPolicy, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyRole(o *rbacv1.Role, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*rbacv1.Role, error), createFunc func(*rbacv1.Role) (*rbacv1.Role, error), patchFunc func(string, types.PatchType, []byte, ...string) (*rbacv1.Role, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyRoleBinding(o *rbacv1.RoleBinding, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*rbacv1.RoleBinding, error), createFunc func(*rbacv1.RoleBinding) (*rbacv1.RoleBinding, error), patchFunc func(string, types.PatchType, []byte, ...string) (*rbacv1.RoleBinding, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyClusterRole(o *rbacv1.ClusterRole, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*rbacv1.ClusterRole, error), createFunc func(*rbacv1.ClusterRole) (*rbacv1.ClusterRole, error), patchFunc func(string, types.PatchType, []byte, ...string) (*rbacv1.ClusterRole, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyClusterRoleBinding(o *rbacv1.ClusterRoleBinding, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*rbacv1.ClusterRoleBinding, error), createFunc func(*rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error), patchFunc func(string, types.PatchType, []byte, ...string) (*rbacv1.ClusterRoleBinding, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(0))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyDeployment(o *appsv1.Deployment, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*appsv1.Deployment, error), createFunc func(*appsv1.Deployment) (*appsv1.Deployment, error), patchFunc func(string, types.PatchType, []byte, ...string) (*appsv1.Deployment, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(60))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyDaemonSet(o *appsv1.DaemonSet, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*appsv1.DaemonSet, error), createFunc func(*appsv1.DaemonSet) (*appsv1.DaemonSet, error), patchFunc func(string, types.PatchType, []byte, ...string) (*appsv1.DaemonSet, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(60))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}

func applyCronJob(o *batchv1beta1.CronJob, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*batchv1beta1.CronJob, error), createFunc func(*batchv1beta1.CronJob) (*batchv1beta1.CronJob, error), patchFunc func(string, types.PatchType, []byte, ...string) (*batchv1beta1.CronJob, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	var curRev int64
	curRevStr, ok := current.Annotations[AnnotationResourceRevision]
	original := current.Annotations[AnnotationResourceOriginal]
	if ok {
		curRev, err = strconv.ParseInt(curRevStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid revision annotation for %s/%s/%s", o.Kind, o.Namespace, o.Name)
		}
	}

	if curRev == rev {
		return nil
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}
	if !ok {
		original = string(modified)
		log.Warn("use modified resource as original for 3-way patch", map[string]interface{}{
			"kind":      o.Kind,
			"namespace": o.Namespace,
			"name":      o.Name,
		})
	}
	currentData, err := encodeToJSON(current)
	if err != nil {
		return err
	}
	pm, err := strategicpatch.NewPatchMetaFromStruct(o)
	if err != nil {
		return err
	}
	patch, err := strategicpatch.CreateThreeWayMergePatch([]byte(original), modified, currentData, pm, true)
	if err != nil {
		return err
	}
	_, err = patchFunc(o.Name, types.StrategicMergePatchType, patch)
	if err == nil {
		return nil
	}

	log.Warn("failed to apply patch", map[string]interface{}{
		"kind":      o.Kind,
		"namespace": o.Namespace,
		"name":      o.Name,
		log.FnError: err,
	})

	err = deleteFunc(o.Name, metav1.NewDeleteOptions(60))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}