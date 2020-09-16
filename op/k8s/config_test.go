package k8s

import (
	"reflect"
	"testing"
	"time"

	"github.com/cybozu-go/cke"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletv1beta1 "k8s.io/kubelet/config/v1beta1"
)

func TestGenerateKubeletConfiguration(t *testing.T) {
	t.Parallel()

	baseExpected := kubeletv1beta1.KubeletConfiguration{
		ReadOnlyPort:          0,
		HealthzBindAddress:    "0.0.0.0",
		OOMScoreAdj:           int32Pointer(-1000),
		FailSwapOn:            boolPointer(true),
		ContainerLogMaxFiles:  int32Pointer(0),
		RuntimeRequestTimeout: metav1.Duration{Duration: 15 * time.Minute},
		TLSCertFile:           "/etc/kubernetes/pki/kubelet.crt",
		TLSPrivateKeyFile:     "/etc/kubernetes/pki/kubelet.key",
		Authentication: kubeletv1beta1.KubeletAuthentication{
			X509:    kubeletv1beta1.KubeletX509Authentication{ClientCAFile: "/etc/kubernetes/pki/ca.crt"},
			Webhook: kubeletv1beta1.KubeletWebhookAuthentication{Enabled: boolPointer(true)},
		},
		Authorization: kubeletv1beta1.KubeletAuthorization{
			Mode: kubeletv1beta1.KubeletAuthorizationModeWebhook,
		},
		ClusterDNS: []string{"1.2.3.4"},
	}

	expected1 := baseExpected.DeepCopy()
	expected1.FailSwapOn = boolPointer(false)
	expected1.ClusterDomain = "foo.local"
	expected1.CgroupDriver = "systemd"
	expected1.ContainerLogMaxSize = "5Mi"
	expected1.ContainerLogMaxFiles = int32Pointer(10)

	expected2 := baseExpected.DeepCopy()
	expected2.CgroupDriver = "systemd"
	expected2.ContainerLogMaxSize = "100Mi"

	cases := []struct {
		Name     string
		Input    cke.KubeletParams
		Expected kubeletv1beta1.KubeletConfiguration
	}{
		{
			Name:     "no config",
			Input:    cke.KubeletParams{},
			Expected: baseExpected,
		},
		{
			Name: "no config_v1beta1",
			Input: cke.KubeletParams{
				AllowSwap:            true,
				Domain:               "foo.local",
				CgroupDriver:         "systemd",
				ContainerLogMaxSize:  "5Mi",
				ContainerLogMaxFiles: 10,
			},
			Expected: *expected1,
		},
		{
			Name: "with config_v1beta1",
			Input: cke.KubeletParams{
				CgroupDriver:        "systemd",
				ContainerLogMaxSize: "5Mi",
				ConfigV1Beta1: &kubeletv1beta1.KubeletConfiguration{
					ContainerLogMaxSize: "100Mi",
				},
			},
			Expected: *expected2,
		},
	}

	for _, c := range cases {
		conf := GenerateKubeletConfiguration(c.Input, "1.2.3.4")
		if !reflect.DeepEqual(conf, c.Expected) {
			t.Errorf("case %q: GenerateKubeletConfiguration() generated unexpected result:\n%s", c.Name, cmp.Diff(conf, c.Expected))
		}
	}
}
