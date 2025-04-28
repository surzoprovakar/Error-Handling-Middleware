package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Operation2Value struct {
	action   string
	operator int
}

func ReadFile(fileName string) []string {

	buf, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var lines []string
	snl := bufio.NewScanner(buf)
	for snl.Scan() {
		lines = append(lines, snl.Text())
	}
	err = snl.Err()
	if err != nil {
		log.Fatal(err)
	}
	return lines
}

func ReadFile_Actions(fileName string) []Operation2Value {
	buf, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = buf.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var OperationsList []Operation2Value
	var lines string
	snl := bufio.NewScanner(buf)

	for snl.Scan() {
		var act string
		var val int
		var opt2val Operation2Value

		lines = snl.Text()
		split := strings.SplitAfter(lines, " ")

		act = strings.Trim(split[0], " ")
		opt2val.action = act

		if len(split) > 1 {

			val, err = strconv.Atoi(split[1])

			opt2val.operator = val
		}
		opt2val.operator = val

		OperationsList = append(OperationsList, opt2val)
	}

	err = snl.Err()
	if err != nil {
		log.Fatal(err)
	}

	return OperationsList
}
