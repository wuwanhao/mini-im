package websocket

// 函数选项模式的实现
type Options func(opt *option)

type option struct {
	Authentication
	pattern string
}

func newOption(opts ...Options) option {
	o := option{
		Authentication: new(authentication),
		pattern:        "/ws",
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithAuthentication(authentication Authentication) Options {
	return func(opt *option) {
		opt.Authentication = authentication
	}
}

func WithHandlerPattern(pattern string) Options {
	return func(opt *option) {
		opt.pattern = pattern
	}
}
