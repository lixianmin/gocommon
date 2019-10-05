/********************************************************************
created:    2018-11-27
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

package dbtools

import (
	"strings"
	"database/sql"
	"fmt"
	"errors"
)

type IExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func lastIndexOfParenthesis(query string) int {
	var count = len(query)
	var leftParenthesisCount, rightParenthesisCount = 0, 0
	for i := count - 1; i >= 0; i-- {
		var c = query[i]
		if c == ')' {
			rightParenthesisCount ++
		} else if c == '(' {
			leftParenthesisCount ++

			if leftParenthesisCount == rightParenthesisCount {
				return i
			}
		}
	}

	return -1
}

func BatchInsert(tx IExec, query string, values [][]interface{}) error {
	var count = len(values)
	if count == 0 {
		return nil
	}

	var lastIndex = lastIndexOfParenthesis(query)
	if lastIndex == -1 {
		var message = fmt.Sprintf("Can not find '(' in the query=%q", query)
		return errors.New(message)
	}

	const stepLength = 10
	var valueText = "," + query[lastIndex:]
	var batchQuery = query + strings.Repeat(valueText, stepLength-1)
	var stopIndex = count / stepLength * stepLength

	var args = make([]interface{}, 0)
	for i := 0; i < stopIndex; i += stepLength {
		args = args[0:0]
		for j := i; j < i+stepLength; j++ {
			args = append(args, values[j]...)
		}

		var _, err = tx.Exec(batchQuery, args...)
		if err != nil {
			return err
		}
	}

	var leftCount = count - stopIndex
	if leftCount >= 1 {
		batchQuery = query + strings.Repeat(valueText, leftCount-1)
		args = args[0:0]
		for j := stopIndex; j < count; j++ {
			args = append(args, values[j]...)
		}

		var _, err = tx.Exec(batchQuery, args...)
		if err != nil{
			return err
		}
	}

	return nil
}
