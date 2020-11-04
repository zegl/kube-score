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

In future releases the `container-security-context` will become *optional*
and replaced by the more detailed checks.