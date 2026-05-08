package ingestion

import (
	"fmt"
	"path/filepath"
	"strings"

	pdfingest "woolf/internal/ingestion/pdf"
)

type Document struct {
	Path    string
	Format  string
	Content string
}

type Ingester interface {
	Ingest(path string) (Document, error)
}

func Ingest(path string) (Document, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return IngestMarkdown(path)
	case ".txt":
		return IngestText(path)
	case ".pdf":
		text, err := pdfingest.ExtractText(path)
		if err != nil {
			return Document{}, err
		}
		return Document{Path: path, Format: "pdf", Content: pdfingest.ToMarkdown(text)}, nil
	default:
		return Document{}, fmt.Errorf("ING-002: unsupported file format %s", filepath.Ext(path))
	}
}
