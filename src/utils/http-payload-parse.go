package utils

import (
	"fmt"
	"strconv"
	"strings"
)

/**
only parse the http content
*/
func parseContent(content []byte, request bool, chunked bool) {
	length := len(content)
	if length == 0 {
		fmt.Println("Data:")
		return
	} else {
		if request {
			fmt.Println("Data:" + string(content))
		} else {
			fmt.Print("Data:")
		}
	}
	if chunked {
		cur := 0
		for i := 0; i < length; i++ {
			if i <= length-1 && content[i] == '\n' {
				var line string
				if content[i-1] == '\r' {
					line = string(content[cur : i-1])
				} else {
					line = string(content[cur:i])
				}
				length, err := strconv.ParseInt(line, 16, 0)
				if err != nil {
					break
				}
				if length > 0 {
					fmt.Println(string(content[i+1 : i+1+int(length)]))
					i = i + 2 + int(length)
					cur = i
				}
			}
		}
	} else {
		fmt.Println(string(content))
	}
}

/**
parse the whole payload(http)
*/
func ParsePayload(payload []byte) {
	request := false
	chunked := false
	cur := 0
	for i, b := range payload {
		if i <= len(payload)-1 && b == '\n' {
			line := string(payload[cur : i+1])
			if strings.HasPrefix(line, "Host:") {
				request = true
			}
			if strings.HasPrefix(line, "Transfer-Encoding: chunked") {
				chunked = true
			}
			cur = i + 1
			if len(line) == 2 {
				parseContent(payload[cur:], request, chunked)
				break
			} else {
				fmt.Print(line)
			}
		}
	}
}
