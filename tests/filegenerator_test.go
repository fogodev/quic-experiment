package tests

import (
	"quic-experiment/internal"
	"strings"
	"testing"
)

func TestFileGenerator(t *testing.T) {
	buf, hash := internal.GenerateFileBytes(internal.File100KB)

	if !strings.HasPrefix(string(buf), "ABCDE") && hash == "886783f2419513e00c1c705fc08c5077" {
		t.Error("Malformed string")
	}
}
