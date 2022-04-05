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

type configuration struct {
	AddAllDefaultChecks            bool     `yaml:"addAllDefaultChecks"`
	AddAllOptionalChecks           bool     `yaml:"addAllOptionalChecks"`
	DisableIgnoreChecksAnnotations bool     `yaml:"disableIgnoreChecksAnnotations"`
	DefaultChecks                  []string `yaml:"defaultChecks"`
	OptionalChecks                 []string `yaml:"optionalChecks"`
	IncludeChecks                  []string `yaml:"include"`
	ExcludeChecks                  []string `yaml:"exclude"`
}

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

	allChecks := score.RegisterAllChecks(parser.Empty(), config.Configuration{})

	var checks configuration

	checks.AddAllDefaultChecks = true
	checks.AddAllOptionalChecks = false
	checks.DisableIgnoreChecksAnnotations = false

	for _, c := range allChecks.All() {
		if c.Optional {
			checks.OptionalChecks = append(checks.OptionalChecks, c.ID)
		} else {
			checks.DefaultChecks = append(checks.DefaultChecks, c.ID)
		}
	}

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

func includeChecks(k *configuration) (checks []string) {
	if k.AddAllOptionalChecks {
		checks = append(checks, k.OptionalChecks...)
	}
	if len(k.IncludeChecks) > 0 {
		checks = append(checks, k.IncludeChecks...)
	}
	return
}

func excludeChecks(k *configuration) (checks []string) {
	if !k.AddAllDefaultChecks {
		checks = append(checks, k.DefaultChecks...)
	}
	if len(k.ExcludeChecks) > 0 {
		checks = append(checks, k.ExcludeChecks...)
	}
	return
}
