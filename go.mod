module github.com/clamoriniere/ddksm

go 1.13

require (
	github.com/DataDog/datadog-go v2.2.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.20.0
	k8s.io/autoscaler/vertical-pod-autoscaler v0.0.0-20200123122250-fa95810cfc1e
	k8s.io/client-go v0.20.0
	k8s.io/klog v1.0.0
	k8s.io/kube-state-metrics v1.9.5
)

replace k8s.io/kube-state-metrics => github.com/clamoriniere/kube-state-metrics v1.8.1-0.20200412142917-6c764c23fffb
