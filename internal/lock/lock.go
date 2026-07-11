package lock

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type LockFile struct {
	Version   string          `json:"version"`
	Generated string          `json:"generated"`
	Artifacts []LockArtifact  `json:"artifacts"`
	Checksum  string          `json:"checksum"`
}

type LockArtifact struct {
	Path     string `json:"path"`
	Size     int    `json:"size"`
	Checksum string `json:"checksum"`
}

type ArtifactInfo struct {
	Path    string
	Content []byte
}

func Generate(artifacts []ArtifactInfo) (*LockFile, error) {
	lock := &LockFile{
		Version:   "1",
		Generated: time.Now().UTC().Format(time.RFC3339),
	}

	for _, a := range artifacts {
		hash := sha256.Sum256(a.Content)
		lock.Artifacts = append(lock.Artifacts, LockArtifact{
			Path:     a.Path,
			Size:     len(a.Content),
			Checksum: hex.EncodeToString(hash[:]),
		})
	}

	sort.Slice(lock.Artifacts, func(i, j int) bool {
		return lock.Artifacts[i].Path < lock.Artifacts[j].Path
	})

	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal lock file: %w", err)
	}

	hash := sha256.Sum256(data)
	lock.Checksum = hex.EncodeToString(hash[:])

	return lock, nil
}

func WriteToFile(lock *LockFile, path string) error {
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lock file: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

func ReadFromFile(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lock file: %w", err)
	}
	var lock LockFile
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("parse lock file: %w", err)
	}
	return &lock, nil
}

func Verify(lock *LockFile, current []ArtifactInfo) ([]string, error) {
	if lock == nil {
		return nil, fmt.Errorf("lock file is nil")
	}

	existing := make(map[string]LockArtifact)
	for _, a := range lock.Artifacts {
		existing[a.Path] = a
	}

	var changes []string
	currentMap := make(map[string]bool)

	for _, a := range current {
		currentMap[a.Path] = true
		hash := sha256.Sum256(a.Content)
		checksum := hex.EncodeToString(hash[:])

		if old, ok := existing[a.Path]; ok {
			if old.Checksum != checksum {
				changes = append(changes, fmt.Sprintf("modified: %s", a.Path))
			}
		} else {
			changes = append(changes, fmt.Sprintf("added: %s", a.Path))
		}
	}

	for _, a := range lock.Artifacts {
		if !currentMap[a.Path] {
			changes = append(changes, fmt.Sprintf("removed: %s", a.Path))
		}
	}

	return changes, nil
}
