package modelpb

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func randomStruct(t *testing.T) (*structpb.Struct, map[string]any) {
	m := map[string]any{
		t.Name() + ".key." + randString(): t.Name() + ".value." + randString(),
	}

	s, err := structpb.NewStruct(m)
	require.NoError(t, err)

	return s, m
}

func uintPtr(i uint32) *uint32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString() string {
	size := 5
	var sb strings.Builder
	sb.Grow(size)
	for i := 0; i < size; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}
