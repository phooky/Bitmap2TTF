package ttflib

import "testing"
import "os"

func TestLoadTableDir(t *testing.T) {
	file, err := os.Open("FreeMono.ttf")
	if err != nil {
		t.Error(err)
	}
	t.Log("TEST")
	o := ReadOffsetSubtable(file)
	ReadTableDir(file,o.NumTables)
}
