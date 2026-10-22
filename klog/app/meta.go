package app

import (
	"regexp"
	"strings"
)

// Meta holds miscellaneous information about the klog binary.
type Meta struct {

	// Specification contains the file format specification in full text.
	Specification string

	// License contains the license text.
	License string

	BuildInfo BuildInfo
}

type BuildInfo struct {
	// Version contains the release version, e.g. `v2.7`.
	Version string

	// SrcHash contains the hash of the sources that the binary was built from.
	SrcHash string

	// Date is the date at which the binary was built, e.g. `2026-10-19`.
	Date string
}

func (b BuildInfo) IsOk() bool {
	return b.Version != "v?.?" && b.SrcHash != "???????"
}

// NewMeta creates a meta info object from input strings. For invalid values, it
// replaces them with ?-placeholders instead of failing. The validation is mostly
// structural, not logical.
func NewMeta(specification string, license string, version string, srcHash string, buildDate string) Meta {
	return Meta{
		Specification: specification,
		License:       license,
		BuildInfo:     NewBuildInfo(version, srcHash, buildDate),
	}
}

func NewBuildInfo(version string, srcHash string, buildDate string) BuildInfo {
	return BuildInfo{
		Version: func() string {
			version = strings.TrimSpace(version)
			pattern := regexp.MustCompile(`^v\d+\.\d+([.-][a-zA-Z0-9.\-_]{1,100})?$`)
			if !pattern.MatchString(version) {
				return "v?.?"
			}
			return version
		}(),
		SrcHash: func() string {
			srcHash = strings.TrimSpace(srcHash)
			pattern := regexp.MustCompile(`^[a-fA-F0-9]{7,}$`)
			if !pattern.MatchString(srcHash) {
				return "???????"
			}
			return strings.ToLower(srcHash[:7])
		}(),
		Date: func() string {
			buildDate = strings.TrimSpace(buildDate)
			pattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}(T.*)?$`)
			if !pattern.MatchString(buildDate) {
				return "????-??-??"
			}
			return buildDate[:10]
		}(),
	}
}
