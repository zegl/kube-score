# kube-score

<p align="center"><img src="https://user-images.githubusercontent.com/47952/56085330-6c0a2480-5e41-11e9-89ba-0cfddd7714a8.png" height="200"></p>

[![Join #kube-score on the Kubernetes Slack](https://img.shields.io/badge/Slack-kubernetes%2Fkube--score-blue.svg)](https://slack.k8s.io/)

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
| Homebrew  (macOS and Linux)                         | `brew install kube-score/tap/kube-score`                                                |
| [Krew](https://krew.sigs.k8s.io/) (macOS and Linux) | `kubectl krew install score`                                                            |


## Checks

For a full list of checks, see [README_CHECKS.md](README_CHECKS.md).

* Container limits (should be set)
* Pod is targeted by a `NetworkPolicy`, both egress and ingress rules are recommended
* Deployments and StatefulSets should have a `PodDisruptionPolicy`
* Deployments and StatefulSets should have host PodAntiAffinity configured
* Container probes, a readiness should be configured, and should not be identical to the liveness probe. Read more in  [README_PROBES.md](README_PROBES.md).
* Container securityContext, run as high number user/group, do not run as root or with privileged root fs
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
docker run -v $(pwd):/project zegl/kube-score:v1.7.0 score my-app/*.yaml
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
      --enable-optional-test strings        Enable an optional test, can be set multiple times
      --exit-one-on-warning                 Exit with code 1 in case of warnings
      --help                                Print help
      --ignore-container-cpu-limit          Disables the requirement of setting a container CPU limit
      --ignore-container-memory-limit       Disables the requirement of setting a container memory limit
      --ignore-test strings                 Disable a test, can be set multiple times
  -o, --output-format string                Set to 'human', 'json' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs. (default "human")
      --output-version string               Changes the version of the --output-format. The 'json' format has version 'v1' (default) and 'v2'. The 'human' and 'ci' formats has only version 'v1' (default). If not explicitly set, the default version for that particular output format will be used.
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
