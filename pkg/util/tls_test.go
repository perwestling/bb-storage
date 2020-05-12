package util_test

import (
	"crypto/tls"
	"testing"

	configuration "github.com/buildbarn/bb-storage/pkg/proto/configuration/tls"
	"github.com/buildbarn/bb-storage/pkg/util"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// Example keypair generated by running:
	// openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:4096 -keyout key -out crt
	exampleCertificate = `
-----BEGIN CERTIFICATE-----
MIIEqDCCApACCQDTWa80U9UhejANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDDAtl
eGFtcGxlLmNvbTAeFw0yMDA0MDYxOTUwMDBaFw0yMTA0MDYxOTUwMDBaMBYxFDAS
BgNVBAMMC2V4YW1wbGUuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKC
AgEA6Cnp1BFSckXMTX+LPWXClN8IjZ3wca4KKsbOCx8WMhI1ziF+le9IEHrLiH8r
9k1E7M40d9wn0YqGLLXEMleor1HMR2s3aDS6A+TARPwrDNWxt3vmgMRXQxHgwuE/
PVyyGNUcP8Fm1nEiT8ngJClMhcIw9meJrUf6kk70CHYrXdVgUNaLlrIVYTCR38ki
dMXzu+DAby0cRoYUDluxJKCXfoochxJteXtO5PdEPQ0OLjj3f9dBSsviMTybYtlm
KSBxiyphn+3P59hYp+w4r0z4Kn23fj/ac3MpFcDFHHS6NP1W5WhxfSrb4O6dxDZI
1hE182xktblgx+tfcqgXsSFLczBIZQuBoRTeuYHTjtakQYIM4asvHnYyfEwau03a
+uM0Hl39lq1jqQNUnSFxWpyWtXA1ac/plMstZWuBd124pORACE1ZjnJS+3rJ+n9G
w25OSvSdmGw19uiFqRZMbAy4VjcokWiMRFkqQ4POuh5VuHr25uSO3JtnqKbJ6K6u
z/C4QdW/Yo+Sid08/+xIK+jFfkTJAStrprNo9Jw0t4RxWn1xGUW6LNIQpG9EjKFU
lA8WS+OwecUBrj+HPZa95heT+aiH9UNZgNXRzjConJOc9GXyH6UL2nBOCQSGfkS+
g9jwqt7NhTB4T3bAX7J1vFS5AOnDI99xujnz5kcv6qq7D40CAwEAATANBgkqhkiG
9w0BAQsFAAOCAgEArUkcdw0wHKx6kfVHOkVXTCU6mq/vuNOtF0lkDMRZ2FTX5yb2
9QGVPVPg7Ypp0QVsPXhfdjgRM+tFHC52zF4LzW/wfN64ywAQIZD0/C0pjsedRqSF
gz+HJhw3+aXCUIQTRpkKMmiuRH0Va2druKsMBRc9S09s2W1KF7Mw5OXA4xO+VpEv
+599F38N0Zv0H7uk761pol/jghcYlzaEhTK8S+oGjUAPF3BAymCZVtFY2kuyIr0J
O5qdcMWLkaiiu9SDCFZbOgrLgN3KbLoL0RiYB2AJLauwW89Y0vCkMM9+xRpjdPGB
r0vfZsgv1Crp3rhGEzdXadpgU8u/rqOzmCW11k78OSFiLkwzOxqZbvlNcZXSDcx2
pqckdNDUFVMdu1m7OjH4hW4theTaxK4++OURE1G/h9MNaKQAJjO+A9ozvuz7Aklv
bezRPmz/G2op75hQoSABNfZHavDw62WSXpDNguXW39GB7wWmAnOmbSb64DXddgAf
spTilKqMOzXPPJLviLphqNnHQse3mIm0be2hgqlBYnc96n+LoJm6sq0MZ7sukap8
pugFyoWaEHx90ECEfX3WXQYig4FZ6r5qmUlnrXFLzlgt/VrQOVdb7AuPc4XywGOg
qZcgecOjvl935nm4udYIBueBGhf6VnOgwbN7i8hr4IN2vd4+eF0tD7dOsGA=
-----END CERTIFICATE-----`
	examplePrivateKey = `
-----BEGIN PRIVATE KEY-----
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
ecgKVpPvVNRL4/3RQYXPEdErkwCshVk=
-----END PRIVATE KEY-----
`
)

