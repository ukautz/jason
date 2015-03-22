package jason

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type testJsonObjectSuper struct {
	Super string `json:"super"`
}

type testJsonObject struct {
	testJsonObjectSuper
	Foo   int                    `json:"foo"`
	Bar   float64                `json:"bar"`
	Baz   string                 `json:"baz"`
	Boing bool                   `json:"boing"`
	Zoing time.Time              `json:"zoing"`
	Arg   map[string]interface{} `json:"arg"`
}

var testsRelaxedUnmarshalJson = []struct {
	obj *testJsonObject
	raw string
}{
	/*
		-------
		SIMPLE
		-------
	*/
	{
		raw: `{"foo":123}`,
		obj: &testJsonObject{
			Foo: 123,
		},
	},
	{
		raw: `{"bar":123.23}`,
		obj: &testJsonObject{
			Bar: 123.23,
		},
	},
	{
		raw: `{"baz":"bla"}`,
		obj: &testJsonObject{
			Baz: "bla",
		},
	},
	{
		raw: `{"boing":true}`,
		obj: &testJsonObject{
			Boing: true,
		},
	},
	{
		raw: `{"arg":{"bla": "blub"}}`,
		obj: &testJsonObject{
			Arg: map[string]interface{}{
				"bla": "blub",
			},
		},
	},
	{
		raw: `{"zoing": "2015-03-12"}`,
		obj: &testJsonObject{
			Zoing: func() time.Time {
				t, _ := time.Parse("2006-01-02", "2015-03-12");
				return t
			}(),
		},
	},
	{
		raw: `{"zoing": "2015-03-12 20:33:21"}`,
		obj: &testJsonObject{
			Zoing: func() time.Time {
				t, _ := time.Parse("2006-01-02 15:04:05", "2015-03-12 20:33:21");
				return t
			}(),
		},
	},
	{
		raw: `{"zoing": "2015-03-22T20:11:00+01:00"}`,
		obj: &testJsonObject{
			Zoing: func() time.Time {
				t, _ := time.Parse(time.RFC3339, "2015-03-22T20:11:00+01:00");
				return t
			}(),
		},
	},
	{
		raw: `{"zoing": "2015-03-12T20:33:21.379118462+01:00"}`,
		obj: &testJsonObject{
			Zoing: func() time.Time {
				t, _ := time.Parse(time.RFC3339Nano, "2015-03-12T20:33:21.379118462+01:00");
				return t
			}(),
		},
	},
	{
		raw: `{"zoing": "2015-03-12 20:33:21.379118462 +0100 CET"}`,
		obj: &testJsonObject{
			Zoing: func() time.Time {
				t, _ := time.Parse(`2006-01-02 15:04:05.999999999 -0700 MST`, "2015-03-12 20:33:21.379118462 +0100 CET");
				return t
			}(),
		},
	},
	/*
		-------
		ALL TOGETHER
		-------
	*/
	{
		raw: `{"foo":123, "bar": 123.23, "baz": "bla", "boing": true, "arg":{"bla": "blub"}}`,
		obj: &testJsonObject{
			Foo:   123,
			Bar:   123.23,
			Baz:   "bla",
			Boing: true,
			Arg: map[string]interface{}{
				"bla": "blub",
			},
		},
	},
	/*
		-------
		THE HARD ONES
		-------
	*/
	{
		raw: `{"foo":"123"}`,
		obj: &testJsonObject{
			Foo: 123,
		},
	},
	{
		raw: `{"foo":"123.23"}`,
		obj: &testJsonObject{
			Foo: 123,
		},
	},
	{
		raw: `{"bar":"123.23"}`,
		obj: &testJsonObject{
			Bar: 123.23,
		},
	},
	{
		raw: `{"bar":"123"}`,
		obj: &testJsonObject{
			Bar: 123.0,
		},
	},
	{
		raw: `{"boing":"true"}`,
		obj: &testJsonObject{
			Boing: true,
		},
	},
	{
		raw: `{"boing":"1"}`,
		obj: &testJsonObject{
			Boing: true,
		},
	},
	{
		raw: `{"boing":1}`,
		obj: &testJsonObject{
			Boing: true,
		},
	},
	{
		raw: `{"super":"trooper"}`,
		obj: &testJsonObject{
			testJsonObjectSuper: testJsonObjectSuper{
				Super: "trooper",
			},
		},
	},
}

func TestRelaxedUnmarshalJson(t *testing.T) {
	Convey("Unmarshalling JSON to object very relaxed..", t, func() {
		for idx, test := range testsRelaxedUnmarshalJson {
			Convey(fmt.Sprintf("%d) from %s", idx, test.raw), func() {
				o := new(testJsonObject)
				err := RelaxedUnmarshalJSONMap(o, []byte(test.raw))
				So(err, ShouldBeNil)
				So(o, ShouldResemble, test.obj)
			})
		}
	})
}
