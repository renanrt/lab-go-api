package main

import (
	"fmt"
	"os"

	"github.com/renanrt/lab-go-api/api"
)

func main() {

	mode := "api"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	switch mode {
	case "api":
		RunAPI()
	case "migrate":
		fmt.Printf("migrate mode")

	default:
		fmt.Printf("Unknown mode")
	}
}

func RunAPI() {
	err := api.Serve()
	if err != nil {
		fmt.Printf("Error serving API")
	}
}
