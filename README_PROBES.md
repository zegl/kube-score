# Readiness and Liveness Probes

Kubernetes supports three types of probes, that all serve different purposes, and have different pros and cons.

This article is here to describe what recommendations kube-score will make in what cituations.

## Different probe types

### `readinessProbe`

This probe roughly answers the question

> Is it a good idea to send traffic to this Pod right now?

A common *misunderstanding* is that, since Kubernetes manages the Pods, you don't need to do graceful draining of Pods during shutdown.

Without a readinessProbe you're risking that:

* Traffic is sent to the Pod before the server has started.
* Traffic is still sent to the Pod _after_ the Pod has stopped. 

**kube-score recommends**:

Every application that is targeted by a Service to:

* Setup a readinessProbe that responds with a healthy status when: The application is fully booted, and ready to serve traffic.
* Handle shutdowns gracefully: Applications should start to fail the readinessProbe after receiving `SIGTERM`, and _wait_ until the service is unregistered from all Load Balancers, before shutting down.
* Set `interval` (default: 10s), `timeout` (default: 1s), `successThreshold` (default: 1), `failureThreshold` (default: 3) to your needs. In the default configuration, your application will fail for 30s (+ the time it takes for the network to react), for clients to stop sending traffic to your application.

### `livenessProbe`

This probe roughly answers the question

> Is the _container_ healthy right now, or do we need to restart it?

**kube-score recommends**:

Only touch this configuration if you know what you're doing.

It should _never_, be the same as your `readinessProbe`.

// TODO

### `startupProbe` (alpha since v1.16, beta since v1.17)

// TODO

## Further reading

* [Pod Lifecycle, kubernetes.io](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)
* [Configure Liveness, Readiness and Startup Probes, kubernetes.io](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
* [Liveness probes are dangerous, srcco.de](https://srcco.de/posts/kubernetes-liveness-probes-are-dangerous.html)
