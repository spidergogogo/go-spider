package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type accountBalance struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

type balanceResponse struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}

type balanceResponseError struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
	Result string `json:"result"`
}

type balanceResponseSuccess struct {
	Status string           `json:"status"`
	Msg    string           `json:"message"`
	Result []accountBalance `json:"result"`
}

type coinBalance struct {
	coin    string
	balance string
}

func (c *coinBalance) surprise() bool {
	if c.balance != "Err" && c.balance != "0" {
		return true
	}
	return false
}

type coinBalanceSlice []*coinBalance

func (c coinBalanceSlice) surprise() bool {
	for _, coinBal := range c {
		if coinBal.surprise() {
			return true
		}
	}
	return false
}

const CoinEth = "ETH"
const CoinBsc = "BSC"

type funcBalance func(string) (map[string]string, error)

var coinSlice = []string{CoinEth, CoinBsc}
var coinFuncMap = map[string]funcBalance{
	CoinEth: ethBalance,
	CoinBsc: bscBalance,
}

func allBalance(address []string) map[string]coinBalanceSlice {
	addrStr := strings.Join(address, ",")
	addrCoinBalMap := make(map[string]coinBalanceSlice)
	for _, coin := range coinSlice {
		funcBal, ok := coinFuncMap[coin]
		if !ok {
			log.Error(fmt.Sprintf("*%s* Not have funcBalance\n", coin))
		} else {
			addrBalMap, err := funcBal(addrStr)
			if err != nil {
				log.Error(fmt.Sprintf("*%s* funcBal Error:%s", coin, err))
				for _, addr := range address {
					addrCoinBalMap[addr] = append(addrCoinBalMap[addr], &coinBalance{
						coin:    coin,
						balance: "Err",
					})
				}
			} else {
				for _, addr := range address {
					addrCoinBalMap[addr] = append(addrCoinBalMap[addr], &coinBalance{
						coin:    coin,
						balance: addrBalMap[addr],
					})
				}
			}
		}
	}
	return addrCoinBalMap
}

func bscBalance(address string) (map[string]string, error) {
	//reqUrlTemplate := "https://api.bscscan.com/api?module=account&action=balance&address=%s"
	reqUrlTemplate := "https://api.bscscan.com/api?module=account&action=balancemulti&address=%s&tag=latest&apikey=%s"
	//reqUrlTemplate := "https://api.bscscan.com/api?module=account&action=balancemulti&address=%s&tag=latest"
	reqUrl := fmt.Sprintf(reqUrlTemplate, address, cfg.ApiKey.Bsc)
	rMap, err := balance(reqUrl)
	return rMap, err
}
func ethBalance(address string) (map[string]string, error) {
	//reqUrlTemplate := "https://api.etherscan.com/api?module=account&action=balance&address=%s"
	reqUrlTemplate := "https://api.etherscan.com/api?module=account&action=balancemulti&address=%s&tag=latest&apikey=%s"
	//reqUrlTemplate := "https://api.etherscan.com/api?module=account&action=balancemulti&address=%s&tag=latest"
	reqUrl := fmt.Sprintf(reqUrlTemplate, address, cfg.ApiKey.Eth)
	rMap, err := balance(reqUrl)
	return rMap, err
}

// return addr:balance
func balance(reqUrl string) (map[string]string, error) {
	log.Info(reqUrl)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := proxyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err1 != nil {
			log.Error(fmt.Sprintf("Body Close Error: %s, %s", err1, reqUrl))
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("status code:%d", resp.StatusCode)
	} else {
		balResp := &balanceResponse{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, balResp)
		if err != nil {
			return nil, err
		}
		if balResp.Status != "1" {
			balRespErr := &balanceResponseError{}
			err = json.Unmarshal(body, balRespErr)
			if err != nil {
				return nil, errors.Wrap(err, "balRespErr")
			}
			return nil, errors.New(fmt.Sprintf("%s,%s,%s", balResp.Status,
				balResp.Msg, balRespErr.Result))
		} else {
			balRespSuccess := &balanceResponseSuccess{}
			err = json.Unmarshal(body, balRespSuccess)
			if err != nil {
				return nil, errors.Wrap(err, "balRespSuccess")
			}
			rMap := make(map[string]string)
			for _, ab := range balRespSuccess.Result {
				rMap[ab.Account] = ab.Balance
			}
			return rMap, nil
		}
	}
}
