package main

import (
	"log"
	"time"
)

func ReleaseUpdateLoop(manager *ReleaseManager) {
	logger := CreateLogger("ReleaseManager", DEBUG)
	sleepDuration, err := time.ParseDuration(manager.config.UpdateInterval)

	if err != nil {
		// Parsing failed, fallback to 5 minutes
		logger.Warning("Failed to parse update interval:", err)
		sleepDuration = 5 * time.Minute
	}

	for {
		err := manager.DownloadAndUpdateLatestRelease()
		if err != nil {
			logger.Error("Failed to update release:", err)
			time.Sleep(sleepDuration)
			continue
		}
		logger.Infof("Updated to release '%s'", manager.LatestRelease.Name)
		time.Sleep(sleepDuration)
	}
}

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	manager, err := NewReleaseManager(config)
	if err != nil {
		log.Fatal("Failed to initialize release manager:", err)
	}
	go ReleaseUpdateLoop(manager)

	server := NewServer(config, manager)
	server.Serve()
}
