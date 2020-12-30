# Server

To quickly get a grpc server (in this example, with non-system TLS certs) you'd simply do this:

Note: this example uses the [grpc example protobufs](https://google.golang.org/grpc/examples/features/proto/echo)

```go
func startTestServer() error {
    srv, err := server.New(
        server.WithCertificate(tls.LocalhostTLSConfig().Certificates[0]),
    )
    if err != nil {
        return err
    }

    echosrv := grpc.New()

    echo.RegisterEchoServer(
        srv.Transport(),
        echosrv,
    )

    srv.Start()
    <-srv.Notify(server.Ready)
    return nil
}
```

This expects the following environment variables to be in place when calling WithViper() without arguments

| Environment Variable | Description | Default |
| ------------ | ------------ | ------------ |
| DDL_RPC_CLIENT_AUTHENTICATION  | Values: [NoClientCert, RequestClientCert, RequireAnyClientCert, VerifyClientCertIfGiven, RequireAndVerifyClientCert]   | NoClientCert  |
| DDL_RPC_TLS_CERTIFICATE  | Sets the certificate to use for transport encryption   |   |
| DDL_RPC_KEY  | Private key that pairs with tls certificate   |   |
| DDL_RPC_TLS_CA  | Load a custom CA pool instead of using the system CA  | empty  |
| DDL_RPC_PORT  | Sets the listener port.   | 0 |
| DDL_RPC_TLS_CA  | Sets the certificate authority for client verification | empty |
| DDL_RPC_INSECURE  | disable TLS verification | false  |

A full list of options can be found in [options.go](options.go)
