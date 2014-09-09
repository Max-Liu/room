package room

import (
	"encoding/json"

	"github.com/funny/link"
)

type Stream struct {
	link.Message
	Content []byte
}

func NewStream(data interface{}) Stream {
	var msg link.Message
	jsonByte, _ := json.Marshal(data)
	return Stream{msg, jsonByte}
}
func (msg Stream) RecommendPacketSize() uint {
	return uint(len(msg.Content))
}

func (msg Stream) AppendToPacket(packet []byte) ([]byte, error) {
	var err error
	return append(packet, msg.Content...), err
}

type Box struct {
	Object interface{}
	Kind   string
}

//func NewBox(obj interface{}, kind string) Box {
//return Box{obj, kind}
//}
