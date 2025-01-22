package hello

import (
	"encoding/json"
	"fmt"
	"os"
)

type PluginMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--describe" {
		metadata := PluginMetadata{
			Name:        "hello",
			Version:     "1.0.0",
			Description: "A simple hello world plugin for Nexlayer CLI",
			Usage:       "nexlayer hello [name]",
		}
		json.NewEncoder(os.Stdout).Encode(metadata)
		return
	}

	name := "World"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	fmt.Printf("Hello, %s!"
", name)"
}
