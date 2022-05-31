package formater

import (
	"strings"
)

type tableptions struct {
	fields map[string]struct{}
}

// TableOption - опция таблицы
type TableOption func(o *tableptions)

// Fields - выводить только указанные поля
func Fields(names ...string) TableOption {
	return func(o *tableptions) {
		if o.fields == nil {
			o.fields = make(map[string]struct{})
		}

		for _, n := range names {
			o.fields[strings.ToUpper(n)] = struct{}{}
		}
	}
}

func processOptions(opts ...TableOption) *tableptions {
	options := &tableptions{}

	for i := 0; i < len(opts); i++ {
		opts[i](options)
	}

	return options
}
