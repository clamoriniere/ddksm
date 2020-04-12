package options

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog"
	"k8s.io/kube-state-metrics/pkg/options"
)

// Options are the configurable parameters for ddksm.
type Options struct {
	options.Options

	StatsdHost string
	StatsdPort int

	flags *pflag.FlagSet
}

// NewOptions returns a new instance of `Options`.
func NewOptions() *Options {
	return &Options{
		Options: *options.NewOptions(),
	}
}

// AddFlags populated the Options struct from the command line arguments passed.
func (o *Options) AddFlags() {
	o.flags = pflag.NewFlagSet("", pflag.ExitOnError)
	// add klog flags
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	o.flags.AddGoFlagSet(klogFlags)
	o.flags.Lookup("logtostderr").Value.Set("true")
	o.flags.Lookup("logtostderr").DefValue = "true"
	o.flags.Lookup("logtostderr").NoOptDefVal = "true"

	o.flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		o.flags.PrintDefaults()
	}

	o.flags.StringVar(&o.Apiserver, "apiserver", "", `The URL of the apiserver to use as a master`)
	o.flags.StringVar(&o.Kubeconfig, "kubeconfig", "", "Absolute path to the kubeconfig file")
	o.flags.BoolVarP(&o.Help, "help", "h", false, "Print Help text")
	o.flags.IntVar(&o.Port, "port", 8080, `Port to expose metrics on.`)
	o.flags.StringVar(&o.Host, "host", "0.0.0.0", `Host to expose metrics on.`)
	o.flags.IntVar(&o.TelemetryPort, "telemetry-port", 8081, `Port to expose kube-state-metrics self metrics on.`)
	o.flags.StringVar(&o.TelemetryHost, "telemetry-host", "0.0.0.0", `Host to expose kube-state-metrics self metrics on.`)
	o.flags.Var(&o.Resources, "resources", fmt.Sprintf("Comma-separated list of Resources to be enabled. Defaults to %q", &options.DefaultResources))
	o.flags.Var(&o.Namespaces, "namespaces", fmt.Sprintf("Comma-separated list of namespaces to be enabled. Defaults to %q", &options.DefaultNamespaces))
	o.flags.Var(&o.MetricAllowlist, "metric-allowlist", "Comma-separated list of metrics to be exposed. This list comprises of exact metric names and/or regex patterns. The allowlist and denylist are mutually exclusive.")
	o.flags.Var(&o.MetricDenylist, "metric-denylist", "Comma-separated list of metrics not to be enabled. This list comprises of exact metric names and/or regex patterns. The allowlist and denylist are mutually exclusive.")
	o.flags.Int32Var(&o.Shard, "shard", int32(0), "The instances shard nominal (zero indexed) within the total number of shards. (default 0)")
	o.flags.IntVar(&o.TotalShards, "total-shards", 1, "The total number of shards. Sharding is disabled when total shards is set to 1.")

	autoshardingNotice := "When set, it is expected that --pod and --pod-namespace are both set. Most likely this should be passed via the downward API. This is used for auto-detecting sharding. If set, this has preference over statically configured sharding. This is experimental, it may be removed without notice."

	o.flags.StringVar(&o.Pod, "pod", "", "Name of the pod that contains the kube-state-metrics container. "+autoshardingNotice)
	o.flags.StringVar(&o.Namespace, "pod-namespace", "", "Name of the namespace of the pod specified by --pod. "+autoshardingNotice)
	o.flags.BoolVarP(&o.Version, "version", "", false, "kube-state-metrics build version information")
	o.flags.BoolVar(&o.EnableGZIPEncoding, "enable-gzip-encoding", false, "Gzip responses when requested by clients via 'Accept-Encoding: gzip' header.")

	o.flags.IntVar(&o.StatsdPort, "statsd-port", 80, `Statsd port to send metrics.`)
	o.flags.StringVar(&o.StatsdHost, "statsd-host", "0.0.0.0", `Statsd host to send metrics on.`)
}

// Parse parses the flag definitions from the argument list.
func (o *Options) Parse() error {
	err := o.flags.Parse(os.Args)
	return err
}

// Usage is the function called when an error occurs while parsing flags.
func (o *Options) Usage() {
	o.flags.Usage()
}
