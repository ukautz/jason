package jason

import (
	"encoding/json"
	"fmt"
	"reflect"
	. "github.com/ukautz/reflekt"
	"time"
)

var (
	TimeFormats = []string{
		`2006-01-02 15:04:05.999999999 -0700 MST`,
		time.RFC3339Nano,
		time.RFC3339,
		`2006-01-02`,
		`2006-01-02 15:04:05`,
	}
)

// Unmarshal JSON in a relaxed manner: integers can be in strings or encoded as float. Same for
// bool, float and string.
func RelaxedUnmarshalJSONMap(obj interface{}, b []byte) error {
	return RelaxedUnmarshalJSONMapCustom(obj, b, nil)
}

func RelaxedUnmarshalJSONMapCustom(obj interface{}, b []byte, cb func(val reflect.Value, typ reflect.StructField, input interface{}) bool) error {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("Could not unmarshall JSON relaxed: %v", r))
		}
	}()
	var x map[string]interface{}
	if e := json.Unmarshal(b, &x); e != nil {
		return e
	}
	var r reflect.Value
	if reflect.TypeOf(obj) == reflect.TypeOf(reflect.Value{}) {
		r = obj.(reflect.Value)
	} else {
		r = reflect.ValueOf(obj)
	}
	ok := false
	if r.Kind() == reflect.Ptr {
		ok = true
		r = r.Elem()
	} else if r.Kind() == reflect.Struct {
		ok = true
	}
	if !ok {
		return fmt.Errorf("Object %s must be pointer or struct, but is %s", reflect.TypeOf(obj).String(), r.Kind().String())
	} else if r.Kind() == reflect.Invalid {
		return fmt.Errorf("Object is nil!")
	}

	t := r.Type()
	for i := 0; i < r.NumField(); i++ {
		fv := r.Field(i)
		ft := t.Field(i)
		if ft.Anonymous {
			if fv.Kind() == reflect.Struct {
				if err := RelaxedUnmarshalJSONMapCustom(fv, b, cb); err != nil {
					return err
				}
			} else {
				continue
			}
		}
		tag := ft.Tag
		n := tag.Get("json")
		if n == "" {
			continue
		}
		v, ok := x[n]

		// apply custom callback to build data
		if cb != nil && cb(fv, ft, v) {
			continue
		} else if !ok {
			continue
		}
		k := fv.Kind()
		vv := reflect.ValueOf(v)
		switch {
		case IsIntKind(k):
			fv.SetInt(int64(AsInt(v)))
		case IsFloatKind(k):
			fv.SetFloat(AsFloat(v))
		case k == reflect.Bool:
			fv.SetBool(AsBool(v))
		case k == reflect.String:
			fv.SetString(AsString(v))
		case k == reflect.Map && vv.Kind() == reflect.Map:
			fv.Set(vv)
		case isTimeV(fv):
			vt := reflect.TypeOf(v)
			vk := vt.Kind()
			switch {
				case vt == fv.Type():
					fv.Set(vv)
				case vk == reflect.String:
					if t, err := extractTime(v.(string)); err != nil {
						return err
					} else {
						fv.Set(reflect.ValueOf(t))
					}
			}
		default:
			fmt.Printf("DUNNO %s - %s\n", fv.Type(), reflect.TypeOf(v))
		}
	}

	return nil
}

func extractTime(v string) (time.Time, error) {
	var t time.Time
	var err error
	for _, f := range TimeFormats {
		if t, err = time.Parse(f, v); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("None of the time formats matched. \"%s\" is not time!", v)
}

func isTime(v interface{}) bool {
	return reflect.TypeOf(time.Time{}) == reflect.TypeOf(v)
}

func isTimeV(v reflect.Value) bool {
	return reflect.TypeOf(time.Time{}) == v.Type()
}
