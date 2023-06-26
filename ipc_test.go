package gocuda

import (
	"gorgonia.org/cu"
	"io/fs"
	"log"
	"os"
	"testing"
	"time"
	"unsafe"
)

func TestGetIpcMemHandle(t *testing.T) {
	var size int64 = 64
	ctx, _ := cu.Device(0).MakeContext(cu.SchedAuto)
	defer ctx.Destroy()
	ptr, err := cu.MemAlloc(int64(size))
	if err != nil {
		t.Fatal(err)
	}
	hostMem := make([]byte, size)
	for i := 0; i < int(size); i++ {
		hostMem[i] = byte(i)
	}
	h, err := GetIpcMemHandle(ptr)
	if err != nil {
		t.Fatal(err)
	}
	defer CloseIpcMemHandle(ptr)
	cu.MemcpyHtoD(ptr, unsafe.Pointer(&hostMem[0]), size)
	log.Println(h.handle)
	log.Println(h.GoBytes())

	log.Println(ptr.IsCUDAMemory())
	log.Println(ptr.AddressRange())

	os.WriteFile("ipc.txt", h.GoBytes(), fs.ModePerm)
	for {
		err = cu.MemcpyDtoH(unsafe.Pointer(&hostMem[0]), ptr, size)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(hostMem)
		time.Sleep(2 * time.Second)
	}
}

func TestOpenIpcMemHandle(t *testing.T) {
	var size int64 = 64
	ctx, _ := cu.Device(0).MakeContext(cu.SchedAuto)
	defer ctx.Destroy()
	data, err := os.ReadFile("ipc.txt")
	if err != nil {
		t.Fatal(err, "run TestGetIpcMemHandle first")
	}
	var ptr cu.DevicePtr
	handle := IpcHandleFromBytes(data)
	log.Println(handle.GoBytes())
	ptr, err = OpenIpcMemHandle(handle)
	if err != nil {
		t.Fatal()
	}
	if !ptr.IsCUDAMemory() {
		t.Fatal(ptr)
	}
	log.Println(ptr.AddressRange())
	hostMem := make([]byte, size)
	err = cu.MemcpyDtoH(unsafe.Pointer(&hostMem[0]), ptr, size)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(hostMem)
	for i := 0; i < len(hostMem); i++ {
		hostMem[i] += 1
	}
	err = cu.MemcpyHtoD(ptr, unsafe.Pointer(&hostMem[0]), size)
	log.Println(hostMem)
}
