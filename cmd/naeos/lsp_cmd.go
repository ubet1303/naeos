package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/lsp"
)

func newLSPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lsp",
		Short: "Start LSP (Language Server Protocol) server for NEIR specs",
		Long: `Start a language server that provides IDE support for NAEOS specification files.

Features:
  - Real-time diagnostics (parse errors, missing fields, resolution issues)
  - Autocompletion for spec keywords, module fields, service kinds, etc.
  - Hover information for all spec fields

The server communicates via stdin/stdout using the LSP protocol.

Example:
  naeos lsp
  naeos lsp --stdio`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			server := lsp.NewServer(os.Stdout)
			stdio := lsp.NewStdio(os.Stdin, server)
			return stdio.Run()
		},
	}

	return cmd
}
