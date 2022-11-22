package main

import (
	"context"
	"log"
	"os"

	"github.com/fujiwara/riex"
)

func main() {
	ctx := context.TODO()
	err := riex.RunCLI(ctx, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
