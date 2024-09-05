package main

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"
)

func TestSSL(t *testing.T) {
	conn, err := tls.Dial("tcp", "api.oomph.ac:443", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := conn.VerifyHostname("api.oomph.ac"); err != nil {
		t.Fatalf("Bad hostname: %v", err)
	}

	for _, cert := range conn.ConnectionState().PeerCertificates {
		if time.Now().Before(cert.NotBefore) || time.Now().After(cert.NotAfter) {
			t.Fatalf(
				"certificate is invalid: %s (notBefore: %s, notAfter: %s)",
				time.Now().Format(time.RFC850),
				cert.NotBefore.Format(time.RFC850),
				cert.NotAfter.Format(time.RFC850),
			)
		}

		fmt.Println(cert.Issuer)
	}
}
