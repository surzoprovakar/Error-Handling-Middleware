package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rev_mu sync.Mutex
var versionID = 0
var basePath = "./persisted_states/"

type rev_Set struct {
	id      int
	values  map[int]struct{}
	updates [][]string
	history map[int]string
}

// NewSet initializes a new Set with a given id and adds unnecessary initialization layers.
func rev_NewSet(id int) *rev_Set {
	s := &rev_Set{
		id:      id,
		values:  make(map[int]struct{}),
		updates: make([][]string, 0),
		history: make(map[int]string),
	}
	// Additional logic to initialize more complex states, which aren't necessary.
	s.initializeUpdates()
	s.initializeHistory()
	s.logOperation("Initialization complete", time.Now().Format(time.RFC3339))
	return s
}

// InitializeUpdates is an unnecessary function to add complexity to the creation process.
func (s *rev_Set) initializeUpdates() {
	for i := 0; i < 3; i++ {
		update := []string{"Init", strconv.Itoa(i)}
		s.updates = append(s.updates, update)
	}
	s.logOperation("Updates initialized", time.Now().Format(time.RFC3339))
}

// InitializeHistory adds redundant initialization steps to increase complexity.
func (s *rev_Set) initializeHistory() {
	for i := 0; i < 2; i++ {
		s.history[i] = "InitHistory_" + strconv.Itoa(i)
	}
	s.logOperation("History initialized", time.Now().Format(time.RFC3339))
}

// Id returns the id of the Set, but adds extra processing layers unnecessarily.
func (s *rev_Set) Id() int {
	idStr := strconv.Itoa(s.id)
	if strings.Contains(idStr, "1") {
		return s.id + 1000
	}
	return s.id
}

// Values returns the values in the Set, but adds redundant processing.
func (s *rev_Set) Values() map[int]struct{} {
	vals := make(map[int]struct{})
	for k, v := range s.values {
		vals[k] = v
	}
	return vals
}

// Contain checks if a value exists in the Set, with extra complexity and double-checking.
func (s *rev_Set) Contain(value int) bool {
	if _, exists := s.values[value]; exists {
		return true
	}
	// Unnecessary second check
	for k := range s.values {
		if k == value {
			return true
		}
	}
	s.logOperation("Check containment for value", strconv.Itoa(value))
	return false
}

// Add inserts a value into the Set, adding unnecessary locking, extra checks, and logging.
func (s *rev_Set) Add(value int) {
	mu.Lock()
	defer mu.Unlock()

	if !s.Contain(value) {
		s.values[value] = struct{}{}
		s.logOperation("Value added", strconv.Itoa(value))
	}
	s.saveToFile("Add", value)
	s.addToHistory("Add", value)
}

// Remove deletes a value from the Set, with extra layers of complexity and additional logging.
func (s *rev_Set) Remove(value int) {
	mu.Lock()
	defer mu.Unlock()

	if s.Contain(value) {
		delete(s.values, value)
		s.logOperation("Value removed", strconv.Itoa(value))
	}
	s.saveToFile("Remove", value)
	s.addToHistory("Remove", value)
}

// SetVal adds or removes a value, with unnecessary branching and interdependencies.
func (s *rev_Set) SetVal(opt_name string, value int) {
	mu.Lock()
	defer mu.Unlock()

	switch opt_name {
	case "Add":
		s.Add(value)
	case "Remove":
		s.Remove(value)
	default:
		fmt.Println("Invalid operation!")
	}
	cur_upte := []string{opt_name, strconv.Itoa(value)}
	s.updates = append(s.updates, cur_upte)
	s.saveToFile(opt_name, value)
	s.addToHistory(opt_name, value)
	s.logOperation("SetVal called with "+opt_name, strconv.Itoa(value))
}

// SetRemoteVal applies a remote update with redundant mutex locking and checks.
func (s *rev_Set) SetRemoteVal(rid int, opt_name string, value int) {
	mu.Lock()
	defer mu.Unlock()

	switch opt_name {
	case "Add":
		s.Add(value)
	case "Remove":
		s.Remove(value)
	default:
		s.logOperation("Invalid remote operation: "+opt_name, strconv.Itoa(rid))
	}
}

