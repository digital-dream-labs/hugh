package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/digital-dream-labs/hugh/config"
)

// viperize augments options based on viper config
//
// args are a k, v scheme. Special keys are:

const (
	// EnvironmentPrefix sets the prefix to strip from environment variables when resolving keys. Default: "DDL_RPC".
	EnvironmentPrefix = "env-prefix"
	// EnvironmentReplace takes a comma separated list of old,new. Default: .,_,-,_
	EnvironmentReplace = "env-replace"
)

// The remaining k,v's are treated as viper BindEnv args where k is the viper key and v is the env variable. Note that prefix is not used when specifically binding vars.

func (o *options) viperize(args ...string) error {
	v, err := config.New("DDL_RPC", args...)
	if err != nil {
		return err
	}

	if x := "target"; v.IsSet(x) {
		o.target = v.GetString(x)
	}

	if x := "insecure"; v.IsSet(x) {
		o.insecure = v.GetBool(x)
	}

	if v.IsSet("disable-tls") {
		o.disableTLS = true
	}

	if v.IsSet("tls-certificate") && v.IsSet("tls-key") {
		crt, err := tls.X509KeyPair(
			[]byte(v.GetString("tls-certificate")),
			[]byte(v.GetString("tls-key")),
		)
		if err != nil {
			return err
		}
		o.certificates = append(o.certificates, crt)
		o.log.Debug("RPC::tls-certificate: ", v.GetString("tls-certificate"))
		o.log.Debug("RPC::tls-key: ", v.GetString("tls-key"))
	}

	if v.IsSet("tls-ca") {
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(v.GetString("tls-ca"))); !ok {
			return fmt.Errorf("CA is not a valid pem file")
		}

		o.certPool = pool
	}

	return nil
}
