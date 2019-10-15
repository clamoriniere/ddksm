module github.com/clamoriniere/ddksm

go 1.13

require (
	github.com/DataDog/datadog-go v2.2.0+incompatible
	github.com/campoy/embedmd v1.0.0
	github.com/dgryski/go-jump v0.0.0-20170409065014-e1f439676b57
	github.com/gogo/protobuf v1.2.0
	github.com/google/go-jsonnet v0.14.0
	github.com/jsonnet-bundler/jsonnet-bundler v0.1.1-0.20190930114713-10e24cb86976
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/robfig/cron/v3 v3.0.0
	github.com/spf13/pflag v1.0.3
	golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e
	k8s.io/api v0.0.0-20190620084959-7cf5895f2711
	k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/autoscaler v0.0.0-20190607113959-1b4f1855cb8e
	k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
	k8s.io/klog v0.4.0
	k8s.io/kube-state-metrics v0.0.0-00010101000000-000000000000
)

replace k8s.io/kube-state-metrics => github.com/clamoriniere/kube-state-metrics v1.6.1-0.20191015135933-3f6d9d9f6e93
