package main

import (
	"log"
	"os"

	"github.com/miyanokomiya/mm-code/server"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("[MM]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server.Run()
}
