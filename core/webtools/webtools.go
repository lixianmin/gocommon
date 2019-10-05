/********************************************************************
created:    2018-09-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package webtools

import (
	"io/ioutil"
	"net/http"
	"time"
)

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func Get(url string, initRequest func(request *http.Request)) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if nil != initRequest {
		initRequest(request)
	}

	var client = http.Client{
		Timeout: time.Second * 20, // 控制从链接建立到返回的整个生命周期的时间
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var body = response.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	return bodyBytes, err
}
