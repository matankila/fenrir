package config

import (
	"os"
)

var (
	LogLvl = getEnv("LOG_LVL", "info")
	Port   = getEnv("PORT", "8080")
	Output = getEnv("OUTPUT", "stdout")
	LoggerName = getEnv("LOGGER_NAME", "github.com.matankila.fenrir.logger")
	PodLivenessAndReadiness = getEnv("POD_LIVENESS_READINESS_CHECK", "true")
	PodRestrictedNs = getEnv("POD_RESTRICTED_NS_CHECK", "true")
	PodSecurityContext = getEnv("POD_RUN_AS_NON_ROOT_CHECK", "false")
	PodLatestImageTag = getEnv("POD_LATEST_IMAGE_TAG_CHECK", "false")
)

type RequestInfo struct {
	Method string
	Url string
	Ip string
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}