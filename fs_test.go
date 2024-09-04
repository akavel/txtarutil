package txtarutil

import (
	"fmt"
	"os"
	"strings"
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
	if d := diff.Diff(have, want); d != "" {
		t.Errorf("bad archive contents\n%s", d)
	}
}

func TestToDirAndFromFS(t *testing.T) {
	const script = `
-- empty.txt --
-- non empty --
non empty
-- sample subdir/file 1.txt --
file1 line 1
-- sample subdir/file2 --
file2  1
file2  2
-- spec!char --
spec char in name!
`

	dir := t.TempDir()
	err := ToDir(dir, txtar.Parse([]byte(script)))
	if err != nil {
		t.Error(err)
	}

	// Check created directory tree - cheat somewhat by reusing FromFS
	a, err := FromFS(os.DirFS(dir))
	if err != nil {
		t.Fatal(err)
	}
	have := "\n" + string(txtar.Format(a))
	d := diff.Diff(have, script)
	if d != "" {
		t.Errorf("bad disk contents loaded\n%s", d)
	}
}

func TestToDir_Fail_NotValid(t *testing.T) {
	tests := []struct {
		script      string
		wantErrPath string
	}{
		{`-- /absolute/path --`, `/absolute/path`},
	}
	for _, tt := range tests {
		dir := t.TempDir()
		err := ToDir(dir, txtar.Parse([]byte(tt.script)))
		errMsg := fmt.Sprint(err)
		const wantFluff = "is not a valid local path"
		if !strings.Contains(errMsg, wantFluff) {
			t.Errorf("missing %q text in error %q for script:\n%s",
				wantFluff, errMsg, tt.script)
		}
		wantSuffix := ": " + tt.wantErrPath
		if !strings.HasSuffix(errMsg, wantSuffix) {
			t.Errorf("missing suffix %q in error %q for script:\n%s",
				wantSuffix, errMsg, tt.script)
		}
	}
}
