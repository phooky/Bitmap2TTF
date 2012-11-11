package ttflib

import "testing"
import "os"

func TestLoad(t *testing.T) {
	file, err := os.Open("FreeMono.ttf")
	if err != nil {
		t.Error(err)
	}
	loaded := ReadOffsetSubtable(file)
	verify := loaded.Verify()
	if verify != nil {
		t.Error(verify)
	}
}

func TestNew(t *testing.T) {
	var err error
	var o OffsetSubtable
	o = NewOffsetSubtable(0x12)
	err = o.Verify()
	if err != nil {
		t.Errorf("Unexpected error on new subtable (%s)",err.Error())
	}
	o = NewOffsetSubtable(0x10)
	err = o.Verify()
	if err != nil {
		t.Errorf("Unexpected error on new subtable (%s)",err.Error())
	}
}
	
func TestVerifyGood(t *testing.T) {
	// Good sample from FreeMono.ttf
	good := OffsetSubtable { 0x10000, 0x12, 0x100, 0x4, 0x20 }
	err := good.Verify()
	if err != nil {
		t.Errorf("Unexpected error on good subtable (%s)",err.Error())
	}
}

		