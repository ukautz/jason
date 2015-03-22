[![Build Status](https://travis-ci.org/ukautz/jason.svg?branch=master)](https://travis-ci.org/ukautz/jason)

jason
=====

Provides JSON unmarshall for Go in a very relaxed manner, using [reflekt](https://github.com/ukautz/reflekt).

Documentation
-------------

GoDoc can be [found here](http://godoc.org/github.com/ukautz/jason)

Examples
--------

``` go
type Me struct {
    Foo   int    `json:"foo"`
    Bar   float  `json:"bar"`
    Baz   bool   `json:"baz"`
    Zoing string `json:"zoing"`
}

func (this *Me) UnmarshalJSON(b []byte) error {
    return jason.RelaxedUnmarshalJSONMap(this, b)
}

// -%<-

// all the following are unmarshalled to an equal object
from := []string{
    `{"foo":1, "bar": 1.1, "baz": true, "zoing": "bla"}`,
    `{"foo":"1", "bar": "1.1", "baz": "true", "zoing": "bla"}`,
    `{"foo":"1.1", "bar": "1.1", "baz": 1, "zoing": "bla"}`,
}

for _, s := range from {
    o := &Me{}
    json.Unmarshal([]byte(s), o)
}
```
