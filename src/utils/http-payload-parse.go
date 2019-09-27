package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func parseContent(content []byte, request bool) {
	if len(content) == 0 {
		fmt.Println("Data:")
		return
	} else {
		if request {
			fmt.Println("Data:" + string(content))
		} else {
			fmt.Print("Data:")
		}
	}
	cur := 0
	chunkIndex := 0
	for i, b := range content {
		if i < len(content)-1 && b == '\r' && content[i+1] == '\n' {
			line := string(content[cur:i])
			if chunkIndex%2 == 0 {
				length, err := strconv.ParseInt(line, 16, 0)
				if err != nil {
					log.Fatal(err)
				}
				if length == 0 {
					break
				}
			} else {
				fmt.Println(line)
			}
			cur = i + 2
			chunkIndex++
		}
	}

}

func ParsePayload(payload []byte) {
	request := false
	cur := 0
	for i, b := range payload {
		if i < len(payload)-1 && b == '\r' && payload[i+1] == '\n' {
			line := string(payload[cur : i+2])
			if strings.HasPrefix(line, "Host:") {
				request = true
			}
			cur = i + 2
			if len(line) == 2 {
				parseContent(payload[cur:], request)
				break
			} else {
				fmt.Print(line)
			}
		}
	}
}
