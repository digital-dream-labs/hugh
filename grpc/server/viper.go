package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"

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

	if v.IsSet("insecure") {
		o.insecure = true
	}

	if v.IsSet("client-authentication") {
		a := v.GetString("client-authentication")
		switch a {
		case "NoClientCert":
			o.clientAuth = tls.NoClientCert
		case "RequestClientCert":
			o.clientAuth = tls.RequestClientCert
		case "RequireAnyClientCert":
			o.clientAuth = tls.RequireAnyClientCert
		case "VerifyClientCertIfGiven":
			o.clientAuth = tls.VerifyClientCertIfGiven
		case "RequireAndVerifyClientCert":
			o.clientAuth = tls.RequireAndVerifyClientCert
		default:
			return fmt.Errorf("invalid client-authentication value %q, valid values [NoClientCert, RequestClientCert, RequireAnyClientCert, VerifyClientCertIfGiven, RequireAndVerifyClientCert]", a)
		}
	}

	if v.IsSet("tls-certificate") && v.IsSet("tls-key") {
		crt, err := tls.X509KeyPair(
			[]byte(v.GetString("tls-certificate")),
			[]byte(v.GetString("tls-key")),
		)
		if err != nil {
			log.Fatal(err)
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
		o.log.Debug("RPC::tls-ca: ", v.GetString("tls-ca"))

		o.certPool = pool
	}

	if v.IsSet("port") {
		o.port = v.GetInt("port")
		o.log.Debugf("RPC::port: %d", v.GetInt("port"))
	}

	return nil
}
