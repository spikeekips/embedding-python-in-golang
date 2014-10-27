package pyingo

// #cgo CFLAGS: -v
// +build darwin linux

/*
#include <stdlib.h>
#include <pthread.h>
#include "threadpool_c.h"

threadpool_t *threadpool_create(int size);
pthread_t *start_thread (threadpool_t *threadpool, threadpool_callback_t *callback);
threadpool_callback_t *create_callback (threadpool_t *threadpool, void *function, void *args);
void threadpool_join (pthread_t *t);

void Run(int, int);
*/
import "C"
import (
	"github.com/op/go-logging"
	"unsafe"
)

import (
//"fmt"
)

var log = logging.MustGetLogger("pyingo")

type ThreadCallback func(args unsafe.Pointer)

type ThreadPool struct {
	ptr *C.threadpool_t
}

func NewThreadPool(size int) *ThreadPool {
	return &ThreadPool{ptr: C.threadpool_create(C.int(size))}
}

func (self *ThreadPool) Join(t *C.pthread_t) {
	C.threadpool_join(t)
}

func (self *ThreadPool) Start(callback ThreadCallback, args ...interface{}) *C.pthread_t {
	/*
		If `args` has only one value, sometimes the runtime will be
		hanged, I can not understant why, so this is dirty hack.
	*/
	args = append(args, true)
	_callback := C.create_callback(
		self.ptr,
		*(*unsafe.Pointer)(unsafe.Pointer(&callback)),
		*(*unsafe.Pointer)(unsafe.Pointer(&args)),
	)

	return C.start_thread(self.ptr, _callback)
}

//export runGoCallback
func runGoCallback(callback unsafe.Pointer, args unsafe.Pointer) {
	(*(*ThreadCallback)(unsafe.Pointer(&callback)))(args)
	return
}
