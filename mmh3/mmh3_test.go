package mmh3

import (
	"testing"
)

func TestAll(t *testing.T) {
	s := "hello"
	if Hash32(s) != 0x248bfa47 {
		t.Fail()
	}
	h1, h2 := Hash128(s)
	if h1 != 0xcbd8a7b341bd9b02 || h2 != 0x5b1e906a48ae1d19 {
		t.Fail()
	}
	s = "Winter is coming"
	if Hash32(s) != 0x43617e8f {
		t.Fail()
	}
	h1, h2 = Hash128(s)

	if h1 != 0x6c373b5d61dced95 || h2 != 0xc549d8ea0c0bfb13 {
		t.Fail()
	}
}
