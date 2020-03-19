/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| io/encoding/time_encoder.go                              |
|                                                          |
| LastModified: Mar 19, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"reflect"
	"time"

	"github.com/modern-go/reflect2"
)

// TimeEncoder is the implementation of ValueEncoder for time.Time/*time.Time.
type TimeEncoder struct{}

// Encode writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as reference
func (valenc TimeEncoder) Encode(enc *Encoder, v interface{}) (err error) {
	if reflect.TypeOf(v).Kind() == reflect.Struct {
		return valenc.Write(enc, v)
	}
	if reflect2.IsNil(v) {
		return WriteNil(enc.Writer)
	}
	var ok bool
	if ok, err = enc.WriteReference(v); !ok && err == nil {
		err = valenc.Write(enc, v)
	}
	return
}

// Write writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as value
func (TimeEncoder) Write(enc *Encoder, v interface{}) (err error) {
	t := reflect.TypeOf(v)
	var dt time.Time
	if t.Kind() == reflect.Ptr {
		enc.SetReference(v)
		dt = *(v.(*time.Time))
	} else {
		enc.AddReferenceCount(1)
		dt = v.(time.Time)
	}
	return WriteTime(enc.Writer, dt)
}

func init() {
	RegisterEncoder((*time.Time)(nil), TimeEncoder{})
}
