package main

import (
	"flag"
	"os/exec"
	"testing"
)

var signatureExePath = flag.String("exe", "", "path to the signed .exe to verify")

func TestSignatureVerification(t *testing.T) {
	if *signatureExePath == "" {
		t.Skip("missing -exe flag for signature verification test")
		return
	}

	cfg, err := LoadConfig("config.json")
	if err != nil {
		t.Skip("config not available for signature verification test: ", err)
		return
	}
	if _, err := exec.LookPath(cfg.SignatureVerification.OsslsigncodePath); err != nil {
		t.Skip("osslsigncode not available: ", err)
		return
	}

	// Force signature verification to be enabled for this test
	cfg.SignatureVerification.Enabled = true

	verifier, err := NewSignatureVerifier(cfg)
	if err != nil {
		t.Fatalf("initialize signature verifier: %v", err)
	}
	if !verifier.ShouldVerify(*signatureExePath) {
		t.Fatalf("%q does not match the configured extensions", *signatureExePath)
	}

	if err := verifier.Verify(*signatureExePath); err != nil {
		t.Fatalf("verify signature of %q: %v", *signatureExePath, err)
	}
}
