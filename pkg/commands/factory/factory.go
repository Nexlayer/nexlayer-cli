// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package factory

import (
	"context"

	"github.com/Nexlayer/nexlayer-cli/pkg/errors"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/spf13/cobra"
)

// CommandMiddleware represents a function that wraps a command's execution
type CommandMiddleware func(next RunFunc) RunFunc

// RunFunc represents a command's execution function
type RunFunc func(ctx context.Context, cmd *cobra.Command, args []string) error

// CommandFactory handles the creation and configuration of cobra commands
type CommandFactory struct {
	logger      *observability.Logger
	middlewares []CommandMiddleware
}

// NewCommandFactory creates a new command factory
func NewCommandFactory(logger *observability.Logger) *CommandFactory {
	return &CommandFactory{
		logger:      logger,
		middlewares: make([]CommandMiddleware, 0),
	}
}

// AddMiddleware adds a middleware to the factory
func (f *CommandFactory) AddMiddleware(middleware CommandMiddleware) {
	f.middlewares = append(f.middlewares, middleware)
}

// CreateCommand creates a new cobra command with the given configuration
func (f *CommandFactory) CreateCommand(cfg *CommandConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:     cfg.Use,
		Short:   cfg.Short,
		Long:    cfg.Long,
		Example: cfg.Example,
	}

	if cfg.Run != nil {
		run := cfg.Run
		// Apply middlewares in reverse order
		for i := len(f.middlewares) - 1; i >= 0; i-- {
			run = f.middlewares[i](run)
		}

		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return run(ctx, cmd, args)
		}
	}

	return cmd
}

// CommandConfig represents the configuration for creating a command
type CommandConfig struct {
	Use     string
	Short   string
	Long    string
	Example string
	Run     RunFunc
}

// Common middlewares

// LoggingMiddleware adds logging to command execution
func LoggingMiddleware(logger *observability.Logger) CommandMiddleware {
	return func(next RunFunc) RunFunc {
		return func(ctx context.Context, cmd *cobra.Command, args []string) error {
			logger.Info(ctx, "Executing command %s with args %v", cmd.Name(), args)
			err := next(ctx, cmd, args)
			if err != nil {
				logger.Error(ctx, "Command %s failed: %v", cmd.Name(), err)
			}
			return err
		}
	}
}

// ErrorHandlingMiddleware adds structured error handling
func ErrorHandlingMiddleware() CommandMiddleware {
	return func(next RunFunc) RunFunc {
		return func(ctx context.Context, cmd *cobra.Command, args []string) error {
			err := next(ctx, cmd, args)
			if err != nil {
				// Convert regular errors to our structured error type
				if _, ok := err.(*errors.Error); !ok {
					err = errors.InternalError("command execution failed", err)
				}
			}
			return err
		}
	}
}

// RecoveryMiddleware adds panic recovery
func RecoveryMiddleware(logger *observability.Logger) CommandMiddleware {
	return func(next RunFunc) RunFunc {
		return func(ctx context.Context, cmd *cobra.Command, args []string) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(ctx, "Panic recovered in command %s: %v", cmd.Name(), r)
					err = errors.InternalError("panic during command execution", nil)
				}
			}()
			return next(ctx, cmd, args)
		}
	}
}
