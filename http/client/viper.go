package client

import "github.com/digital-dream-labs/hugh/config"

// viperize augments options based on viper config
//
// args are a k, v scheme. Special keys are:
// * env-prefix - sets the prefix to strip from environment variables when resolving keys. Default: "DDL_RPC".
// * env-replace - takes a comma separated list of old,new. Default: .,_,-,_
//
// The remaining k,v's are treated as viper BindEnv args where k is the viper key and v is the env variable. Note that prefix is not used when specifically binding vars.
func (o *options) viperize(args ...string) error {
	v, err := config.New("DDL_HTTP", args...)
	if err != nil {
		return err
	}

	if x := "target"; v.IsSet(x) {
		o.target = v.GetString(x)
	}

	return nil
}
