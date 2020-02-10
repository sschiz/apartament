package apartment

type Options struct {
	limit  int
	offset int
}

type Option interface {
	Apply(*Options)
}

type optionFunc func(*Options)

func (f optionFunc) Apply(o *Options) {
	f(o)
}

func WithLimit(limit int) Option {
	return optionFunc(func(o *Options) {
		o.limit = limit
	})
}

func WithOffset(offset int) Option {
	return optionFunc(func(o *Options) {
		o.offset = offset
	})
}

const WithoutLimit int = -1 // default limit
