package main

import (
	"bufio"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
	"go-spider/hd"
	"os"
	"strings"
	"time"
)

func process() {
	initCoinConfig()
	initProxyClient(cfg.Proxy)
	//processV1()
	processV2()
}

func processV2() {
	chainMap := make(map[string]chan *addrMnEntity)
	for _, coin := range coinSlice {
		chanAddrMn := make(chan *addrMnEntity)
		chainMap[coin] = chanAddrMn
		go dispatchHandler(chanAddrMn)
	}
	//for i := 0; i < 1; i++ {
	mnMap := make(map[string]struct{})
	for {
		addrMnMap := make(map[string]string)
		var addrSlice []string
		for j := 0; j < 6; j++ {
			addr, mnemonic, err := genNewAddress()
			if err != nil {
				log.Error(fmt.Sprintf("genNewAddress Error: %s", err))
				continue
			}
			if _, ok := mnMap[mnemonic]; ok {
				log.Error(fmt.Sprintf("重复的词:%s", mnemonic))
			} else {
				mnMap[mnemonic] = struct{}{}
			}
			addrMnMap[addr] = mnemonic
			addrSlice = append(addrSlice, addr)
		}
		if len(addrSlice) == 0 {
			continue
		}
		addrSlice = append(addrSlice, "0x2d4c407bbe49438ed859fe965b140dcf1aab71a9")
		for _, coin := range coinSlice {
			ame := &addrMnEntity{
				coin:      coin,
				addrMnMap: addrMnMap,
				addr:      addrSlice,
			}
			chainMap[coin] <- ame
		}
		time.Sleep(240 * time.Millisecond)
	}
	for {
	}
}

func processV1() {
	//for i := 0; i < 1; i++ {
	for {
		addrMnMap := make(map[string]string)
		var addrSlice []string
		for j := 0; j < 6; j++ {
			addr, mnemonic, err := genNewAddress()
			if err != nil {
				log.Error(fmt.Sprintf("genNewAddress Error: %s", err))
				continue
			}
			addrMnMap[addr] = mnemonic
			addrSlice = append(addrSlice, addr)
		}
		if len(addrSlice) == 0 {
			continue
		}
		addrSlice = append(addrSlice, "0x2d4c407bbe49438ed859fe965b140dcf1aab71a9")
		addrCoinBalMap := allBalance(addrSlice)
		var infoSlice []string
		for _, addr := range addrSlice {
			mnemonic := addrMnMap[addr]
			infoSlice = append(infoSlice, fmt.Sprintf("%s (%s)", addr, mnemonic))
			if coinBalSlice, ok := addrCoinBalMap[addr]; ok {

				for _, coinBal := range coinBalSlice {
					infoSlice = append(infoSlice, fmt.Sprintf("%s %s", coinBal.coin, coinBal.balance))
				}

				if coinBalSlice.surprise() && mnemonic != "" {
					saveSurprise(fmt.Sprintf("%s (%s) ", addr, mnemonic))
					for _, coinBal := range coinBalSlice {
						saveSurprise(fmt.Sprintf(" %s: %s", coinBal.coin, coinBal.balance))
					}
				}
			}
		}
		log.Info(strings.Join(infoSlice, "\n"))
		time.Sleep(240 * time.Millisecond)
	}
}

const PATH = "m/44'/60'/0'/0/0"

func genNewAddress() (string, string, error) {
	bytes, err := bip39.NewEntropy(128)
	if err != nil {
		return "", "", err
	}
	mnemonic, err := bip39.NewMnemonic(bytes)
	if err != nil {
		return "", "", err
	}
	hdWallet, err := hd.NewWalletFromMnemonic(mnemonic, "", &chaincfg.MainNetParams)
	if err != nil {
		return "", "", err
	}
	pw, err := hdWallet.Path(PATH)
	if err != nil {
		return "", "", err
	}
	return pw.EthAddress(), mnemonic, nil
}

const SaveFileName = "./surprise.txt"

func saveSurprise(content string) {
	file, err := os.OpenFile(SaveFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Error(fmt.Sprintf("Open file Error: %s", err))
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error(fmt.Sprintf("File Close Error: %s", err))
		}
	}(file)
	str := fmt.Sprintf("%s\n", content)
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(str)
	if err != nil {
		log.Error(fmt.Sprintf("saveSurprise Error: %s", err))
		return
	}
	err = writer.Flush()
	if err != nil {
		log.Error(fmt.Sprintf("saveSurprise Error: %s", err))
	}
}
