package object_test

import (
	"testing"

	"github.com/sam8helloworld/uwscgo/object"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &object.String{Value: "Hello World"}
	hello2 := &object.String{Value: "Hello World"}
	diff1 := &object.String{Value: "My name is johnny"}
	diff2 := &object.String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with same content have same hash keys")
	}
}

func TestHashTableGeyPairByIndex(t *testing.T) {
	ht := &object.HashTable{
		Pairs: map[object.HashKey]object.HashPair{
			{
				Type:  object.STRING_OBJ,
				Value: uint64(12638189399578898418),
			}: {
				Key: &object.String{
					Value: "c",
				},
				Value: &object.Integer{
					Value: 3,
				},
			},
			{
				Type:  object.STRING_OBJ,
				Value: uint64(12638187200555641996),
			}: {
				Key: &object.String{
					Value: "a",
				},
				Value: &object.Integer{
					Value: 1,
				},
			},
			{
				Type:  object.STRING_OBJ,
				Value: uint64(12638190499090526629),
			}: {
				Key: &object.String{
					Value: "b",
				},
				Value: &object.Integer{
					Value: 2,
				},
			},
		},
		IsSort: true,
	}

	pair := ht.GetPairByIndex(1)
	key, ok := pair.Key.(*object.String)
	if !ok {
		t.Fatalf("pair.Key is not *object.String. got=%T", pair.Key)
	}
	if key.Value != "b" {
		t.Errorf("key.Value is not 'b'. got=%s", key.Value)
	}
}
