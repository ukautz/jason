package jason

import (
	"encoding/json"
	"fmt"
	"reflect"
	. "github.com/ukautz/reflekt"
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
		}
	}

	return nil
}
