package config

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Namespace string

type AutoConfigure struct {
	config Config
	sync.RWMutex
}

type Config struct {
	Pod     Pod     `json:"pod"`
	Service Service `json:"service"`
}

type Service struct {
	PolicyEnforcement bool                        `json:"policy_enforcement"`
	DefaultPolicy     ServicePolicy               `json:"default_policy"`
	CustomPolicies    map[Namespace]ServicePolicy `json:"custom_policies"`
}

type Pod struct {
	PolicyEnforcement bool                    `json:"policy_enforcement"`
	DefaultPolicy     PodPolicy               `json:"default_policy"`
	CustomPolicies    map[Namespace]PodPolicy `json:"custom_policies"`
}

type ServicePolicy struct {
	LoadBalancer bool `json:"load_balancer"`
	DefaultNs    bool `json:"default_ns"`
}

// pod policy settings, if not set default is false
type PodPolicy struct {
	ReadinessLiveness bool `json:"readiness_liveness"`
	Resources         bool `json:"resources"`
	DefaultNs         bool `json:"default_ns"`
	LatestImageTag    bool `json:"latest_image_tag"`
	RunAsNonRoot      bool `json:"run_as_non_root"`
}

type Configuration interface {
	Load(path string) error
}

func (c *AutoConfigure) Load(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// reset config settings
	temp := Config{}
	if err := json.Unmarshal(f, &temp); err != nil {
		return err
	}

	c.Lock()
	c.config = temp
	c.Unlock()

	return nil
}

func (c *AutoConfigure) Get() Config {
	c.RLock()
	cc := c.config
	c.RUnlock()

	return cc
}
