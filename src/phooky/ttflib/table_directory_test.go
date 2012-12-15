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


func TestLoadBitmapTableDir(t *testing.T) {
	file, err := os.Open("bitmap.ttf")
	if err != nil {
		t.Error(err)
	}
	t.Log("TEST")
	o := ReadOffsetSubtable(file)
	dir := ReadTableDir(file,o.NumTables)
	for key,value := range(dir) {
		keytrans := FromTag(key)

		t.Log(keytrans,value)
	}
}

func TestLoadSimpleTableDir(t *testing.T) {
	file, err := os.Open("simple.ttf")
	if err != nil {
		t.Error(err)
	}
	t.Log("TEST2")
	o := ReadOffsetSubtable(file)
	dir := ReadTableDir(file,o.NumTables)
	for key,value := range(dir) {
		keytrans := FromTag(key)
		t.Log(keytrans,value)
	}
}
