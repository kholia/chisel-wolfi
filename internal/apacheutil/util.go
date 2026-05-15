// SPDX-License-Identifier: Apache-2.0

package apacheutil

import (
	"fmt"
	"regexp"
)

type SliceKey struct {
	Package string
	Slice   string
}

func (s SliceKey) String() string { return s.Package + "_" + s.Slice }

// FnameExp matches the slice definition file basename. APK package names may
// contain uppercase letters, so package matching is intentionally broader than
// Debian package name matching.
var FnameExp = regexp.MustCompile(`^([A-Za-z0-9][_.+A-Za-z0-9-]*[_.+A-Za-z0-9])\.yaml$`)

// SnameExp matches only the slice name, without the leading package name.
var SnameExp = regexp.MustCompile(`^([a-z](?:-?[a-z0-9]){2,})$`)

// knameExp matches the slice full name in pkg_slice format. APK package names
// may contain underscores, so the package part is intentionally greedy and the
// slice part remains underscore-free.
var knameExp = regexp.MustCompile(`^([A-Za-z0-9](?:-?[_.A-Za-z0-9+]){1,})_([a-z](?:-?[a-z0-9]){2,})$`)

func ParseSliceKey(sliceKey string) (SliceKey, error) {
	match := knameExp.FindStringSubmatch(sliceKey)
	if match == nil {
		return SliceKey{}, fmt.Errorf("invalid slice reference: %q", sliceKey)
	}
	return SliceKey{match[1], match[2]}, nil
}
