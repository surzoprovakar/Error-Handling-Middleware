package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

var mu sync.Mutex

type Set struct {
	id      int
	values  map[int]struct{}
	updates [][]string
}

func NewSet(id int) *Set {
	return &Set{id: id, values: make(map[int]struct{})}
}

func (s *Set) Id() int {
	return s.id
}

func (s *Set) Values() map[int]struct{} {
	return s.values
}

func (s *Set) Contain(value int) bool {
	_, c := s.values[value]
	return c
}

func (s *Set) Add(value int) {
	mu.Lock()
	if !s.Contain(value) {
		s.values[value] = struct{}{}
	}
	mu.Unlock()
}

func (s *Set) Remove(value int) {
	mu.Lock()
	if s.Contain(value) {
		delete(s.values, value)
	}
	mu.Unlock()
}

func (s *Set) SetVal(opt_name string, value int) {
	if opt_name == "Add" {
		s.Add(value)
	} else if opt_name == "Remove" {
		s.Remove(value)
	}
	cur_upte := []string{opt_name, strconv.Itoa(value)}
	s.updates = append(s.updates, cur_upte)
}

func (s *Set) SetRemoteVal(rid int, opt_name string, value int) {
	if opt_name == "Add" {
		s.Add(value)
	} else if opt_name == "Remove" {
		s.Remove(value)
	}
}

func (s *Set) Union(o *Set) {

	mu.Lock()
	// fmt.Println("Starting to merge:")
	// s.Print()
	// o.Print()
	m := make(map[int]bool)

	for item := range s.values {
		m[item] = true
	}

	for item := range o.values {
		if _, ok := m[item]; !ok {
			s.values[item] = struct{}{}
		}
	}
	// fmt.Print("Merged ")
	// s.Print()
	mu.Unlock()
}

// func (s *Set) Intersection(o *Set) {

// 	mu.Lock()
// 	m := make(map[int]bool)

// 	for item := range s.values {
// 		m[item] = true
// 	}

// 	for item := range o.values {
// 		if _, ok := m[item]; ok {
// 			s.values[item] = struct{}{}
// 		}
// 	}
// 	mu.Unlock()
// }

// func (s *Set) ToByteArray() []byte {

// 	a1 := make([]byte, 64)

// 	binary.LittleEndian.PutUint64(a1, uint64(s.Id()))

// 	values := get_keys_from_map(s.values)
// 	for i := range values {
// 		tmp := make([]byte, 64)
// 		binary.LittleEndian.PutUint64(tmp, uint64(values[i]))
// 		a1 = append(a1, tmp...)
// 	}
// 	return a1
// }

// func FromByteArray(bytes []byte) *Set {

// 	var length = len(bytes)
// 	var div = length / 64

// 	var r1 = binary.LittleEndian.Uint64(bytes[0:(len(bytes) / div)])

// 	i := 1
// 	var tmp_list []int
// 	for i < div {
// 		var r2 = binary.LittleEndian.Uint64(bytes[i*len(bytes)/div : (i+1)*len(bytes)/div])
// 		tmp_list = append(tmp_list, int(r2))
// 		i++
// 	}

// 	id := int64(r1)
// 	s := NewSet(int(id))

// 	for j := range tmp_list {
// 		s.values[tmp_list[j]] = struct{}{}
// 	}

// 	return s
// }

func (s *Set) Merge(rid int, r_updates [][]string) {
	// mu.Lock()
	// res := fmt.Sprintf("%s%d:%d", "Counter_", rid, rval)
	fmt.Println("Starting to merge req from replica_", rid)

	tmpSet := new(Set)
	tmpSet.id = rid
	tmpSet.values = make(map[int]struct{})
	if len(r_updates) > 0 {
		for i := 0; i < len(r_updates); i++ {
			req_val, _ := strconv.Atoi(r_updates[i][1])
			if r_updates[i][0] == "Add" {
				tmpSet.Add(req_val)
			} else {
				tmpSet.Remove(req_val)
			}
		}
	}

	// if len(r_updates) > 0 {
	// 	for i := 0; i < len(r_updates); i++ {
	// 		req_val, _ := strconv.Atoi(r_updates[i][1])
	// 		s.SetRemoteVal(rid, r_updates[i][0], req_val)
	// 	}
	// }

	s.Union(tmpSet)
	tmpSet = nil
	fmt.Print("Merged :-->")
	s.Print()
	// mu.Unlock()
}

func (s *Set) Print() {
	fmt.Print("Set:", s.id, " ")
	values := get_keys_from_map(s.values)
	fmt.Println(values)
}

func get_keys_from_map(mapp map[int]struct{}) []int {
	keys := make([]int, 0, len(mapp))
	for key := range mapp {
		keys = append(keys, key)
	}
	sort.Ints(keys[:])
	return keys
}

func (s *Set) ToMarshal() []byte {
	id_2d := []string{strconv.Itoa(s.Id()), ""}
	s.updates = append(s.updates, id_2d)

	jsonData, err := json.Marshal(s.updates)
	if err != nil {
		fmt.Println("Error while marshaling updates")
		return nil
	}
	s.updates = [][]string{}
	return jsonData
}

func FromMarshalData(bytes []byte) (int, [][]string) {

	var remote_updates [][]string
	err := json.Unmarshal(bytes, &remote_updates)

	if err != nil {
		fmt.Println("Error while unmarshaling ", err)
		return -1, nil
	}

	rid, _ := strconv.Atoi(remote_updates[len(remote_updates)-1][0])
	// fmt.Println("rid is: ", rid)
	if len(remote_updates) == 1 {
		return rid, nil
	}
	r_updates := remote_updates[0 : len(remote_updates)-1]
	// fmt.Println(r_updates)
	return rid, r_updates
}
