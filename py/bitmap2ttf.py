#!/usr/bin/python
from struct import pack, unpack, calcsize

class VerificationError(StandardError):
    pass

class OffsetSubtable:
    fmt = ">IHHHH"
    def read(self, f):
        data = f.read(calcsize(OffsetSubtable.fmt))
        self.scalerType, self.numTables, self.searchRange, self.entrySelector, self.rangeShift = unpack(OffsetSubtable.fmt,data)

    SCALER_T_WIN = 0x00010000
    SCALER_T_MAC = 0x74727565

    def verify(self):
        if self.scalerType != OffsetSubtable.SCALER_T_WIN and self.scalerType != OffsetSubtable.SCALER_T_MAC:
            raise VerificationError("Bad scaler type {0:#x}".format(self.scalerType))
	maxPower = 1
	log = 0
        while maxPower*2 < self.numTables:
            maxPower = maxPower * 2
            log = log + 1
        if self.searchRange != maxPower*16:
            raise VerificationError("Bad search range (expected {0}, got {1})",
                                    self.searchRange,
                                    maxPower*16)
        if self.entrySelector != log:
            raise VerificationError("Bad entry selector (expected {0}, got {1})",
                                    self.entrySelector,
                                    log)
        expectedRangeShift = (self.numTables-maxPower)*16
        if self.rangeShift != expectedRangeShift:
            raise VerificationError("Bad range shift (expected {0}, got {1})",
                                    self.rangeShift,
                                    expectedRangeShift)

class TableDirEntry:
    fmt = ">IIII"
    def read(self, f):
        data = f.read(calcsize(TableDirEntry.fmt))
        self.tag, self.checksum, self.offset, self.length = unpack(TableDirEntry.fmt,data)

def strToTag(s):
    v = 0
    for c in s:
        v = v << 8
        v = v + ord(c)
    return v

def tagToStr(t):
    s = ['\0']*4
    for i in range(4):
        s[3-i] = chr(t & 0xff)
        t = t >> 8
    return "".join(s)

class TTFFile:
    def __init__(self):
        pass

    def read(self,path):
        f = open(path,'r')
        self.offset = OffsetSubtable()
        self.offset.read(f)
        self.offset.verify()
        self.tableDirectory = {}
        for i in range(self.offset.numTables):
            entry = TableDirEntry()
            entry.read(f)
            self.tableDirectory[entry.tag] = entry
	self.tables = {}
	for tag, entry in self.tableDirectory.items():
            f.seek(entry.offset)
            self.tables[tag] = f.read(entry.length)

    def getTable(self,tag):
        print "seeking {0:#x}".format(tag)
        return self.tables[tag]

    def write(self,path):
        pass

import sys

if __name__ == '__main__':
    path = sys.argv[1]
    ttf = TTFFile()
    ttf.read(path)
    fout = open('glyf.dat','w')
    fout.write(ttf.getTable(strToTag('glyf')))
    fout.close()
