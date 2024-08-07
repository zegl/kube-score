# Readiness and Liveness Probes

Kubernetes supports three types of probes, that all serve different purposes, and have different pros and cons.

This article is here to describe what recommendations kube-score will make in what situations.

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
* Set `interval` (default: 10s), `timeout` (default: 1s), `successThreshold` (default: 1), `failureThreshold` (default: 3) to your needs. In the default configuration, your application will fail for 30s (+ the time it takes for the network to react), for clients to stop sending traffic to your application.
* Don't depend on downstream dependencies, such as other services or databases in your probe. If the dependency has a hickup, or for example, a database is restarted, removing your Pods from the Load Balancers rotation will likely only make the downtime worse.
* Handle shutdowns gracefully: Applications should start to fail the readinessProbe after receiving `SIGTERM`, and _wait_ until the service is unregistered from all Load Balancers, before shutting down. As an alternative to this, a [`preStop`](https://kubernetes.io/docs/tasks/configure-pod-container/attach-handler-lifecycle-event/#define-poststart-and-prestop-handlers) hook with a sleep can be used. 

### `livenessProbe`

This probe roughly answers the question

> Is the _container_ healthy right now, or do we need to restart it?

It can be used to let Kubernetes know if your application is deadlocked, and needs to be restarted. Only the container with the failing probe will be restarted, other containers in the same Pod will be unaffected.

**kube-score recommends**:

* If you don't know why you need a livenessProbe, don't configure it.
* It should _never_, be the same as your `readinessProbe`.
* The livenessProbe should *never* depend on downstream dependencies, such as databases or other services.


### `startupProbe` (alpha since v1.16, beta since v1.17)

This probe roughly answers the question

> Should we start running the livenessProbe now?

For applications that take a longer time to boot than the livenessProbes `initialDelaySeconds` + `periodSeconds` * `failureThreshold`, a `startupProbe` can be configured.

The startupProbe allows you to decrease the liveness probes `initialDelaySeconds`, and catch application deadlocks earlier.

As soon as the startupProbe has succeeded once, the livenessProbe will start to be executed.

**kube-score recommends**:

* Configure a startupProbe if you have a livenessProbe configured. 

## Further reading

* [Pod Lifecycle, kubernetes.io](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)
* [Configure Liveness, Readiness and Startup Probes, kubernetes.io](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
* [Liveness probes are dangerous, srcco.de](https://srcco.de/posts/kubernetes-liveness-probes-are-dangerous.html)
