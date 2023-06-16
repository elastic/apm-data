package modeldecoderutil

import (
	"google.golang.org/protobuf/types/known/structpb"
)

func ToStruct(m map[string]any) *structpb.Struct {
	if str, err := structpb.NewStruct(m); err == nil {
		return str
	}
	return nil
}

func ToValue(a any) *structpb.Value {
	if v, err := structpb.NewValue(a); err == nil {
		return v
	}
	return nil
}

func NormalizeMap(m map[string]any) map[string]any {
	v := NormalizeHTTPRequestBody(m)
	return v.(map[string]any)
}
