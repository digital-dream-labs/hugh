# Client

To quickly get a grpc connection (in this example, with non-system TLS certs) you'd simply do this:

Note: this example uses the [grpc example protobufs](https://google.golang.org/grpc/examples/features/proto/echo)

```go
func getconn() echo.EchoClient {
    cli, err := New(
        WithCertPool(tls.LocalhostTLSConfig().RootCAs),
        WithTarget("hostame:1234")
    )
    if err != nil {
        log.Fatal(err)
    }

    err := cli.Connect(); err != nil {
        log.Fatal(err)
    }

    return echo.NewEchoClient(cli.Conn())
}
```

This expects the following environment variables to be in place when calling WithViper() without arguments

| Environment Variable | Description | 
| ------------ | ------------ |
| DDL_RPC_TARGET  | Sets the target host for an RPC client   |  
| DDL_RPC_TLS_CERTIFICATE  | Sets the client authentication cert for mutual tls   | 
| DDL_RPC_KEY  | Private key that pairs with tls certificate   | 
| DDL_RPC_TLS_CA  | Sets the certificate authority for client verification | 
| DDL_RPC_APP_KEY | populates the RPC context with an app key |

A full list of options can be found in [options.go](options.go)
