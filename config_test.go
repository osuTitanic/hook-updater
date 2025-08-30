package main

import (
	"testing"
)

func TestConfig(t *testing.T) {
	err := EnsureConfigExists()
	if err != nil {
		t.Fatal(err)
		return
	}

	cfg, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Logf("GitHub: %s/%s", cfg.GitHub.Owner, cfg.GitHub.Repository)
	t.Logf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	t.Logf("Data Folder: %s\n", cfg.DataFolder)

	for _, release := range cfg.ReleaseTypes {
		t.Logf("Release: %s -> %s\n", release.Identifier, release.Filename)
	}
}
