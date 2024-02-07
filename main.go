package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func forward2Request(LD string) {
	posturl := "http://127.0.0.1:7000"
	body := []byte(LD)

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("forward LD", LD, "to", posturl, res)
}

func main() {
	fmt.Println("start MESDD/Consul event forwarder")
	reader, err := os.Open("consul.log")
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		//os.Exit(1)
	}
	//var stdout string
	go func() {
		buf := bufio.NewReader(reader)

	buf := bufio.NewReader(reader)
	for {
		line, err := buf.ReadString('\n')
		if len(line) > 0 {
			if strings.Contains(line, "service.consul") == true {

				re := regexp.MustCompile(`name=(.*)type`)
				match := re.FindStringSubmatch(line)

				if len(match) > 1 {
					fmt.Println("match found -", match[1][:len(match[1])-2], len(match[1][:len(match[1])-2]))
					forward2Request(match[1][:len(match[1])-2])
				}
				//fmt.Println("new elem", line)
			}

			//stdout = stdout + line + "\n"
		}

		if err != nil {
			//	fmt.Println(err)
		}
	}
}
