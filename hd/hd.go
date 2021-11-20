package hd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

var ErrKeyPathFormat = errors.New("wallet path error")

type HD struct {
	masterKey  *hdkeychain.ExtendedKey
	net        *chaincfg.Params
	pathKeyMap map[string]*hdkeychain.ExtendedKey
}

type purpose int

const (
	pathP2PKH        purpose = 44
	pathNestedSegwit purpose = 49
	pathNativeSegwit purpose = 84
)

// m/purpose'/coin_type'/account'/change/address_index
const pathTemplate = "m/%d'/%d'/%d'/%d/%d"

// PathP2PKH P2PKH Start 1, 44
func (h *HD) PathP2PKH(account, change, addressIdx int) (*PathWrapper, error) {
	path := fmt.Sprintf(pathTemplate, pathP2PKH, h.net.HDCoinType, account, change, addressIdx)
	return h.Path(path)
}

// PathNestedSegwit P2SH(P2WPKH) Start 3, 49
func (h *HD) PathNestedSegwit(account, change, addressIdx int) (*PathWrapper, error) {
	path := fmt.Sprintf(pathTemplate, pathNestedSegwit, h.net.HDCoinType, account, change, addressIdx)
	return h.Path(path)
}

// PathNativeSegwit P2WPKH Start bc, 84
func (h *HD) PathNativeSegwit(account, change, addressIdx int) (*PathWrapper, error) {
	path := fmt.Sprintf(pathTemplate, pathNativeSegwit, h.net.HDCoinType, account, change, addressIdx)
	return h.Path(path)
}

func validPath(path []string) error {
	if path[0] != "m" {
		return ErrKeyPathFormat
	}
	for i := 1; i < len(path); i++ {
		childNumStr := path[i]
		if strings.HasSuffix(childNumStr, "'") {
			childNumStr = strings.Replace(childNumStr, "'", "", -1)
		}
		childNum, err := strconv.Atoi(childNumStr)
		if err != nil {
			return ErrKeyPathFormat
		}
		if childNum >= hdkeychain.HardenedKeyStart || childNum < 0 {
			return ErrKeyPathFormat
		}
	}
	return nil
}

// Path create path info
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
func (h *HD) Path(path string) (*PathWrapper, error) {
	cKey, ok := h.pathKeyMap[path]
	if ok {
		pathWrapper, err := NewPathWrapper(path, cKey, h.net)
		if err != nil {
			return nil, err
		}
		return pathWrapper, nil
	}
	pathSlice := strings.Split(path, "/")
	err := validPath(pathSlice)
	if err != nil {
		return nil, err
	}
	var tmpPath []string
	var tmpPathStr string
	var tmpParentKey *hdkeychain.ExtendedKey
	for _, childNumStr := range pathSlice {
		tmpPath = append(tmpPath, childNumStr)
		tmpPathStr = strings.Join(tmpPath, "/")
		cKey, ok = h.pathKeyMap[tmpPathStr]
		if !ok {
			if tmpPathStr == "m" {
				cKey = h.masterKey
			} else {
				isHardenedChild := false
				if strings.HasSuffix(childNumStr, "'") {
					childNumStr = strings.Replace(childNumStr, "'", "", -1)
					isHardenedChild = true
				}
				childNum, _ := strconv.Atoi(childNumStr)
				var err error
				if isHardenedChild {
					childNum = hdkeychain.HardenedKeyStart + childNum
					cKey, err = tmpParentKey.Derive(uint32(childNum))
				} else {
					cKey, err = tmpParentKey.Derive(uint32(childNum))
				}
				if err != nil {
					return nil, err
				}
			}
			h.pathKeyMap[tmpPathStr] = cKey
		}
		tmpParentKey = cKey
	}
	pathWrapper, err := NewPathWrapper(path, cKey, h.net)
	if err != nil {
		return nil, err
	}
	return pathWrapper, nil
}

func NewWalletFromMnemonic(mnemonic, passwd string, net *chaincfg.Params) (*HD, error) {
	seed := bip39.NewSeed(mnemonic, passwd)
	masterKey, err := hdkeychain.NewMaster(seed, net)
	if err != nil {
		return nil, err
	}
	return &HD{
		masterKey:  masterKey,
		net:        net,
		pathKeyMap: make(map[string]*hdkeychain.ExtendedKey),
	}, nil

}

type PathWrapper struct {
	path string
	key  *hdkeychain.ExtendedKey
	net  *chaincfg.Params
	eth  *Eth
}

func NewPathWrapper(path string, key *hdkeychain.ExtendedKey, net *chaincfg.Params) (*PathWrapper, error) {
	privKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	eth := NewEth(privKey)
	return &PathWrapper{
		path: path,
		key:  key,
		net:  net,
		eth:  eth,
	}, nil
}

func (p *PathWrapper) ExtendedPubKey() (string, error) {
	key, err := p.key.Neuter()
	if err != nil {
		return "", err
	}
	return key.String(), nil
}

func (p *PathWrapper) ExtendedPrivKey() (string, error) {
	if !p.key.IsPrivate() {
		return "", errors.New("this is PubKey")
	}
	return p.key.String(), nil
}

func (p *PathWrapper) PrivKeyImportable() (string, error) {
	ecPrivKey, err := p.key.ECPrivKey()
	if err != nil {
		return "", err
	}
	wif, err := btcutil.NewWIF(ecPrivKey, p.net, true)
	if err != nil {
		return "", err
	}
	return wif.String(), nil
}

// AddressNormal P2PKH Start 1, 44
func (p *PathWrapper) AddressNormal() (string, error) {
	addr, err := p.key.Address(p.net)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

// AddressNativeSegwit P2WPKH Start bc, 84
func (p *PathWrapper) AddressNativeSegwit() (string, error) {
	pubKey, err := p.key.ECPubKey()
	if err != nil {
		return "", err
	}

	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

	p2wkhAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, p.net)
	if err != nil {
		return "", err
	}
	return p2wkhAddr.EncodeAddress(), nil
}

// AddressNestedSegwit P2SH(P2WPKH) Start 3, 49
func (p *PathWrapper) AddressNestedSegwit() (string, error) {
	pubKey, err := p.key.ECPubKey()
	if err != nil {
		return "", err
	}
	witnessAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), p.net)
	data := witnessAddr.ScriptAddress()
	dataLen := len(data)
	witnessProg := []byte{0x00, byte(dataLen)}
	witnessProg = append(witnessProg, data...)
	addr, err := btcutil.NewAddressScriptHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	return addr.EncodeAddress(), nil
}

func (p *PathWrapper) EthAddress() string {
	return p.eth.EthEncodeAddress()
}

func (p *PathWrapper) EthImportablePrivKey() string {
	return p.eth.EthPrivKeyImportable()
}
