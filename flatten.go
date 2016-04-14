package i18n

import (
	"fmt"
	"reflect"
	"strconv"
)

func flatten(in map[string]interface{}) (out map[string]string, err error) {
	out = make(map[string]string)

	for k, v := range in {
		if err = flattenValue(out, k, reflect.ValueOf(v)); err != nil {
			return
		}
	}

	return
}

func flattenValue(out map[string]string, p string, v reflect.Value) (err error) {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			out[p] = "true"
		} else {
			out[p] = "false"
		}

	case reflect.Int:
		out[p] = strconv.Itoa(int(v.Int()))

	case reflect.Map:
		err = flattenMap(out, p, v)

	case reflect.Slice:
		err = flattenSlice(out, p, v)

	case reflect.String:
		out[p] = v.String()

	default:
		err = fmt.Errorf("unsupported type %v", v)
	}

	return
}

func flattenMap(out map[string]string, p string, v reflect.Value) (err error) {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}

		if k.Kind() != reflect.String {
			return fmt.Errorf("%s: map key is not string: %s", p, k)
		}

		if err = flattenValue(out, p+"."+k.String(), v.MapIndex(k)); err != nil {
			break
		}
	}
	return
}

func flattenSlice(out map[string]string, p string, v reflect.Value) (err error) {
	p += "."
	for i := 0; i < v.Len(); i++ {
		if err = flattenValue(out, p+strconv.Itoa(i), v.Index(i)); err != nil {
			break
		}
	}
	return
}
