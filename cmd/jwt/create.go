package jwt

import (
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jose-util/generator"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Generate a JWKS",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return gen()
		},
	}
	return cmd
}

// https://github.com/go-jose/go-jose/blob/fdc2ceb0bbe2a29c582edfe07ea914c8dacd7e1b/jose-util/generate.go#L29
func gen() error {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	generateUseFlag := fs.String("use", "", "Desired public key usage (use header), one of [enc sig]")
	generateAlgFlag := fs.String("alg", "", "Desired key pair algorithm (alg header)")
	generateKeySizeFlag := fs.Int("size", 0, "Key size in bits (e.g. 2048 if generating an RSA key)")
	generateKeyIdentFlag := fs.String("kid", "", "Optional Key ID (kid header, generate random kid if not set)")

	var privKey crypto.PrivateKey
	var pubKey crypto.PublicKey
	var err error
	*generateAlgFlag = "RS256"
	*generateUseFlag = "sig"
	switch *generateUseFlag {
	case "sig":
		pubKey, privKey, err = generator.NewSigningKey(jose.SignatureAlgorithm(*generateAlgFlag), *generateKeySizeFlag)
	case "enc":
		pubKey, privKey, err = generator.NewEncryptionKey(jose.KeyAlgorithm(*generateAlgFlag), *generateKeySizeFlag)
	default:
		// According to RFC 7517 section-8.2.  This is unlikely to change in the
		// near future. If it were, new values could be found in the registry under
		// "JSON Web Key Use": https://www.iana.org/assignments/jose/jose.xhtml
		return fmt.Errorf("invalid key use '%s'.  Must be \"sig\" or \"enc\"", *generateUseFlag)
	}
	if err != nil {
		return fmt.Errorf("unable to generate key: %w", err)
	}

	pubData, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return err
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubData,
	})
	fmt.Printf("%s\n", string(pubBytes))

	privData, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return err
	}
	privBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privData,
	})
	fmt.Printf("%s\n", string(privBytes))

	kid := *generateKeyIdentFlag

	priv := jose.JSONWebKey{Key: privKey, KeyID: kid, Algorithm: *generateAlgFlag, Use: *generateUseFlag}

	// Generate a canonical kid based on RFC 7638
	if kid == "" {
		thumb, err := priv.Thumbprint(crypto.SHA256)
		if err != nil {
			return fmt.Errorf("unable to compute thumbprint: %w", err)
		}

		kid = base64.URLEncoding.EncodeToString(thumb)
		priv.KeyID = kid
	}

	sa := &ServiceAccount{
		Type:         "service_account",
		PrivateKeyId: kid,
		PrivateKey:   string(privBytes),
	}
	sbytes, err := json.MarshalIndent(sa, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(sbytes))
	// I'm not sure why we couldn't use `pub := priv.Public()` here as the private
	// key should contain the public key.  In case for some reason it doesn't,
	// this builds a public JWK from scratch.
	pub := jose.JSONWebKey{Key: pubKey, KeyID: kid, Algorithm: *generateAlgFlag, Use: *generateUseFlag}

	if priv.IsPublic() || !pub.IsPublic() || !priv.Valid() || !pub.Valid() {
		// This should never happen
		panic("invalid keys were generated")
	}

	privJSON, err := priv.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal private key to JSON: %w", err)
	}
	pubJSON, err := pub.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal public key to JSON: %w", err)
	}

	name := fmt.Sprintf("jwk-%s-%s", *generateUseFlag, kid)
	pubFile := fmt.Sprintf("%s-pub.json", name)
	privFile := fmt.Sprintf("%s-priv.json", name)

	err = writeNewFile(pubFile, pubJSON)
	if err != nil {
		return fmt.Errorf("error on write to file %s: %w", pubFile, err)
	}

	err = writeNewFile(privFile, privJSON)
	if err != nil {
		return fmt.Errorf("error on write to file %s: %w", privFile, err)
	}

	return nil
}

// Write new file to current dir
func writeNewFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

type ServiceAccount struct {
	Type                    string `json:"type,omitempty"`
	ProjectId               string `json:"project_id,omitempty"`
	PrivateKeyId            string `json:"private_key_id,omitempty"`
	PrivateKey              string `json:"private_key,omitempty"`
	ClientEmail             string `json:"client_email,omitempty"`
	ClientId                string `json:"client_id,omitempty"`
	AuthUri                 string `json:"auth_uri,omitempty"`
	TokenUri                string `json:"token_uri,omitempty"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientX509CertUrl       string `json:"client_x509_cert_url,omitempty"`
	UniverseDomain          string `json:"universe_domain,omitempty"`
}
