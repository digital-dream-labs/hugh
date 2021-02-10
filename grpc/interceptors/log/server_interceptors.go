package log

// SafeLoggingModifier provides functions needed for this
// middleware to safely modify copies of data on the chain
type SafeLoggingModifier interface {
	Handles(interface{}) bool
	Copy(interface{}) interface{}
	Modify(interface{}) interface{}
}

// InterceptOption is used to configure interceptors
type InterceptOption func(*interceptConfig)

// InterceptWithReqModifier returns an interceptOption to pass to Intercept* functions
func InterceptWithReqModifier(method string, modifier SafeLoggingModifier) InterceptOption {
	return func(cfg *interceptConfig) {
		cfg.reqModifiers[method] = modifier
	}
}

// InterceptWithRespModifier returns an interceptOption to pass to Intercept* functions
func InterceptWithRespModifier(method string, modifier SafeLoggingModifier) InterceptOption {
	return func(cfg *interceptConfig) {
		cfg.respModifiers[method] = modifier
	}
}

type interceptConfig struct {
	reqModifiers  map[string]SafeLoggingModifier
	respModifiers map[string]SafeLoggingModifier
}

func newInterceptConfig() *interceptConfig {
	return &interceptConfig{
		reqModifiers:  make(map[string]SafeLoggingModifier),
		respModifiers: make(map[string]SafeLoggingModifier),
	}
}

func prepareForLog(modifier SafeLoggingModifier, src interface{}) interface{} {
	cp := src
	if modifier != nil && modifier.Handles(src) {
		cp = modifier.Copy(src)
		if src == cp {
			return src
		}
		cp = modifier.Modify(cp)
	}
	return cp
}
