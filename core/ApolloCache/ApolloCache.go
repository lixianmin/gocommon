/********************************************************************
created:    2018-11-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package ApolloCache

import (
	"coinbene.com/gocommon/core/console"
	"fmt"
	"github.com/lixianmin/gocore/loom"
	"github.com/philchia/agollo"
	"strings"
	"sync"
	"time"
)

type updateValueTable struct {
	sync.RWMutex
	m map[string]interface{}
}

var updateValues = &updateValueTable{
	m: make(map[string]interface{}),
}

func init() {
	loom.Repeat(time.Minute, func() {
		// 这里并没有修改updateValues，所以用只读锁
		updateValues.RLock()
		defer updateValues.RUnlock()

		for key, val := range updateValues.m {
			switch val := val.(type) {
			case *loom.String:
				var oldText = val.Load()
				var newText = innerGetStringValue(key, oldText)
				if newText != oldText {
					val.Store(newText)
					console.Notice("[ApolloCache.loom.Repeat()] key=%q, oldText=%q, newText=%q", key, oldText, newText)
				}
			}
		}
	})
}

// 这个方法的返回值是*loom.String而不是string，原因是它会在后台每分钟自动更新
func GetStringValue(key string, defaultValue string) *loom.String {
	updateValues.RLock()
	var oldValue, ok = updateValues.m[key]
	updateValues.RUnlock()

	if ok {
		oldText, ok := oldValue.(*loom.String)
		if ok {
			return oldText
		}
		var message = fmt.Sprintf("[GetStringValue()] Try to call GetStringValue() with different value type, key=%q", key)
		panic(message)
	}

	var text = innerGetStringValue(key, defaultValue)
	var newValue = new(loom.String)
	newValue.Store(text)
	console.Notice("[GetStringValue()] key=%q, text=%q", key, text)

	updateValues.Lock()
	updateValues.m[key] = newValue
	updateValues.Unlock()
	return newValue
}

func innerGetStringValue(key string, defaultValue string) string {
	return strings.TrimSpace(agollo.GetStringValue(key, defaultValue))
}
