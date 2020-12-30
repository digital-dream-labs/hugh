package sql

import "github.com/digital-dream-labs/hugh/config"

// viperize augments options based on viper config
func (o *options) viperize(args ...string) error {
	v, err := config.New("DDL_DB", args...)
	if err != nil {
		return err
	}

	if x := "type"; v.IsSet(x) {
		o.databaseType = v.GetString(x)
	}

	if x := "username"; v.IsSet(x) {
		o.username = v.GetString(x)
	}

	if x := "password"; v.IsSet(x) {
		o.password = v.GetString(x)
	}

	if x := "name"; v.IsSet(x) {
		o.name = v.GetString(x)
	}

	if x := "host"; v.IsSet(x) {
		o.host = v.GetString(x)
	}

	if x := "port"; v.IsSet(x) {
		o.port = v.GetInt(x)
	}

	if x := "tls-mode"; v.IsSet(x) {
		o.tlsMode = v.GetString(x)
	}

	return nil
}
