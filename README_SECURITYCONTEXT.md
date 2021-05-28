# Security Context

The default `container-security-context` check checks the `SecurityContext`
for 

* Containers with writeable root filesystems
* Containers that run with user ID or group ID < 10000
* Privileged containers

If you do not want all of this checks you can disable `container-security-context`
and enable one or more of the following optional checks:

* `container-security-context-user-group-id`
* `container-security-context-privileged`
* `container-security-context-readonlyrootfilesystem`

## Removal timeline of `container-security-context`

`container-security-context` has been deprecated (see [#204](https://github.com/zegl/kube-score/pull/204), [#325](https://github.com/zegl/kube-score/pull/325), [#326](https://github.com/zegl/kube-score/pull/326)).

The checks that has container-security-context preformed has been split into three different checks, which where all introduced in v1.10.

* v1.10: Introduce the three new checks (opt-in), and officially deprecate `container-security-context`.
* v1.12: Make `container-security-context` optional (opt-in), and make the three new checks run by default.
* v1.13: Remove `container-security-context`.

In v1.10, run kube-score with the following flags to ensure compatability with v1.12 and later:

```bash
kube-score score \
    --enable-optional-test container-security-context-user-group-id \
    --enable-optional-test container-security-context-privileged \
    --enable-optional-test container-security-context-readonlyrootfilesystem \
    --ignore-test container-security-context
```

----

_Note:_ The "flip" and the deletion of the tests where originally scheduled to happen in v1.11 and v1.12. This did not happend, and the migration is now scheduled for v1.12 and v1.13 instead.
