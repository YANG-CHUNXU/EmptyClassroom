package snapshot

import (
	"EmptyClassroom/service/model"
	"EmptyClassroom/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	BlobTokenEnvKey       = "BLOB_READ_WRITE_TOKEN"
	DefaultSnapshotPath   = "snapshots/today.json"
	DefaultLocalStorePath = ".empty-classroom-snapshot.json"

	blobAPIURL = "https://vercel.com/api/blob"
)

var ErrSnapshotNotFound = errors.New("snapshot not found")

type Store interface {
	Load(ctx context.Context) (*model.ClassInfo, error)
	Save(ctx context.Context, classInfo *model.ClassInfo) error
}

type blobStore struct {
	token    string
	pathname string
	client   *http.Client
}

type fileStore struct {
	path string
}

func NewDefaultStore() Store {
	token := os.Getenv(BlobTokenEnvKey)
	if token != "" {
		return NewBlobStore(token, DefaultSnapshotPath)
	}
	return NewFileStore(DefaultLocalStorePath)
}

func NewBlobStore(token string, pathname string) Store {
	return &blobStore{
		token:    token,
		pathname: pathname,
		client:   utils.OutboundHTTPClient(),
	}
}

func NewFileStore(path string) Store {
	return &fileStore{path: path}
}

func (s *blobStore) Load(ctx context.Context) (*model.ClassInfo, error) {
	blobURL, err := s.blobURL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, blobURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrSnapshotNotFound
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(strings.TrimSpace(string(body)))
	}

	return decodeSnapshot(resp.Body)
}

func (s *blobStore) Save(ctx context.Context, classInfo *model.ClassInfo) error {
	body, err := json.Marshal(classInfo)
	if err != nil {
		return err
	}

	apiURL := blobAPIURL + "/?pathname=" + url.QueryEscape(s.pathname)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-vercel-blob-access", "private")
	req.Header.Set("x-add-random-suffix", "0")
	req.Header.Set("x-allow-overwrite", "1")
	req.Header.Set("x-content-type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return errors.New(strings.TrimSpace(string(respBody)))
	}

	return nil
}

func (s *blobStore) blobURL() (string, error) {
	parts := strings.Split(s.token, "_")
	storeID := ""
	if len(parts) >= 4 {
		storeID = parts[3]
	}
	if storeID == "" {
		return "", errors.New("invalid blob token")
	}
	return "https://" + storeID + ".private.blob.vercel-storage.com/" + s.pathname, nil
}

func (s *fileStore) Load(_ context.Context) (*model.ClassInfo, error) {
	raw, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrSnapshotNotFound
	}
	if err != nil {
		return nil, err
	}

	return decodeSnapshot(bytes.NewReader(raw))
}

func (s *fileStore) Save(_ context.Context, classInfo *model.ClassInfo) error {
	body, err := json.Marshal(classInfo)
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(s.path, body, 0o644)
}

func decodeSnapshot(reader io.Reader) (*model.ClassInfo, error) {
	classInfo := new(model.ClassInfo)
	if err := json.NewDecoder(reader).Decode(classInfo); err != nil {
		return nil, err
	}
	return classInfo, nil
}
