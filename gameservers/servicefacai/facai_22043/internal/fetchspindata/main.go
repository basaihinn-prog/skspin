package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"serve/servicejdb/jdbcomm"
)

func main() {
	if len(os.Args) > 1 {
		bytes, _ := base64.StdEncoding.DecodeString(os.Args[1])
		fmt.Printf("%x", bytes)
		got, err := jdbcomm.NewFromBinaryData(bytes)
		fmt.Println(got, err)
	} else {
		fmt.Println("input base64 code")
	}

}
