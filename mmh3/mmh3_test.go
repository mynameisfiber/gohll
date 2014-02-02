package mmh3

import (
	"testing"
	"fmt"
)

func TestAll(t *testing.T) {
  s := []byte("hello")
  if Hash32(s) != 0x248bfa47 {
    t.Fail()
  }
  if fmt.Sprintf("%x", Hash128(s)) != "029bbd41b3a7d8cb191dae486a901e5b" {
    t.Fail()
  }
  s = []byte("Winter is coming")
  if Hash32(s) != 0x43617e8f {
    t.Fail()
  }
  if fmt.Sprintf("%x", Hash128(s)) != "95eddc615d3b376c13fb0b0cead849c5" {
    t.Fail()
  }
}
