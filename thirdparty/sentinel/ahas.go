package sentinel

import (
	config2 "github.com/alibaba/sentinel-golang/core/config"
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
