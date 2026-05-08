package exporter

import "woolf/internal/session"

type Exporter interface {
	Export(session.Session) ([]byte, error)
}
