package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex

type Set struct {
	id      int
	values  map[int]struct{}
	updates [][]string
}

func NewSet(id int) *Set {
	//fmt.Println("after creating ", id)
	createFile(id)
	//go c.check()
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

	if opt_name == "Add" {
		rejct_up := opt_name + "_" + strconv.Itoa(value)
		addhoc_reject_local_update = make(map[string]bool)
		addhoc_reject_local_update[rejct_up] = true

		// fmt.Println("reject_local_update")
		// fmt.Println(reject_local_update)
	}
	val := addhoc_keys4Persist(s.Values())
	record(s.Id(), val, opt_name, value)
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

// Wrapper's functionality
var addhoc_updates = 0
var addhoc_undo_opt_not_arrived string
var addhoc_reject_local_update map[string]bool

func addhoc_keys4Persist(mapp map[int]struct{}) string {
	keys := make([]int, 0, len(mapp))
	for key := range mapp {
		keys = append(keys, key)
	}
	sort.Ints(keys[:])

	jsonStr, _ := json.Marshal(keys)
	return string(jsonStr)
}

func (s *Set) addhoc_setRemoteWrap(rid int, opt_name string, value int) {
	// s.SetRemoteVal(rid, opt_name, value)
	// val := keys4Persist(s.Values())
	remoteRecord(s.Id(), rid, opt_name, value)
}

func (s *Set) addhoc_mergeWrap(rid int, r_updates [][]string) {
	if len(r_updates) > 0 {
		// fmt.Println(r_updates)
		if r_updates[0][0] == "undo" {
			fmt.Println("undo requested from replica:", rid)
			// fmt.Println("undo req from others")
			// fmt.Println(r_updates)
			undo_info := strings.Split(r_updates[0][1], "_")
			undo_val, _ := strconv.Atoi(undo_info[1])
			remote_Lt, _ := strconv.Atoi(undo_info[2])
			if undo(rid, undo_info[0], undo_val, remote_Lt) {
				fmt.Println("Undo executing")
				if undo_info[0] == "Add" {
					s.SetRemoteVal(rid, "Remove", undo_val)
					s.addhoc_setRemoteWrap(rid, "Remove", undo_val)
				} else if undo_info[0] == "Remove" {
					s.SetRemoteVal(rid, "Add", undo_val)
					s.addhoc_setRemoteWrap(rid, "Add", undo_val)
				}
				val := addhoc_keys4Persist(s.Values())
				fmt.Println("Set:", s.Id(), val)
				record(s.Id(), val, "Undo done from replica", rid)

				// Revisit the previous rejected local updates
				fmt.Println("Revisit previous updates")
				check_opt := undo_info[0] + "_" + undo_info[1]
				fmt.Println("revisit opt:", check_opt)

				if addhoc_reject_local_update[check_opt] == true {
					fmt.Println("Execute previous rejected local update due to undone now")
					s.SetVal(undo_info[0], undo_val)
					revist_val := addhoc_keys4Persist(s.Values())
					record(s.Id(), revist_val, undo_info[0]+"_revisit", undo_val)
				}

			} else {
				fmt.Println("Undo not required as this operation has not merged yet")
				fmt.Println("Recording this undo info for future rejection")
				addhoc_undo_opt_not_arrived += strconv.Itoa(rid) + "_" + undo_info[0] + "_" + undo_info[1]
				// fmt.Println("undo_opt_not_arrived:", undo_opt_not_arrived)
			}
		} else {
			// fmt.Println("req replica:", rid)
			// fmt.Println(r_updates)
			for i := 0; i < len(r_updates); i++ {

				check := strconv.Itoa(rid) + "_" + r_updates[i][0] + "_" + r_updates[i][1]
				// fmt.Println("check:", check)
				if check == addhoc_undo_opt_not_arrived {
					fmt.Println("No need to merge as already undo req by same replica")
				} else {
					req_val, _ := strconv.Atoi(r_updates[i][1])
					s.addhoc_setRemoteWrap(rid, r_updates[i][0], req_val)
				}
			}
			s.Merge(rid, r_updates)

			val := addhoc_keys4Persist(s.Values())
			record(s.Id(), val, "Merge from replica", rid)
		}
	}
}

// C++ portion
type LamportClock struct {
	time int
	mu   sync.Mutex
}

func (lc *LamportClock) localEvent() int {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.time++
	return lc.time
}

func (lc *LamportClock) getTime() int {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.time
}

var id2ver = make(map[int]int)
var updateList []Update
var cClock LamportClock
var rid2LT = make(map[int]*LamportClock)

type Update struct {
	ID    int
	Ver   int
	Value string
}

func createFile(id int) {
	dirPath := "DBs/"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0700)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}
	filename := dirPath + strconv.Itoa(id) + ".txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	fmt.Println("Storage file created for Replica", id)
	file.Close()

	id2ver[id] = 0
}

func record(id int, val string, optName string, upValue int) {
	// Write to physical file
	filename := "DBs/" + strconv.Itoa(id) + ".txt"
	id2ver[id]++
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	logEntry := fmt.Sprintf("R_ID:%d Ver:%d St:%s LT:%d Opt:%s_%d Type:Local\n", id, id2ver[id], val, cClock.localEvent(), optName, upValue)
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Update list
	updateList = append(updateList, Update{id, id2ver[id], val})
}

func remoteRecord(id int, rid int, optName string, upValue int) {
	filename := "DBs/" + strconv.Itoa(id) + ".txt"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	remoteLT := rid2LT[rid].localEvent()
	logEntry := fmt.Sprintf("R_ID:%d Opt:%s_%d Type:{Remote, rid:%d local_LT:%d}\n", id, optName, upValue, rid, remoteLT)
	if _, err := file.WriteString(logEntry); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func undo(rID int, undoUpdate string, undoVal int, rLT int) bool {
	if rid2LT[rID].getTime() > rLT {
		return true
	}
	return false
}
