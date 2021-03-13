package config

import "errors"

const (
	DefaultLoggerName = "github.com.matankila.fenrir.default.logger"
	WatcherLoggerName = "github.com.matankila.fenrir.watcher.logger"
	DefaultNs         = "default"
)

var (
	FallBackConf = AutoConfigure{
		config: Config{
			Pod: Pod{
				PolicyEnforcement: true,
				CustomPolicies:    map[Namespace]PodPolicy{},
				DefaultPolicy: PodPolicy{
					Resources:         true,
					ReadinessLiveness: true,
					DefaultNs:         true,
					LatestImageTag:    false,
					RunAsNonRoot:      false,
				},
			},
			Service: Service{
				PolicyEnforcement: true,
				DefaultPolicy: ServicePolicy{
					LoadBalancer: true,
					DefaultNs:    true,
				},
				CustomPolicies: map[Namespace]ServicePolicy{},
			},
		},
	}

	EmptyRequest        = errors.New("admission request is empty")
	RestrictedNamespace = errors.New("deployment in 'default' namespace is restricted")
	NoLiveness          = errors.New("deployment without liveness probe is prohibited")
	NoReadiness         = errors.New("deployment without readiness probe is prohibited")
	NoResources         = errors.New("deployment without resource declared is prohibited")
	EmptyResources      = errors.New("deployment without empty resource declared is prohibited")
	RunAsRoot           = errors.New("deployment is able to run as root, please fix set pod.spec.securityContext.runAsNonRoot to true")
	LatestImageTag      = errors.New("container image tag is latest, this might lead to unexpected behaviours, please set it to valid version")

	LoadBalancer = errors.New("LoadBalancer type in service, is not recommended, use ingress instead")
)
