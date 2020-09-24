package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Vps struct {
	Host string
	Port uint16
	User string
}

var (
	GITLAB_ADDR           = "github.com"
	DOCKER_REGISTRY       = "docker.pkg.github.com"
	DOCKER_USERNAME       = "FTwOoO"
	DOCKER_PASSWORD       = "f3df9275a3931f4fc72c11f0bcaef5f1fa8e7364"
	K8sDOCKER_PULL_SECRET = "docker-github"
	GOPRIVATE             = `github.com/FTwOoO/im_grpc,github.com/FTwOoO/im_common,gitlab.livedev.shika2019.com/*,github.com/FTwOoO/*`

	HOSTS = map[string][]Vps{

		"kube-master": {
			{"120.77.2.33", 10002, "root"},
		},
		"dev":  {{"47.106.124.124", 10002, "root"}},
		"test": {{"39.108.235.76", 10002, "root"}},
		"prod": {
			{"120.77.2.33", 10002, "root"},
			{"120.24.71.80", 10002, "root"},
			{"120.24.79.207", 10002, "root"},
		},
	}
)

func buildGoProject(projectDir string, buildTarget string, gitlabHost string) error {
	err := AddKnownHostToSSH(gitlabHost)
	if err != nil {
		return err
	}

	err = AddPrivateGitRepo(gitlabHost)
	if err != nil {
		return err
	}

	err = SetupEnvs(map[string]string{
		"GOPROXY":     "goproxy.io,direct",
		"GOPRIVATE":   GOPRIVATE,
		"CGO_ENABLED": "0",
		"GO111MODULE": "on",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	})
	if err != nil {
		return err
	}

	absPath := filepath.Join(projectDir, "main/main.go")
	if strings.Contains(absPath, " ") {
		absPath = `'` + absPath + `'`
	}

	buildCommand := fmt.Sprintf(`go build -o %s %s`, buildTarget, absPath)
	return RunCommand(buildCommand)
}

type ProjectType string

const (
	ProjectTypeGo     ProjectType = "go"
	ProjectTypePython ProjectType = "python"
)

func getProjectTypeForDirectory(dir string) (projectType ProjectType, err error) {
	path, _ := filepath.Abs(filepath.Join(dir, "main/main.go"))

	isExists, err := PathExists(path)
	if err != nil {
		return "", err
	}

	if isExists {
		return ProjectTypeGo, nil
	} else {
		return ProjectTypePython, nil
	}
}

func generateDockerFile(dockerFileTemplatePath string, tempDockerFile string, projectName string, env string, execBuildTarget string) (err error) {

	tpl, err := template.ParseFiles(dockerFileTemplatePath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(tempDockerFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	err = tpl.Execute(f, map[string]interface{}{
		"env":             env,
		"projectName":     projectName,
		"execBuildTarget": execBuildTarget,
	})

	_ = f.Close()
	return err
}

func buildAndPushDockerImage(dockerFilePath string, dockerImageTarget string) error {

	homeDir, _ := os.UserHomeDir()
	dockerConfig := filepath.Join(homeDir, ".docker/config.json")
	exists, err := PathExists(dockerConfig)
	if err != nil {
		return err
	}
	var needCreateDockerConfig bool = true
	if exists {
		needCreateDockerConfig = false
	}

	if needCreateDockerConfig {
		dir, _ := filepath.Split(dockerConfig)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", DOCKER_USERNAME, DOCKER_PASSWORD)))
		authContent := fmt.Sprintf(`{"auths":{"https://%s":{"auth":"%s"}}}`, DOCKER_REGISTRY, auth)
		f, err := os.OpenFile(dockerConfig, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write([]byte(authContent))
		if err != nil {
			return err
		}
	}

	build_image_cmd := fmt.Sprintf(`docker build -f %s --tag %s .`, dockerFilePath, dockerImageTarget)
	push_image_cmd := fmt.Sprintf("docker push %s", dockerImageTarget)

	dir, _ := filepath.Split(dockerFilePath)
	err = os.Chdir(dir)
	if err != nil {
		return err
	}

	err = RunCommand(build_image_cmd)
	if err != nil {
		return err
	}

	err = RunCommand(push_image_cmd)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	/*	env := os.Args[1]
		projectDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}*/

	env := "prod"
	projectDir := "/Users/ganxiangle/Desktop/workspace/k8sdemo"

	_, projectName := filepath.Split(projectDir)

	relExecBuildTarget := fmt.Sprintf("main/%s.linux", projectName)
	execBuildTarget := filepath.Join(projectDir, relExecBuildTarget)

	projectType, err := getProjectTypeForDirectory(projectDir)
	if err != nil {
		panic(err)
	}

	if projectType == ProjectTypeGo {
		err := buildGoProject(projectDir, execBuildTarget, GITLAB_ADDR)
		if err != nil {
			panic(err)
		}
	}

	projectGitVersion, err := GetProjectGitVersion(projectDir)
	if err != nil {
		panic(err)
	}

	dockerImageTarget := fmt.Sprintf(`docker.pkg.github.com/FTwOoO/%s/%s-%s:%s`, projectName, projectName, env, projectGitVersion)
	//dockerImageTarget := fmt.Sprintf("%s/go/%s%s:%s", DOCKER_REGISTRY, projectName, env, projectGitVersion)

	dockerFileTemplatePath := filepath.Join(projectDir, "Dockerfile_tpl")
	tempDockerfile := filepath.Join(projectDir, "Dockerfile_"+env)
	err = generateDockerFile(dockerFileTemplatePath, tempDockerfile, projectName, env, relExecBuildTarget)
	if err != nil {
		panic(err)
	}

	err = buildAndPushDockerImage(tempDockerfile, dockerImageTarget)
	if err != nil {
		panic(err)
	}

	_ = os.Remove(execBuildTarget)

	configFile := filepath.Join(projectDir, fmt.Sprintf("main/%s.json", env))

	yamlFile := filepath.Join(projectDir, fmt.Sprintf("k8s_%s.yml", env))
	yamlContent, err := generateK8sYaml(configFile, dockerImageTarget)
	f, err := os.Create(yamlFile)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(yamlContent)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	for _, vps := range HOSTS["kube-master"] {
		cmd := fmt.Sprintf("scp -P %d %s %s@%s:~/%s.yml", vps.Port, yamlFile, vps.User, vps.Host, projectName)
		err = RunCommand(cmd)
		if err != nil {
			panic(err)
		}

		kubeCmd := fmt.Sprintf("kubectl apply -f ~/%s.yml", projectName)
		err = RunCommandsAtHost(
			vps.User,
			fmt.Sprintf("%s:%d", vps.Host, vps.Port),
			kubeCmd)
	}

}
