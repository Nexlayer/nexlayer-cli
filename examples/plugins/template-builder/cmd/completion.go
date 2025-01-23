package cmd

import (
	"os"
	
	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script for your shell",
		Long: `To load completions:

Bash:
  $ source <(nexlayer completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ nexlayer completion bash > /etc/bash_completion.d/nexlayer
  # macOS:
  $ nexlayer completion bash > /usr/local/etc/bash_completion.d/nexlayer

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ nexlayer completion zsh > "${fpath[1]}/_nexlayer"

  # You will need to start a new shell for this setup to take effect.

fish:
  $ nexlayer completion fish | source

  # To load completions for each session, execute once:
  $ nexlayer completion fish > ~/.config/fish/completions/nexlayer.fish

PowerShell:
  PS> nexlayer completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> nexlayer completion powershell > nexlayer.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
}
