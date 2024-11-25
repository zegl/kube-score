package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mattn/go-isatty"
	flag "github.com/spf13/pflag"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/renderer/ci"
	"github.com/zegl/kube-score/renderer/human"
	"github.com/zegl/kube-score/renderer/json_v2"
	"github.com/zegl/kube-score/renderer/junit"
	"github.com/zegl/kube-score/renderer/sarif"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	"golang.org/x/term"
)

func main() {
	helpName := execName(os.Args[0])

	fs := flag.NewFlagSet(helpName, flag.ExitOnError)
	setDefault(fs, helpName, "", true)

	cmds := map[string]cmdFunc{
		"score": func(helpName string, args []string) {
			if err := scoreFiles(helpName, args); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to score files: %v\n", err)
				os.Exit(1)
			}
		},

		"list": func(helpName string, args []string) {
			if err := listChecks(helpName, args); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to list checks: %v\n", err)
				os.Exit(1)
			}
		},

		"version": func(helpName string, args []string) {
			cmdVersion()
		},

		"help": func(helpName string, args []string) {
			fs.Usage()
			os.Exit(1)
		},
	}

	command, cmdArgsOffset, err := parse(os.Args, cmds)
	if err != nil {
		fs.Usage()
		os.Exit(1)
	}

	// Execute the command, or the help command if no matching command is found
	if ex, ok := cmds[command]; ok {
		ex(helpName, os.Args[cmdArgsOffset:])
	} else {
		cmds["help"](helpName, os.Args[cmdArgsOffset:])
	}
}

func setDefault(fs *flag.FlagSet, binName, actionName string, displayForMoreInfo bool) {
	fs.Usage = func() {
		usage := fmt.Sprintf(`Usage of %s:
%s [action] --flags

Actions:
	score	Checks all files in the input, and gives them a score and recommendations
	list	Prints a CSV list of all available score checks
	version	Print the version of kube-score
	help	Print this message`+"\n\n", binName, binName)

		if displayForMoreInfo {
			usage += fmt.Sprintf(`Run "%s [action] --help" for more information about a particular command`, binName)
		}

		if len(actionName) > 0 {
			usage += "Flags for " + actionName + ":"
		}

		fmt.Println(usage)

		if len(actionName) > 0 {
			fs.PrintDefaults()
		}
	}
}

