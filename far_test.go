package simsfar

import (
	"fmt"
	"os"
	"testing"
)

var fileName = "test.bmp"
var path = "test.far"

func TestParseFar(t *testing.T) {
	f := far{Path: path}
	f.ParseFar()

	if f.Signature != "FAR!byAZ" {
		t.Errorf("got %s, wanted %s", f.Signature, "FAR!byAZ")
	}
	if f.Version != 1 {
		t.Errorf("got %d, wanted %d", f.Version, 1)
	}
	if f.ManifestOffset != 160 {
		t.Errorf("got %d, wanted %d", f.ManifestOffset, 160)
	}
	if f.Manifest.NumberOfFiles != 1 {
		t.Errorf("got %d, wanted %d", f.Manifest.NumberOfFiles, 1)
	}
	if len(f.Manifest.ManifestEntries) != 1 {
		t.Errorf("got %d, wanted %d", len(f.Manifest.ManifestEntries), 1)
	}
	if f.Manifest.ManifestEntries[0].FileLength1 != 144 {
		t.Errorf("got %d, wanted %d", f.Manifest.ManifestEntries[0].FileLength1, 144)
	}
	if f.Manifest.ManifestEntries[0].FileLength2 != 144 {
		t.Errorf("got %d, wanted %d", f.Manifest.ManifestEntries[0].FileLength1, 144)
	}
	if f.Manifest.ManifestEntries[0].Filename != fileName {
		t.Errorf("got %s, wanted %s", f.Manifest.ManifestEntries[0].Filename, fileName)
	}
	if f.Manifest.ManifestEntries[0].FileNameLength != len(fileName) {
		t.Errorf("got %d, wanted %d", f.Manifest.ManifestEntries[0].FileNameLength, len(fileName))
	}
	if f.Manifest.ManifestEntries[0].FileOffset != 16 {
		t.Errorf("got %d, wanted %d", f.Manifest.ManifestEntries[0].FileOffset, 16)
	}
}

func TestFar_GetBytesByFileName(t *testing.T) {
	f := far{Path: path}
	f.ParseFar()
	b := f.GetBytesByFileName(fileName)
	fmt.Println(b)

	p, err := os.Open(path)
	check(err)
	defer p.Close()
	p.Seek(16, 0)
	oneByte := make([]byte, 1)
	for i := 0; i < 144; i++ {
		p.Read(oneByte)
		if oneByte[0] != b[i] {
			t.Errorf("got %d, wanted %d", b[i], oneByte)
		}
	}
}
