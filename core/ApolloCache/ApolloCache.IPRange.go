/********************************************************************
created:    2018-11-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package ApolloCache

import (
	"github.com/philchia/agollo"
	"coinbene.com/gocommon/core/console"
	"coinbene.com/gocommon/core/iptools"
	"coinbene.com/gocommon/core/convert"
	"sync/atomic"
	"math"
	"github.com/lixianmin/gocore/loom"
	"time"
)

type numRange struct {
	left  int64
	right int64
}

type IPRangeConfig struct {
	ipRanges atomic.Value
}

func CreateIPRangeConfig(key string) *IPRangeConfig {
	var config = &IPRangeConfig{}

	var ranges = createNumRanges(key)
	config.ipRanges.Store(ranges)

	loom.Repeat(time.Minute, func() {
		var newRanges = createNumRanges(key)
		config.ipRanges.Store(newRanges)
	})

	return config
}

func createNumRanges(key string) []numRange {
	var text = agollo.GetStringValue(key, "")
	if text != "" {
		var rawRangeList [][2]string

		var err = convert.FromJson([]byte(text), &rawRangeList)
		if err == nil {
			var list = make([]numRange, 0)
			for _, item := range rawRangeList {
				left, err := iptools.GetIPNum(item[0])
				if err != nil {
					console.Error(err)
					continue
				}
				right, err := iptools.GetIPNum(item[1])
				if err != nil {
					console.Error(err)
					continue
				}

				list = append(list, numRange{left: left, right: right})
			}

			return list
		} else {
			console.Error("text=%s, err=%q", text, err)
		}
	}

	// 默认所有ip都合法
	var defaultIPRanges = make([]numRange, 1)
	defaultIPRanges[0] = numRange{left: 0, right: math.MaxInt64}
	return defaultIPRanges
}

func (config *IPRangeConfig) IsAllowedIP(ip string) bool {
	var ipNum, err = iptools.GetIPNum(ip)
	if err != nil {
		console.Error(err)
		return false
	}

	var ranges = config.ipRanges.Load().([]numRange)
	for _, numRange := range ranges {
		if ipNum >= numRange.left && ipNum <= numRange.right {
			return true
		}
	}

	return false
}
