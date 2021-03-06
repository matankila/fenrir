package service

import (
	"encoding/json"
	"errors"
	"github.com/matankila/fenrir/config"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	"strings"
)

type service struct {
	Log *zap.Logger
}

type Service interface {
	Validate(req v1beta1.AdmissionReview) error
	Health() error
	GetLogger() *zap.Logger
}

var (
	defaultNs = "default"
	emptyRequest = errors.New("admission request is empty")
	restrictedNamespace = errors.New("deployment in 'default' namespace is restricted")
	noLiveness = errors.New("deployment without liveness probe is prohibited")
	noReadiness = errors.New("deployment without readiness probe is prohibited")
	noResources = errors.New("deployment without resource declared is prohibited")
	emptyResources = errors.New("deployment without empty resource declared is prohibited")
	runAsRoot = errors.New("deployment is able to run as root, please fix set pod.spec.securityContext.runAsNonRoot to true")
	latestImageTag = errors.New("container image tag is latest, this might lead to unexpected behaviours, please set it to valid version")
)

func NewService(log *zap.Logger) Service {
	return &service{
		Log: log,
	}
}

func (s service) GetLogger() *zap.Logger {
	return s.Log
}

func (s service) Validate(req v1beta1.AdmissionReview) error {
	if req.Request == nil {
		return emptyRequest
	}

	rawObj := req.Request.Object.Raw
	if err := isPodValid(rawObj); err != nil {
		return err
	}

	return nil
}

func (s service)Health() error {
	// TODO: add business logic
	return nil
}

func isPodValid(rawObj []byte) error {
	var pod v1.Pod

	// if cannot be serialized to pod, skip.
	if err := json.Unmarshal(rawObj, &pod); err != nil {
		return nil
	}

	// check ns settings
	if config.PodRestrictedNs == "true" {
		if pod.Namespace == defaultNs {
			return restrictedNamespace
		}
	}

	// check containers settings
	for _, c := range pod.Spec.Containers {
		if config.PodLivenessAndReadiness == "true" {
			if c.LivenessProbe == nil {
				return noLiveness
			} else if c.ReadinessProbe == nil {
				return noReadiness
			} else if c.Resources.Limits == nil || c.Resources.Requests == nil {
				return noResources
			} else if len(c.Resources.Limits) == 0 || len(c.Resources.Requests) == 0 {
				return emptyResources
			}
		}

		if config.PodLatestImageTag == "true" {
			s := strings.Split(c.Image, ":")
			if s[1] == "latest" {
				return latestImageTag
			}
		}

	}


	// check security context settings
	if config.PodSecurityContext == "true" {
		if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.RunAsNonRoot != nil && *(pod.Spec.SecurityContext.RunAsNonRoot) == false {
			return runAsRoot
		}
	}


	return nil
}