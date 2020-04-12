package main

import (
	"fmt"
	"path"
	"strings"
)

type cmdFunc func(string, []string)

func parse(args []string, cmds map[string]cmdFunc) (command string, cmdArgsOffset int, err error) {
	helpName := execName(args[0])

	// When executing kube-score as a kubectl plugin, default to the "score" sub-command to avoid stuttering
	// "kubectl score" is equivalent to "kubectl score score"
	if isKubectlPlugin(helpName) {
		command = "score"
		cmdArgsOffset = 1
	}

	// No command, flag, or file has been specified
	if len(args) <= cmdArgsOffset {
		err = fmt.Errorf("No command, flag or file")
		return
	}

	// If arg 1 is set and is a valid command, always use it as the command to execute, instead of the default
	if _, ok := cmds[args[1]]; ok {
		command = args[1]
		cmdArgsOffset = 2
	}
	return
}

func execName(args0 string) string {
	// Detect name of the binary
	binName := path.Base(args0)

	// If executed as a kubectl plugin, replace dash with a space
	// "kubectl-score" -> "kubectl score"
	if strings.HasPrefix(binName, "kubectl-") {
		binName = strings.Replace(binName, "kubectl-", "kubectl ", 1)
	}

	return binName
}

func isKubectlPlugin(helpName string) bool {
	return execName(helpName) == "kubectl score"
}
