package formater

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type tab struct {
	w table.Writer
}

func newTab(data interface{}, opts ...TableOption) Table {
	options := processOptions(opts...)

	tb := &tab{w: table.NewWriter()}

	headers, rows := headerAndRows(data, options.fields)
	if headers != nil {
		tb.w.AppendHeader(headers)
	}

	if rows != nil {
		tb.w.AppendRows(rows)
	}

	return tb
}

func (t *tab) CSV() string {
	if t != nil {
		return t.w.RenderCSV()
	}

	return ""
}

func (t *tab) HTML() string {
	if t != nil {
		return t.w.RenderHTML()
	}

	return ""
}

func (t *tab) Markdown() string {
	if t != nil {
		return t.w.RenderMarkdown()
	}

	return ""
}

func (t *tab) Render() string {
	if t != nil {
		return t.w.Render()
	}

	return ""
}

func headerAndRows(data interface{}, fields map[string]struct{}) (h Row, r []Row) {
	v := reflect.ValueOf(data)

	var (
		t, elemT reflect.Type
	)

	if t = v.Type(); t.Kind() != reflect.Slice {
		t1 := t.Elem()
		if t1.Kind() != reflect.Slice {
			return nil, nil
		}

		t = t1
	}

	if elemT = t.Elem(); elemT.Kind() != reflect.Struct {
		elemT = t.Elem().Elem()
	}

	h = make(Row, 0, elemT.NumField())

	excludes := make(map[int]struct{})
	allFields := len(fields) == 0

	for fi := 0; fi < elemT.NumField(); fi++ {
		name := strings.ToUpper(elemT.Field(fi).Name)

		if _, ok := fields[name]; allFields || ok {
			h = append(h, name)
		} else {
			excludes[fi] = struct{}{}
		}
	}

	vs := v
	if v.Kind() == reflect.Ptr && v.CanInterface() {
		vs = vs.Elem()
	}

	r = make([]Row, vs.Len())

	for ri := 0; ri < vs.Len(); ri++ {
		rr := vs.Index(ri)
		if rr.Kind() == reflect.Ptr {
			rr = rr.Elem()
		}

		r[ri] = make(Row, 0, len(h))

		for fi := 0; fi < elemT.NumField(); fi++ {
			if _, exclude := excludes[fi]; !exclude {
				rrVal := rr.Field(fi)
				if rrVal.Kind() == reflect.Slice {
					r[ri] = append(r[ri], renderSliceVal(rrVal))
				} else {
					if rrVal.Kind() == reflect.Ptr && !rrVal.IsZero() {
						r[ri] = append(r[ri], rrVal.Elem().Interface())
					} else {
						r[ri] = append(r[ri], rrVal.Interface())
					}
				}
			}
		}
	}

	return h, r
}

func renderSliceVal(val reflect.Value) Row {
	slice := make([]interface{}, val.Len())

	for si := 0; si < val.Len(); si++ {
		elem := val.Index(si)

		var data []byte

		if elem.Kind() == reflect.Ptr && elem.CanInterface() {
			data, _ = json.MarshalIndent(val.Index(si).Elem().Interface(), "", "  ")
		} else {
			data, _ = json.MarshalIndent(val.Index(si).Interface(), "", "  ")
		}

		slice[si] = string(data)
	}

	return slice
}
