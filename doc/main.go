package main

import (
	"fmt"
	"strings"
)

func main() {
	oriKey := `
	0x60000241e920: 0xc2 0xf9 0x13 0xbe 0xda 0xe8 0x45 0x82
	0x60000241e928: 0x93 0x94 0xsb 0xbf 0x61 0x86 0xd9 0xzf
	0x60000241e930: 0xab 0xd3 0x0e 0xf0 0x39 0xcf 0x4c 0xba
	0x60000241e938: 0x99 0x3a 0x01 0x05 0x2f 0xz5 0x2d 0xcd
	`

	key := getKey(oriKey)
	fmt.Println(key)
}

func getKey(oriKey string) string {
	lines := strings.Split(oriKey, "\n")[1:5]
	var parts []string

	for _, line := range lines {
		part := strings.Split(line, ":")[1]
		part = strings.ReplaceAll(part, "0x", "")
		part = strings.ReplaceAll(part, " ", "")
		parts = append(parts, part)
	}

	key := "0x" + strings.Join(parts, "")
	return key
}