func scoreFiles(binName string, args []string) error {
	fs := flag.NewFlagSet(binName, flag.ExitOnError)
	exitOneOnWarning := fs.Bool("exit-one-on-warning", false, "Exit with code 1 in case of warnings")
	ignoreContainerCpuLimit := fs.Bool("ignore-container-cpu-limit", false, "Disables the requirement of setting a container CPU limit")
	ignoreContainerMemoryLimit := fs.Bool("ignore-container-memory-limit", false, "Disables the requirement of setting a container memory limit")
	verboseOutput := fs.CountP("verbose", "v", "Enable verbose output, can be set multiple times for increased verbosity.")
	printHelp := fs.Bool("help", false, "Print help")
	outputFormat := fs.StringP("output-format", "o", "human", "Set to 'human', 'json', 'ci' or 'sarif'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs. Sarif output allows for easier integration with CI platforms.")
	outputVersion := fs.String("output-version", "", "Changes the version of the --output-format. The 'json' format has version 'v2' (default) and 'v1' (deprecated, will be removed in v1.7.0). The 'human' and 'ci' formats has only version 'v1' (default). If not explicitly set, the default version for that particular output format will be used.")
	color := fs.String("color", "auto", "If the output should be colored. Set to 'always', 'never' or 'auto'. If set to 'auto', kube-score will try to detect if the current terminal / platform supports colors. If set to 'never', kube-score will not output any colors. If set to 'always', kube-score will output colors even if the current terminal / platform does not support colors.")
	optionalTests := fs.StringSlice("enable-optional-test", []string{}, "Enable an optional test, can be set multiple times")
	ignoreTests := fs.StringSlice("ignore-test", []string{}, "Disable a test, can be set multiple times")
	disableIgnoreChecksAnnotation := fs.Bool("disable-ignore-checks-annotations", false, "Set to true to disable the effect of the 'kube-score/ignore' annotations")
	disableOptionalChecksAnnotation := fs.Bool("disable-optional-checks-annotations", false, "Set to true to disable the effect of the 'kube-score/enable' annotations")
	allDefaultOptional := fs.Bool("all-default-optional", false, "Set to true to enable all tests")
	kubernetesVersion := fs.String("kubernetes-version", "v1.18", "Setting the kubernetes-version will affect the checks ran against the manifests. Set this to the version of Kubernetes that you're using in production for the best results.")
	minReplicasDeployment := fs.Int("min-replicas-deployment", 2, "Minimum required number of replicas for a deployment")
	minReplicasHPA := fs.Int("min-replicas-hpa", 2, "Minimum required number of replicas for a horizontal pod autoscaler")
	setDefault(fs, binName, "score", false)

	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("failed to parse files: %w", err)
	}

	if *printHelp {
		fs.Usage()
		return nil
	}

	if *outputFormat != "human" && *outputFormat != "ci" && *outputFormat != "json" && *outputFormat != "sarif" && *outputFormat != "junit" {
		fs.Usage()
		return fmt.Errorf("Error: --output-format must be set to: 'human', 'json', 'sarif', 'junit' or 'ci'")
	}

	acceptedColors := map[string]bool{
		"auto":   true,
		"always": true,
		"never":  true,
	}
	if !acceptedColors[*color] {
		fs.Usage()
		return fmt.Errorf("Error: --color must be set to: 'auto', 'always' or 'never'")
	}

	filesToRead := fs.Args()
	if len(filesToRead) == 0 {
		return fmt.Errorf(`Error: No files given as arguments.

Usage: %s score [--flag1 --flag2] file1 file2 ...

Use "-" as filename to read from STDIN.`, execName(binName))
	}

	var allFilePointers []ks.NamedReader

	for _, file := range filesToRead {
		var fp io.Reader
		var filename string

		if file == "-" {
			fp = os.Stdin
			filename = "STDIN"
		} else {
			var err error
			fp, err = os.Open(file)
			if err != nil {
				return err
			}
			filename, _ = filepath.Abs(file)
		}
		allFilePointers = append(allFilePointers, namedReader{Reader: fp, name: filename})
	}

	if len(*ignoreTests) > 0 && *allDefaultOptional {
		return errors.New("Invalid argument combination. --all-default-optional and --ignore-tests cannot be used together")
	}

	if *allDefaultOptional {
		var addOptionalChecks []string
		for _, c := range score.RegisterAllChecks(parser.Empty(), nil, nil).All() {
			if c.Optional {
				addOptionalChecks = append(addOptionalChecks, c.ID)
			}
		}
		optionalTests = &addOptionalChecks
	}

	ignoredTests := listToStructMap(ignoreTests)
	enabledOptionalTests := listToStructMap(optionalTests)

	kubeVer, err := config.ParseSemver(*kubernetesVersion)
	if err != nil {
		return errors.New("Invalid --kubernetes-version. Use on format \"vN.NN\"")
	}

	runConfig := &config.RunConfiguration{
		IgnoreContainerCpuLimitRequirement:    *ignoreContainerCpuLimit,
		IgnoreContainerMemoryLimitRequirement: *ignoreContainerMemoryLimit,
		EnabledOptionalTests:                  enabledOptionalTests,
		UseIgnoreChecksAnnotation:             !*disableIgnoreChecksAnnotation,
		UseOptionalChecksAnnotation:           !*disableOptionalChecksAnnotation,
		KubernetesVersion:                     kubeVer,
		MinReplicasDeployment:                 *minReplicasDeployment,
		MinReplicasHPA:                        *minReplicasHPA,
	}

	p, err := parser.New(&parser.Config{
		VerboseOutput: *verboseOutput,
	})
	if err != nil {
		return fmt.Errorf("failed to initializer parser: %w", err)
	}

	parsedFiles, err := p.ParseFiles(allFilePointers)
	if err != nil {
		return fmt.Errorf("failed to parse files: %w", err)
	}

	checks := score.RegisterAllChecks(parsedFiles, &checks.Config{IgnoredTests: ignoredTests}, runConfig)

	scoreCard, err := score.Score(parsedFiles, checks, runConfig)
	if err != nil {
		return err
	}

	var exitCode int
	switch {
	case scoreCard.AnyBelowOrEqualToGrade(scorecard.GradeCritical):
		exitCode = 1
	case *exitOneOnWarning && scoreCard.AnyBelowOrEqualToGrade(scorecard.GradeWarning):
		exitCode = 1
	default:
		exitCode = 0
	}

	var r io.Reader

	version := getOutputVersion(*outputVersion, *outputFormat)

	switch {
	case *outputFormat == "json" && version == "v1":
		d, _ := json.MarshalIndent(scoreCard, "", "    ")
		w := bytes.NewBufferString("")
		w.WriteString(string(d))
		r = w
	case *outputFormat == "json" && version == "v2":
		r = json_v2.Output(scoreCard)
	case *outputFormat == "human" && version == "v1":
		termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
		// Assume a width of 80 if it can't be detected
		if err != nil {
			termWidth = 80
		}
		r, err = human.Human(scoreCard, *verboseOutput, termWidth, useColor(*color))
		if err != nil {
			return err
		}
	case *outputFormat == "ci" && version == "v1":
		r = ci.CI(scoreCard)
	case *outputFormat == "sarif":
		r = sarif.Output(scoreCard)
	case *outputFormat == "junit":
		r = junit.JUnit(scoreCard)
	default:
		return fmt.Errorf("error: Unknown --output-format or --output-version")
	}

	output, _ := io.ReadAll(r)
	fmt.Print(string(output))
	os.Exit(exitCode)
	return nil
}

