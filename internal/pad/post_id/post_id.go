// Package postid provides base-62 encoded identifiers for pads.
// IDs are derived from millisecond timestamps and formatted as "XXXXXXXX".
package postid

import (
	"fmt"
	"strings"
	"time"

	paderrors "github.com/guionardo/gs-dev/internal/errors"
)

// PostID is a base-62 encoded identifier string for a pad.
type PostID string

const (
	// PostIDCharset is the character set used for base-62 encoding:
	// digits 0-9, then lowercase a-z, then uppercase A-Z.
	PostIDCharset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// MaxPostID is the maximum value that can be encoded in 8 base-62 digits (62^8 - 1).
	MaxPostID = 62*62*62*62*62*62*62*62 - 1
)

// NewPostID encodes n as base-62 using PostIDCharset (digits 0–9, then a–z, then A–Z).
func NewPostID(n uint64) PostID {
	if n > MaxPostID {
		panic(fmt.Errorf("post ID is too large: %d", n))
	}

	base := uint64(len(PostIDCharset))
	if n == 0 {
		return PostID(PostIDCharset[:1])
	}
	// uint64 fits in at most 11 base-62 digits (62^11 > 2^64).
	var buf [11]byte

	i := len(buf)
	for n > 0 {
		i--
		buf[i] = PostIDCharset[n%base]
		n /= base
	}

	return PostID(buf[i:])
}

// NewPostIDFromNow generates a PostID from the current Unix millisecond timestamp.
func NewPostIDFromNow() PostID {
	return NewPostID(uint64(time.Now().UnixMilli()))
}

// String returns the formatted string representation "XXXX-XXXX".
func (p PostID) String() string {
	s := fmt.Sprintf("%08s", string(p))
	return s[:4] + "-" + s[4:]
}

// ParsePostID parses a formatted "XXXXXXXX" or "XXXX-XXXX" string into a PostID.
// Returns ErrPostIDInvalid if the format is incorrect.
func ParsePostID(s string) (postID PostID, err error) {
	var messageError string

	if len(s) != 9 || s[4] != '-' {
		messageError = "post ID must be 9 characters long"
	} else {
		for i, c := range s {
			if i != 4 && !strings.Contains(PostIDCharset, string(c)) {
				messageError = "post ID must contain only valid characters"
			}
		}
	}

	if messageError != "" {
		return "", paderrors.ErrPostIDInvalid
	}

	return PostID(s[:4] + s[5:]), nil
}
