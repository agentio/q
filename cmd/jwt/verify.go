package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/agent-kit/q/pkg/jws"
	"github.com/spf13/cobra"
)

type KeySet struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	E   string `json:"e"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kty string `json:"rsa"`
	N   string `json:"n"`
	Kid string `json:"kid"`
}

func verifyCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a JWT signed by a Google Service Account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			token := args[0]
			claims, err := jws.Decode(token)
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(claims, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s", string(b))

			// now get the private key from Google's service accounts service
			keyUrl := "https://www.googleapis.com/service_accounts/v1/jwk/" + claims.Sub
			response, err := http.Get(keyUrl)
			if err != nil {
				return err
			}
			responseBytes, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			var keySet KeySet
			err = json.Unmarshal(responseBytes, &keySet)
			if err != nil {
				return err
			}
			b, err = json.MarshalIndent(keySet, "", "  ")
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			//slices.Reverse(keySet.Keys)
			for i, k := range keySet.Keys {

				N := k.N
				decN, err := base64.RawURLEncoding.DecodeString(N)
				if err != nil {
					return err
				}
				n := big.NewInt(0)
				n.SetBytes(decN)

				E := k.E
				decE, err := base64.RawURLEncoding.DecodeString(E)
				if err != nil {
					return err
				}
				e := big.NewInt(0)
				e.SetBytes(decE)

				key := &rsa.PublicKey{
					N: n,
					E: int(e.Int64()),
				}
				fmt.Printf("verifying with key %d\n", i)
				err = jws.Verify(token, key)
				if err == nil {
					return nil
				}
			}
			return errors.New("unable to verify signature")
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "output format")
	return cmd
}
