# NAEOS NEIR VS Code Extension

Language support for NAEOS NEIR specification files.

## Features

- Syntax highlighting for `.naeos.yaml`, `.naeos.yml`, and `.neir.yaml` files
- LSP-powered real-time validation and diagnostics
- Autocompletion for spec keywords, service kinds, HTTP methods, and more
- Hover information for all spec fields
- Commands for validation and AI suggestions

## Commands

- `NEIR: Validate Specification` — Run validation on the current spec
- `NEIR: Get AI Suggestions` — Get AI-powered improvement suggestions

## Configuration

- `naeos.executablePath` — Path to the `naeos` executable (default: `naeos`)
- `naeos.lsp.enabled` — Enable LSP server for real-time validation (default: `true`)

## Requirements

- NAEOS CLI installed and available in PATH
- VS Code 1.85 or later
