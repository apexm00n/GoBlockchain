package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Wallet_Server: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5001", "Blockchain Gateway")
	flag.Parse()
	app := NewWalletServer(uint16(*port), *gateway)
	log.Print("Server starts, port ", *port)
	app.Run()
}