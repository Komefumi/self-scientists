package main

import (
	"fmt"
	"os"
	"regexp"
)

var envValueRegexp *regexp.Regexp = regexp.MustCompile("\\w+=(\\w+)")

func main() {
	res, _ := os.ReadFile(".env")
	matches := envValueRegexp.FindAllStringSubmatch(string(res), -1)
	for _, v := range matches {
		fmt.Println(v[1])
	}
}
