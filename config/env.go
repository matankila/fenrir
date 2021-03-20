package config

import (
	"os"
)

var (
	LogLvl           = getEnv("LOG_LVL", "info")
	Port             = getEnv("PORT", "8080")
	Output           = getEnv("OUTPUT", "stdout")
	ConfigPolicyPath = getEnv("CONFIG_POLICY_PATH", "./conf.json")
)

type RequestInfo struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Ip     string `json:"ip"`
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}
