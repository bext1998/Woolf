package ingestion

import "os"

func IngestText(path string) (Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	return Document{Path: path, Format: "txt", Content: string(data)}, nil
}
