package parser

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"
)

// Asset represents a discovered cryptographic file
type Asset struct {
	Type    string
	Status  string
	Expires string
}

// AuditData analyzes raw file bytes to extract certificates or private keys
func AuditData(data []byte) ([]Asset, bool, bool) {
	var assets []Asset
	var isCert, isKey bool

	// 1. Try Text-Based (PEM) Format
	rest := data
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				status, expiryStr := evaluateExpiry(cert.NotAfter)
				assets = append(assets, Asset{Type: "PEM Cert", Status: status, Expires: expiryStr})
				isCert = true
			}
		} else if strings.Contains(block.Type, "PRIVATE KEY") {
			assets = append(assets, Asset{Type: "⚠️ PRIVATE KEY", Status: "EXPOSED KEY", Expires: "N/A"})
			isKey = true
		}
	}

	if isCert || isKey {
		return assets, isCert, isKey
	}

	// 2. Try Binary (DER) Format Fallback
	if cert, err := x509.ParseCertificate(data); err == nil {
		status, expiryStr := evaluateExpiry(cert.NotAfter)
		assets = append(assets, Asset{Type: "DER Cert", Status: status, Expires: expiryStr})
		return assets, true, false
	}

	return nil, false, false
}

func evaluateExpiry(notAfter time.Time) (string, string) {
	expiryStr := notAfter.Format("2006-01-02")
	if time.Now().After(notAfter) {
		return "❌ EXPIRED", expiryStr
	}
	return "✅ Valid", expiryStr
}
