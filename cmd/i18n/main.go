package main

import (
	"fmt"

	"golang.org/x/text/language"
)

func main() {
	tag, _, err := language.ParseAcceptLanguage("zh")

	fmt.Println(tag, err)
}
