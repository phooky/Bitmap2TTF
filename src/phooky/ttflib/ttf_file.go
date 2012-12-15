package ttflib

import "os"
import "encoding/binary"
import "fmt"

func MakeTable(tag uint32) (TtfTable) {
	switch tag {
	//case NameToTag("cmap"): return new(CmapTable)
	}
	return new(BlobTable)
}

type TtfTable interface {
	FromBlob(data []byte) (uint32)
	ToBlob() ([]byte)
}

type BlobTable struct {
	data []byte
}

func (b *BlobTable) FromBlob(data []byte) (uint32) {
	b.data = data
	return uint32(len(b.data))
}

func (b *BlobTable) ToBlob() ([]byte) {
	return b.data
}

func padSize(size int) (int) {
	remainder := size % 4
	if remainder > 0 {
		return (4-remainder)
	}
	return 0
}

type TtfFile struct {
	Offset OffsetSubtable
	Directory TableDir
	Tables map[uint32]TtfTable
}

func (t *TtfFile) Read(r *os.File) {
	t.Offset = ReadOffsetSubtable(r)
	t.Directory = ReadTableDir(r,t.Offset.NumTables)
	t.Tables = make(map[uint32]TtfTable)
	for tag, entry := range t.Directory {
		var b []byte = make([]byte, entry.Length)
		r.Seek(int64(entry.Offset),0)
		r.Read(b)
		var tab TtfTable = MakeTable(tag)
		tab.FromBlob(b)
		t.Tables[tag]=tab
	}
}

func checksum(b []byte) (uint32) {
	var sum uint32
	var i int
	for i = 0; (i+3) < len(b); i += 4 {
		sum += binary.BigEndian.Uint32(b[i:i+4])
	}
	if (i < len(b)) {
		// pad out last word
		padded := make([]byte,4)
		for j := 0; i < len(b); j, i = j+1, i+1 { padded[j] = b[i] }
		sum += binary.BigEndian.Uint32(padded)
	}
	return sum
}

func NameToTag(s string) (uint32) {
	val := uint32(0)
	for _,c := range s {
		val <<= 8
		val += uint32(c)
	}
	return val
}
		
var RequiredTagNames = [...]string { 
	"cmap", "glyf", "head", "hhea", "hmtx", "loca", "maxp", "name", "post",
}

func (t *TtfFile) Verify() (error) {
	// Verify that required tables are present
	for _, tagName := range RequiredTagNames {
		tag := NameToTag(tagName)
		if _,ok := t.Tables[tag]; !ok {
			return fmt.Errorf("Required tag %s not found!",tagName)
		}

	}
	return nil
}

	

func (t *TtfFile) Write(w *os.File) {
	t.Offset.Write(w)
	offset := binary.Size(t.Offset)
	// Skip space necessary for table
	directoryOffset := offset
	directoryLen := len(t.Tables) * binary.Size(TableDirEntry{})
	// temporarily zero out. golang's docs here are sketchy.
	w.Write( make([]byte, directoryLen) )
	t.Directory = make(TableDir)
	for tag, table := range t.Tables {
		blob := table.ToBlob()
		size := len(blob)
		offset, _ := w.Seek(0,1)
		t.Directory[tag] = TableDirEntry { tag, checksum(blob), uint32(offset), uint32(size) }
		w.Write(blob)
		w.Write( make([]byte, padSize(size)) )
	}
	w.Seek(int64(directoryOffset),0)
	t.Directory.Write(w)
}