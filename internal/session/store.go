package session

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const Version = "1.0"

var (
	ErrNotFound  = errors.New("SES-001: session not found")
	ErrAmbiguous = errors.New("SES-001: session reference is ambiguous")
)

type Store interface {
	Dir() string
	Create(title, draftPath string) (Session, string, error)
	Save(Session) (string, error)
	Load(ref string) (Session, string, error)
	Find(ref string) (Session, string, error)
	Resume(ref string) (Session, string, error)
	List(ListOptions) ([]SessionSummary, error)
}

type ListOptions struct {
	Limit  int
	Status Status
}

type FileStore struct {
	dir string
}

func NewStore(dir string) *FileStore {
	return &FileStore{dir: dir}
}

func (s *FileStore) Dir() string {
	return s.dir
}

func (s *FileStore) Create(title, draftPath string) (Session, string, error) {
	now := time.Now().UTC()
	session := Session{
		SessionID:     s.nextID(now, title),
		Version:       Version,
		Title:         title,
		Status:        StatusActive,
		CreatedAt:     now,
		UpdatedAt:     now,
		AgentsConfig:  []AgentConfig{},
		Rounds:        []Round{},
		Interventions: []Intervention{},
		Summaries:     map[string]string{},
	}

	if draftPath != "" {
		source, err := sourceFromFile(draftPath)
		if err != nil {
			return Session{}, "", err
		}
		session.Source = &source
	}

	path, err := s.Save(session)
	return session, path, err
}

func (s *FileStore) Save(session Session) (string, error) {
	if err := validateSession(session); err != nil {
		return "", err
	}
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return "", fmt.Errorf("SES-002: create session directory: %w", err)
	}

	session.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return "", fmt.Errorf("SES-002: encode session: %w", err)
	}
	data = append(data, '\n')

	path := s.pathFor(session.SessionID)
	tmp, err := os.CreateTemp(s.dir, ".woolf-*.tmp")
	if err != nil {
		return "", fmt.Errorf("SES-002: create temporary session file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", fmt.Errorf("SES-002: write temporary session file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("SES-002: close temporary session file: %w", err)
	}
	if err := os.Chmod(tmpPath, 0o600); err != nil {
		return "", fmt.Errorf("SES-002: secure temporary session file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return "", fmt.Errorf("SES-002: save session file: %w", err)
	}
	return path, nil
}

func (s *FileStore) Load(ref string) (Session, string, error) {
	path, err := s.resolve(ref)
	if err != nil {
		return Session{}, "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Session{}, "", fmt.Errorf("SES-002: read session file: %w", err)
	}
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return Session{}, "", fmt.Errorf("SES-003: parse session file %s: %w", path, err)
	}
	if err := validateSession(session); err != nil {
		return Session{}, "", fmt.Errorf("SES-003: invalid session file %s: %w", path, err)
	}
	return session, path, nil
}

func (s *FileStore) Find(ref string) (Session, string, error) {
	return s.Load(ref)
}

func (s *FileStore) Resume(ref string) (Session, string, error) {
	session, _, err := s.Load(ref)
	if err != nil {
		return Session{}, "", err
	}
	if session.Status == StatusCompleted {
		return Session{}, "", fmt.Errorf("SES-001: session %s is completed", session.SessionID)
	}
	session.Status = StatusActive
	path, err := s.Save(session)
	return session, path, err
}

func (s *FileStore) List(opts ListOptions) ([]SessionSummary, error) {
	paths, err := s.sessionPaths()
	if err != nil {
		return nil, err
	}

	items := make([]SessionSummary, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var sess Session
		if err := json.Unmarshal(data, &sess); err != nil {
			continue
		}
		if opts.Status != "" && sess.Status != opts.Status {
			continue
		}
		items = append(items, SessionSummary{
			SessionID: sess.SessionID,
			Title:     sess.Title,
			Status:    sess.Status,
			Rounds:    sess.Totals.RoundsCompleted,
			Path:      path,
			UpdatedAt: sess.UpdatedAt,
		})
		if opts.Limit > 0 && len(items) == opts.Limit {
			break
		}
	}
	return items, nil
}

func (s *FileStore) resolve(ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", ErrNotFound
	}

	paths, err := s.sessionPaths()
	if err != nil {
		return "", err
	}
	var matches []string
	for _, path := range paths {
		id := strings.TrimSuffix(filepath.Base(path), ".json")
		if id == ref {
			return path, nil
		}
		if strings.HasPrefix(id, ref) {
			matches = append(matches, path)
		}
	}
	switch len(matches) {
	case 0:
		if index, err := strconv.Atoi(ref); err == nil {
			if index < 1 {
				return "", ErrNotFound
			}
			items, err := s.List(ListOptions{})
			if err != nil {
				return "", err
			}
			if index > len(items) {
				return "", ErrNotFound
			}
			return items[index-1].Path, nil
		}
		return "", ErrNotFound
	case 1:
		return matches[0], nil
	default:
		ids := make([]string, 0, len(matches))
		for _, path := range matches {
			ids = append(ids, strings.TrimSuffix(filepath.Base(path), ".json"))
		}
		return "", fmt.Errorf("%w: %s", ErrAmbiguous, strings.Join(ids, ", "))
	}
}

func (s *FileStore) sessionPaths() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, err
	}
	sort.Slice(matches, func(i, j int) bool {
		return strings.TrimSuffix(filepath.Base(matches[i]), ".json") > strings.TrimSuffix(filepath.Base(matches[j]), ".json")
	})
	return matches, nil
}

func (s *FileStore) pathFor(id string) string {
	return filepath.Join(s.dir, id+".json")
}

func (s *FileStore) nextID(now time.Time, title string) string {
	base := now.Format("20060102-150405") + "-" + slugify(title)
	id := base
	for i := 2; ; i++ {
		if _, err := os.Stat(s.pathFor(id)); errors.Is(err, os.ErrNotExist) {
			return id
		}
		id = fmt.Sprintf("%s-%d", base, i)
	}
}

func sourceFromFile(path string) (Source, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Source{}, fmt.Errorf("ING-001: read draft file: %w", err)
	}
	sum := sha256.Sum256(data)
	return Source{
		Type:           "file",
		Path:           path,
		ContentHash:    hex.EncodeToString(sum[:]),
		ContentPreview: preview(string(data), 200),
	}, nil
}

func validateSession(session Session) error {
	if session.SessionID == "" {
		return errors.New("missing session_id")
	}
	if !validSessionID(session.SessionID) {
		return fmt.Errorf("invalid session_id %q", session.SessionID)
	}
	if session.Version != Version {
		return fmt.Errorf("unsupported version %q", session.Version)
	}
	switch session.Status {
	case StatusActive, StatusPaused, StatusCompleted, StatusError:
	default:
		return fmt.Errorf("invalid status %q", session.Status)
	}
	if session.CreatedAt.IsZero() {
		return errors.New("missing created_at")
	}
	if session.AgentsConfig == nil {
		return errors.New("missing agents_config")
	}
	if session.Rounds == nil {
		return errors.New("missing rounds")
	}
	return nil
}

func validSessionID(id string) bool {
	parts := strings.SplitN(id, "-", 3)
	if len(parts) != 3 || len(parts[0]) != 8 || len(parts[1]) != 6 || parts[2] == "" {
		return false
	}
	for _, r := range parts[0] + parts[1] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteByte('-')
			lastDash = true
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "session"
	}
	return slug
}

func preview(value string, maxRunes int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes])
}

type SessionSummary struct {
	SessionID string
	Title     string
	Status    Status
	Rounds    int
	Path      string
	UpdatedAt time.Time
}
