package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/agent-kit/q/pkg/jws"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2/google"
)

func generateCmd() *cobra.Command {
	var audience string
	var serviceAccountFile string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a JWT",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := generateJWT(serviceAccountFile, audience, 3600)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", s)
			return nil
		},
	}
	cmd.Flags().StringVar(&audience, "audience", "", "JWT audience")
	cmd.Flags().StringVar(&serviceAccountFile, "credentials", "", "service account credentials file")
	return cmd
}

// generateJWT creates a signed JSON Web Token using a Google API Service Account.
func generateJWT(saKeyfile, audience string, expiryLength int64) (string, error) {

	// Extract the RSA private key from the service account keyfile.
	sa, err := os.ReadFile(saKeyfile)
	if err != nil {
		return "", fmt.Errorf("could not read service account file: %w", err)
	}
	conf, err := google.JWTConfigFromJSON(sa)
	if err != nil {
		return "", fmt.Errorf("could not parse service account JSON: %w", err)
	}
	block, _ := pem.Decode(conf.PrivateKey)
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("private key parse error: %w", err)
	}
	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("private key failed rsa.PrivateKey type assertion")
	}

	saEmail := conf.Email
	now := time.Now().Unix()

	// Build the JWT payload.
	jwt := &jws.ClaimSet{
		Iat: now,
		// expires after 'expiryLength' seconds.
		Exp: now + expiryLength,
		// Iss must match 'issuer' in the security configuration in your
		// swagger spec (e.g. service account email). It can be any string.
		Iss: saEmail,
		// Aud must be either your Endpoints service name, or match the value
		// specified as the 'x-google-audience' in the OpenAPI document.
		Aud: audience,
		// Sub and Email should match the service account's email address.
		Sub:           saEmail,
		PrivateClaims: map[string]interface{}{"email": saEmail},
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}

	// Sign the JWT with the service account's private key.
	return jws.Encode(jwsHeader, jwt, rsaKey)
}
