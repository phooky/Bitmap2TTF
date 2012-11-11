package ttflib

import "testing"
import "os"


func TestLoadTtfFile(t *testing.T) {
	file, err := os.Open("FreeMono.ttf")
	if err != nil {
		t.Error(err)
	}
	var ttf TtfFile
	ttf.Read(file)
	verificationError := ttf.Verify()
	if verificationError != nil {
		t.Error(verificationError)
	}
	outfile, err := os.OpenFile("CloneMono.ttf",os.O_RDWR|os.O_CREATE|os.O_TRUNC,0666)
	if err != nil {
		t.Error(err)
	}
	ttf.Write(outfile)
}
