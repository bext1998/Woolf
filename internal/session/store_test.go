package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCreateSaveLoadAndList(t *testing.T) {
	dir := t.TempDir()
	draft := filepath.Join(dir, "chapter3.md")
	if err := os.WriteFile(draft, []byte("# Chapter 3\n\nThis is a draft."), 0o600); err != nil {
		t.Fatal(err)
	}

	store := NewStore(filepath.Join(dir, "sessions"))
	sess, path, err := store.Create("Chapter 3", draft)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if !validSessionID(sess.SessionID) {
		t.Fatalf("SessionID %q does not match required rule", sess.SessionID)
	}
	if filepath.Base(path) != sess.SessionID+".json" {
		t.Fatalf("path = %q, want file named from session id", path)
	}

	loaded, loadedPath, err := store.Load(sess.SessionID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loadedPath != path {
		t.Fatalf("loaded path = %q, want %q", loadedPath, path)
	}
	if loaded.Status != StatusActive {
		t.Fatalf("status = %q, want %q", loaded.Status, StatusActive)
	}
	if loaded.Source == nil || loaded.Source.Type != "file" || loaded.Source.ContentHash == "" {
		t.Fatalf("source was not serialized with file metadata: %#v", loaded.Source)
	}

	items, err := store.List(ListOptions{Limit: 20})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 1 || items[0].SessionID != sess.SessionID {
		t.Fatalf("List() = %#v, want created session", items)
	}
}

func TestLoadResolvesPrefixAndListIndex(t *testing.T) {
	store := NewStore(t.TempDir())
	writeSession(t, store, Session{
		SessionID:     "20260507-143022-first",
		Version:       Version,
		Title:         "first",
		Status:        StatusPaused,
		CreatedAt:     time.Date(2026, 5, 7, 14, 30, 22, 0, time.UTC),
		UpdatedAt:     time.Date(2026, 5, 7, 14, 30, 22, 0, time.UTC),
		AgentsConfig:  []AgentConfig{},
		Rounds:        []Round{},
		Interventions: []Intervention{},
		Summaries:     map[string]string{},
	})
	writeSession(t, store, Session{
		SessionID:     "20260508-091500-second",
		Version:       Version,
		Title:         "second",
		Status:        StatusCompleted,
		CreatedAt:     time.Date(2026, 5, 8, 9, 15, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2026, 5, 8, 9, 15, 0, 0, time.UTC),
		AgentsConfig:  []AgentConfig{},
		Rounds:        []Round{},
		Interventions: []Intervention{},
		Summaries:     map[string]string{},
	})

	byPrefix, _, err := store.Load("20260507")
	if err != nil {
		t.Fatalf("Load(prefix) error = %v", err)
	}
	if byPrefix.SessionID != "20260507-143022-first" {
		t.Fatalf("Load(prefix) = %q", byPrefix.SessionID)
	}

	byIndex, _, err := store.Load("1")
	if err != nil {
		t.Fatalf("Load(index) error = %v", err)
	}
	if byIndex.SessionID != "20260508-091500-second" {
		t.Fatalf("Load(index) = %q", byIndex.SessionID)
	}
}

func TestLoadAmbiguousPrefix(t *testing.T) {
	store := NewStore(t.TempDir())
	base := Session{
		Version:       Version,
		Status:        StatusPaused,
		CreatedAt:     time.Date(2026, 5, 8, 9, 15, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2026, 5, 8, 9, 15, 0, 0, time.UTC),
		AgentsConfig:  []AgentConfig{},
		Rounds:        []Round{},
		Interventions: []Intervention{},
		Summaries:     map[string]string{},
	}
	base.SessionID = "20260508-091500-alpha"
	writeSession(t, store, base)
	base.SessionID = "20260508-091501-beta"
	writeSession(t, store, base)

	_, _, err := store.Load("20260508")
	if err == nil || !strings.Contains(err.Error(), ErrAmbiguous.Error()) {
		t.Fatalf("Load(ambiguous prefix) error = %v, want ambiguous error", err)
	}
}

func TestResumePersistsActiveStatus(t *testing.T) {
	store := NewStore(t.TempDir())
	writeSession(t, store, Session{
		SessionID:     "20260508-091500-draft",
		Version:       Version,
		Title:         "draft",
		Status:        StatusPaused,
		CreatedAt:     time.Date(2026, 5, 8, 9, 15, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2026, 5, 8, 9, 15, 0, 0, time.UTC),
		AgentsConfig:  []AgentConfig{},
		Rounds:        []Round{},
		Interventions: []Intervention{},
		Summaries:     map[string]string{},
	})

	resumed, _, err := store.Resume("20260508")
	if err != nil {
		t.Fatalf("Resume() error = %v", err)
	}
	if resumed.Status != StatusActive {
		t.Fatalf("resumed status = %q, want %q", resumed.Status, StatusActive)
	}
	loaded, _, err := store.Load("20260508-091500-draft")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Status != StatusActive {
		t.Fatalf("persisted status = %q, want %q", loaded.Status, StatusActive)
	}
}

func writeSession(t *testing.T, store *FileStore, session Session) {
	t.Helper()
	if err := os.MkdirAll(store.Dir(), 0o700); err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(store.Dir(), session.SessionID+".json"), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}
