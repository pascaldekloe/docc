// Package repo provides repository management.
package repo // import "docc.io/source/repo"

import (
	"errors"
	"unicode"
)

// NameSegSizeMax is the maximum number of UTF-8 bytes for a Name segment.
const NameSegSizeMax = 256

// Name is a repository identifier.
//
//	name :≡ <host> "/" <path>
//	path :≡ <seg> o̅r̅ <seg> "/" <path>
//	seg  :≡ <safe> <more>
//	safe :≡ <unicode letter> o̅r̅ <unicode number> o̅r̅ "_"
//	more :≡ <char> o̅r̅ <char> <more>
//	char :≡ <safe> o̅r̅ "-"  o̅r̅ "."
type Name string

var (
	errNameSegNone  = errors.New("name segment empty")
	errNameSegStart = errors.New("illegal start character in name segment")
	errNameSegChar  = errors.New("illegal character in name segment")
	errNameSegSize  = errors.New("name segment too long")
)

// ValidNameSeg validates a Name segment.
func ValidNameSeg(s string) error {
	if s == "" {
		return errNameSegNone
	}
	if len(s) > NameSegSizeMax {
		return errNameSegSize
	}

	for i, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			continue
		}

		switch r {
		case '_':
			continue
		case '.', '-':
			if i == 0 {
				return errNameSegStart
			}
		default:
			return errNameSegChar
		}
	}

	return nil
}
