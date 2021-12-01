package hd

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/islishude/bip32"
	"github.com/tyler-smith/go-bip39"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestWallet_NormalPubKey(t *testing.T) {
	mnemonic := "govern truth flame pact cool cause dirt left owner behind angry arm"
	w, err := NewWalletFromMnemonic(mnemonic, "", &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(1, err)
		return
	}
	pi, err := w.Path("m/44'/0'/0'/0/0")
	if err != nil {
		fmt.Println(2, err)
		return
	}
	pubKey, err := pi.ExtendedPubKey()
	if err != nil {
		fmt.Println(4, err)
		return
	}
	fmt.Println(pubKey)
	privKey, err := pi.ExtendedPrivKey()
	if err != nil {
		fmt.Println(5, err)
		return
	}
	fmt.Println(privKey)

	privKeyImport, err := pi.PrivKeyImportable()
	if err != nil {
		fmt.Println(6, err)
		return
	}
	fmt.Println(privKeyImport)

	addr, err := pi.AddressNormal()
	if err != nil {
		fmt.Println(7, err)
		return
	}
	fmt.Println("addr1 ", addr)

	addr3, err := pi.AddressNestedSegwit()
	if err != nil {
		fmt.Println(8, err)
		return
	}
	fmt.Println("addr3 ", addr3)

	addrBC, err := pi.AddressNativeSegwit()
	if err != nil {
		fmt.Println(9, err)
		return
	}
	fmt.Println("addrBC", addrBC)
}

func TestPathWrapper_EthAddress(t *testing.T) {

	mnemonic := "govern truth flame pact cool cause dirt left owner behind angry arm"
	w, err := NewWalletFromMnemonic(mnemonic, "", &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(1, err)
		return
	}
	pi, err := w.Path("m/44'/60'/0'/0/0")
	if err != nil {
		fmt.Println(2, err)
		return
	}
	ethAddr := pi.EthAddress()

	targetAddr := "0x037D0c8Fbadff7344F7e67AE63765694F5358bFd"

	fmt.Println(ethAddr, targetAddr == ethAddr)

	privKey := pi.EthImportablePrivKey()

	targetPrivKey := "91fd7c93aad334a8d573fcdc266ddb4c5cf4071d0042f62b3b48aef84d04f334"
	fmt.Println(privKey, targetPrivKey == privKey)

}

func TestWallet_ED25519(t *testing.T) {
	mnemonic := "govern truth flame pact cool cause dirt left owner behind angry arm"
	seed := bip39.NewSeed(mnemonic, "")
	vSeed := "f78fc136f2391e6af447db4cf8dc2738db9dd6ff03577bbfb11d398da68a2ba8e6b10292de0f8db87c1bbc38b91bdb9ef374aa99f99b66b78cf097d127ce0be2"
	fmt.Println(hex.EncodeToString(seed) == vSeed)
	masterKey := []byte("ed25519 seed")
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	lr := hmac512.Sum(nil)
	//fmt.Println(lr)
	l := lr[:32]
	vl := "52566cedddfa3133dd462d6e7323ba9092d802c44eddb8b8b2572cc00270269a"
	fmt.Println("l", l)
	fmt.Println(hex.EncodeToString(l) == vl)
	r := lr[32:]
	vr := "e6d2c65eaa72a83fa83cef8aaa86618998add98013414d59a8f60cc10937aba5"
	fmt.Println(hex.EncodeToString(r) == vr)
	privKey := base58.Encode(l)
	fmt.Println(privKey)
	fmt.Println("seed", l)
	fmt.Println("1", len(privKey))
	privKey2 := ed25519.NewKeyFromSeed(l)
	fmt.Println("2", len(privKey2))
	fmt.Println("p2", privKey2)
	fmt.Println(base58.Encode(privKey2))
	fmt.Println(base58.Encode(privKey2.Public().(ed25519.PublicKey)))
	seri := make([]byte, 4)
	//binary.LittleEndian.PutUint32(seri, 100)
	binary.BigEndian.PutUint32(seri, 100)
	fmt.Println(seri)
	xprv := bip32.NewRootXPrv(l)
	fmt.Println("xprv", xprv.Bytes())
	fmt.Println(base58.Encode(xprv.Bytes()[:64]))
	fmt.Println(base58.Encode(xprv.PublicKey()))
	//fmt.Println(r)
	//fmt.Println(xprv.ChainCode())
}
