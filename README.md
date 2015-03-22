[![Build Status](https://travis-ci.org/ukautz/jason.svg?branch=master)](https://travis-ci.org/ukautz/jason)

jason
=====

Provides JSON unmarshall for Go in a very relaxed manner, using [reflekt](https://github.com/ukautz/reflekt).

My use-case is parsing JSON from multiple inputs, which have slight format differences (eg one
source encodes integer as quoted string, the other source as float, the next as, well, integer..
don't get me started on bool!).

Documentation
-------------

GoDoc can be [found here](http://godoc.org/github.com/ukautz/jason)

Examples
--------

``` go
type Me struct {
    Foo   int     `json:"xfoo"`
    Bar   float64
    Baz   bool
    Zoing string  `json:"xzoing"`
}

func (this *Me) UnmarshalJSON(b []byte) error {
    return jason.RelaxedUnmarshalJSONMap(this, b)
}

// -%<-

// all the following are unmarshalled to an equal object
from := []string{
    `{"xfoo":1, "bar": 1.1, "baz": true, "xzoing": "bla"}`,
    `{"xfoo":"1", "bar": "1.1", "baz": "true", "xzoing": "bla"}`,
    `{"xfoo":"1.1", "bar": "1.1", "baz": 1, "xzoing": "bla"}`,
}

for _, s := range from {
    o := &Me{}
    json.Unmarshal([]byte(s), o)
}
```

### JSON vs Time

If your JSON contains time strings and you want to parse them into `time.Time` objects, check
out `jason.TimeFormats`, which you can change to a shorter list (less recognized, better
performance) or extend (more recognized, worse performance).

