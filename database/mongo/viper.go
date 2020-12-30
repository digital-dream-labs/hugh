package mongo

import "github.com/digital-dream-labs/hugh/config"

// viperize augments options based on viper config
func (o *opts) viperize(args ...string) error {
	v, err := config.New("DDL_DB", args...)
	if err != nil {
		return err
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

	if x := "auth-name"; v.IsSet(x) {
		o.authdbname = v.GetString(x)
	}

	if x := "host"; v.IsSet(x) {
		o.host = v.GetString(x)
	}

	if x := "port"; v.IsSet(x) {
		o.port = v.GetInt(x)
	}

	if x := "direct"; v.IsSet(x) {
		o.direct = v.GetBool(x)
	}

	if x := "cluster"; v.IsSet(x) {
		o.cluster = v.GetBool(x)
	}

	if x := "retry-writes"; v.IsSet(x) {
		o.retryWrites = v.GetBool(x)
	}

	if x := "write-concern"; v.IsSet(x) {
		o.writeConcern = v.GetString(x)
	}

	return nil
}
