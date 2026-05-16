package postid_test

import (
	"testing"
	"time"

	postid "github.com/guionardo/gs-dev/internal/pad/post_id"
	"github.com/stretchr/testify/assert"
)

func TestNewPostID(t *testing.T) {
	t.Parallel()

	got := postid.NewPostID(1)
	assert.Equal(t, "0000-0001", got.String())

	got = postid.NewPostID(62)
	assert.Equal(t, "0000-0010", got.String())

	got = postid.NewPostID(62 * 62)
	assert.Equal(t, "0000-0100", got.String())

	got = postid.NewPostID(62 * 62 * 62)
	assert.Equal(t, "0000-1000", got.String())

	got = postid.NewPostID(62*62*62*62*62*62*62*62 - 1)
	assert.Equal(t, "ZZZZ-ZZZZ", got.String())

	currentID := time.Now().UnixMilli()
	got = postid.NewPostID(uint64(currentID))
	asString := got.String()
	assert.LessOrEqualf(t, len(asString), 9, "post ID length should be less than or equal to 8 for current time %d -> %s", currentID, asString)
}
