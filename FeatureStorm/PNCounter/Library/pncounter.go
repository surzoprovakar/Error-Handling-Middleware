package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

var mu sync.Mutex

type CounterIface interface {
	Inc()
	Dec()
	Value() int
}

type Counter struct {
	id      int
	value   int
	updates []string
}

func NewCounter(id int) *Counter {
	return &Counter{id: id, value: 0}
}

func (c *Counter) SetVal(new_val int, opt_name string) {
	c.value = new_val
	c.updates = append(c.updates, opt_name)
}

func (c *Counter) SetRemoteVal(rid int, opt_name string) {
	if opt_name == "Inc" {
		c.value += 1
	} else if opt_name == "Dec" {
		c.value -= 1
	}

}

func (c *Counter) Inc() {
	mu.Lock()
	new_val := c.value + 1
	c.SetVal(new_val, "Inc")
	mu.Unlock()
}

func (c *Counter) Dec() {
	mu.Lock()
	new_val := c.value - 1
	c.SetVal(new_val, "Dec")
	mu.Unlock()
}

func (c *Counter) Id() int {
	return c.id
}

func (c *Counter) Value() int {
	return c.value
}

func (c *Counter) Merge(rid int, r_updates []string) {
	// mu.Lock()
	// res := fmt.Sprintf("%s%d:%d", "Counter_", rid, rval)
	fmt.Println("Starting to merge req from replica_", rid)
	if len(r_updates) > 0 {
		for i := 0; i < len(r_updates); i++ {
			c.SetRemoteVal(rid, r_updates[i])
		}
	}
	fmt.Println("Merged " + c.Print())
	// mu.Unlock()
}

func (c *Counter) Print() string {
	res := fmt.Sprintf("%s%d:%d", "Counter_", c.Id(), c.Value())
	return res
}

/*
func (c *Counter) ToByteArray() []byte {

	a1 := make([]byte, 64)
	a2 := make([]byte, 64)

	binary.LittleEndian.PutUint64(a1, uint64(c.Id()))
	binary.LittleEndian.PutUint64(a2, uint64(c.Value()))

	return append(a1, a2...)
}

func FromByteArray(bytes []byte) (int, int) {

	var r1 = binary.LittleEndian.Uint64(bytes[0 : len(bytes)/2])
	var r2 = binary.LittleEndian.Uint64(bytes[len(bytes)/2:])

	rid := int(r1)
	rval := int(r2)
	return rid, rval
}
*/

func (c *Counter) ToMarshal() []byte {

	data := append([]string{strconv.Itoa(c.Id())}, c.updates...)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error while marshaling updates")
		return nil
	}
	c.updates = []string{}
	return jsonData
}

func FromMarshalData(bytes []byte) (int, []string) {

	var remote_updates []string
	err := json.Unmarshal(bytes, &remote_updates)

	if err != nil {
		fmt.Println("Error while unmarshaling ", err)
		return -1, nil
	}

	rid, _ := strconv.Atoi(remote_updates[0])
	if len(remote_updates) == 1 {
		return rid, nil
	}
	r_updates := remote_updates[1:]
	return rid, r_updates
}
