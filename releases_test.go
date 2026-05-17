package main

import "testing"

func TestReleases(t *testing.T) {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
		return
	}

	manager, err := NewReleaseManager(cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = manager.DownloadAndUpdateLatestRelease()
	if err != nil {
		t.Fatal(err)
		return
	}
}
