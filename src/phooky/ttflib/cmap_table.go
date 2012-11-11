package ttflib

import "encoding/binary"

type CmapFormat0 struct {
}

type CmapTable struct {
	subtables []TtfTable
}

func (tab *CmapTable) FromBlob(b []byte) (uint32) {
	
}

func (tab *CmapTable) ToBlob() (b []byte) {
	return nil
}