package ingestion

type Document struct {
	Path    string
	Format  string
	Content string
}

type Ingester interface {
	Ingest(path string) (Document, error)
}
