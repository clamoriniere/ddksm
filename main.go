package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	vpaclientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"k8s.io/kube-state-metrics/pkg/allowdenylist"
	ksmoptions "k8s.io/kube-state-metrics/pkg/options"
	"k8s.io/kube-state-metrics/pkg/version"

	ddbuilder "github.com/clamoriniere/ddksm/pkg/builder"
	"github.com/clamoriniere/ddksm/pkg/options"
	ddstore "github.com/clamoriniere/ddksm/pkg/store"
)

func main() {
	opts := options.NewOptions()
	opts.AddFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := opts.Parse()
	if err != nil {
		klog.Fatalf("Error: %s", err)
	}

	if opts.Version {
		fmt.Printf("%#v\n", version.GetVersion())
		os.Exit(0)
	}

	if opts.Help {
		opts.Usage()
		os.Exit(0)
	}

	// Datadog statsd creation
	client, err := statsd.New(fmt.Sprintf("%s:%d", opts.StatsdHost, opts.StatsdPort),
		statsd.WithNamespace("ddksm."), // prefix every metric with the app name
	)
	if err != nil {
		log.Fatal(err)
	}
	klog.Infof("statsd server : %s:%d", opts.StatsdHost, opts.StatsdPort)
	builder := ddbuilder.New(client)

	var resources []string
	if len(opts.Resources) == 0 {
		klog.Info("Using default resources")
		resources = ksmoptions.DefaultResources.AsSlice()
	} else {
		klog.Infof("Using collectors %s", opts.Resources.String())
		resources = opts.Resources.AsSlice()
	}

	if err := builder.WithEnabledResources(resources); err != nil {
		klog.Fatalf("Failed to set up collectors: %v", err)
	}

	if len(opts.Namespaces) == 0 {
		klog.Info("Using all namespace")
		builder.WithNamespaces(ksmoptions.DefaultNamespaces)
	} else {
		if opts.Namespaces.IsAllNamespaces() {
			klog.Info("Using all namespace")
		} else {
			klog.Infof("Using %s namespaces", opts.Namespaces)
		}
		builder.WithNamespaces(opts.Namespaces)
	}

	allowdenylist, err := allowdenylist.New(opts.MetricAllowlist, opts.MetricDenylist)
	if err != nil {
		klog.Fatal(err)
	}

	err = allowdenylist.Parse()
	if err != nil {
		klog.Fatalf("error initializing the whiteblack list : %v", err)
	}

	klog.Infof("metric white-blacklisting: %v", allowdenylist.Status())

	builder.WithAllowDenyList(allowdenylist)

	kubeClient, vpaClient, err := createKubeClient(opts.Apiserver, opts.Kubeconfig)
	if err != nil {
		klog.Fatalf("Failed to create client: %v", err)
	}
	builder.WithKubeClient(kubeClient)
	builder.WithVPAClient(vpaClient)
	builder.WithContext(ctx)
	builder.WithGenerateStoreFunc(builder.GenerateStore)

	// Finally build the cache.Store
	stores := builder.Build()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	klog.Infof("Start processing")
	for {
		select {
		case <-ctx.Done():
			klog.Infof("Done")
			return
		case t := <-ticker.C:
			klog.Infof("ticker: %v", t)
			for _, store := range stores {
				store.(*ddstore.MetricsStore).Push()
			}
		}
	}
}

func createKubeClient(apiserver string, kubeconfig string) (clientset.Interface, vpaclientset.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	kubeClient, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	vpaClient, err := vpaclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	// Informers don't seem to do a good job logging error messages when it
	// can't reach the server, making debugging hard. This makes it easier to
	// figure out if apiserver is configured incorrectly.
	klog.Infof("Testing communication with server")
	v, err := kubeClient.Discovery().ServerVersion()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error while trying to communicate with apiserver")
	}
	klog.Infof("Running with Kubernetes cluster version: v%s.%s. git version: %s. git tree state: %s. commit: %s. platform: %s",
		v.Major, v.Minor, v.GitVersion, v.GitTreeState, v.GitCommit, v.Platform)
	klog.Infof("Communication with server successful")

	return kubeClient, vpaClient, nil
}
