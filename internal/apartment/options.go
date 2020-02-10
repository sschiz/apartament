package apartment

type Options struct {
	Limit  int
	Offset int
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
		o.Limit = limit
	})
}

func WithOffset(offset int) Option {
	return optionFunc(func(o *Options) {
		o.Offset = offset
	})
}

// Options
const (
	DefaultLimit  int = -1
	DefaultOffset int = 0
)
