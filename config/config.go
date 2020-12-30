package config

import (
	"strings"

	"github.com/spf13/viper"
)

// New returns a configured viper instance
func New(prefix string, args ...string) (*viper.Viper, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()

	for i, j := 0, 1; j < len(args); i, j = i+2, j+2 {
		key, val := args[i], args[j]
		switch key {
		case "env-prefix":
			v.SetEnvPrefix(val)
		case "env-replace":
			x := strings.Split(val, ",")
			v.SetEnvKeyReplacer(strings.NewReplacer(x...))
		default:
			if err := v.BindEnv(key, val); err != nil {
				return nil, err
			}
		}
	}

	return v, nil
}
