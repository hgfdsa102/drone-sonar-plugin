package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
    "log"
)

type (
	Config struct {
		Key   string
		Name  string
		Host  string
		Token string

		Version        string
		Branch         string
		Sources        string
		Timeout        string
		Inclusions     string
		Exclusions     string
		Level          string
		ShowProfiling  string
		BranchAnalysis bool
		UsingProperties bool

        PullrequestKey string
        PullrequestBranch string
        PullrequestBase string
	}
	Plugin struct {
		Config Config
	}
)

func (p Plugin) Exec() error {
    /*
    sonar.projectName=myproject
    sonar.projectKey=myproject

    sonar.sources=src/main => .
    sonar.issuesReport.console.enable=true
    // sonar.tests=src/test
    // sonar.language=java
    sonar.java.source=1.8
    sonar.java.binaries=target/classes
    sonar.junit.reportPaths=target/surefire-reports
    sonar.sourceEncoding=UTF-8

    sonar.pullrequest.key=${ghprbPullId}
    sonar.pullrequest.base=${ghprbTargetBranch}
    sonar.pullrequest.branch=${ghprbSourceBranch}
    */

	args := []string{
		"-Dsonar.host.url=" + p.Config.Host,
		"-Dsonar.login=" + p.Config.Token,
	}

	if !p.Config.UsingProperties {
		argsParameter := []string{
			"-Dsonar.projectKey=" + strings.Replace(p.Config.Key, "/", ":", -1),
			"-Dsonar.projectName=" + p.Config.Name,
			"-Dsonar.projectVersion=" + p.Config.Version,
			"-Dsonar.sources=" + p.Config.Sources,
			"-Dsonar.ws.timeout=" + p.Config.Timeout,
			"-Dsonar.inclusions=" + p.Config.Inclusions,
			"-Dsonar.exclusions=" + p.Config.Exclusions,
			"-Dsonar.log.level=" + p.Config.Level,
			"-Dsonar.showProfiling=" + p.Config.ShowProfiling,
			"-Dsonar.scm.provider=git",
            // sonar pr
            "-Dsonar.pullrequest.key=" + p.Config.PullrequestKey,
            "-Dsonar.pullrequest.branch=" + p.Config.PullrequestBranch,
            "-sonar.pullrequest.base=" + p.Config.PullrequestBase,

		}
		args = append(args, argsParameter...)
	}


	if p.Config.BranchAnalysis {
		args = append(args, "-Dsonar.branch.name=" + p.Config.Branch)
	}

    log.Println("=== ARGS ===")
    for _, e := range args {
        log.Printf("\t - %v \n", e)
    }
    log.Println("===========")

	cmd := exec.Command("sonar-scanner", args...)
	// fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("==> Code Analysis Result:\n")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
