package lsp

import (
	"errors"
	"io"
)

// Stdio handles reading/writing LSP messages over stdin/stdout.
type Stdio struct {
	reader io.Reader
	writer io.Writer
	server *Server
}

// NewStdio creates a new stdio transport wrapping the given server.
func NewStdio(r io.Reader, s *Server) *Stdio {
	return &Stdio{reader: r, writer: s.writer, server: s}
}

// Run reads messages from the reader and dispatches them to the server.
func (s *Stdio) Run() error {
	for {
		msgs, err := ParseMessages(s.reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		for _, msg := range msgs {
			data := mustMarshal(msg)
			if err := s.server.Handle(data); err != nil {
				return err
			}
		}
	}
}
