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

type kubescorechecks struct {
	AddAllDefaultChecks            bool     `yaml:"addAllDefaultChecks"`
	AddAllOptionalChecks           bool     `yaml:"addAllOptionalChecks"`
	DisableIgnoreChecksAnnotations bool     `yaml:"disableIgnoreChecksAnnotations"`
	DefaultChecks                  []string `yaml:"defaultChecks"`
	OptionalChecks                 []string `yaml:"optionalChecks"`
	IncludeChecks                  []string `yaml:"include"`
	ExcludeChecks                  []string `yaml:"exclude"`
}

func mkConfigFile(binName string, args []string) {
	fs := flag.NewFlagSet(binName, flag.ExitOnError)
	printHelp := fs.Bool("help", false, "Print help")
	setDefault(fs, binName, "mkconfig", false)
	cfgFile := fs.String("config", ".kube-score.yml", "Optional kube-score configuration file")
	cfgForce := fs.Bool("force", false, "Force overwrite of existing .kube-score.yml file")

	err := fs.Parse(args)

	if err != nil {
		panic("Failed to parse mkconfig arguments")
	}

	if *printHelp {
		fs.Usage()
		return
	}

	if _, err := os.Stat(*cfgFile); err == nil {
		if !*cfgForce {
			errmsg := fmt.Errorf("File %s exists. Use --force flag to overwrite\n", *cfgFile)
			fmt.Println(errmsg)
			fs.Usage()
			return
		}
	}

	allChecks := score.RegisterAllChecks(parser.Empty(), config.Configuration{})

	var checks kubescorechecks

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
		err := fmt.Errorf("Failed to marshal checks")
		fmt.Println(err.Error())
	} else {
		if err := os.WriteFile(*cfgFile, []byte(o), 0600); err != nil {
			panic(err)
		}
		fmt.Println("Created kube-score configuration file ", *cfgFile)
	}
}

func loadConfigFile(fp string) (config kubescorechecks) {

	content, err := os.ReadFile(fp)

	// if the file does not exist, create it
	if err != nil {
		mkConfigFile("mkconfig", []string{fp})
	}

	err2 := yaml.Unmarshal(content, &config)
	if err2 != nil {
		panic(err2)
	}

	return config
}

func includeChecks(k *kubescorechecks) (checks []string) {
	if k.AddAllOptionalChecks {
		checks = append(checks, k.OptionalChecks...)
	}
	if len(k.IncludeChecks) > 0 {
		checks = append(checks, k.IncludeChecks...)
	}
	return
}

func excludeChecks(k *kubescorechecks) (checks []string) {
	if !k.AddAllDefaultChecks {
		checks = append(checks, k.DefaultChecks...)
	}
	if len(k.ExcludeChecks) > 0 {
		checks = append(checks, k.ExcludeChecks...)
	}
	return
}
