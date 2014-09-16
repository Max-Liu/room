package room

import (
	"bytes"

	"github.com/ugorji/go/codec"
)

type Box struct {
	Object interface{} `codec:"Object"`
	Kind   string      `codec:"Kind"`
}

func Encode(ob interface{}) ([]byte, error) {
	w := new(bytes.Buffer)
	enc := codec.NewEncoder(w, &codec.MsgpackHandle{RawToString: true, WriteExt: true})
	err := enc.Encode(ob)
	return w.Bytes(), err
}

func Decode(b []byte, ob interface{}) error {
	dec := codec.NewDecoderBytes(b, &codec.MsgpackHandle{RawToString: true, WriteExt: true})
	err := dec.Decode(ob)
	return err
}
