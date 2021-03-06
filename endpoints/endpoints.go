package endpoints

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matankila/fenrir/config"
	"github.com/matankila/fenrir/service"
	"go.uber.org/zap"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Endpoint func(*fiber.Ctx) error

type Endpoints struct {
	Health   Endpoint
	Validate Endpoint
}

type ValidationResp struct {
	Message string `json:"message"`
}

type HealthResp struct {
	Ok bool `json:"ok"`
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Health:   makeHealthEndpoint(s),
		Validate: makeValidationEndpoint(s),
	}
}

var (
	respValid = v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Allowed: true,
		},
	}
	pool = sync.Pool{
		New: func() interface{} {
			return &HealthResp{Ok: true}
		},
	}
	errPool = sync.Pool{
		New: func() interface{} {
			return &v1beta1.AdmissionResponse{
				Allowed:          false,
				Result:           &metav1.Status{},
			}
		},
	}
)

func makeHealthEndpoint(s service.Service) Endpoint {
	return func(c *fiber.Ctx) error {
		resp := pool.Get().(*HealthResp)
		defer pool.Put(resp)
		resp.Ok = true

		err := s.Health()
		if err != nil {
			resp.Ok = false
			return c.Status(http.StatusInternalServerError).JSON(resp)
		}

		return c.JSON(resp)
	}
}

func makeValidationEndpoint(s service.Service) Endpoint {
	return func(c *fiber.Ctx) error {
		var request v1beta1.AdmissionReview

		if err := c.BodyParser(&request); err != nil {
			return err
		}

		err := s.Validate(request)
		if err != nil {
			var req types.UID
			r := config.RequestInfo{
				Method: c.Method(),
				Url:    string(c.Request().RequestURI()),
				Ip:     c.IP(),
			}

			if request.Request != nil {
				req = request.Request.UID
			}

			s.GetLogger().Error(err.Error(),
				zap.String("uid", c.Get(fiber.HeaderXRequestID)),
				zap.Any("requestInfo", r))
			res := errPool.Get().(*v1beta1.AdmissionResponse)
			defer errPool.Put(res)
			res.UID = req
			res.Result.Message = err.Error()
			return c.Status(http.StatusBadRequest).JSON(res)
		}

		return c.JSON(respValid)
	}
}

