package main

import (
	"log"
	"os"

	"github.com/MohamedBassem/DSSS/server"
)

func main() {
	logger := log.New(os.Stdout, "Server", log.LstdFlags)
	server.Main(8082, 8083, 8081, logger)
}
