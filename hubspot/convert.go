package hubspot

import (
	"reflect"
	"time"
	"unicode"

	"github.com/spf13/cast"
)

var timetype reflect.Type = reflect.TypeOf(time.Time{})

func convert(value interface{}, t reflect.Type) interface{} {
	if reflect.TypeOf(value) == t {
		return value
	}

	if t == timetype {
		switch v := value.(type) {
		case string:
			if len(v) >= 0 {
				alldigits := true
				for _, char := range v {
					if !unicode.IsDigit(char) {
						alldigits = false
						break
					}
				}

				if alldigits {
					// hubspot sends unix time in milliseconds
					return cast.ToTime(cast.ToInt64(v) / 1000).UTC()
				}
			}
		}
		return cast.ToTime(value)
	}

	switch t.Kind() {
	case reflect.Bool:
		return cast.ToBool(value)
	case reflect.Int:
		return cast.ToInt(value)
	case reflect.Int8:
		return cast.ToInt8(value)
	case reflect.Int16:
		return cast.ToInt16(value)
	case reflect.Int32:
		return cast.ToInt32(value)
	case reflect.Int64:
		return cast.ToInt64(value)
	case reflect.Uint:
		return cast.ToUint(value)
	case reflect.Uint8:
		return cast.ToUint8(value)
	case reflect.Uint16:
		return cast.ToUint16(value)
	case reflect.Uint32:
		return cast.ToUint32(value)
	case reflect.Uint64:
		return cast.ToUint64(value)
	case reflect.Float32:
		return cast.ToFloat32(value)
	case reflect.Float64:
		return cast.ToFloat64(value)
	case reflect.String:
		return cast.ToString(value)
	case reflect.Slice:
		elementtype := t.Elem()
		array := reflect.MakeSlice(t, 0, 8)
		sourcevalue := reflect.ValueOf(value)
		for i := 0; i < sourcevalue.Len(); i++ {
			array = reflect.Append(array, reflect.ValueOf(convert(sourcevalue.Index(i).Interface(), elementtype)))
		}
		return array.Interface()
	}

	return nil
}
