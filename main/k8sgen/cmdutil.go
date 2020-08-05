package main

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RunCommand(command string) error {
	arr := strings.Split(command, " ")
	ctx, _ := context.WithTimeout(context.TODO(), 120*time.Second)

	var cmd *exec.Cmd
	if len(arr) == 0 {
		cmd = exec.CommandContext(ctx,
			arr[0],
		)
	} else {
		cmd = exec.CommandContext(ctx,
			arr[0],
			arr[1:]...,
		)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunCommandWithResult(command string) (out []byte, err error) {
	arr := strings.Split(command, " ")
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)

	var cmd *exec.Cmd
	if len(arr) == 0 {
		cmd = exec.CommandContext(ctx,
			arr[0],
		)
	} else {
		cmd = exec.CommandContext(ctx,
			arr[0],
			arr[1:]...,
		)
	}

	b := bytes.NewBuffer(nil)
	cmd.Stdout = b
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	out = b.Bytes()
	return
}

func RunCommandsAtHost(user, hostAndPort string, commands ...string) error {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".ssh/id_rsa")

	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	key, err := ssh.ParseRawPrivateKey(fileContent)
	if err != nil {
		return err
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", hostAndPort, sshConfig)
	if err != nil {
		return err
	}

	defer client.Close()

	for _, command := range commands {
		session, err := client.NewSession()
		if err != nil {
			return err
		}
		session.Stderr = os.Stderr
		session.Stdout = os.Stdout
		fmt.Println("===>运行:" + command)
		err = session.Run(command)
		if err != nil {

			continue
		}
	}

	return nil
}

func SetupEnvs(envs map[string]string) error {
	for key, value := range envs {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddKnownHostToSSH(host string) error {
	out, err := RunCommandWithResult("ssh-keyscan" + " " + host)
	if err != nil {
		return err
	}

	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".ssh/known_hosts")
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if strings.Contains(string(fileContent), host) {
		return nil
	}

	f, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	_, err = f.Write(out)
	_ = f.Close()
	return err
}

func AddPrivateGitRepo(gitlabHost string) error {
	command := fmt.Sprintf(`git config --global url.git@%s:.insteadOf https://%s/`, gitlabHost, gitlabHost)
	return RunCommand(command)
}

func GetProjectGitVersion(dir string) (string, error) {
	var out []byte
	err := os.Chdir(dir)
	if err != nil {
		return "", err
	}

	out, err = RunCommandWithResult(`git describe --tags`)
	if err == nil {
		gitTag := strings.TrimSuffix(string(out), "\n")
		if gitTag != "" {
			return gitTag, nil
		}
	}

	out, err = RunCommandWithResult(`git rev-parse --short HEAD`)
	if err != nil {
		return "", err
	}

	version := strings.TrimSuffix(string(out), "\n")
	return version, nil
}
