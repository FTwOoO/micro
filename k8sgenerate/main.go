package main

import (
	"bytes"
	"github.com/FTwOoO/micro/cfg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

func generateDockerFile(dockerFileTemplatePath string, tempDockerFile string, projectName string, env string, execBuildTarget string) (err error) {

	tpl, err := template.ParseFiles(dockerFileTemplatePath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(tempDockerFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	err = tpl.Execute(f, map[string]interface{}{
		"env":             env,
		"serviceName":     projectName,
		"execBuildTarget": execBuildTarget,
	})

	_ = f.Close()
	return err
}

func main() {
	env := os.Args[1]
	projectDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	projectDir, err = filepath.Abs(projectDir)
	if err != nil {
		panic(err)
	}

	cfgPath := filepath.Join(projectDir, "main/"+env+".json")

	configFile := &cfg.ConfigurationImp{}
	err = cfg.NewConfiguration(cfgPath, configFile)
	if err != nil {
		panic(err)
	}

	projectName := configFile.Name
	relExecBuildTarget := "main/main"

	dockerFileTemplatePath := filepath.Join(projectDir, "Dockerfile_tpl")
	tempDockerfile := filepath.Join(projectDir, "Dockerfile_"+env)
	err = generateDockerFile(dockerFileTemplatePath, tempDockerfile, projectName, env, relExecBuildTarget)
	if err != nil {
		panic(err)
	}

	k8sTemplateFile := filepath.Join(projectDir, "k8s_tpl")
	f, err := os.OpenFile(k8sTemplateFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	tempK8sFile := filepath.Join(projectDir, "k8s_"+env+".yml")
	arr := strings.Split(configFile.GetHttp().Addr, ":")
	port, _ := strconv.ParseUint(arr[1], 10, 64)
	out := templateText(string(fileContent), map[string]interface{}{
		"domain":       configFile.GetHttp().Route.Host,
		"path":         configFile.GetHttp().Route.PathPrefix[0],
		"serviceName":  configFile.Name,
		"instancePort": port,
		"dockerSecret": "docker-github",
		"image":        `{{ .image }}`,
	})
	f, _ = os.Create(tempK8sFile)
	_, err = f.WriteString(out + "\n")
	if err != nil {
		panic(err)
	}
	f.Close()

}

func templateText(tpl string, tplArgv map[string]interface{}) string {
	tpl1, err := template.New("x").Parse(tpl)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	_ = tpl1.Execute(&b, tplArgv)
	return b.String()
}
