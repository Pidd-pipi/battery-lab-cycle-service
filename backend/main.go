package main

import (
	"log"
)

func main() {
	config := LoadConfig()
	service := NewCycleService(NewCellStore())
	log.Printf("battery lab cycle service listening on :%s", config.Port)
	log.Fatal(serveAddress(":"+config.Port, NewRouter(service)))
}
