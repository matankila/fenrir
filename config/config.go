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
	Pod Pod `json:"pod"`
}

type Pod struct {
	PodPolicyEnforcement     bool                            `json:"pod_policy_enforcement"`
	DefaultPodPolicySettings PodPolicySettings               `json:"default_pod_policy_settings"`
	CustomPodPolicies        map[Namespace]PodPolicySettings `json:"custom_pod_policies"`
}

// pod policy settings, if not set default is false
type PodPolicySettings struct {
	ReadinessLiveness bool `json:"readiness_liveness"`
	Resources         bool `json:"resources"`
	DefaultNs         bool `json:"default_ns"`
	LatestImageTag    bool `json:"latest_image_tag"`
	RunAsNonRoot      bool `json:"run_as_non_root"`
}

type Policy interface {
	policy()
}

func (p *PodPolicySettings) policy() {}

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