func getOutputVersion(flagValue, format string) string {
	if len(flagValue) > 0 {
		return flagValue
	}

	switch format {
	case "json":
		return "v2"
	default:
		return "v1"
	}
}

func listChecks(binName string, args []string) error {
	fs := flag.NewFlagSet(binName, flag.ExitOnError)
	printHelp := fs.Bool("help", false, "Print help")
	setDefault(fs, binName, "list", false)
	err := fs.Parse(args)
	if err != nil {
		return nil
	}

	if *printHelp {
		fs.Usage()
		return nil
	}

	allChecks := score.RegisterAllChecks(parser.Empty(), nil, nil)

	output := csv.NewWriter(os.Stdout)
	for _, c := range allChecks.All() {
		optionalString := "default"
		if c.Optional {
			optionalString = "optional"
		}
		err := output.Write([]string{c.ID, c.TargetType, c.Comment, optionalString})
		if err != nil {
			return nil
		}
	}
	output.Flush()

	return nil
}

func listToStructMap(items *[]string) map[string]struct{} {
	structMap := make(map[string]struct{})
	for _, testID := range *items {
		structMap[testID] = struct{}{}
	}
	return structMap
}

type namedReader struct {
	io.Reader
	name string
}

func (n namedReader) Name() string {
	return n.name
}

func useColor(colorArg string) bool {
	// Respect user preference
	switch colorArg {
	case "always":
		return true
	case "never":
		return false
	}

	// If running on Github Actions, use colors
	if _, ok := os.LookupEnv("GITHUB_ACTIONS"); ok {
		return true
	}

	// If NO_COLOR is set, don't use color
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}

	// Dont use color if not a terminal
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return false
	}

	// Use colors
	return true
}
