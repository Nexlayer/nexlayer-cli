// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package main is the entry point for the Nexlayer CLI application.
// It initializes and executes the root command, which sets up all
// subcommands and their respective functionality.
package main

import "github.com/Nexlayer/nexlayer-cli/cmd"

// main is the entry point of the Nexlayer CLI.
// It delegates to cmd.Execute() which handles command-line parsing,
// configuration loading, and command execution.
func main() {
	cmd.Execute()
}
