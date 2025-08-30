package main

import "testing"

func TestReleases(t *testing.T) {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
		return
	}

	manager := NewReleaseManager(cfg)
	err = manager.DownloadAndUpdateLatestRelease()
	if err != nil {
		t.Fatal(err)
		return
	}
}
