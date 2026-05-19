package git

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const gitStatsOutput = `commit 576d54578837644118341a14f4fcb140dfe18642 (HEAD -> feature/v0.4)
Author: Guionardo Furlan <guionardo@gmail.com>
Date:   2026-04-08T21:50:31-03:00

    refactor: remove and relocate code

 16 files changed, 175 insertions(+), 367 deletions(-)

commit 663545ca03c87d4a8f0685cad84c14c6cf130657 (origin/feature/v0.3)
Author: Guionardo Furlan <guionardo@gmail.com>
Date:   2026-04-08T18:22:16-03:00

    refactor: removed unused code on init_setup

 6 files changed, 52 insertions(+), 68 deletions(-)

commit 3658aa4f51870d23a1eb477250a29d1871b1b03c (origin/feature/v0.2)
Author: Guionardo Furlan <guionardo@gmail.com> 
Date:   2026-04-08T16:23:00-03:00

    refactor(plugin): removed unused pad and base

 5 files changed, 1 insertion(+), 179 deletions(-)
 `

func Test_parseGitStats(t *testing.T) {
	t.Parallel()

	output := []byte(gitStatsOutput)

	got := slices.Collect(parseGitStats(output))
	require.Len(t, got, 3, "expected commits")
	_ = got[0].Comments()

	assert.Equal(t, 2026, got[0].Timestamp.Year())
	assert.Equal(t, 8, got[0].Timestamp.Day())
	assert.Equal(t, 31, got[0].Timestamp.Second())
	assert.Equal(t, GitCommit{
		hash:         "576d54578837644118341a14f4fcb140dfe18642",
		meta:         "(HEAD -> feature/v0.4)",
		author:       "Guionardo Furlan <guionardo@gmail.com>",
		Timestamp:    got[0].Timestamp,
		ChangedFiles: 16,
		Insertions:   175,
		Deletions:    367,
		comments:     []string{"    refactor: remove and relocate code"},
	}, got[0])
}
