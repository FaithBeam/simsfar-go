package simsfar

import (
	"encoding/binary"
	"os"
)

type manifestEntry struct {
	FileLength1 int
	FileLength2 int
	FileOffset int64
	FileNameLength int
	Filename string
}

type manifest struct {
	NumberOfFiles int
	ManifestEntries []manifestEntry
}

type far struct {
	Path string
	Signature string
	Version int
	ManifestOffset int64
	Files    [][]byte
	Manifest manifest
}

func (f *far) ParseFar()  {
	b := make([]byte, 8)
	p, err := os.Open(f.Path)
	check(err)
	defer p.Close()

	_, err = p.Read(b)
	check(err)
	f.Signature = string(b)

	b = make([]byte, 4)

	_, err = p.Read(b)
	check(err)
	f.Version = int(binary.LittleEndian.Uint32(b))

	_, err = p.Read(b)
	check(err)
	f.ManifestOffset = int64(binary.LittleEndian.Uint32(b))

	p.Seek(f.ManifestOffset, 0)
	_, err = p.Read(b)
	check(err)
	f.Manifest = manifest{
		NumberOfFiles:   int(binary.LittleEndian.Uint32(b)),
		ManifestEntries: make([]manifestEntry, 0),
	}

	for i := 0; i < f.Manifest.NumberOfFiles; i++ {
		f.Manifest.ManifestEntries = append(f.Manifest.ManifestEntries, ParseManifestEntry(p))
	}
}

func ParseManifestEntry(p *os.File) (me manifestEntry) {
	b := make([]byte, 4)

	_, err := p.Read(b)
	check(err)
	me.FileLength1 = int(binary.LittleEndian.Uint32(b))

	_, err = p.Read(b)
	check(err)
	me.FileLength2 = int(binary.LittleEndian.Uint32(b))

	_, err = p.Read(b)
	check(err)
	me.FileOffset = int64(binary.LittleEndian.Uint32(b))

	_, err = p.Read(b)
	check(err)
	me.FileNameLength = int(binary.LittleEndian.Uint32(b))

	b = make([]byte, me.FileNameLength)
	_, err = p.Read(b)
	check(err)
	me.Filename = string(b)

	return
}

func (f *far) GetBytesByFileName(fileName string) (b []byte)  {
	var entry manifestEntry
	for _, e := range f.Manifest.ManifestEntries {
		if e.Filename == fileName {
			entry = e
		}
	}

	b = make([]byte, entry.FileLength1)
	p, err := os.Open(f.Path)
	check(err)
	defer p.Close()
	p.Seek(entry.FileOffset, 0)
	p.Read(b)
	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}