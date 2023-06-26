package gocuda

// #include <cuda.h>
// #cgo pkg-config:cuda
import "C"
import "gorgonia.org/cu"

func CloseIpcMemHandle(dptr cu.DevicePtr) (err error) {
	Cdptr := C.CUdeviceptr(dptr)
	return result(C.cuIpcCloseMemHandle(Cdptr))
}

func GetIpcMemHandle(dptr cu.DevicePtr) (pHandle IpcMemHandle, err error) {
	Cdptr := C.CUdeviceptr(dptr)
	var CpHandle C.CUipcMemHandle
	err = result(C.cuIpcGetMemHandle(&CpHandle, Cdptr))
	pHandle = IpcMemHandle{CpHandle}

	return
}

func OpenIpcMemHandle(handle IpcMemHandle) (ptr cu.DevicePtr, err error) {
	var Cdptr C.CUdeviceptr
	Chandle := handle.c()
	err = result(C.cuIpcOpenMemHandle(&Cdptr, Chandle, C.CU_IPC_MEM_LAZY_ENABLE_PEER_ACCESS))
	ptr = cu.DevicePtr(Cdptr)
	return
}
