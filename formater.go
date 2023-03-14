package formater

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/gomisc/errors.v1"
	"gopkg.in/yaml.v3"
)

type (
	// Row - алиас типа строки таблицы
	Row = table.Row

	// Formater - многоформатный форматер вывода
	Formater interface {
		String() error
		Json(writer io.Writer) error
		Yaml(writer io.Writer) error
		Table(opts ...TableOption) Table
		isTable() bool
	}

	// Table - табличный форматер вывода
	Table interface {
		CSV() string
		HTML() string
		Markdown() string
		Render() string
	}
)

// OutputFormat - тип форматировани
type OutputFormat string

// Перечисление типов форматирования
const (
	Json      OutputFormat = "json"
	Yaml      OutputFormat = "yaml"
	HTML      OutputFormat = "html"
	Markdownn OutputFormat = "md"
	CSV       OutputFormat = "csv"
)

type formater struct {
	data interface{}
}

// Format - конструктор форматера
func Format(data interface{}) Formater {
	return &formater{data: data}
}

// Print - вывести объект в указанном формате на стандартный поток вывода
func Print(data interface{}, f OutputFormat, opts ...TableOption) {
	printer := Format(data)

	switch f {
	case Json:
		_ = printer.Json(os.Stdout)
	case Yaml:
		_ = printer.Yaml(os.Stdout)
	case CSV:
		fmt.Println(printer.Table(opts...).CSV())
	case Markdownn:
		fmt.Println(printer.Table(opts...).Markdown())
	case HTML:
		fmt.Println(printer.Table(opts...).HTML())
	default:
		if !printer.isTable() {
			_ = printer.String()

			return
		}

		fmt.Println(printer.Table(opts...).Render())
	}
}

func (f *formater) String() error {
	if _, err := fmt.Fprintf(os.Stdout, "%#v\n", f.data); err != nil {
		return errors.Wrap(err, "print data as string")
	}

	return nil
}

func (f *formater) Json(w io.Writer) error {
	data, err := json.MarshalIndent(f.data, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal data json")
	}

	if _, err = w.Write(append(data, []byte("\n")...)); err != nil {
		return errors.Wrap(err, "write output json")
	}

	return nil
}

func (f *formater) Yaml(w io.Writer) error {
	data, err := yaml.Marshal(f.data)
	if err != nil {
		return errors.Wrap(err, "marshal data yaml")
	}

	if _, err = w.Write(append(data, []byte("\n")...)); err != nil {
		return errors.Wrap(err, "write output yaml")
	}

	return nil
}

func (f *formater) Table(opts ...TableOption) Table {
	return newTab(f.data, opts...)
}

func (f *formater) isTable() bool {
	var (
		v = reflect.ValueOf(f.data)
		t reflect.Type
	)

	if t = v.Type(); t.Kind() != reflect.Slice {
		t1 := t.Elem()
		if t1.Kind() != reflect.Slice {
			return false
		}
	}

	return true
}
