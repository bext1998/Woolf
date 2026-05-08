package exporter

import (
	"fmt"

	"woolf/internal/session"
)

type MarkdownExporter struct{}

func (MarkdownExporter) Export(sess session.Session) ([]byte, error) {
	return []byte(fmt.Sprintf("# %s\n\nSession: %s\n", sess.Title, sess.SessionID)), nil
}
