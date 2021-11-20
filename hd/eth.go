package hd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/sha3"
)

const (
	EthHashLength = 32
)

type (
	EthHash [EthHashLength]byte
)

type Eth struct {
	privKey *btcec.PrivateKey
}

func NewEth(privKey *btcec.PrivateKey) *Eth {
	return &Eth{
		privKey: privKey,
	}
}

func (e *Eth) keccak256Hash(data ...[]byte) (h EthHash) {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

func (e *Eth) EthPrivKeyImportable() string {
	return hex.EncodeToString(e.privKey.Serialize())
}

func (e *Eth) ChecksumAddress(addr string) string {
	addrL := strings.ToLower(addr)
	if strings.HasSuffix(addrL, "0x") {
		addrL = addrL[2:]
		addr = addr[2:]
	}
	var binaryStr string
	addrBytes := []byte(addrL)
	hash256 := e.keccak256Hash(addrBytes)
	for i, l := range addrL {
		if l >= '0' && l <= '9' {
			continue
		} else {
			binaryStr = fmt.Sprintf("%08b", hash256[i/2])
			if binaryStr[4*(i%2)] == '1' {
				addrBytes[i] -= 32
			}
		}
	}
	return "0x" + string(addrBytes)
}

func (e *Eth) EthEncodeAddress() string {
	pubKey := e.privKey.PubKey()
	keccak256Hash := e.keccak256Hash(pubKey.SerializeUncompressed()[1:])
	h := hex.EncodeToString(keccak256Hash[12:])
	return e.ChecksumAddress(h)
}
