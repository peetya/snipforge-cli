package main

import (
	"github.com/peetya/snipforge-cli/cmd"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	cmd.Execute()
}
