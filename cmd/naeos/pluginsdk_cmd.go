package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPluginSDKCommand() *cobra.Command {
	return &cobra.Command{
		Use:        "pluginsdk",
		Short:      "Plugin SDK commands (deprecated)",
		Deprecated: "use 'naeos plugin' instead",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.ErrOrStderr(), "Warning: 'naeos pluginsdk' is deprecated. Use 'naeos plugin' instead.")
			fmt.Fprintln(cmd.ErrOrStderr(), "")
			fmt.Fprintln(cmd.ErrOrStderr(), "Available commands:")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin list        - List installed plugins")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin install     - Install a plugin")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin uninstall   - Uninstall a plugin")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin enable      - Enable a plugin")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin disable     - Disable a plugin")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin info        - Show plugin info")
			fmt.Fprintln(cmd.ErrOrStderr(), "  naeos plugin execute     - Execute a plugin action")
			return nil
		},
	}
}
