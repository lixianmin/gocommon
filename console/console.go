/********************************************************************
created:    2018-09-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package console

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lixianmin/gocore/loom"
	"sync"
	"time"
)

var history = sync.Map{}

func init() {
	// 输出文件名和文件行号
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(4)

	// 设置异步输出，以提升性能
	logs.Async()

	// 周期性清空logHistory
	loom.Repeat(1*time.Second, func() {
		history.Range(func(key, value interface{}) bool {
			history.Delete(key)
			return true
		})
	})
}

func Info(f interface{}, v ...interface{}) {
	logs.Info(f, v...)
}

// 不输出重复的日志，防止日志过于频繁导致server不可用
func InfoNew(format string, v ...interface{}) {
	var message = fmt.Sprintf(format, v...)
	if _, loaded := history.LoadOrStore(message, nil); !loaded {
		logs.Info(message)
	}
}

func Notice(f interface{}, v ...interface{}) {
	logs.Notice(f, v...)
}

// 不输出重复的日志，防止日志过于频繁导致server不可用
func NoticeNew(format string, v ...interface{}) {
	var message = fmt.Sprintf(format, v...)
	if _, loaded := history.LoadOrStore(message, nil); !loaded {
		logs.Notice(message)
	}
}

func Warn(f interface{}, v ...interface{}) {
	logs.Warn(f, v...)
}

// 不输出重复的日志，防止日志过于频繁导致server不可用
func WarnNew(format string, v ...interface{}) {
	var message = fmt.Sprintf(format, v...)
	if _, loaded := history.LoadOrStore(message, nil); !loaded {
		logs.Warn(message)
	}
}

func Error(f interface{}, v ...interface{}) {
	logs.Error(f, v...)
}

// 不输出重复的日志，防止日志过于频繁导致server不可用
func ErrorNew(format string, v ...interface{}) {
	var message = fmt.Sprintf(format, v...)
	if _, loaded := history.LoadOrStore(message, nil); !loaded {
		logs.Error(message)
	}
}

func Flush() {
	logs.GetBeeLogger().Flush()
}
