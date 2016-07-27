package main

import (
	"fmt"
	"log"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
	ARM     string
}

func PlatformFromString(os, arch string) Platform {
	if strings.HasPrefix(arch, "armv") && len(arch) >= 5 {
		return Platform{
			OS:   os,
			Arch: "arm",
			ARM:  arch[4:],
		}
	}
	return Platform{
		OS:   os,
		Arch: arch,
	}
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.GetArch())
}

func (p *Platform) GetArch() string {
	return fmt.Sprintf("%s%s", p.Arch, p.GetARMVersion())
}

func (p *Platform) GetARMVersion() string {
	if len(p.ARM) > 0 {
		return "v" + p.ARM
	}
	return ""
}

var (
	OsList = []string{
		"darwin",
		"dragonfly",
		"linux",
		"android",
		"solaris",
		"freebsd",
		"nacl",
		"netbsd",
		"openbsd",
		"plan9",
		"windows",
	}

	ArchList = []string{
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"mips64",
		"mips64le",
		"ppc64",
		"ppc64le",
	}

	Platforms_1_0 = []Platform{
		{OS: "darwin", Arch: "386", Default: true},
		{OS: "darwin", Arch: "amd64", Default: true},
		{OS: "linux", Arch: "386", Default: true},
		{OS: "linux", Arch: "amd64", Default: true},
		{OS: "linux", Arch: "arm", Default: true},
		{OS: "freebsd", Arch: "386", Default: true},
		{OS: "freebsd", Arch: "amd64", Default: true},
		{OS: "openbsd", Arch: "386", Default: true},
		{OS: "openbsd", Arch: "amd64", Default: true},
		{OS: "windows", Arch: "386", Default: true},
		{OS: "windows", Arch: "amd64", Default: true},
	}

	Platforms_1_1 = append(Platforms_1_0, []Platform{
		{OS: "freebsd", Arch: "arm", Default: true},
		{OS: "linux", Arch: "arm", Default: false, ARM: "5"},
		{OS: "linux", Arch: "arm", Default: false, ARM: "6"},
		{OS: "linux", Arch: "arm", Default: false, ARM: "7"},
		{OS: "netbsd", Arch: "386", Default: true},
		{OS: "netbsd", Arch: "amd64", Default: true},
		{OS: "netbsd", Arch: "arm", Default: true},
		{OS: "plan9", Arch: "386", Default: false},
	}...)

	Platforms_1_3 = append(Platforms_1_1, []Platform{
		{OS: "dragonfly", Arch: "386", Default: false},
		{OS: "dragonfly", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "amd64p32", Default: false},
		{OS: "nacl", Arch: "arm", Default: false},
		{OS: "solaris", Arch: "amd64", Default: false},
	}...)

	Platforms_1_4 = append(Platforms_1_3, []Platform{
		{OS: "android", Arch: "arm", Default: false},
		{OS: "plan9", Arch: "amd64", Default: false},
	}...)

	Platforms_1_5 = append(Platforms_1_4, []Platform{
		{OS: "darwin", Arch: "arm", Default: false},
		{OS: "darwin", Arch: "arm64", Default: false},
		{OS: "linux", Arch: "arm64", Default: false},
		{OS: "linux", Arch: "ppc64", Default: false},
		{OS: "linux", Arch: "ppc64le", Default: false},
	}...)

	Platforms_1_6 = append(Platforms_1_5, []Platform{
		{"android", "386", false},
		{"linux", "mips64", false},
		{"linux", "mips64le", false},
	}...)

	Platforms_1_7 = append(Platforms_1_5, []Platform{
		// While not fully supported s390x is generally useful
		{"linux", "s390x", true},
		{"plan9", "arm", false},
		// Add the 1.6 Platforms, but reflect full support for mips64 and mips64le
		{"android", "386", false},
		{"linux", "mips64", true},
		{"linux", "mips64le", true},
	}...)

	Platforms_1_8 = append(Platforms_1_7, []Platform{
		{"linux", "mips", true},
		{"linux", "mipsle", true},
	}...)

	// no new platforms in 1.9
	Platforms_1_9 = Platforms_1_8

	// no new platforms in 1.10
	Platforms_1_10 = Platforms_1_9

	PlatformsLatest = Platforms_1_10
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	// Use latest if we get an unexpected version string
	if !strings.HasPrefix(v, "go") {
		return PlatformsLatest
	}
	// go-version only cares about version numbers
	v = v[2:]

	current, err := version.NewVersion(v)
	if err != nil {
		log.Printf("Unable to parse current go version: %s\n%s", v, err.Error())

		// Default to latest
		return PlatformsLatest
	}

	var platforms = []struct {
		constraint string
		plat       []Platform
	}{
		{"<= 1.0", Platforms_1_0},
		{">= 1.1, < 1.3", Platforms_1_1},
		{">= 1.3, < 1.4", Platforms_1_3},
		{">= 1.4, < 1.5", Platforms_1_4},
		{">= 1.5, < 1.6", Platforms_1_5},
		{">= 1.6, < 1.7", Platforms_1_6},
		{">= 1.7, < 1.8", Platforms_1_7},
		{">= 1.8, < 1.9", Platforms_1_8},
		{">= 1.9, < 1.10", Platforms_1_9},
		{">=1.10, < 1.11", Platforms_1_10},
	}

	for _, p := range platforms {
		constraints, err := version.NewConstraint(p.constraint)
		if err != nil {
			panic(err)
		}
		if constraints.Check(current) {
			return p.plat
		}
	}

	// Assume latest
	return Platforms_1_9
}
