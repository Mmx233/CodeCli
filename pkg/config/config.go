package config

import (
	"fmt"
	"github.com/Mmx233/CodeCli/internal/global"
	"github.com/Mmx233/CodeCli/internal/util"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
)

func List() error {
	d, e := yaml.Marshal(&global.Config)
	if e != nil {
		return e
	}
	fmt.Printf(string(d))
	return nil
}

func loadField(fieldRaw string) (reflect.Value, error) {
	rv := reflect.ValueOf(&global.Config).Elem()
	for _, field := range strings.Split(fieldRaw, ".") {
		rt := rv.Type()
		var hit bool
		for i := 0; i < rv.NumField(); i++ {
			rtf := rt.Field(i)
			if field == rtf.Name || field == rtf.Tag.Get("yaml") {
				hit = true
				rv = rv.Field(i)
				break
			}
		}
		if !hit {
			return rv, util.ErrIllegalInput
		}
	}
	return rv, nil
}

func Set(field, value string) error {
	raw := strings.Split(field, "=")
	if len(raw) != 2 {
		if len(raw) == 1 && value != "" {
			raw = []string{
				field, value,
			}
		} else {
			return util.ErrIllegalInput
		}
	}

	fieldName, fieldValue := raw[0], raw[1]
	rv, e := loadField(fieldName)
	if e != nil {
		return e
	}
	rv.SetString(fieldValue)
	return global.ConfigLoader.Save()
}

func Unset(field string) error {
	rv, e := loadField(field)
	if e != nil {
		return e
	}
	rv.SetString("")
	return global.ConfigLoader.Save()
}
