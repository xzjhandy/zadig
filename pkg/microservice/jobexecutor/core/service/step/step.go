/*
Copyright 2022 The KodeRover Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package step

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/koderover/zadig/pkg/microservice/jobexecutor/config"
	"github.com/koderover/zadig/pkg/microservice/jobexecutor/core/service/cmd"
	"github.com/koderover/zadig/pkg/microservice/jobexecutor/core/service/meta"
	"github.com/koderover/zadig/pkg/tool/log"
	"github.com/koderover/zadig/pkg/util"
)

type Step interface {
	Run(ctx context.Context) error
}

func RunSteps(ctx context.Context, steps []*meta.Step, workspace, paths string, envs, secretEnvs []string) error {
	hasFailed := false
	var respErr error
	for _, stepInfo := range steps {
		if hasFailed && !stepInfo.Onfailure {
			continue
		}
		if err := runStep(ctx, stepInfo, workspace, paths, envs, secretEnvs); err != nil {
			hasFailed = true
			respErr = err
		}
	}
	return respErr
}

func runStep(ctx context.Context, step *meta.Step, workspace, paths string, envs, secretEnvs []string) error {
	var stepInstance Step
	var err error

	switch step.StepType {
	case "shell":
		stepInstance, err = NewShellStep(step.Spec, workspace, paths, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "git":
		stepInstance, err = NewGitStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "docker_build":
		stepInstance, err = NewDockerBuildStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "tools":
		stepInstance, err = NewToolInstallStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "archive":
		stepInstance, err = NewArchiveStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "junit_report":
		stepInstance, err = NewJunitReportStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "tar_archive":
		stepInstance, err = NewTararchiveStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	case "sonar_check":
		stepInstance, err = NewSonarCheckStep(step.Spec, workspace, envs, secretEnvs)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("step type: %s does not match any known type", step.StepType)
		log.Error(err)
		return err
	}
	if err := stepInstance.Run(ctx); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func prepareScriptsEnv() []string {
	scripts := []string{}
	scripts = append(scripts, "eval $(ssh-agent -s) > /dev/null")
	// $HOME/.ssh/id_rsa 为 github 私钥
	scripts = append(scripts, fmt.Sprintf("ssh-add %s/.ssh/id_rsa.github &> /dev/null", config.Home()))
	scripts = append(scripts, fmt.Sprintf("rm %s/.ssh/id_rsa.github &> /dev/null", config.Home()))
	// $HOME/.ssh/gitlab 为 gitlab 私钥
	scripts = append(scripts, fmt.Sprintf("ssh-add %s/.ssh/id_rsa.gitlab &> /dev/null", config.Home()))
	scripts = append(scripts, fmt.Sprintf("rm %s/.ssh/id_rsa.gitlab &> /dev/null", config.Home()))

	return scripts
}

func handleCmdOutput(pipe io.ReadCloser, needPersistentLog bool, logFile string, secretEnvs []string) {
	reader := bufio.NewReader(pipe)

	for {
		lineBytes, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Errorf("Failed to read log when processing cmd output: %s", err)
			break
		}

		fmt.Printf("%s", maskSecretEnvs(string(lineBytes), secretEnvs))

		if needPersistentLog {
			err := util.WriteFile(logFile, lineBytes, 0700)
			if err != nil {
				log.Warnf("Failed to write file when processing cmd output: %s", err)
			}
		}
	}
}

const (
	secretEnvMask = "********"
)

func maskSecret(secrets []string, message string) string {
	out := message

	for _, val := range secrets {
		if len(val) == 0 {
			continue
		}
		out = strings.Replace(out, val, "********", -1)
	}
	return out
}

func maskSecretEnvs(message string, secretEnvs []string) string {
	out := message

	for _, val := range secretEnvs {
		if len(val) == 0 {
			continue
		}
		sl := strings.Split(val, "=")

		if len(sl) != 2 {
			continue
		}

		if len(sl[0]) == 0 || len(sl[1]) == 0 {
			// invalid key value pair received
			continue
		}
		out = strings.Replace(out, strings.Join(sl[1:], "="), secretEnvMask, -1)
	}
	return out
}

func isDirEmpty(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer f.Close()

	_, err = f.Readdir(1)
	return err == io.EOF
}

func setCmdsWorkDir(dir string, cmds []*cmd.Command) {
	for _, c := range cmds {
		c.Cmd.Dir = dir
	}
}
