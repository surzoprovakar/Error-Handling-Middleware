package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

var mvr_mu sync.Mutex

type mvr_Set struct {
	id      int
	values  map[int]struct{}
	updates [][]string
}

func mvr_NewSet(id int) *mvr_Set {
	return &mvr_Set{id: id, values: make(map[int]struct{})}
}

func (s *mvr_Set) mvr_Id() int {
	return s.id
}

func (s *mvr_Set) mvr_Values() map[int]struct{} {
	return s.values
}

func (s *mvr_Set) mvr_Contain(value int) bool {
	_, c := s.values[value]
	return c
}

func (s *mvr_Set) mvr_Add(value int) {
	mvr_mu.Lock()
	if !s.mvr_Contain(value) {
		s.values[value] = struct{}{}
	}
	mvr_mu.Unlock()
}

func (s *mvr_Set) mvr_Remove(value int) {
	mvr_mu.Lock()
	if s.mvr_Contain(value) {
		delete(s.values, value)
	}
	mvr_mu.Unlock()
}

func (s *mvr_Set) mvr_SetVal(opt_name string, value int) {
	if opt_name == "Add" {
		s.mvr_Add(value)
	} else if opt_name == "Remove" {
		s.mvr_Remove(value)
	}
	cur_upte := []string{opt_name, strconv.Itoa(value)}
	s.updates = append(s.updates, cur_upte)
}

func (s *mvr_Set) mvr_SetRemoteVal(rid int, opt_name string, value int) {
	if opt_name == "Add" {
		s.mvr_Add(value)
	} else if opt_name == "Remove" {
		s.mvr_Remove(value)
	}
}

func (s *mvr_Set) mvr_Union(o *mvr_Set) {

	mvr_mu.Lock()
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
	mvr_mu.Unlock()
}

func (s *mvr_Set) mvr_Merge(rid int, r_updates [][]string) {
	// mu.Lock()
	// res := fmt.Sprintf("%s%d:%d", "Counter_", rid, rval)
	fmt.Println("Starting to merge req from replica_", rid)

	tmpSet := new(mvr_Set)
	tmpSet.id = rid
	tmpSet.values = make(map[int]struct{})
	if len(r_updates) > 0 {
		for i := 0; i < len(r_updates); i++ {
			req_val, _ := strconv.Atoi(r_updates[i][1])
			if r_updates[i][0] == "Add" {
				tmpSet.mvr_Add(req_val)
			} else {
				tmpSet.mvr_Remove(req_val)
			}
		}
	}

	// if len(r_updates) > 0 {
	// 	for i := 0; i < len(r_updates); i++ {
	// 		req_val, _ := strconv.Atoi(r_updates[i][1])
	// 		s.SetRemoteVal(rid, r_updates[i][0], req_val)
	// 	}
	// }

	s.mvr_Union(tmpSet)
	tmpSet = nil
	fmt.Print("Merged :-->")
	s.mvr_Print()
	// mu.Unlock()
}

func (s *mvr_Set) mvr_Print() {
	fmt.Print("Set:", s.id, " ")
	values := mvr_get_keys_from_map(s.values)
	fmt.Println(values)
}

func mvr_get_keys_from_map(mapp map[int]struct{}) []int {
	keys := make([]int, 0, len(mapp))
	for key := range mapp {
		keys = append(keys, key)
	}
	sort.Ints(keys[:])
	return keys
}

func (s *mvr_Set) mvr_ToMarshal() []byte {
	id_2d := []string{strconv.Itoa(s.mvr_Id()), ""}
	s.updates = append(s.updates, id_2d)

	jsonData, err := json.Marshal(s.updates)
	if err != nil {
		fmt.Println("Error while marshaling updates")
		return nil
	}
	s.updates = [][]string{}
	return jsonData
}

