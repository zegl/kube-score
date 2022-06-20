# kube-score

<p align="center"><img src="https://user-images.githubusercontent.com/47952/56085330-6c0a2480-5e41-11e9-89ba-0cfddd7714a8.png" height="200"></p>

[![Go Report Card](https://goreportcard.com/badge/github.com/zegl/kube-score?)](https://goreportcard.com/report/github.com/zegl/kube-score)
[![Build Status](https://circleci.com/gh/zegl/kube-score.svg?style=shield)](https://app.circleci.com/pipelines/github/zegl/kube-score?branch=master)
[![Releases](https://img.shields.io/github/release-pre/zegl/kube-score.svg)](https://github.com/zegl/kube-score/releases)
![GitHub stars](https://img.shields.io/github/stars/zegl/kube-score.svg?label=github%20stars)
![Downloads](https://img.shields.io/github/downloads/zegl/kube-score/total.svg)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/zegl/kube-score/blob/master/LICENSE)

---

`kube-score` is a tool that performs static code analysis of your Kubernetes object definitions.

The output is a list of recommendations of what you can improve to make your application more secure and resilient.

You can test kube-score out in the browser with the [online demo](https://kube-score.com) ([source](https://github.com/kube-score/web)).

## Installation

kube-score is easy to install, and is available from the following sources:

| Distribution                                        | Command / Link                                                                          |
|-----------------------------------------------------|-----------------------------------------------------------------------------------------|
| Pre-built binaries for macOS, Linux, and Windows    | [GitHub releases](https://github.com/zegl/kube-score/releases)                          |
| Docker                                              | `docker pull zegl/kube-score` ([Docker Hub)](https://hub.docker.com/r/zegl/kube-score/) |
| Homebrew  (macOS and Linux)                         | `brew install kube-score`                                                |
| [Krew](https://krew.sigs.k8s.io/) (macOS and Linux) | `kubectl krew install score`                                                            |


## Checks

For a full list of checks, see [README_CHECKS.md](README_CHECKS.md).

* Container limits (should be set)
* Pod is targeted by a `NetworkPolicy`, both egress and ingress rules are recommended
* Deployments and StatefulSets should have a `PodDisruptionPolicy`
* Deployments and StatefulSets should have host PodAntiAffinity configured
* Container probes, a readiness should be configured, and should not be identical to the liveness probe. Read more in  [README_PROBES.md](README_PROBES.md).
* Container securityContext, run as high number user/group, do not run as root or with privileged root fs. Read more in [README_SECURITYCONTEXT.md](README_SECURITYCONTEXT.md).
* Stable APIs, use a stable API if available (supported: Deployments, StatefulSets, DaemonSet)

## Example output

![](https://user-images.githubusercontent.com/47952/63225706-5b90fe80-c1d3-11e9-8b9d-fad7e723afad.png)

## Usage in CI

`kube-score` can run in your CI/CD environment and will exit with exit code 1 if a critical error has been found.
The trigger level can be changed to warning with the `--exit-one-on-warning` argument.

The input to `kube-score` should be all applications that you deploy to the same namespace for the best result.

### Example with Helm

```bash
helm template my-app | kube-score score -
```

### Example with Kustomize

```bash
kustomize build . | kube-score score -
```

### Example with static YAMLs

```bash
kube-score score my-app/*.yaml
```

```bash
kube-score score my-app/deployment.yaml my-app/service.yaml
```

### Example with an existing cluster

```bash
kubectl api-resources --verbs=list --namespaced -o name \
  | xargs -n1 -I{} bash -c "kubectl get {} --all-namespaces -oyaml && echo ---" \
  | kube-score score -
```

### Example with Docker

```bash
docker run -v $(pwd):/project zegl/kube-score:latest score my-app/*.yaml
```

## Configuration

```
Usage of kube-score:
kube-score [action] --flags

Actions:
	score	Checks all files in the input, and gives them a score and recommendations
	list	Prints a CSV list of all available score checks
	version	Print the version of kube-score
	help	Print this message

Flags for score:
      --disable-ignore-checks-annotations   Set to true to disable the effect of the 'kube-score/ignore' annotations
      --disable-enable-checks-annotations   Set to true to disable the effect of the 'kube-score/enable' annotations
      --enable-test strings                 Enable a test, can be set multiple times
      --exit-one-on-warning                 Exit with code 1 in case of warnings
      --help                                Print help
      --ignore-container-cpu-limit          Disables the requirement of setting a container CPU limit
      --ignore-container-memory-limit       Disables the requirement of setting a container memory limit
      --ignore-test strings                 Disable a test, can be set multiple times
      --kubernetes-version string           Setting the kubernetes-version will affect the checks ran against the manifests. Set this to the version of Kubernetes that you're using in production for the best results. (default "v1.18")
  -o, --output-format string                Set to 'human', 'json', 'ci' or 'sarif'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs. Sarif output allows for easier integration with CI platforms. (default "human")
      --output-version string               Changes the version of the --output-format. The 'json' format has version 'v2' (default) and 'v1' (deprecated, will be removed in v1.7.0). The 'human' and 'ci' formats has only version 'v1' (default). If not explicitly set, the default version for that particular output format will be used.
  -v, --verbose count                       Enable verbose output, can be set multiple times for increased verbosity.
```

### Ignoring a test

Tests can be ignored in the whole run of the program, with the `--ignore-test` flag.

A test can also be ignored on a per-object basis, by adding the annotation `kube-score/ignore` to the object.
The value should be a comma separated string of the [test IDs](README_CHECKS.md).

Example:

Testing this object will temporarily disable the `service-type` test, which warns against using services of type NodePort.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: node-port-service-with-ignore
  namespace: foospace
  annotations:
    kube-score/ignore: service-type
spec:
  selector:
    app: my-app
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: NodePort
```

### Enabling a optional test

Optional tests can be enabled in the whole run of the program, with the `--enable-test` flag.

A test can also be enabled on a per-object basis, by adding the annotation `kube-score/enable` to the object.
The value should be a comma separated string of the [test IDs](README_CHECKS.md).

Example:

Testing this object will enable the `container-seccomp-profile` test.
Also multiple tests defined by `kube-score/ignore` are also ignored at the same.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: optional-test-manifest-deployment
  labels:
    app: optional-test-manifest
  annotations:
    kube-score/ignore: pod-networkpolicy,container-resources,container-image-pull-policy,container-security-context-privileged,container-security-context-user-group-id,container-security-context-readonlyrootfilesystem,container-ephemeral-storage-request-and-limit
    kube-score/enable: container-seccomp-profile
spec:
  replicas: 1
  selector:
    matchLabels:
      app: optional-test-manifest
  template:
    metadata:
      labels:
        app: optional-test-manifest
    spec:
      containers:
      - name: optional-test-manifest
        image: busybox:1.34
        command:
        - /bin/sh
        - -c
        - date; env; tail -f /dev/null
```

## Building from source

`kube-score` requires [Go](https://golang.org/) `1.11` or later to build. Clone this repository, and then:

```bash
# Build the project
go build github.com/zegl/kube-score/cmd/kube-score

# Run all tests
go test -v github.com/zegl/kube-score/...
```

## Contributing?

Do you want to help out? Take a look at the [Contributing Guidelines](./.github/CONTRIBUTING.md) for more info. 🤩

## Sponsors

The development of kube-score is proudly sponsored by [Sturdy](https://github.com/sturdy-dev/sturdy). 🐥

<p align="center"><a href="https://getsturdy.com/?ref=kube-score"><img src="https://getsturdy.com/img/Sturdy-Logotype-Transparent.png" height="200"></a></p>

## Made by

<a href="https://github.com/zegl/kube-score/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=zegl/kube-score" />
</a>
