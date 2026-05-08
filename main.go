package main

import (
	Logger "23102081/logging_middleware"
	"fmt"
	"log"
)

func main() {
	if err := Logger.Log("backend", "debug", "middleware", "testing"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("log sent successfully")
}
