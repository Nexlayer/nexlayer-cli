package main

import "fmt"

type HelloPlugin struct{}

func (p *HelloPlugin) Name() string {
	return "hello"
}

func (p *HelloPlugin) Description() string {
	return "A simple hello world plugin"
}

func (p *HelloPlugin) Run(opts map[string]interface{}) error {
	name, _ := opts["name"].(string)
	if name == "" {
		name = "World"
	}
	fmt.Printf("Hello, %s!\n", name)
	return nil
}

// Plugin is the exported symbol that Nexlayer will look for
var Plugin HelloPlugin
