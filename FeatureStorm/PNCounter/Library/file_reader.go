package main

import (
	"bufio"
	"log"
	"os"
)

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
