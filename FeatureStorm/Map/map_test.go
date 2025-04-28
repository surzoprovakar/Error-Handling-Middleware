package main

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {

	map1, map2 := NewMap(1), NewMap(2)

	map1.Add("a", 1)
	map1.Add("b", 3)
	map1.Add("c", 5)

	map2.Add("p", 2)
	map2.Add("q", 4)
	map2.Add("r", 6)

	map1.Delete("d")
	map2.Delete("r")

	map1.Update("c", 7)
	map2.Update("r", 8)

	map1.Print()
	map2.Print()

	b1 := map1.ToMarshal()
	// fmt.Println("b1: ", b1)
	b2 := map2.ToMarshal()

	rid1, updates1 := FromMarshalData(b1)
	rid2, updates2 := FromMarshalData(b2)

	map3, map4 := NewMap(3), NewMap(4)

	map3.Merge(rid1, updates1)
	map4.Merge(rid2, updates2)

	map3.Print()
	map4.Print()

	if !reflect.DeepEqual(map1.Values(), map3.Values()) {
		t.Errorf("map1 and map3 values are not the same")
	}

	if !reflect.DeepEqual(map2.Values(), map4.Values()) {
		t.Errorf("map2 and map4 values are not the same")
	}

}
