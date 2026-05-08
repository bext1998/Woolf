package ingestion

import "os"

func IngestMarkdown(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	return Document{Path: path, Format: "md", Content: string(data)}, nil
}
