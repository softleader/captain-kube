package ctx

import (
	"github.com/spf13/pflag"
	"reflect"
)

// Flaggable 代表可以被加入 flags
type Flaggable interface {
	AddFlags(f *pflag.FlagSet)
}

// 支援無限層的 Flaggable
func addFlags(x interface{}, f *pflag.FlagSet) {
	v := valueOf(x)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			addFlags(x, f)
		}
		if flaggable, ok := field.Interface().(Flaggable); ok {
			flaggable.AddFlags(f)
		}
	}
}

func valueOf(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		return val.Elem()
	}
	return val
}
