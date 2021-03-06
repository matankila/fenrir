package service

import (
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/validation"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
)

type service struct {
	Log *zap.Logger
}

type Service interface {
	Validate(req v1beta1.AdmissionReview) error
	Health() error
}

func NewService() Service {
	return &service{}
}

func (s service) Validate(req v1beta1.AdmissionReview) error {
	if req.Request == nil {
		return config.EmptyRequest
	}

	rawObj := req.Request.Object.Raw
	validation.Init()
	v, ok := validation.Get(req.Kind)
	if ok {
		if err := v.IsValid(rawObj, req.Request.Namespace); err != nil {
			return err
		}
	}

	return nil
}

func (s service) Health() error {
	return nil
}
