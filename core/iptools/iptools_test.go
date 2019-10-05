/********************************************************************
created:    2018-10-14
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
package iptools

import (
	"testing"
	"fmt"
)

func TestGetIPNum(t *testing.T) {
	var ip = "172.20.20.503"
	fmt.Println(GetIPNum(ip))
}