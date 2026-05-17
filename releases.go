package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

type ReleaseManager struct {
	LatestRelease *ReleaseMetadata

	context  context.Context
	client   *github.Client
	config   *Config
	verifier *SignatureVerifier
}

func NewReleaseManager(config *Config) (*ReleaseManager, error) {
	verifier, err := NewSignatureVerifier(config)
	if err != nil {
		return nil, err
	}

	return &ReleaseManager{
		context:  context.Background(),
		client:   github.NewClient(nil),
		config:   config,
		verifier: verifier,
	}, nil
}

func (manager *ReleaseManager) DownloadAndUpdateLatestRelease() error {
	release, err := manager.fetchLatestRelease()
	if err != nil {
		return err
	}

	targetFolder := filepath.Join(
		manager.config.ReleaseFolder(),
		release.TagName,
	)
	if err := release.DownloadAll(targetFolder, manager.verifier); err != nil {
		return err
	}
	manager.LatestRelease = release
	return nil
}

func (manager *ReleaseManager) fetchLatestRelease() (*ReleaseMetadata, error) {
	repoRelease, _, err := manager.client.Repositories.GetLatestRelease(
		manager.context,
		manager.config.GitHub.Owner,
		manager.config.GitHub.Repository,
	)
	if err != nil {
		return nil, err
	}

	metadata := &ReleaseMetadata{
		Name:    repoRelease.GetName(),
		Commit:  repoRelease.GetTargetCommitish(),
		TagName: repoRelease.GetTagName(),
		Items:   make([]*ReleaseItem, len(repoRelease.Assets)),
	}
	for i, asset := range repoRelease.Assets {
		metadata.Items[i] = &ReleaseItem{
			Filename: asset.GetName(),
			Size:     asset.GetSize(),
			Url:      asset.GetBrowserDownloadURL(),
		}
	}
	return metadata, nil
}

type ReleaseItem struct {
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
}

func (item *ReleaseItem) Download(targetFolder string, verifier *SignatureVerifier) error {
	filePath := filepath.Join(targetFolder, item.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	response, err := http.Get(item.Url)
	if err != nil {
		out.Close()
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		out.Close()
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	if err := verifier.Verify(filePath); err != nil {
		os.Remove(filePath)
		return fmt.Errorf("signature verification failed for %s: %w", item.Filename, err)
	}

	checksum, err := calculateFileChecksum(filePath)
	if err != nil {
		return err
	}
	item.Checksum = checksum
	return nil
}

func (item *ReleaseItem) DownloadIfNotExists(targetFolder string, verifier *SignatureVerifier) error {
	filePath := filepath.Join(targetFolder, item.Filename)
	if _, err := os.Stat(filePath); err != nil {
		// Download file from GitHub
		return item.Download(targetFolder, verifier)
	}

	if err := verifier.Verify(filePath); err != nil {
		return fmt.Errorf("signature verification failed for cached %s: %w", item.Filename, err)
	}

	checksum, err := calculateFileChecksum(filePath)
	if err != nil {
		return err
	}
	item.Checksum = checksum
	return nil
}

func calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

type ReleaseMetadata struct {
	Name    string         `json:"name"`
	Commit  string         `json:"commit"`
	TagName string         `json:"tag_name"`
	Items   []*ReleaseItem `json:"items"`
}

func (metadata *ReleaseMetadata) DownloadAll(targetFolder string, verifier *SignatureVerifier) error {
	// Ensure target folder exists
	if err := os.MkdirAll(targetFolder, 0755); err != nil {
		return err
	}
	for _, item := range metadata.Items {
		if err := item.DownloadIfNotExists(targetFolder, verifier); err != nil {
			return err
		}
	}
	return nil
}

func (metadata *ReleaseMetadata) GetItemByFilename(filename string) *ReleaseItem {
	for _, item := range metadata.Items {
		if item.Filename == filename {
			return item
		}
	}
	return nil
}
