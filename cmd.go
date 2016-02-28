package main

import (
	"log"
	"os"

	"github.com/MohamedBassem/DSSS/server"
)

func main() {
	logger := log.New(os.Stdout, "Server", log.LstdFlags)
	server.Main(0, 8082, 8081, logger)
}