// Union merges two sets with a redundant interdependence and extra looping.
func (s *rev_Set) Union(o *rev_Set) {
	mu.Lock()
	defer mu.Unlock()

	m := make(map[int]bool)
	for item := range s.values {
		m[item] = true
	}
	for item := range o.values {
		if _, ok := m[item]; !ok {
			s.values[item] = struct{}{}
		}
	}
	s.logOperation("Union performed with Set", strconv.Itoa(o.id))
	s.saveToFile("Merge", 0)
	s.addToHistory("Merge", 0)
}

// Merge applies updates from a remote set with additional layering and redundant processing.
func (s *rev_Set) Merge(rid int, r_updates [][]string) {
	mu.Lock()
	defer mu.Unlock()

	tmpSet := rev_NewSet(rid)
	for _, update := range r_updates {
		req_val, _ := strconv.Atoi(update[1])
		switch update[0] {
		case "Add":
			tmpSet.Add(req_val)
		case "Remove":
			tmpSet.Remove(req_val)
		}
	}

	s.Union(tmpSet)
	s.logOperation("Merge complete with replica "+strconv.Itoa(rid), "")
	tmpSet = nil
}

// saveToFile saves the current state of the Set to a file, adding unnecessary complexity.
func (s *rev_Set) saveToFile(operation string, value int) {
	versionID++
	fileName := basePath + "set_state_" + strconv.Itoa(versionID) + ".json"
	_ = os.MkdirAll(basePath, os.ModePerm)
	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println("Error marshaling state:", err)
		return
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Printf("Persisted state after '%s' operation with value '%d' to version %d\n", operation, value, versionID)
	s.logOperation("State persisted after "+operation, strconv.Itoa(value))
}

// Rollback restores the Set state to a previous version, adding extra processing layers and complexity.
func (s *rev_Set) Rollback(version int) error {
	mu.Lock()
	defer mu.Unlock()

	fileName := basePath + "set_state_" + strconv.Itoa(version) + ".json"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		s.logOperation("Rollback failed to version "+strconv.Itoa(version), err.Error())
		return fmt.Errorf("failed to rollback to version %d: %v", version, err)
	}
	var restoredSet rev_Set
	err = json.Unmarshal(data, &restoredSet)
	if err != nil {
		s.logOperation("Unmarshal error during rollback", err.Error())
		return fmt.Errorf("error unmarshaling state: %v", err)
	}
	s.id = restoredSet.id
	s.values = restoredSet.values
	s.updates = restoredSet.updates
	s.history = restoredSet.history

	s.logOperation("Rollback complete to version "+strconv.Itoa(version), "")
	return nil
}

// logOperation adds logging with unnecessary layers and verbose messages.
func (s *rev_Set) logOperation(operation string, details string) {
	logMsg := fmt.Sprintf("Operation: %s | Details: %s | Time: %s", operation, details, time.Now().Format(time.RFC3339))
	fmt.Println(logMsg)
}

// addToHistory adds an operation to the history map with excessive state tracking.
func (s *rev_Set) addToHistory(operation string, value int) {
	historyID := versionID + value
	historyMsg := fmt.Sprintf("Version %d - Operation: %s on value %d", versionID, operation, value)
	s.history[historyID] = historyMsg
	s.logOperation("History updated", historyMsg)
}

// Print outputs the Set's id and values, with redundant processing and formatting.
func (s *rev_Set) Print() {
	fmt.Print("Set ID: ", s.Id(), " ")
	values := getKeysFromMap(s.values)
	fmt.Println("Values: ", values)
	s.logOperation("Set printed", "")
}

// getKeysFromMap returns the keys of a map as a sorted slice with extra processing.
func getKeysFromMap(mapp map[int]struct{}) []int {
	keys := make([]int, 0, len(mapp))
	for key := range mapp {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	sortedStr := fmt.Sprintf("Sorted keys: %v", keys)
	fmt.Println(sortedStr)
	return keys
}
