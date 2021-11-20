package hd

import (
	"fmt"
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
