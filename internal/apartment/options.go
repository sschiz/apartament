package apartment

import (
	"github.com/sschiz/apartament/models"
	"reflect"
	"strings"
)

type Options struct {
	Limit      int
	Offset     int
	OrderField string
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

func WithOrder(field string) Option {
	return optionFunc(func(o *Options) {
		house := reflect.ValueOf(&models.House{}).Elem()
		apt := reflect.ValueOf(&models.Apartment{}).Elem()

		for i := 0; i < house.NumField(); i++ {
			if field == strings.ToLower(house.Type().Field(i).Name) {
				o.OrderField = field
			}
		}

		for i := 0; i < apt.NumField(); i++ {
			if field == strings.ToLower(apt.Type().Field(i).Name) {
				o.OrderField = field
			}
		}

		if field == "min_apartment_number" || field == "max_apartment_number" {
			o.OrderField = field
		}
	})
}

// Options
const (
	DefaultLimit      int    = -1
	DefaultOffset     int    = 0
	DefaultOrderField string = ""
)
