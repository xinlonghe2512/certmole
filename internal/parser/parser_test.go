package parser

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestAuditDataPEMCertificate(t *testing.T) {
	tests := []struct {
		name   string
		expiry time.Time
		status string
	}{
		{
			name:   "valid certificate",
			expiry: time.Date(2030, 1, 15, 12, 0, 0, 0, time.Local),
			status: "✅ Valid",
		},
		{
			name:   "expired certificate",
			expiry: time.Date(2020, 1, 15, 12, 0, 0, 0, time.Local),
			status: "❌ EXPIRED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := generatePEMCertificate(t, tt.expiry)

			assets, isCert, isKey := AuditData(data)

			if !isCert {
				t.Fatal("expected isCert to be true")
			}

			if isKey {
				t.Fatal("expected isKey to be false")
			}

			if len(assets) != 1 {
				t.Fatalf("expected 1 asset, got %d", len(assets))
			}

			asset := assets[0]

			if asset.Type != "PEM Cert" {
				t.Errorf(
					"expected type %q, got %q",
					"PEM Cert",
					asset.Type,
				)
			}

			if asset.Status != tt.status {
				t.Errorf(
					"expected status %q, got %q",
					tt.status,
					asset.Status,
				)
			}

			expectedExpiry := tt.expiry.Format("2006-01-02")

			if asset.Expires != expectedExpiry {
				t.Errorf(
					"expected expiry %q, got %q",
					expectedExpiry,
					asset.Expires,
				)
			}
		})
	}
}

func TestAuditDataDERCertificate(t *testing.T) {
	expiry := time.Date(2030, 1, 15, 12, 0, 0, 0, time.Local)
	data := generateDERCertificate(t, expiry)

	assets, isCert, isKey := AuditData(data)

	if !isCert {
		t.Fatal("expected isCert to be true")
	}

	if isKey {
		t.Fatal("expected isKey to be false")
	}

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets))
	}

	if assets[0].Type != "DER Cert" {
		t.Errorf("expected type %q, got %q", "DER Cert", assets[0].Type)
	}

	if assets[0].Status != "✅ Valid" {
		t.Errorf("expected valid certificate, got %q", assets[0].Status)
	}
}

func TestAuditDataPrivateKey(t *testing.T) {
	data := []byte(`-----BEGIN PRIVATE KEY-----
dGVzdC1wcml2YXRlLWtleQ==
-----END PRIVATE KEY-----`)

	assets, isCert, isKey := AuditData(data)

	if isCert {
		t.Fatal("expected isCert to be false")
	}

	if !isKey {
		t.Fatal("expected isKey to be true")
	}

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets))
	}

	asset := assets[0]

	if asset.Type != "⚠️ PRIVATE KEY" {
		t.Errorf("expected private key type, got %q", asset.Type)
	}

	if asset.Status != "EXPOSED KEY" {
		t.Errorf("expected status %q, got %q", "EXPOSED KEY", asset.Status)
	}

	if asset.Expires != "N/A" {
		t.Errorf("expected expiry %q, got %q", "N/A", asset.Expires)
	}
}

func TestAuditDataInvalidData(t *testing.T) {
	assets, isCert, isKey := AuditData([]byte("not a certificate or key"))

	if assets != nil {
		t.Fatalf("expected nil assets, got %v", assets)
	}

	if isCert {
		t.Fatal("expected isCert to be false")
	}

	if isKey {
		t.Fatal("expected isKey to be false")
	}
}

func TestAuditDataMultiplePEMBlocks(t *testing.T) {
	expiry := time.Date(2030, 1, 15, 12, 0, 0, 0, time.Local)
	cert := generatePEMCertificate(t, expiry)

	key := []byte(`-----BEGIN PRIVATE KEY-----
dGVzdC1wcml2YXRlLWtleQ==
-----END PRIVATE KEY-----`)

	data := append(cert, key...)

	assets, isCert, isKey := AuditData(data)

	if !isCert {
		t.Fatal("expected isCert to be true")
	}

	if !isKey {
		t.Fatal("expected isKey to be true")
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(assets))
	}

	if assets[0].Type != "PEM Cert" {
		t.Errorf("expected first asset to be PEM Cert, got %q", assets[0].Type)
	}

	if assets[1].Type != "⚠️ PRIVATE KEY" {
		t.Errorf(
			"expected second asset to be PRIVATE KEY, got %q",
			assets[1].Type,
		)
	}
}

func TestEvaluateExpiry(t *testing.T) {
	tests := []struct {
		name   string
		time   time.Time
		status string
	}{
		{
			name:   "future",
			time:   time.Now().Add(time.Hour),
			status: "✅ Valid",
		},
		{
			name:   "past",
			time:   time.Now().Add(-time.Hour),
			status: "❌ EXPIRED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, expiry := evaluateExpiry(tt.time)

			if status != tt.status {
				t.Errorf(
					"expected status %q, got %q",
					tt.status,
					status,
				)
			}

			expectedExpiry := tt.time.Format("2006-01-02")
			if expiry != expectedExpiry {
				t.Errorf(
					"expected expiry %q, got %q",
					expectedExpiry,
					expiry,
				)
			}
		})
	}
}

func generatePEMCertificate(t *testing.T, expiry time.Time) []byte {
	t.Helper()

	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: generateDERCertificate(t, expiry),
	})
}

func generateDERCertificate(t *testing.T, expiry time.Time) []byte {
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

	notBefore := time.Now().Add(-time.Hour)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "certmole-test",
		},
		NotBefore: notBefore,
		NotAfter:  expiry,
		KeyUsage:  x509.KeyUsageDigitalSignature,
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
