package main

import (
	"bytes"
	"github.com/rexue2019/micro/cfg"
	"text/template"
)

const getUserIdFuncTpl = `
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ .serviceName }}
spec:
  rules:
    - host: {{ .httpHort }}
      http:
        paths:
          - path: {{ .httpPathPrefix }}
            backend:
              serviceName: {{ .serviceName }}
              servicePort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .serviceName }}
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: {{ .serviceName }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .serviceName }}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{ .serviceName }}
  template:
    metadata:
      labels:
        app: {{ .serviceName }}
    spec:
      containers:
        - name: {{ .serviceName }}
          image: {{ .imageUrl }}
          ports:
            - containerPort: 8080
      imagePullSecrets:
        - name: {{ .k8sDockerPullSecret }}
`

func writeWithTpl(tpl string, tplArgv map[string]interface{}) string {
	tpl1, err := template.New("x").Parse(tpl)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	_ = tpl1.Execute(&b, tplArgv)
	return b.String()
}

func generateK8sYaml(configFile string, imageUrl string) (string, error) {
	configPointer := &cfg.ConfigurationImp{}
	err := cfg.NewConfiguration(configFile, configPointer)
	if err != nil {
		panic(err)
	}

	cf := map[string]interface{}{
		"serviceName":         configPointer.Name,
		"httpHort":            configPointer.GetHttp().Route.Host,
		"httpPathPrefix":      configPointer.GetHttp().Route.PathPrefix[0],
		"imageUrl":            imageUrl,
		"k8sDockerPullSecret": K8sDOCKER_PULL_SECRET,
	}

	return writeWithTpl(getUserIdFuncTpl, cf), nil
}
