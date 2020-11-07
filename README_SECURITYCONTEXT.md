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

## Release Timeline

* v1.10: Introduce the three new checks.
* v1.11: Make `container-security-context` optional, and make the three new checks run by default.
* v1.12: Remove `container-security-context`.
