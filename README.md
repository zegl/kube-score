# kube-score

`kube-score` is a tool that does static code analysis of your Kubernetes object definitions.
The output is a list of recommendations of what you can improbe to make your application more secure and resiliant.


## Checks

* Container limits (should be set)
* Container image tag (should not be `:latest`)
* Container image pull policy (should be `Always`)
* Pod is targeted by a `NetworkPolicy`, both egress and ingress rules are recommended.
* Container probes, both readiness andd liveness checks should be configured, and should not be identical.
