package txtarutil

import (
	"os"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"golang.org/x/tools/txtar"
)

func TestFromFS(t *testing.T) {
	a, err := FromFS(os.DirFS("testdata/"))
	if err != nil {
		t.Fatal(err)
	}

	have := string(txtar.Format(a))
	want := `-- dir1/empty-file --
-- file1 --
file1 line 1
file1 line 2
`
	d := diff.Diff(have, want)
	if d != "" {
		t.Errorf("bad archive contents\n%s", d)
	}
}