func TestTLSConfigFromClientConfiguration(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		// When the TLS configuration is nil, TLS should be left
		// disabled.
		tlsConfig, err := util.NewTLSConfigFromClientConfiguration(nil)
		require.NoError(t, err)
		require.Nil(t, tlsConfig)
	})

	t.Run("Default", func(t *testing.T) {
		// The default configuration should enforce the use of
		// TLS 1.2 or higher.
		tlsConfig, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{})
		require.NoError(t, err)
		require.Equal(t, &tls.Config{
			MinVersion: tls.VersionTLS12,
		}, tlsConfig)
	})

	t.Run("ClientCertificate", func(t *testing.T) {
		tlsConfig, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{
				ClientCertificate: exampleCertificate,
				ClientPrivateKey:  examplePrivateKey,
			})
		require.NoError(t, err)
		require.Len(t, tlsConfig.Certificates, 1)
	})

	t.Run("InvalidClientCertificate", func(t *testing.T) {
		_, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{
				ClientCertificate: "This is an invalid certificate",
				ClientPrivateKey:  examplePrivateKey,
			})
		require.Equal(t, status.Error(codes.InvalidArgument, "Invalid client certificate or private key: tls: failed to find any PEM data in certificate input"), err)
	})

	t.Run("ServerCertificateAuthorities", func(t *testing.T) {
		tlsConfig, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{
				ServerCertificateAuthorities: exampleCertificate,
			})
		require.NoError(t, err)
		require.Len(t, tlsConfig.RootCAs.Subjects(), 1)
	})

	t.Run("InvalidServerCertificateAuthorities", func(t *testing.T) {
		// Because CertPool.AppendCertsFromPEM() does not return
		// a rich error message, we have no choice but to return
		// a simple error message in case of CA parsing failures.
		// https://github.com/golang/go/issues/23711#issuecomment-363322424
		_, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{
				ServerCertificateAuthorities: "This is an invalid certificate",
			})
		require.Equal(t, status.Error(codes.InvalidArgument, "Invalid server certificate authorities"), err)
	})

	t.Run("CustomCipherSuites", func(t *testing.T) {
		// Custom cipher suites should be respected.
		tlsConfig, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{
				CipherSuites: []string{
					"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
					"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
					"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
					"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				},
			})
		require.NoError(t, err)
		require.Equal(t, &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}, tlsConfig)
	})

	t.Run("InvalidCipherSuite", func(t *testing.T) {
		_, err := util.NewTLSConfigFromClientConfiguration(
			&configuration.ClientConfiguration{
				CipherSuites: []string{
					"TLS_ECDHE_ECDSA_WITH_AES_257_GCM_SHA385",
				},
			})
		require.Equal(t, status.Error(codes.InvalidArgument, "Unsupported cipher suite: \"TLS_ECDHE_ECDSA_WITH_AES_257_GCM_SHA385\""), err)
	})
}

func TestTLSConfigFromServerConfiguration(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		// When the TLS configuration is nil, TLS should be left
		// disabled.
		tlsConfig, err := util.NewTLSConfigFromServerConfiguration(nil)
		require.NoError(t, err)
		require.Nil(t, tlsConfig)
	})

	t.Run("Default", func(t *testing.T) {
		// The default configuration should enforce the use of
		// TLS 1.2 or higher.
		tlsConfig, err := util.NewTLSConfigFromServerConfiguration(
			&configuration.ServerConfiguration{
				ServerCertificate: exampleCertificate,
				ServerPrivateKey:  examplePrivateKey,
			})
		require.NoError(t, err)
		require.Len(t, tlsConfig.Certificates, 1)
		tlsConfig.Certificates = nil
		require.Equal(t, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientAuth: tls.RequestClientCert,
		}, tlsConfig)
	})

	t.Run("InvalidServerCertificate", func(t *testing.T) {
		_, err := util.NewTLSConfigFromServerConfiguration(
			&configuration.ServerConfiguration{
				ServerCertificate: "This is an invalid certificate",
				ServerPrivateKey:  examplePrivateKey,
			})
		require.Equal(t, status.Error(codes.InvalidArgument, "Invalid server certificate or private key: tls: failed to find any PEM data in certificate input"), err)
	})

	t.Run("CustomCipherSuites", func(t *testing.T) {
		// Custom cipher suites should be respected.
		tlsConfig, err := util.NewTLSConfigFromServerConfiguration(
			&configuration.ServerConfiguration{
				ServerCertificate: exampleCertificate,
				ServerPrivateKey:  examplePrivateKey,
				CipherSuites: []string{
					"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
					"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
					"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
					"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				},
			})
		require.NoError(t, err)
		require.Len(t, tlsConfig.Certificates, 1)
		tlsConfig.Certificates = nil
		require.Equal(t, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientAuth: tls.RequestClientCert,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}, tlsConfig)
	})

	t.Run("InvalidCipherSuite", func(t *testing.T) {
		_, err := util.NewTLSConfigFromServerConfiguration(
			&configuration.ServerConfiguration{
				ServerCertificate: exampleCertificate,
				ServerPrivateKey:  examplePrivateKey,
				CipherSuites: []string{
					"TLS_ECDHE_ECDSA_WITH_AES_257_GCM_SHA385",
				},
			})
		require.Equal(t, status.Error(codes.InvalidArgument, "Unsupported cipher suite: \"TLS_ECDHE_ECDSA_WITH_AES_257_GCM_SHA385\""), err)
	})
}
