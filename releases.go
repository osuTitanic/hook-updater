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
	context context.Context
	client  *github.Client
	config  *Config
}

func NewReleaseManager(config *Config) *ReleaseManager {
	return &ReleaseManager{
		context: context.Background(),
		client:  github.NewClient(nil),
		config:  config,
	}
}

func (manager *ReleaseManager) UpdateReleases() ([]*ReleaseMetadata, error) {
	releases, err := manager.FetchReleases()
	if err != nil {
		return nil, err
	}

	for _, release := range releases {
		targetFolder := filepath.Join(manager.config.ReleaseFolder(), release.TagName)
		release.DownloadAll(targetFolder)
	}

	return releases, nil
}

func (manager *ReleaseManager) FetchReleases() ([]*ReleaseMetadata, error) {
	repoReleases, _, err := manager.client.Repositories.ListReleases(
		manager.context,
		manager.config.GitHub.Owner,
		manager.config.GitHub.Repository,
		&github.ListOptions{PerPage: 100},
	)
	if err != nil {
		return nil, err
	}

	releases := make([]*ReleaseMetadata, 0, len(repoReleases))
	for _, release := range repoReleases {
		releaseMetadata := &ReleaseMetadata{
			Name:    release.GetName(),
			Commit:  release.GetTargetCommitish(),
			TagName: release.GetTagName(),
			Items:   make([]*ReleaseItem, len(release.Assets)),
		}

		for i, asset := range release.Assets {
			releaseMetadata.Items[i] = &ReleaseItem{
				Filename: asset.GetName(),
				Size:     asset.GetSize(),
				Url:      asset.GetBrowserDownloadURL(),
			}
		}

		releases = append(releases, releaseMetadata)
	}

	return releases, nil
}

type ReleaseItem struct {
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
	Size     int    `json:"size"`
	Url      string `json:"url"`
}

func (item *ReleaseItem) Download(targetFolder string) error {
	filePath := filepath.Join(targetFolder, item.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	response, err := http.Get(item.Url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	// Update checksum
	hash := sha256.New()
	hash.Write([]byte(item.Filename))
	item.Checksum = fmt.Sprintf("%x", hash.Sum(nil))
	return nil
}

func (item *ReleaseItem) DownloadIfNotExists(targetFolder string) error {
	filePath := filepath.Join(targetFolder, item.Filename)
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}
	return item.Download(targetFolder)
}

type ReleaseMetadata struct {
	Name    string         `json:"name"`
	Commit  string         `json:"commit"`
	TagName string         `json:"tag_name"`
	Items   []*ReleaseItem `json:"items"`
}

func (metadata *ReleaseMetadata) DownloadAll(targetFolder string) error {
	if err := os.MkdirAll(targetFolder, 0755); err != nil {
		return err
	}
	for _, item := range metadata.Items {
		if err := item.DownloadIfNotExists(targetFolder); err != nil {
			return err
		}
	}
	return nil
}
