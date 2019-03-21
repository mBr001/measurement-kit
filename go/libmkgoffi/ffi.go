package main

import "C"

import (
	"sync"

	"github.com/measurement-kit/measurement-kit/go/task"
)

// mutex protects table from concurrent access
var mutex sync.Mutex

// table contains the state of running tasks. Note that the first
// entry of the table is always nil.
var table [256]*task.State

// getstate returns the state bound to a handle or nil
func getstate(handle int) *task.State {
	if handle >= len(table) || table[handle] == nil {
		return nil
	}
	return table[handle]
}

//export mkgo_ffi_task_start
func mkgo_ffi_task_start(settings *C.char) int {
	mutex.Lock()
	defer mutex.Unlock()
	const minhandle = 1 // first entry must be nil
	for handle := minhandle; settings != nil && handle < len(table); handle += 1 {
		if table[handle] == nil {
			table[handle] = task.New(C.GoString(settings))
			return handle
		}
	}
	return 0
}

//export mkgo_ffi_task_wait_for_next_event
func mkgo_ffi_task_wait_for_next_event(handle int) *C.char {
	mutex.Lock()
	defer mutex.Unlock()
	state := getstate(handle)
	if state == nil {
		return nil
	}
	return C.CString(state.WaitForNextEvent())
}

//export mkgo_ffi_task_is_done
func mkgo_ffi_task_is_done(handle int) int {
	mutex.Lock()
	defer mutex.Unlock()
	state := getstate(handle)
	if state == nil || state.IsDone() {
		return 1
	}
	return 0
}

//export mkgo_ffi_task_interrupt
func mkgo_ffi_task_interrupt(handle int) {
	mutex.Lock()
	defer mutex.Unlock()
	state := getstate(handle)
	if state != nil {
		state.Interrupt()
	}
}

//export mkgo_ffi_task_destroy
func mkgo_ffi_task_destroy(handle int) {
	mutex.Lock()
	defer mutex.Unlock()
	state := getstate(handle)
	if state == nil {
		return
	}
	state.Interrupt()
	for !state.IsDone() {
		state.WaitForNextEvent() // drain
	}
	table[handle] = nil
}

func main() {
}
