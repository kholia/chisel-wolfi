package apk

import (
	"strconv"
	"strings"
	"unicode"
)

// CompareVersions compares two apk package versions. It follows the apk
// version shape closely enough for repository selection: main numeric/letter
// components, apk suffix ordering, optional hash, and package revision.
func CompareVersions(a, b string) int {
	pa := parseVersion(a)
	pb := parseVersion(b)
	if cmp := compareMain(pa.main, pb.main); cmp != 0 {
		return cmp
	}
	if cmp := compareSuffixes(pa.suffixes, pb.suffixes); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(pa.hash, pb.hash); cmp != 0 {
		if cmp < 0 {
			return -1
		}
		return 1
	}
	return compareInt(pa.revision, pb.revision)
}

type version struct {
	main     string
	suffixes []string
	hash     string
	revision int
}

func parseVersion(s string) version {
	v := version{main: s}
	if i := strings.LastIndex(s, "-r"); i != -1 {
		if revision, err := strconv.Atoi(s[i+2:]); err == nil {
			v.revision = revision
			s = s[:i]
		}
	}
	if i := strings.IndexByte(s, '~'); i != -1 {
		v.hash = s[i+1:]
		s = s[:i]
	}
	parts := strings.Split(s, "_")
	v.main = parts[0]
	if len(parts) > 1 {
		v.suffixes = parts[1:]
	}
	return v
}

func compareMain(a, b string) int {
	ta := versionTokens(a)
	tb := versionTokens(b)
	max := len(ta)
	if len(tb) > max {
		max = len(tb)
	}
	for i := 0; i < max; i++ {
		var aa, bb string
		if i < len(ta) {
			aa = ta[i]
		}
		if i < len(tb) {
			bb = tb[i]
		}
		if aa == bb {
			continue
		}
		if aa == "" {
			return -1
		}
		if bb == "" {
			return 1
		}
		aNum := isNumber(aa)
		bNum := isNumber(bb)
		if aNum && bNum {
			if cmp := compareNumericString(aa, bb); cmp != 0 {
				return cmp
			}
			continue
		}
		if aNum != bNum {
			if aNum {
				return 1
			}
			return -1
		}
		if aa < bb {
			return -1
		}
		return 1
	}
	return 0
}

func versionTokens(s string) []string {
	var tokens []string
	for i := 0; i < len(s); {
		r := rune(s[i])
		switch {
		case unicode.IsDigit(r):
			j := i + 1
			for j < len(s) && unicode.IsDigit(rune(s[j])) {
				j++
			}
			tokens = append(tokens, s[i:j])
			i = j
		case unicode.IsLetter(r):
			j := i + 1
			for j < len(s) && unicode.IsLetter(rune(s[j])) {
				j++
			}
			tokens = append(tokens, s[i:j])
			i = j
		default:
			i++
		}
	}
	return tokens
}

var suffixOrder = map[string]int{
	"alpha": 0,
	"beta":  1,
	"pre":   2,
	"rc":    3,
	"":      4,
	"cvs":   5,
	"svn":   6,
	"git":   7,
	"hg":    8,
	"p":     9,
}

func compareSuffixes(a, b []string) int {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	for i := 0; i < max; i++ {
		var aa, bb string
		if i < len(a) {
			aa = a[i]
		}
		if i < len(b) {
			bb = b[i]
		}
		if cmp := compareSuffix(aa, bb); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func compareSuffix(a, b string) int {
	nameA, numA := splitSuffix(a)
	nameB, numB := splitSuffix(b)
	orderA, okA := suffixOrder[nameA]
	orderB, okB := suffixOrder[nameB]
	if okA && okB {
		if orderA != orderB {
			return compareInt(orderA, orderB)
		}
		return compareInt(numA, numB)
	}
	if okA != okB {
		if okA {
			return -1
		}
		return 1
	}
	if nameA != nameB {
		if nameA < nameB {
			return -1
		}
		return 1
	}
	return compareInt(numA, numB)
}

func splitSuffix(s string) (name string, number int) {
	for i, r := range s {
		if unicode.IsDigit(r) {
			n, _ := strconv.Atoi(s[i:])
			return s[:i], n
		}
	}
	return s, 0
}

func compareNumericString(a, b string) int {
	a = strings.TrimLeft(a, "0")
	b = strings.TrimLeft(b, "0")
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	if len(a) != len(b) {
		return compareInt(len(a), len(b))
	}
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return s != ""
}
