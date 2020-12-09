package server

import (
	"github.com/FTwOoO/micro/thirdparty/sentinel"
	"github.com/FTwOoO/util/logging"
	sentinel_api "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"net/http"
)

type HttpMiddleware func(handlerFunc http.HandlerFunc) http.HandlerFunc

func AHASMiddleware(licenseKey string, serviceName string) HttpMiddleware {
	err := sentinel.AHASInit(licenseKey, serviceName)
	if err != nil {
		logging.Log.FatalError(err)
	}

	resourceName := serviceName + "-qps-limit"
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resourceName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		logging.Log.FatalError(err)
	}

	return func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			e, err := sentinel_api.Entry(resourceName, sentinel_api.WithTrafficType(base.Inbound))
			if err != nil {
				writer.WriteHeader(http.StatusServiceUnavailable)
			} else {
				e.Exit()
				handlerFunc(writer, request)
			}
		}
	}
}
