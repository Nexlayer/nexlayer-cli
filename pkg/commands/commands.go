package commands

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ci"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/service"
)

// CICmd is the root command for CI operations
var CICmd = ci.CICmd

// ServiceCmd is the root command for service operations
var ServiceCmd = service.ServiceCmd
