# kube-score

`kube-score` is a tool that does static code analysis of your Kubernetes object definitions.
The output is a list of recommendations of what you can improve to make your application more secure and resiliant.

## Installation

### Download

Pre-built releases can be downloaded from the [Github Releases page](https://github.com/zegl/kube-score/releases), or from [Docker Hub](https://hub.docker.com/r/zegl/kube-score/).

### Building from source

`kube-score` requires [go](https://golang.org/) in version `1.11.+` with [go modules](https://github.com/golang/go/wiki/Modules). To install `kube-score` into you local gobin path run the following commands:

```bash
go get github.com/zegl/kube-score
cd $GOPATH/src/github.com/zegl/kube-score/
GO111MODULE=on go install github.com/zegl/kube-score/cmd/kube-score
```

## Checks

* Container limits (should be set)
* Container image tag (should not be `:latest`)
* Container image pull policy (should be `Always`)
* Pod is targeted by a `NetworkPolicy`, both egress and ingress rules are recommended
* Container probes, both readiness and liveness checks should be configured, and should not be identical
* Container securityContext, run as high number user/group, do not run as root or with privileged root fs
* Stable APIs, use a stable API if available (supported: Deployments, StatefulSets, DaemonSet)

## Example output

![](https://i.imgur.com/zETNJNS.png)

## Usage in CI

`kube-score` can run in your CI/CD environment and will exit with exit code 1 if a critical error has been found.
The trigger level can be changed to warning with the `--exit-one-on-warning` argument.

The input to `kube-score` should be all applications that you deploy to the same namespace for the best result.

### Example with Helm

```bash
helm template my-app | kube-score -
```

### Example with static yamls

```bash
kube-score my-app/*.yaml
```

```bash
kube-score my-app/deployment.yaml my-app/service.yaml
```

## Configuration

```
Usage: kube-score [--flag1 --flag2] file1 file2 ...

Use "-" as filename to read from STDIN.

Usage of ./kube-score:
  -exit-one-on-warning
    	Exit with code 1 in case of warnings
  -help
    	Print help
  -ignore-container-cpu-limit
    	Disables the requirement of setting a container CPU limit
  -output-format string
    	Set to 'human' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs. (default "human")
  -threshold-ok int
    	The score threshold for treating an score as OK. Must be between 1 and 10 (inclusive). Scores graded below this threshold are WARNING or CRITICAL. (default 10)
  -threshold-warning int
    	The score threshold for treating a score as WARNING. Grades below this threshold are CRITICAL. Must be between 1 and 10 (inclusive). (default 5)
  -v	Verbose output
```
