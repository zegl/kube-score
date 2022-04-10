package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/score"
	"gopkg.in/yaml.v3"
)

// Start with an empty enable and disable list
// If enable-all is true, add all tests to the list. If it's false, add only the default tests to the list.
// If disable-all is true, add all default tests to the disable list. If it's false, do nothing.
// If enable is set, use it as the enable list.
// If disable is set, use it as the disable list.
// If --enable-optional-test is set, add the test(s) to the enable list
// If --ignore-test is set, add the test(s) to the disable list
type configuration struct {
	DisableAll    bool     `yaml:"disable-all"`
	EnableChecks  []string `yaml:"enable"`
	EnableAll     bool     `yaml:"enable-all"`
	DisableChecks []string `yaml:"disable"`
}

// Start with an empty enable and disable list
func mkConfigFile(binName string, args []string) error {
	fs := flag.NewFlagSet(binName, flag.ExitOnError)
	printHelp := fs.Bool("help", false, "Print help")
	setDefault(fs, binName, "mkconfig", false)
	cfgFile := fs.String("config", ".kube-score.yml", "Optional kube-score configuration file")
	cfgForce := fs.Bool("force", false, "Force overwrite of existing .kube-score.yml file")

	err := fs.Parse(args)

	if err != nil {
		return fmt.Errorf("Failed to parse mkconfig arguments: %w", err)
	}

	if *printHelp {
		fs.Usage()
		return nil
	}

	if _, err := os.Stat(*cfgFile); err == nil {
		if !*cfgForce {
			return fmt.Errorf("File %s exists. Use --force flag to overwrite\n", *cfgFile)
		}
	}

	var checks configuration

	checks.DisableAll = false
	checks.EnableAll = false

	if o, err := yaml.Marshal(&checks); err != nil {
		return fmt.Errorf("Failed to marshal checks %w", err)
	} else {
		if err := os.WriteFile(*cfgFile, []byte(o), 0600); err != nil {
			return fmt.Errorf("Failed to write file %w,", err)
		}
		fmt.Println("Created kube-score configuration file ", *cfgFile)
		return nil
	}
}

func loadConfigFile(fileName string) (config configuration, err error) {

	content, err := os.ReadFile(fileName)

	if err != nil {
		return config, fmt.Errorf("Failed to read file %s, error %w", fileName, err)
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		return config, fmt.Errorf("Failed to unmarshal yaml %w", err)
	}

	return config, nil
}

func registeredChecks() (requiredChecks []string, optionalChecks []string) {

	allChecks := score.RegisterAllChecks(parser.Empty(), config.Configuration{})

	for _, c := range allChecks.All() {
		if c.Optional {
			optionalChecks = append(optionalChecks, c.ID)
		} else {
			requiredChecks = append(requiredChecks, c.ID)
		}
	}

	return
}

func allRegisteredChecks() (allChecks []string) {

	registeredChecks := score.RegisterAllChecks(parser.Empty(), config.Configuration{})

	for _, c := range registeredChecks.All() {
		allChecks = append(allChecks, c.ID)
	}

	return
}

// If enable is set, use it as the enable list.
// If enable-all is true, add all tests to the list.
// By default add only the non-optional tests to the list.
func includeChecks(k *configuration) (checks []string) {

	switch {
	case len(k.EnableChecks) > 0:
		checks = append(checks, k.EnableChecks...)
		return
	case k.EnableAll:
		checks = append(checks, allRegisteredChecks()...)
		return
	}

	// default case
	defaultChecks, _ := registeredChecks()
	checks = append(checks, defaultChecks...)

	return
}

// If disable is set, use it as the disable list.
// If disable-all is true, add all tests to the disable list.
func excludeChecks(k *configuration) (checks []string) {

	switch {
	case len(k.DisableChecks) > 0:
		checks = append(checks, k.DisableChecks...)
		return
	case k.DisableAll:
		checks = append(checks, allRegisteredChecks()...)
		return
	}

	return
}
