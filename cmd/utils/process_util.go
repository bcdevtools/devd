package utils

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"os"
	"os/exec"
	"strings"
	"time"
)

func LaunchAppWithDirectStd(appName string, args []string, envVars []string) int {
	return LaunchAppWithSetup(appName, args, func(launchCmd *exec.Cmd) {
		if len(envVars) > 0 {
			launchCmd.Env = envVars
		}
		launchCmd.Stdin = os.Stdin
		launchCmd.Stdout = os.Stdout
		launchCmd.Stderr = os.Stderr
	})
}

func LaunchAppWithSetup(appName string, args []string, setup func(launchCmd *exec.Cmd)) int {
	launchCmd := exec.Command(appName, args...)
	setup(launchCmd)
	err := launchCmd.Run()
	if err != nil {
		libutils.PrintfStdErr("problem when running process %s: %s\n", appName, err.Error())
		return 1
	}
	return 0
}

func WarnIfNotRunningUnderScreenSession(sampleSessionName string) {
	if !libutils.IsBlank(os.Getenv("STY")) {
		return
	}

	envTerm := os.Getenv("TERM")
	if envTerm == "screen" || strings.HasPrefix(envTerm, "screen.") {
		return
	}

	fmt.Printf(`
WARN: This command is marked as an Important or Long-Living or Heavy operation!
WARN: For safety purpose, it should be executed under a detach-able terminal session
WARN: with ability to resume after incident like internet/SSH session disconnected,...
WARN: It is recommended to use "screen" or "tmux" first to create a session, then run this command inside that session.
> screen -S %s # start a session
> # Ctrl A+D to detach from the session
> screen -r %s # resume a session
> screen -x %s # force resume a session
`, sampleSessionName, sampleSessionName, sampleSessionName)
	time.Sleep(15 * time.Second)
	fmt.Println()
}
