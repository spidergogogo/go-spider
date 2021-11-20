package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

var proxyClient *http.Client

func initProxyClient(proxy string) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	}}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		} else {
			log.Error(fmt.Sprintf("url.parse Error: %s", err))
		}
	}
	proxyClient = &http.Client{
		Transport: tr,
	}
}
