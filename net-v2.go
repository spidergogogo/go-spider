package main

import (
	"fmt"
	"strings"
)

type addrMnEntity struct {
	coin      string
	addrMnMap map[string]string
	addr      []string
}

func dispatchHandler(chanAddrMn chan *addrMnEntity) {
	for {
		select {
		case am := <-chanAddrMn:
			coinAddrBalanceHandler(am)
		}
	}
}

func coinAddrBalanceHandler(am *addrMnEntity) {
	addrBalMap := reqCoinAddrBalance(am.coin, am.addr)
	var infoSlice []string
	for _, addr := range am.addr {
		mn := am.addrMnMap[addr]
		bal := addrBalMap[addr]
		info := fmt.Sprintf("%s(%s)\n %s:%s", addr, mn, am.coin, bal)
		infoSlice = append(infoSlice, info)
		if mn != "" && bal != "Err" && bal != "0" {
			saveSurprise(info)
		}
	}
	log.Info(strings.Join(infoSlice, "\n"))
}

func reqCoinAddrBalance(coin string, address []string) map[string]string {
	addrStr := strings.Join(address, ",")
	funcBal, ok := coinFuncMap[coin]
	rAddrBalMap := make(map[string]string)
	if !ok {
		log.Error(fmt.Sprintf("%s Not have funcBalance", coin))
	} else {
		addrBalMap, err := funcBal(addrStr)
		if err != nil {
			log.Error(fmt.Sprintf("%s funcBall Error:%s", coin, err))
			for _, addr := range address {
				rAddrBalMap[addr] = "Err"
			}
		} else {
			for _, addr := range address {
				rAddrBalMap[addr] = addrBalMap[addr]
			}
		}
	}
	return rAddrBalMap
}
