package bip32ed25519

import "testing"

func Test_validPath(t *testing.T) {
	validPath("m/1'/2'/3'/")
}
