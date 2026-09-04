package crawler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	root := t.TempDir()

	certDir := filepath.Join(root, "certs")
	nestedDir := filepath.Join(certDir, "nested")

	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("create test directories: %v", err)
	}

	validCert := generateTestCertificate(t, time.Now().Add(24*time.Hour))
	expiredCert := generateTestCertificate(t, time.Now().Add(-24*time.Hour))

	writeTestFile(
		t,
		filepath.Join(certDir, "valid.pem"),
		pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: validCert,
		}),
	)

	writeTestFile(
		t,
		filepath.Join(nestedDir, "expired.crt"),
		pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: expiredCert,
		}),
	)

	writeTestFile(
		t,
		filepath.Join(nestedDir, "private.key"),
		[]byte(`-----BEGIN PRIVATE KEY-----
dGVzdC1wcml2YXRlLWtleQ==
-----END PRIVATE KEY-----`),
	)

	writeTestFile(
		t,
		filepath.Join(root, "ignored.txt"),
		[]byte("not a certificate"),
	)

	result, err := Run(root)
	if err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	if result.CertificateCount != 2 {
		t.Errorf(
			"expected 2 certificates, got %d",
			result.CertificateCount,
		)
	}

	if result.PrivateKeyCount != 1 {
		t.Errorf(
			"expected 1 private key, got %d",
			result.PrivateKeyCount,
		)
	}

	if result.ValidCount != 1 {
		t.Errorf(
			"expected 1 valid certificate, got %d",
			result.ValidCount,
		)
	}

	if result.ExpiredCount != 1 {
		t.Errorf(
			"expected 1 expired certificate, got %d",
			result.ExpiredCount,
		)
	}

	if len(result.Results) != 3 {
		t.Errorf(
			"expected 3 results, got %d",
			len(result.Results),
		)
	}
}

func TestHasScanExtension(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"certificate.crt", true},
		{"certificate.CRT", true},
		{"certificate.pem", true},
		{"certificate.der", true},
		{"certificate.cer", true},
		{"private.key", true},
		{"document.txt", false},
		{"certificate", false},
		{"key.pem.bak", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := hasScanExtension(tt.path)

			if got != tt.want {
				t.Errorf(
					"hasScanExtension(%q) = %v, want %v",
					tt.path,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestRunMissingDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "does-not-exist")

	result, err := Run(root)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Results) != 0 {
		t.Fatalf(
			"expected no results, got %d",
			len(result.Results),
		)
	}

	if result.CertificateCount != 0 {
		t.Errorf(
			"expected 0 certificates, got %d",
			result.CertificateCount,
		)
	}

	if result.PrivateKeyCount != 0 {
		t.Errorf(
			"expected 0 private keys, got %d",
			result.PrivateKeyCount,
		)
	}
}

func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write test file %q: %v", path, err)
	}
}

func generateTestCertificate(t *testing.T, expiry time.Time) []byte {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}

	serialNumber, err := rand.Int(
		rand.Reader,
		new(big.Int).Lsh(big.NewInt(1), 128),
	)
	if err != nil {
		t.Fatalf("generate serial number: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "certmole-test",
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  expiry,
	}

	data, err := x509.CreateCertificate(
		rand.Reader,
		template,
		template,
		&privateKey.PublicKey,
		privateKey,
	)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	return data
}
