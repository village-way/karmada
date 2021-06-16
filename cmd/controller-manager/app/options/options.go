package options

import (
	"time"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/klog/v2"
)

var (
	defaultElectionLeaseDuration = metav1.Duration{Duration: 15 * time.Second}
	defaultElectionRenewDeadline = metav1.Duration{Duration: 10 * time.Second}
	defaultElectionRetryPeriod   = metav1.Duration{Duration: 2 * time.Second}
)

const (
	defaultBindAddress = "0.0.0.0"
	defaultPort        = 10357
)

// Options contains everything necessary to create and run controller-manager.
type Options struct {
	HostNamespace  string
	LeaderElection componentbaseconfig.LeaderElectionConfiguration
	// BindAddress is the IP address on which to listen for the --secure-port port.
	BindAddress string
	// SecurePort is the port that the the server serves at.
	// Note: We hope support https in the future once controller-runtime provides the functionality.
	SecurePort int
	// ClusterStatusUpdateFrequency is the frequency that controller computes and report cluster status.
	// It must work with ClusterMonitorGracePeriod(--cluster-monitor-grace-period) in karmada-controller-manager.
	ClusterStatusUpdateFrequency metav1.Duration
	// ClusterLeaseDuration is a duration that candidates for a lease need to wait to force acquire it.
	// This is measure against time of last observed lease RenewTime.
	ClusterLeaseDuration metav1.Duration
	// ClusterLeaseRenewIntervalFraction is a fraction coordinated with ClusterLeaseDuration that
	// how long the current holder of a lease has last updated the lease.
	ClusterLeaseRenewIntervalFraction float64
}

// NewOptions builds an empty options.
func NewOptions() *Options {
	return &Options{}
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (o *Options) Complete() {
	if len(o.HostNamespace) == 0 {
		o.HostNamespace = "default"
		klog.Infof("Set default value: Options.HostNamespace = %s", "default")
	}

	if len(o.LeaderElection.ResourceLock) == 0 {
		o.LeaderElection.ResourceLock = resourcelock.EndpointsLeasesResourceLock
		klog.Infof("Set default value: Options.LeaderElection.ResourceLock = %s", resourcelock.EndpointsLeasesResourceLock)
	}

	if o.LeaderElection.LeaseDuration.Duration.Seconds() == 0 {
		o.LeaderElection.LeaseDuration = defaultElectionLeaseDuration
		klog.Infof("Set default value: Options.LeaderElection.LeaseDuration = %s", defaultElectionLeaseDuration.Duration.String())
	}

	if o.LeaderElection.RenewDeadline.Duration.Seconds() == 0 {
		o.LeaderElection.RenewDeadline = defaultElectionRenewDeadline
		klog.Infof("Set default value: Options.LeaderElection.RenewDeadline = %s", defaultElectionRenewDeadline.Duration.String())
	}

	if o.LeaderElection.RetryPeriod.Duration.Seconds() == 0 {
		o.LeaderElection.RetryPeriod = defaultElectionRetryPeriod
		klog.Infof("Set default value: Options.LeaderElection.RetryPeriod = %s", defaultElectionRetryPeriod.Duration.String())
	}
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.BindAddress, "bind-address", defaultBindAddress,
		"The IP address on which to listen for the --secure-port port.")
	flags.IntVar(&o.SecurePort, "secure-port", defaultPort,
		"The secure port on which to serve HTTPS.")
	flags.DurationVar(&o.ClusterStatusUpdateFrequency.Duration, "cluster-status-update-frequency", 10*time.Second,
		"Specifies how often karmada-controller-manager posts cluster status to karmada-apiserver.")
	flags.BoolVar(&o.LeaderElection.LeaderElect, "leader-elect", true, "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.")
	flags.DurationVar(&o.ClusterLeaseDuration.Duration, "cluster-lease-duration", 40*time.Second,
		"Specifies the expiration period of a cluster lease.")
	flags.Float64Var(&o.ClusterLeaseRenewIntervalFraction, "cluster-lease-renew-interval-fraction", 0.25,
		"Specifies the cluster lease renew interval fraction.")
}
