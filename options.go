package bexpr

// getOpts - iterate the inbound Options and return a struct
func getOpts(opt ...Option) options {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o != nil {
			o(&opts)
		}
	}
	return opts
}

// Option - how Options are passed as arguments
type Option func(*options)

// options = how options are represented
type options struct {
	withMaxExpressions uint64
}

func WithMaxExpressions(maxExprCnt uint64) Option {
	return func(o *options) {
		o.withMaxExpressions = maxExprCnt
	}
}

func getDefaultOptions() options {
	return options{
		withMaxExpressions: 0,
	}
}
