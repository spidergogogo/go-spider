package bip32ed25519

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
)

const HardenedKeyStart = 0x80000000 // 2^31

var masterKeySeed = []byte("ed25519 seed")

type PrivateKey struct {
	key       []byte //org key
	privKey   ed25519.PrivateKey
	chainCode []byte
}

func NewMasterKeyFromSeed(seed []byte) *PrivateKey {
	hmac512 := hmac.New(sha512.New, masterKeySeed)
	hmac512.Write(seed)
	I := hmac512.Sum(nil)
	IL := I[:32]
	IR := I[32:]
	key := append([]byte(nil), IL...)
	return &PrivateKey{
		key:       key,
		privKey:   ed25519.NewKeyFromSeed(key),
		chainCode: append([]byte(nil), IR...),
	}
}

func (k PrivateKey) KeyHex() string {
	return hex.EncodeToString(k.key)
}

func (k PrivateKey) PrivKeyEncode() string {
	return base58.Encode(k.privKey)
}

func (k PrivateKey) PublicKey() []byte {
	return k.privKey.Public().(ed25519.PublicKey)
}

func (k PrivateKey) PublicKeyHex() string {
	return hex.EncodeToString(k.PublicKey())
}

func (k PrivateKey) PublicKeyEncode() string {
	return base58.Encode(k.PublicKey())
}

func (k PrivateKey) ChainCode() []byte {
	return k.chainCode
}

func (k PrivateKey) ChainCodeHex() string {
	return hex.EncodeToString(k.chainCode)
}

func (k PrivateKey) Derive(index uint32) (*PrivateKey, error) {
	if index < HardenedKeyStart {
		return nil, errors.New("Just Hard Derive")
	}
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	// 0 key...(32)  indexBytes...(4)
	data := make([]byte, 37)
	copy(data[1:], k.key)
	copy(data[33:], indexBytes)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	hmac512.Write(data)
	I := hmac512.Sum(nil)
	IL := I[:32]
	IR := I[32:]
	key := append([]byte(nil), IL...)
	return &PrivateKey{
		key:       key,
		privKey:   ed25519.NewKeyFromSeed(key),
		chainCode: append([]byte(nil), IR...),
	}, nil
}

func (k PrivateKey) DerivePath(path string) (*PrivateKey, error) {
	if !validPath(path) {
		return nil, errors.Errorf("Invalid derivation path: %s", path)
	}

	var segments []uint32
	for _, segment := range strings.Split(path, "/")[1:] {
		segmentN := strings.ReplaceAll(segment, "'", "")
		i64, err := strconv.ParseUint(segmentN, 10, 32)
		if err != nil {
			return nil, errors.Errorf("segment error: %s", segment)
		}
		segments = append(segments, uint32(i64))
	}
	childKey := &k
	var err error
	for _, segment := range segments {
		childKey, err = childKey.Derive(HardenedKeyStart + segment)
		if err != nil {
			return nil, errors.Errorf("Devire error:%s", err)
		}
	}
	return childKey, nil
}
