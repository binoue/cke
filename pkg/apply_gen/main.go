package main

import (
	"fmt"
	"os"
	"text/template"
)

var tmpl = template.Must(template.New("").Parse(`// Code generated by apply_gen.go. DO NOT EDIT.
//go:generate go run ./pkg/apply_gen

package cke

import (
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
{{- range . }}

func apply{{ .Kind }}(o *{{ .API }}.{{ .Kind }}, data []byte, rev int64, getFunc func(string, metav1.GetOptions) (*{{ .API }}.{{ .Kind }}, error), createFunc func(*{{ .API }}.{{ .Kind }}) (*{{ .API }}.{{ .Kind }}, error), patchFunc func(string, types.PatchType, []byte, ...string) (*{{ .API }}.{{ .Kind }}, error), deleteFunc func(string, *metav1.DeleteOptions) error) error {
	annotate(&o.ObjectMeta, rev, data)
	current, err := getFunc(o.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = createFunc(o)
		return err
	}
	if err != nil {
		return err
	}

	modified, err := encodeToJSON(o)
	if err != nil {
		return err
	}

	original, ok := current.Annotations[AnnotationResourceOriginal]
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

	err = deleteFunc(o.Name, metav1.NewDeleteOptions({{ .GracefulSeconds }}))
	if err != nil {
		return err
	}
	_, err = createFunc(o)

	return err
}
{{- end }}
`))

func main() {
	err := subMain()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func subMain() error {
	f, err := os.OpenFile("resource_apply.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, []struct {
		API             string
		Kind            string
		GracefulSeconds int64
	}{
		{"corev1", "Namespace", 60},
		{"corev1", "ServiceAccount", 0},
		{"corev1", "ConfigMap", 0},
		{"corev1", "Service", 0},
		{"policyv1beta1", "PodSecurityPolicy", 0},
		{"networkingv1", "NetworkPolicy", 0},
		{"rbacv1", "Role", 0},
		{"rbacv1", "RoleBinding", 0},
		{"rbacv1", "ClusterRole", 0},
		{"rbacv1", "ClusterRoleBinding", 0},
		{"appsv1", "Deployment", 60},
		{"appsv1", "DaemonSet", 60},
		{"batchv1beta1", "CronJob", 60},
	})
	if err != nil {
		return err
	}
	return f.Sync()
}
