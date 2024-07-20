package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	longrunning "cloud.google.com/go/longrunning/autogen"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func newConn(host string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if host != "" {
		opts = append(opts, grpc.WithAuthority(host))
	}
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	return grpc.NewClient(host, opts...)
}

type Credentials struct {
	QuotaProjectID string `json:"quota_project_id"`
}

func getADC(ctx context.Context) (string, string, error) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", "", err
	}
	var credentials Credentials
	json.Unmarshal(creds.JSON, &credentials)
	t, err := creds.TokenSource.Token()
	if err != nil {
		return "", "", err
	}
	return t.AccessToken, credentials.QuotaProjectID, nil
}

func ApiKeysClient(ctx context.Context) (apikeyspb.ApiKeysClient, context.Context, error) {
	token, project, err := getADC(ctx)
	if err != nil {
		return nil, nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-goog-user-project", project)
	conn, err := newConn("apikeys.googleapis.com:443")
	if err != nil {
		return nil, nil, err
	}
	c := apikeyspb.NewApiKeysClient(conn)
	return c, ctx, nil
}

func ApiKeysLROClient(ctx context.Context) (*longrunning.OperationsClient, context.Context, error) {
	token, project, err := getADC(ctx)
	if err != nil {
		return nil, nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-goog-user-project", project)
	conn, err := newConn("apikeys.googleapis.com:443")
	if err != nil {
		return nil, nil, err
	}
	c, err := longrunning.NewOperationsClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		return nil, nil, err
	}
	return c, ctx, nil
}
