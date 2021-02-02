package client

import (
	"context"
	"testing"

	"github.com/digital-dream-labs/hugh/grpc/server"
	"github.com/digital-dream-labs/hugh/internal/testdata/grpcecho"
	"github.com/digital-dream-labs/hugh/internal/testdata/tls"
)

func TestNewClient(t *testing.T) {
	type args struct {
		clientopts []Option
		serveropts []server.Option
	}
	type serverargs struct {
	}

	tests := []struct {
		name        string
		args        args
		want        *Client
		wantErr     bool
		appKeyCheck bool
	}{
		{
			name: "pass with tls",
			args: args{
				clientopts: []Option{
					WithCertPool(tls.LocalhostTLSConfig().RootCAs),
				},
				serveropts: []server.Option{
					server.WithCertificate(tls.LocalhostTLSConfig().Certificates[0]),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := startTestServer(tt.args.serveropts)
			if err != nil {
				t.Errorf("server start error = %v", err)
			}

			tt.args.clientopts = append(tt.args.clientopts, WithTarget(addr))

			cli, err := New(tt.args.clientopts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := cli.Connect(); err != nil {
				t.Fatal(err)
			}

			c := grpcecho.NewEchoServiceClient(cli.Conn())

			ctx := context.Background()

			resp, err := c.Echo(
				ctx,
				&grpcecho.EchoMessage{Value: "test"},
			)
			if err != nil {
				t.Fatal(err)
			}

			if resp.Value != "test" {
				t.Fatalf("message does not match")
			}

			if err := cli.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func startTestServer(opts []server.Option) (string, error) {
	srv, err := server.New(opts...)
	if err != nil {
		return "", err
	}

	echosrv := grpcecho.Echo{}

	grpcecho.RegisterEchoServiceServer(
		srv.Transport(),
		echosrv,
	)

	srv.Start()
	<-srv.Notify(server.Ready)
	return srv.Address().String(), nil
}
