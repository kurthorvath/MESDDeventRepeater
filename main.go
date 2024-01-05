package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func forward2Request(LD string) {
	posturl := "http://127.0.0.1:7000"
	body := []byte(LD)

	_, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	fmt.Println("forward LD", LD, "to", posturl)
}

func main() {
	ctx, _ := client.NewClientWithOpts(client.FromEnv)
	reader, err := ctx.ContainerLogs(context.Background(), "259e221e47c0", types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})

	if err != nil {
		log.Fatal(err)
	}

	//var stdout string
	go func() {
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
					//fmt.Println("new elem", out)
				}

				//stdout = stdout + line + "\n"
			}
			if err != nil {
				return
			}
		}
	}()

	fmt.Println(reflect.TypeOf(reader))
	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, reader)
	if err != nil && err != io.EOF {
		//	log.Fatal(err)
	}
	sBUF := buf.String()

	if strings.Contains(sBUF, "dc1") {
		//	fmt.Println("new elem", sBUF)
	}

}

/*
func main() {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			fmt.Printf("%s %s %s\n", container.ID[:10], container.Image, container.Names)
		}


	containerID := "259e221e47c0"

	reader, err := cli.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("done read")
	defer reader.Close()

	//read the first 8 bytes to ignore the HEADER part from docker container logs
	p := make([]byte, 8)
	reader.Read(p)
	content, _ := ioutil.ReadAll(reader)
	fmt.Println("read")
	//var codeOutput MyJSONStruct

		if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&codeOutput); err != nil {
			//handle error
		}
		//set some other value in struct
		codeOutput.ContainerID = containerID


	fmt.Println(content)
}
*/
