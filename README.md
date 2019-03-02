# kube-score ― Kubernetes object analysis

![](https://i.imgur.com/BTtsu06.jpg)

`kube-score` is a tool that does static code analysis of your Kubernetes object definitions.
The output is a list of recommendations of what you can improve to make your application more secure and resiliant.

## Installation

### Download

Pre-built releases can be downloaded from the [Github Releases page](https://github.com/zegl/kube-score/releases), or from [Docker Hub](https://hub.docker.com/r/zegl/kube-score/).

### Building from source

`kube-score` requires [Go](https://golang.org/) 1.11 or later to build . Run `go install github.com/zegl/kube-score/cmd/kube-score` and the latest version will be installed automatically. You can also clone the repo, and run `go build github.com/zegl/kube-score/cmd/kube-score` from the checked out directory.

## Checks

For a full list of checks, see [README_CHECKS.md](README_CHECKS.md).

* Container limits (should be set)
* Pod is targeted by a `NetworkPolicy`, both egress and ingress rules are recommended
* Deployments and StatefulSets should have a `PodDisruptionPolicy`
* Deployments and StatefulSets should have host PodAntiAffinity configured
* Container probes, both readiness and liveness checks should be configured, and should not be identical
* Container securityContext, run as high number user/group, do not run as root or with privileged root fs
* Stable APIs, use a stable API if available (supported: Deployments, StatefulSets, DaemonSet)

## Example output

![](https://i.imgur.com/zETNJNS.png)

## Usage in CI

`kube-score` can run in your CI/CD environment and will exit with exit code 1 if a critical error has been found.
The trigger level can be changed to warning with the `--exit-one-on-warning` argument.

The input to `kube-score` should be all applications/objects that you deploy to the same namespace for the best result.

### Example with Helm

```bash
helm template my-app | kube-score score -
```

### Example with static yamls

```bash
kube-score score my-app/*.yaml
```

```bash
kube-score score my-app/deployment.yaml my-app/service.yaml
```

## Configuration

```
Usage of kube-score:
kube-score [action] --flags

Actions:
	score 	Checks all files in the input, and gives them a score and recommendations
	list	Prints a cvs list of all available score checks

Flags for score:
      --exit-one-on-warning          Exit with code 1 in case of warnings
      --help                         Print help
      --ignore-container-cpu-limit   Disables the requirement of setting a container CPU limit
      --ignore-test strings          Disable a test, can be set multiple times
      --output-format string         Set to 'human' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs. (default "human")
      --threshold-ok int             The score threshold for treating an score as OK. Must be between 1 and 10 (inclusive). Scores graded below this threshold are WARNING or CRITICAL. (default 10)
      --threshold-warning int        The score threshold for treating a score as WARNING. Grades below this threshold are CRITICAL. Must be between 1 and 10 (inclusive). (default 5)
      --v                            Verbose output
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
