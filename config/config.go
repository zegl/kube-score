package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RunConfiguration struct {
	IgnoreContainerCpuLimitRequirement    bool
	IgnoreContainerMemoryLimitRequirement bool
	EnabledOptionalTests                  map[string]struct{}
	UseIgnoreChecksAnnotation             bool
	UseOptionalChecksAnnotation           bool
	KubernetesVersion                     Semver
	MinReplicasDeployment                 int
	MinReplicasHPA                        int
}

type Semver struct {
	Major int
	Minor int
}

var errInvalidSemver = errors.New("invalid semver")

func ParseSemver(s string) (Semver, error) {
	if len(s) == 0 {
		return Semver{}, errInvalidSemver
	}
	start := 0
	if s[0] == 'v' {
		start = 1
	}

	// Separate by .
	parts := strings.Split(s[start:], ".")
	if len(parts) != 2 {
		return Semver{}, errInvalidSemver
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Semver{}, errInvalidSemver
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Semver{}, errInvalidSemver
	}

	return Semver{
		Major: major,
		Minor: minor,
	}, nil
}

func (s Semver) LessThan(other Semver) bool {
	if s.Major < other.Major {
		return true
	}
	if s.Major == other.Major && s.Minor < other.Minor {
		return true
	}
	return false
}

func (s Semver) String() string {
	return fmt.Sprintf("v%d.%d", s.Major, s.Minor)
}
