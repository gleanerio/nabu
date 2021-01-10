package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"
)

func main() {

	// desc := gjson.Get(string(b), "description")
	dat, err := ioutil.ReadFile("test.json")
	if err != nil {
		fmt.Println(err)
	}

	result := gjson.Get(string(dat), "description")

	if result.String() == "" {
		fmt.Println("found no description")
	} else {
		fmt.Println(result.String())
	}
}
