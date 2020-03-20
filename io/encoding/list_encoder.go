/*--------------------------------------------------------*\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: https://hprose.com                     |
|                                                          |
| io/encoding/list_encoder.go                              |
|                                                          |
| LastModified: Mar 21, 2020                               |
| Author: Ma Bingyao <andot@hprose.com>                    |
|                                                          |
\*________________________________________________________*/

package encoding

import (
	"container/list"

	"github.com/hprose/hprose-golang/v3/io"
	"github.com/modern-go/reflect2"
)

// ListEncoder is the implementation of ValueEncoder for list.List/*list.List.
type ListEncoder struct{}

// Encode writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as reference
func (valenc ListEncoder) Encode(enc *Encoder, v interface{}) (err error) {
	return ReferenceEncode(valenc, enc, v)
}

// Write writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as value
func (ListEncoder) Write(enc *Encoder, v interface{}) (err error) {
	SetReference(enc, v)
	return writeList(enc, (*list.List)(reflect2.PtrOf(v)))
}

func writeList(enc *Encoder, lst *list.List) (err error) {
	writer := enc.Writer
	count := lst.Len()
	if count == 0 {
		_, err = writer.Write(emptySlice)
		return
	}
	err = WriteHead(writer, count, io.TagList)
	for e := lst.Front(); e != nil && err == nil; e = e.Next() {
		err = enc.Encode(e.Value)
	}
	if err == nil {
		err = WriteFoot(writer)
	}
	return
}

// ElementEncoder is the implementation of ValueEncoder for list.Element/*list.Element.
type ElementEncoder struct{}

// Encode writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as reference
func (valenc ElementEncoder) Encode(enc *Encoder, v interface{}) (err error) {
	e := (*list.Element)(reflect2.PtrOf(v))
	if e == nil {
		return WriteNil(enc.Writer)
	}
	return enc.Encode(e.Value)
}

// Write writes the hprose encoding of v to stream
// if v is already written to stream, it will writes it as value
func (ElementEncoder) Write(enc *Encoder, v interface{}) (err error) {
	return enc.Write((*list.Element)(reflect2.PtrOf(v)).Value)
}

func init() {
	RegisterEncoder((*list.List)(nil), ListEncoder{})
	RegisterEncoder((*list.Element)(nil), ElementEncoder{})
}
