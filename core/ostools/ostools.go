/********************************************************************
created:    2018-09-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package ostools

import (
	"os"
	"net/http"
)

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// 这一句需要注释掉，因为consul的健康检查不能使用json格式
	//w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"UP"}`))
}