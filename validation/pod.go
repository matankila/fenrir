package validation

import (
	"encoding/json"
	"github.com/matankila/fenrir/config"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type Pod struct{}

// this function checks if pod settings are valid according to set policy.
func (p Pod) IsValid(rawObj []byte, ns string) error {
	var pod v1.Pod
	conf := config.FallBackConf.Get()
	customPodPolicies := conf.Pod.CustomPolicies
	policy := conf.Pod.DefaultPolicy
	if v, ok := customPodPolicies[config.Namespace(ns)]; ok {
		policy = v
	}

	// if pod policy is off\ false, skip.
	if !conf.Pod.PolicyEnforcement {
		return nil
	}

	// if cannot be serialized to pod, skip.
	if err := json.Unmarshal(rawObj, &pod); err != nil {
		return nil
	}

	// check ns settings
	if policy.DefaultNs {
		if pod.Namespace == config.DefaultNs {
			return config.RestrictedNamespace
		}
	}

	// check containers settings
	for _, c := range pod.Spec.Containers {
		if policy.ReadinessLiveness {
			if c.LivenessProbe == nil {
				return config.NoLiveness
			} else if c.ReadinessProbe == nil {
				return config.NoReadiness
			}
		}

		if policy.Resources {
			if c.Resources.Limits == nil || c.Resources.Requests == nil {
				return config.NoResources
			} else if len(c.Resources.Limits) == 0 || len(c.Resources.Requests) == 0 {
				return config.EmptyResources
			}
		}

		if policy.LatestImageTag {
			s := strings.Split(c.Image, ":")
			if s[1] == "latest" {
				return config.LatestImageTag
			}
		}
	}

	// check security context settings
	if policy.RunAsNonRoot {
		if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.RunAsNonRoot != nil && *(pod.Spec.SecurityContext.RunAsNonRoot) == false {
			return config.RunAsRoot
		}
	}

	return nil
}
