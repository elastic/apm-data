package vtproto

import "fmt"

type vtprotoMessage interface {
	MarshalVT() ([]byte, error)
	UnmarshalVT([]byte) error
}

// Codec is a composite of Encoder and Decoder.
type Codec struct {
}

// Encode encodes vtprotoMessage type into byte slice
func (Codec) Encode(in any) ([]byte, error) {
	vt, ok := in.(vtprotoMessage)
	if !ok {
		return nil, fmt.Errorf("failed to encode, message is %T (missing vtprotobuf helpers)", in)
	}
	return vt.MarshalVT()
}

// Decode decodes a byte slice into vtprotoMessage type.
func (Codec) Decode(in []byte, out any) error {
	vt, ok := out.(vtprotoMessage)
	if !ok {
		return fmt.Errorf("failed to decode, message is %T (missing vtprotobuf helpers)", out)
	}
	return vt.UnmarshalVT(in)
}
