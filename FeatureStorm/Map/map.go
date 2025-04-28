package main

import (
	"fmt"
	_map "map/SyncMsg"
	"math"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/protobuf/proto"
)

var mu sync.Mutex

type Map struct {
	id      int
	values  map[string]int
	updates []string
}

func NewMap(id int) *Map {
	return &Map{id: id, values: make(map[string]int)}
}

func (m *Map) Id() int {
	return m.id
}

func (m *Map) Values() map[string]int {
	return m.values
}

func (m *Map) Add(key string, val int) {
	mu.Lock()
	if _, ok := m.values[key]; !ok {
		m.values[key] = val
		cur_update := "Add:" + key + ":" + strconv.Itoa(val)
		m.updates = append(m.updates, cur_update)
	}
	mu.Unlock()
}

func (m *Map) Delete(key string) {
	mu.Lock()
	if _, ok := m.values[key]; ok {
		delete(m.values, key)
		cur_update := "Delete:" + key
		m.updates = append(m.updates, cur_update)
	}
	mu.Unlock()
}

func (m *Map) Update(key string, newVal int) {
	mu.Lock()
	if _, ok := m.values[key]; ok {
		m.values[key] = newVal
		cur_update := "Update:" + key + ":" + strconv.Itoa(newVal)
		m.updates = append(m.updates, cur_update)
	}
	mu.Unlock()
}

func (m *Map) SetRemoteVal(rid int, opt_name string, key string, value int) {
	if opt_name == "Add" {
		m.Add(key, value)
	} else if opt_name == "Delete" {
		m.Delete(key)
	} else if opt_name == "Update" {
		m.Update(key, value)
	}
}

func (m *Map) Merge(rid int, r_updates []string) {
	fmt.Println("Starting to merge req from replica_", rid)

	if len(r_updates) > 0 {
		for i := 0; i < len(r_updates); i++ {
			reqs := strings.Split(r_updates[i], ":")
			r_opt := reqs[0]
			r_key := reqs[1]
			if len(reqs) > 2 {
				r_value, _ := strconv.Atoi(reqs[2])
				m.SetRemoteVal(rid, r_opt, r_key, r_value)
			} else {
				m.SetRemoteVal(rid, r_opt, r_key, math.MaxInt)
			}
		}
	}

	fmt.Print("Merged: ")
	m.Print()
}

func (m *Map) Print() {
	fmt.Print("Map:", m.id, " ")
	fmt.Println(m.Values())
}

// func (m *Map) ToMarshal() []byte {
// 	data := append([]string{strconv.Itoa(m.Id())}, m.updates...)
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		fmt.Println("Error while marshaling updates")
// 		return nil
// 	}
// 	m.updates = []string{}
// 	return jsonData
// }

func (m *Map) ToMarshal() []byte {

	msg := &_map.SyncMsg{
		Id:      int32(m.Id()),
		Updates: []string{},
	}

	msg.Updates = append(msg.Updates, m.updates...)
	m.updates = []string{}

	serialized, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	return serialized
}

// func FromMarshalData(bytes []byte) (int, []string) {
// 	var remote_updates []string
// 	err := json.Unmarshal(bytes, &remote_updates)

// 	if err != nil {
// 		fmt.Println("Error while unmarshaling ", err)
// 		return -1, nil
// 	}

// 	rid, _ := strconv.Atoi(remote_updates[0])
// 	if len(remote_updates) == 1 {
// 		return rid, nil
// 	}
// 	r_updates := remote_updates[1:]
// 	return rid, r_updates
// }

func FromMarshalData(bytes []byte) (int, []string) {
	var msg _map.SyncMsg
	err := proto.Unmarshal(bytes, &msg)
	if err != nil {
		fmt.Println("unmarshalling error")
		return -1, nil
	}

	// fmt.Println("id: ", msg.GetId())
	// fmt.Println("Updates: ", msg.GetUpdates())

	return int(msg.GetId()), msg.GetUpdates()
}
