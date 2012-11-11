package ttflib

import "fmt"
import "encoding/binary"
import "io"

type OffsetSubtable struct {
	ScalerType uint32
	NumTables uint16
	SearchRange uint16
	EntrySelector uint16
	RangeShift uint16
}

func NewOffsetSubtable(numTables uint16) (OffsetSubtable) {
	o := OffsetSubtable { 
		ScalerType:0x00010000,
		NumTables:numTables,
	}
	maxPower := uint16(1)
	log := uint16(0)
	for ; (maxPower*2) < numTables ; maxPower *= 2 { log++ }
	o.SearchRange = 16*maxPower
	o.EntrySelector = log
	o.RangeShift = 16*(numTables-maxPower)
	return o
}

func ReadOffsetSubtable(r io.Reader) (OffsetSubtable) {
	var o OffsetSubtable
	binary.Read(r, binary.BigEndian, &o)
	return o
}

func (subtable *OffsetSubtable) Write(w io.Writer) {
	binary.Write(w, binary.BigEndian, subtable)
}

func (subtable *OffsetSubtable) Verify() (error) {
	switch subtable.ScalerType {
	case 0x00010000: ; // OK for windows and mac
	case 0x74727565: ; // OK for mac
	default: return fmt.Errorf("Unrecognized scaler type %#08X",subtable.ScalerType)
	}
	maxPower := uint16(1)
	log := uint16(0)
	for ; (maxPower*2) < subtable.NumTables ; maxPower *= 2 { log++ }

	expected := OffsetSubtable {
		NumTables:subtable.NumTables,
		SearchRange:maxPower*16,
		EntrySelector:log,
		RangeShift:(subtable.NumTables-maxPower)*16,
	}

	if subtable.SearchRange != expected.SearchRange {
		return fmt.Errorf("Bad SearchRange (expected %d, got %d)",
			expected.SearchRange,
			subtable.SearchRange)
	}
	if subtable.EntrySelector != expected.EntrySelector {
		return fmt.Errorf("Bad EntrySelector (expected %d, got %d)",
			expected.EntrySelector,
			subtable.EntrySelector)
	}
	if subtable.RangeShift != expected.RangeShift {
		return fmt.Errorf("Bad RangeShift (expected %d, got %d)",
			expected.RangeShift,
			subtable.RangeShift)
	}
	return nil
}
