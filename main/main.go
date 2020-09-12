package main

import (
	"os"

	"log"

	"fmt"

	"github.com/makasim/ymltoerl"
)

func main() {
	b, err := ymltoerl.ConvertFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
