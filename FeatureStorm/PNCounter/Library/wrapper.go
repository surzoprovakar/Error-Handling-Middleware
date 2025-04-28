package main

/*
#cgo LDFLAGS: -L. -lexample
#include <stdbool.h>
extern void Create_File(int id);
extern void Record(int id, const char* val, const char* opt_name);
extern bool undo(int r_id, const char* undo_update, int r_LT);
extern void Remote_Record(int id, int rid, const char* val, const char *opt_name);
*/
import "C"

import (
	"fmt"
	"strconv"
)

var check_undo_interval = 7
var updates = 0

func (c *Counter) setValueWrap(v int, opt_name string) {
	//fmt.Println("Before setValue:", v)
	c.SetVal(v, opt_name)
	val := strconv.Itoa(c.Value())
	C.Record(C.int(c.Id()), C.CString(val), C.CString(opt_name))
	//fmt.Println("After setValue", c.Value())
	updates++

	// if updates%check_undo_interval == 0 {
	// 	c.undo(c.Id())
	// }
}

func (c *Counter) setRemoteWrap(rid int, opt_name string) {
	c.SetRemoteVal(rid, opt_name)
	val := strconv.Itoa(c.Value())
	C.Remote_Record(C.int(c.Id()), C.int(rid), C.CString(val), C.CString(opt_name))
}

func createWrap(id int) *Counter {
	c := NewCounter(id)
	//fmt.Println("after creating ", id)
	C.Create_File(C.int(id))
	//go c.check()
	return c
}

func (c *Counter) mergeWrap(rid int, r_updates []string) {
	if len(r_updates) > 0 {
		// fmt.Println("requested merge:", r_updates)
		if r_updates[0] == "undo" {
			fmt.Println("undo requested from replica:", rid)
			remote_LT, _ := strconv.Atoi(r_updates[2])
			if C.undo(C.int(rid), C.CString(r_updates[1]), C.int(remote_LT)) {
				fmt.Println("Undo executing")
				if r_updates[1] == "Inc" {
					c.setRemoteWrap(rid, "Dec")
				} else if r_updates[1] == "Dec" {
					c.setRemoteWrap(rid, "Inc")
				}
				// c.Merge(rid, r_updates[1:2]) //undo only one operation
			} else {
				fmt.Println("Undo not required as this operation has not merged yet")
				fmt.Println("Recording this undo info for future rejection")
			}
		} else {
			c.Merge(rid, r_updates)
		}
	}
}

// func (c *Counter) undo(id int) {
// 	fmt.Println("check undo")

// 	res := C.undo(C.int(id), C.int(1))
// 	if res != nil {
// 		undo_val, _ := strconv.Atoi(C.GoString(res))
// 		fmt.Println("undo val is: ", undo_val)
// 		c.setValueWrap(undo_val, "undo")
// 	}
// 	time.Sleep(100 * time.Millisecond)
// }
