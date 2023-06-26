package gocuda

/*
#include "cuda.h"
*/
import "C"
import "unsafe"

type IpcEventHandle struct {
	handle C.CUipcEventHandle
}

func (h *IpcEventHandle) c() C.CUipcEventHandle {
	return h.handle
}

type IpcMemHandle struct {
	handle C.CUipcMemHandle
}

func (h *IpcMemHandle) c() C.CUipcMemHandle {
	return h.handle
}

func (h *IpcMemHandle) GoBytes() []byte {
	return C.GoBytes(unsafe.Pointer(&(h.handle.reserved)), 64)
}

func IpcHandleFromBytes(data []byte) IpcMemHandle {
	var handle C.CUipcMemHandle
	for i := 0; i < 64; i++ {
		handle.reserved[i] = C.char(data[i])
	}
	return IpcMemHandle{handle: handle}
}
