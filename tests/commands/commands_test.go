package commands_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/feedback"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// CommandTestSuite contains all command tests
type CommandTestSuite struct {
	suite.Suite
	client *commands.MockAPIClient
	buffer *bytes.Buffer
}

func (s *CommandTestSuite) SetupTest() {
	s.client = &commands.MockAPIClient{
		GetDeploymentInfoFunc: func(ctx context.Context, namespace, appID string) (*schema.APIResponse[schema.Deployment], error) {
			return &schema.APIResponse[schema.Deployment]{
				Message: "Success",
				Data: schema.Deployment{
					Namespace: namespace,
					Status:    "Running",
				},
			}, nil
		},
		ListDeploymentsFunc: func(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
			return &schema.APIResponse[[]schema.Deployment]{
				Message: "Success",
				Data: []schema.Deployment{
					{
						Namespace: "test",
						Status:    "Running",
					},
				},
			}, nil
		},
	}
	s.buffer = new(bytes.Buffer)
}

// TestAICommand tests the AI command and its subcommands
func (s *CommandTestSuite) TestAICommand() {
	cmd := ai.NewCommand()
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "ai [subcommand]", cmd.Use)
	assert.Equal(s.T(), "AI-powered features for Nexlayer", cmd.Short)
	assert.NotEmpty(s.T(), cmd.Long)

	// Test subcommands
	generateCmd, _, err := cmd.Find([]string{"generate"})
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), generateCmd)
	assert.Equal(s.T(), "generate <app-name>", generateCmd.Use)

	detectCmd, _, err := cmd.Find([]string{"detect"})
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), detectCmd)
	assert.Equal(s.T(), "detect", detectCmd.Use)
}

// TestDeployCommand tests the deploy command
func (s *CommandTestSuite) TestDeployCommand() {
	cmd := deploy.NewCommand(s.client)
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "deploy [applicationID]", cmd.Use)
	assert.Contains(s.T(), cmd.Short, "Deploy")
}

// TestDomainCommand tests the domain command and its subcommands
func (s *CommandTestSuite) TestDomainCommand() {
	cmd := domain.NewDomainCommand(s.client)
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "domain", cmd.Use)
	assert.Contains(s.T(), cmd.Short, "domain")

	// Test domain set command
	setCmd, _, err := cmd.Find([]string{"set"})
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), setCmd)
}

// TestFeedbackCommand tests the feedback command
func (s *CommandTestSuite) TestFeedbackCommand() {
	cmd := feedback.NewFeedbackCommand(s.client)
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "feedback send", cmd.Use)
	assert.Contains(s.T(), cmd.Short, "feedback")

	// Test required message flag
	msgFlag := cmd.Flags().Lookup("message")
	assert.NotNil(s.T(), msgFlag)
}

// TestInfoCommand tests the info command
func (s *CommandTestSuite) TestInfoCommand() {
	cmd := info.NewInfoCommand(s.client)
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "info <namespace> <applicationID>", cmd.Use)
	assert.Contains(s.T(), cmd.Short, "info")
}

// TestListCommand tests the list command
func (s *CommandTestSuite) TestListCommand() {
	cmd := list.NewListCommand(s.client)
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "list [applicationID]", cmd.Use)
	assert.Contains(s.T(), cmd.Short, "List")
}

// TestLoginCommand tests the login command
func (s *CommandTestSuite) TestLoginCommand() {
	cmd := login.NewLoginCommand(s.client)
	assert.NotNil(s.T(), cmd)
	assert.Equal(s.T(), "login", cmd.Use)
	assert.Equal(s.T(), "Log in to Nexlayer", cmd.Short)
	assert.NotEmpty(s.T(), cmd.Long)

	cmd.SetOut(s.buffer)
	err := cmd.Execute()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "not yet implemented")
}

// TestCommandSuite runs all command tests
func TestCommandSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
