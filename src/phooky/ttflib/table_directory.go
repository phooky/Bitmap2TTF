package ttflib

import "encoding/binary"
import "io"

type TableDirEntry struct {
	Tag uint32
	CheckSum uint32
	Offset uint32
	Length uint32
}


type TableDir map[uint32]TableDirEntry

func MakeTag(t string) (uint32) {
	return uint32(t[3]) |
		uint32(t[2])<<8 |
		uint32(t[1])<<16 |
		uint32(t[0])<<24;
}

func FromTag(tag uint32) (string) {
	var runes [5]rune
	runes[0] = rune(tag >> 24) & 0xff
	runes[1] = rune(tag >> 16) & 0xff
	runes[2] = rune(tag >> 8 ) & 0xff
	runes[3] = rune(tag >> 0 ) & 0xff
	return string(runes[:])
}

func ReadTableDir(r io.Reader, numTables uint16) (TableDir) {
	t := make(TableDir)
	for i := uint16(0); i < numTables; i++ {
		var entry TableDirEntry
		binary.Read(r, binary.BigEndian, &entry)
		t[entry.Tag] = entry
	}
	return t
}

func (dir *TableDir) Write(w io.Writer) {
	for _,entry := range *dir {
		binary.Write(w, binary.BigEndian, &entry)
	}
}
