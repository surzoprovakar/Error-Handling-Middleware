package main

/*
#cgo LDFLAGS: -L. -lexample
#include <stdbool.h>
extern void Create_File(int id);
extern void Record(int id, const char* val, const char* opt_name, int up_value);
extern void Remote_Record(int id, int rid, const char *opt_name, int up_value);
extern bool undo(int r_id, const char *undo_update, int undo_val, int r_LT);
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// var muw sync.Mutex

var updates = 0
var undo_opt_not_arrived string
var reject_local_update map[string]bool

func keys4Persist(mapp map[int]struct{}) string {
	keys := make([]int, 0, len(mapp))
	for key := range mapp {
		keys = append(keys, key)
	}
	sort.Ints(keys[:])

	jsonStr, _ := json.Marshal(keys)
	return string(jsonStr)
}

func createSetWrap(id int) *Set {
	s := NewSet(id)
	//fmt.Println("after creating ", id)
	C.Create_File(C.int(id))
	//go c.check()
	return s
}

func (s *Set) setValueWrap(opt_name string, value int) {

	if s.Contain(value) && opt_name == "Add" {
		rejct_up := opt_name + "_" + strconv.Itoa(value)
		reject_local_update = make(map[string]bool)
		reject_local_update[rejct_up] = true

		// fmt.Println("reject_local_update")
		// fmt.Println(reject_local_update)
	}
	s.SetVal(opt_name, value)
	val := keys4Persist(s.Values())
	C.Record(C.int(s.Id()), C.CString(val), C.CString(opt_name), C.int(value))
}

func (s *Set) setRemoteWrap(rid int, opt_name string, value int) {
	// s.SetRemoteVal(rid, opt_name, value)
	// val := keys4Persist(s.Values())
	C.Remote_Record(C.int(s.Id()), C.int(rid), C.CString(opt_name), C.int(value))
}

func (s *Set) mergeWrap(rid int, r_updates [][]string) {
	if len(r_updates) > 0 {
		// fmt.Println(r_updates)
		if r_updates[0][0] == "undo" {
			fmt.Println("undo requested from replica:", rid)
			// fmt.Println("undo req from others")
			// fmt.Println(r_updates)
			undo_info := strings.Split(r_updates[0][1], "_")
			undo_val, _ := strconv.Atoi(undo_info[1])
			remote_Lt, _ := strconv.Atoi(undo_info[2])
			if C.undo(C.int(rid), C.CString(undo_info[0]), C.int(undo_val), C.int(remote_Lt)) {
				fmt.Println("Undo executing")
				if undo_info[0] == "Add" {
					s.SetRemoteVal(rid, "Remove", undo_val)
					s.setRemoteWrap(rid, "Remove", undo_val)
				} else if undo_info[0] == "Remove" {
					s.SetRemoteVal(rid, "Add", undo_val)
					s.setRemoteWrap(rid, "Add", undo_val)
				}
				val := keys4Persist(s.Values())
				fmt.Println("Set:", s.Id(), val)
				C.Record(C.int(s.Id()), C.CString(val), C.CString("Undo done from replica"), C.int(rid))

				// Revisit the previous rejected local updates
				fmt.Println("Revisit previous updates")
				check_opt := undo_info[0] + "_" + undo_info[1]
				fmt.Println("revisit opt:", check_opt)

				if reject_local_update[check_opt] == true {
					fmt.Println("Execute previous rejected local update due to undone now")
					s.SetVal(undo_info[0], undo_val)
					revist_val := keys4Persist(s.Values())
					C.Record(C.int(s.Id()), C.CString(revist_val), C.CString(undo_info[0]+"_revisit"), C.int(undo_val))
				}

			} else {
				fmt.Println("Undo not required as this operation has not merged yet")
				fmt.Println("Recording this undo info for future rejection")
				undo_opt_not_arrived += strconv.Itoa(rid) + "_" + undo_info[0] + "_" + undo_info[1]
				// fmt.Println("undo_opt_not_arrived:", undo_opt_not_arrived)
			}
		} else {
			// fmt.Println("req replica:", rid)
			// fmt.Println(r_updates)
			for i := 0; i < len(r_updates); i++ {

				check := strconv.Itoa(rid) + "_" + r_updates[i][0] + "_" + r_updates[i][1]
				// fmt.Println("check:", check)
				if check == undo_opt_not_arrived {
					fmt.Println("No need to merge as already undo req by same replica")
				} else {
					req_val, _ := strconv.Atoi(r_updates[i][1])
					s.setRemoteWrap(rid, r_updates[i][0], req_val)
				}
			}
			s.Merge(rid, r_updates)

			val := keys4Persist(s.Values())
			C.Record(C.int(s.Id()), C.CString(val), C.CString("Merge from replica"), C.int(rid))
		}
	}
}
