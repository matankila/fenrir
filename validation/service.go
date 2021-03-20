package validation

import (
	"encoding/json"
	"github.com/matankila/fenrir/config"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type Service struct{}

func (s Service) IsValid(rawObj []byte, ns string) error {
	var service v1.Service
	conf := config.FallBackConf.Get()
	customPodPolicies := conf.Service.CustomPolicies
	policy := conf.Service.DefaultPolicy
	if v, ok := customPodPolicies[config.Namespace(ns)]; ok {
		policy = v
	}

	if !conf.Service.PolicyEnforcement {
		return nil
	}

	// if cannot be serialized to service, skip.
	if err := json.Unmarshal(rawObj, &service); err != nil {
		return err
	}

	// check ns settings
	if policy.DefaultNs {
		if service.Namespace == config.DefaultNs {
			return config.RestrictedNamespace
		}
	}

	// load balancers are costs alot of money, not always secure and tightly link to cloud provider.
	// use ingress instead
	// links:
	// https://medium.com/sainsburys-engineering/why-we-dont-use-the-loadbalancer-k8s-service-type-5f5403d42dfd
	// https://www.weave.works/blog/kubernetes-best-practices
	if policy.LoadBalancer {
		if strings.EqualFold(string(service.Spec.Type), string(v1.ServiceTypeLoadBalancer)) {
			return config.LoadBalancer
		}
	}

	return nil
}
