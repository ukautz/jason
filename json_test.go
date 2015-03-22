package jason

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
