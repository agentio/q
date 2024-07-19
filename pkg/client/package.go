package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"cloud.google.com/go/apikeys/apiv2/apikeyspb"
	"github.com/agentio/q/pkg/gcloud"
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

func ApiKeysClient(ctx context.Context, project string) (apikeyspb.ApiKeysClient, context.Context, error) {
	token, err := gcloud.GetADCToken(false)
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