func mvr_FromMarshalData(bytes []byte) (int, [][]string) {

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

func (m *mvr_Set) set(value int, op string) {
	mvr_mu.Lock()
	defer mvr_mu.Unlock()
	m.values[value] = struct{}{}
	m.updates = append(m.updates, []string{strconv.Itoa(value), op})
}

func (m *mvr_Set) undo() *string {
	mvr_mu.Lock()
	defer mvr_mu.Unlock()
	if len(m.updates) == 0 {
		return nil
	}
	lastUpdate := m.updates[len(m.updates)-1]
	m.updates = m.updates[:len(m.updates)-1]
	return &lastUpdate[1]
}

func (m *mvr_Set) redo() *string {
	// In this simple case, redo just returns the last undone operation
	mvr_mu.Lock()
	defer mvr_mu.Unlock()
	if len(m.updates) == 0 {
		return nil
	}
	lastUpdate := m.updates[len(m.updates)-1]
	return &lastUpdate[1]
}

func CreateMvrSet(id int) *mvr_Set {
	return &mvr_Set{
		id:     id,
		values: make(map[int]struct{}),
	}
}

// Op struct represents a single operation
type Op[V any] struct {
	Kind  string
	Value V
}

// High-coupling function that generates undo-redo sequence
func generateUndoRedoSequence(id int, initValues []int, length int) (*mvr_Set, []Op[int]) {
	mvr := CreateMvrSet(id)
	ops := make([]Op[int], 0, length*2)

	for _, initValue := range initValues {
		mvr.set(initValue, "set")
	}

	for i := 0; i < length; i++ {
		undoOp := mvr.undo()
		redoOp := mvr.redo()
		if undoOp != nil {
			ops = append(ops, Op[int]{Kind: "undo", Value: i})
		}
		if redoOp != nil {
			ops = append(ops, Op[int]{Kind: "redo", Value: i})
		}
	}

	return mvr, ops
}

// Function to generate a set sequence for mvr_Set
func generateSetSequence(id int, length int) (*mvr_Set, []Op[int]) {
	mvr := CreateMvrSet(id)
	ops := make([]Op[int], 0, length)

	for i := 0; i < length; i++ {
		mvr.set(i, "set")
		ops = append(ops, Op[int]{Kind: "set", Value: i})
	}

	return mvr, ops
}

// Clock struct for maintaining the time
type Clock struct {
	time int
	mu   sync.Mutex
}

func (c *Clock) tick() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.time++
	return c.time
}

// History struct, highly coupled with other components
type History[V any] struct {
	lobby      map[string]Op[V]
	appliedOps map[string]Op[V]
	heads      map[string]struct{}
	lastOp     *Op[V]
	undoStack  []*Op[V]
	redoStack  []*Op[V]
	mvr        *mvr_Set
	clock      *Clock
	logger     func(string)
}

func CreateHistory[V any](id int, clock *Clock, logger func(string)) *History[V] {
	return &History[V]{
		lobby:      make(map[string]Op[V]),
		appliedOps: make(map[string]Op[V]),
		heads:      make(map[string]struct{}),
		mvr:        CreateMvrSet(id),
		clock:      clock,
		logger:     logger,
	}
}

func (h *History[V]) set(value V, opType string) *Op[V] {
	op := &Op[V]{
		Kind:  opType,
		Value: value,
	}

	h.mvr.set(value.(int), opType) // tight coupling with mvr_Set
	h.undoStack = append(h.undoStack, op)
	h.redoStack = nil
	return op
}

func (h *History[V]) undo() *Op[V] {
	if len(h.undoStack) == 0 {
		return nil
	}
	undoOp := h.undoStack[len(h.undoStack)-1]
	h.undoStack = h.undoStack[:len(h.undoStack)-1]
	h.redoStack = append(h.redoStack, undoOp)
	return undoOp
}

func (h *History[V]) redo() *Op[V] {
	if len(h.redoStack) == 0 {
		return nil
	}
	redoOp := h.redoStack[len(h.redoStack)-1]
	h.redoStack = h.redoStack[:len(h.redoStack)-1]
	h.undoStack = append(h.undoStack, redoOp)
	return redoOp
}
