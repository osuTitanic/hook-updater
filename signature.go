package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultOsslsigncodePath = "osslsigncode"
	defaultVerifyTimeout    = 30 * time.Second
)

type SignatureVerifier struct {
	Enabled            bool
	OsslsigncodePath   string
	RequiredLeafSha256 string
	RequiredExtensions map[string]bool
}

func NewSignatureVerifier(config *Config) (*SignatureVerifier, error) {
	cfg := config.SignatureVerification
	path := cfg.OsslsigncodePath
	if path == "" {
		path = defaultOsslsigncodePath
	}

	extensions := cfg.RequiredExtensions
	if len(extensions) == 0 {
		extensions = []string{".exe"}
	}

	verifier := &SignatureVerifier{
		Enabled:            cfg.Enabled,
		OsslsigncodePath:   path,
		RequiredLeafSha256: normalizeFingerprint(cfg.RequiredLeafSha256),
		RequiredExtensions: make(map[string]bool, len(extensions)),
	}

	for _, extension := range extensions {
		extension = strings.ToLower(strings.TrimSpace(extension))
		if extension == "" {
			continue
		}
		if !strings.HasPrefix(extension, ".") {
			extension = "." + extension
		}
		verifier.RequiredExtensions[extension] = true
	}

	if verifier.Enabled && verifier.RequiredLeafSha256 == "" {
		return nil, fmt.Errorf("signature verification requires RequiredLeafSha256")
	}
	if verifier.Enabled {
		if _, err := exec.LookPath(verifier.OsslsigncodePath); err != nil {
			return nil, fmt.Errorf("signature verification requires osslsigncode at %q: %w", verifier.OsslsigncodePath, err)
		}
	}

	return verifier, nil
}

func (verifier *SignatureVerifier) ShouldVerify(filename string) bool {
	if !verifier.Enabled {
		return false
	}
	return verifier.RequiredExtensions[strings.ToLower(filepath.Ext(filename))]
}

func (verifier *SignatureVerifier) Verify(path string) error {
	if !verifier.ShouldVerify(path) {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultVerifyTimeout)
	defer cancel()
	var output bytes.Buffer

	command := exec.CommandContext(
		ctx,
		verifier.OsslsigncodePath,
		"verify",
		"-require-leaf-hash",
		"sha256:"+verifier.RequiredLeafSha256,
		"-in",
		path,
	)
	command.Stdout = &output
	command.Stderr = &output

	if err := command.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("osslsigncode timed out after %s", defaultVerifyTimeout)
		}
		outputTrimmed := trimVerifierOutput(output.String())
		return fmt.Errorf("osslsigncode verify failed: %w: %s", err, outputTrimmed)
	}

	return nil
}

func normalizeFingerprint(fingerprint string) string {
	fingerprint = strings.ReplaceAll(fingerprint, ":", "")
	fingerprint = strings.ReplaceAll(fingerprint, " ", "")
	return strings.ToUpper(strings.TrimSpace(fingerprint))
}

func trimVerifierOutput(output string) string {
	output = strings.TrimSpace(output)
	if len(output) <= 4096 {
		return output
	}
	return output[:4096] + "... (truncated)"
}
