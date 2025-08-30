package main

import "testing"

func TestReleases(t *testing.T) {
	cfg, err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
		return
	}

	manager := NewReleaseManager(cfg)
	_, err = manager.UpdateReleases()
	if err != nil {
		t.Fatal(err)
		return
	}
}
