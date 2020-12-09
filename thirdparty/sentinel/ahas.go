package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	config2 "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	ahas "github.com/aliyun/aliyun-ahas-go-sdk"
	"github.com/aliyun/aliyun-ahas-go-sdk/config"
	"os"
)

func AHASInit(license string, appName string) error {
	_ = os.Setenv(config.LicenseEnvKey, license)
	_ = os.Setenv(config.EnvironmentEnvKey, config.DeployEnvProd)
	_ = os.Setenv(config2.AppNameEnvKey, appName)
	return ahas.InitAhasDefault()
}

func SentinelInit() error {
	conf := config2.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	return err
}
