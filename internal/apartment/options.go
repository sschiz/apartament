package apartment

type options struct {
	limit  int
	offset int
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithLimit(limit int) Option {
	return optionFunc(func(o *options) {
		o.limit = limit
	})
}

func WithOffset(offset int) Option {
	return optionFunc(func(o *options) {
		o.offset = offset
	})
}
