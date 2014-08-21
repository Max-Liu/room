package room

import "encoding/json"

type Stream struct {
	Content []byte
}

func NewStream(data interface{}) Stream {
	jsonByte, _ := json.Marshal(data)
	return Stream{jsonByte}
}

func (msg Stream) RecommendPacketSize() uint {
	return uint(len(msg.Content))
}

func (msg Stream) AppendToPacket(packet []byte) []byte {
	return append(packet, msg.Content...)
}

type Box struct {
	Object interface{}
	Kind   string
}

//func NewBox(obj interface{}, kind string) Box {
//return Box{obj, kind}
//}
