// +build windows

/*
Copyright 2018 Ryan Dahl <ry@tinyclouds.org>. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to
deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
IN THE SOFTWARE.
*/
package v8worker2

import "C"

import "sync"
import "runtime"

type workerTableIndex int

var workerTableLock sync.Mutex


// This table will store all pointers to all active workers. Because we can't safely
// pass pointers to Go objects to C, we instead pass a key to this table.
var workerTable = make(map[workerTableIndex]*worker)

// Keeps track of the last used table index. Incremeneted when a worker is created.
var workerTableNextAvailable workerTableIndex = 0

// To receive messages from javascript.
type ReceiveMessageCallback func(msg []byte) []byte

// To resolve modules from javascript.
type ModuleResolverCallback func(moduleName, referrerName string) int


// Internal worker struct which is stored in the workerTable.
// Weak-ref pattern https://groups.google.com/forum/#!topic/golang-nuts/1ItNOOj8yW8/discussion
type worker struct {
	cWorker    *C.worker
	cb         ReceiveMessageCallback
	tableIndex workerTableIndex
}

// This is a golang wrapper around a single V8 Isolate.
type Worker struct {
	*worker
	disposed bool
}

// Return the V8 version E.G. "6.6.164-v8worker2"
func Version() string {
	return "6.6.164-v8worker2"
}

// Sets V8 flags. Returns the input args but with the V8 flags removed.
// Use --help to print a list of flags to stdout.
func SetFlags(args []string) []string {
	return args
}


//export ResolveModule
func ResolveModule(moduleSpecifier *C.char, referrerSpecifier *C.char, resolverToken int) C.int {
	return 0
}

// Creates a new worker, which corresponds to a V8 isolate. A single threaded
// standalone execution context.
func New(cb ReceiveMessageCallback) *Worker {
	workerTableLock.Lock()
	w := &worker{
		cb:         cb,
		tableIndex: workerTableNextAvailable,
	}

	externalWorker := &Worker{
		worker:   w,
		disposed: false,
	}

	runtime.SetFinalizer(externalWorker, func(final_worker *Worker) {
		final_worker.Dispose()
	})
	return externalWorker
}

// Forcefully frees up memory associated with worker.
// GC will also free up worker memory so calling this isn't strictly necessary.
func (w *Worker) Dispose() {
	if w.disposed {
		panic("worker already disposed")
	}
	w.disposed = true
}

// Load and executes a javascript file with the filename specified by
// scriptName and the contents of the file specified by the param code.
func (w *Worker) Load(scriptName string, code string) error {

	return nil
}

// LoadModule loads and executes a javascript module with filename specified by
// scriptName and the contents of the module specified by the param code.
// All `import` dependencies must be loaded before a script otherwise it will error.
func (w *Worker) LoadModule(scriptName string, code string, resolve ModuleResolverCallback) error {

	return nil
}

// Same as Send but for []byte. $recv callback will get an ArrayBuffer.
func (w *Worker) SendBytes(msg []byte) error {
		return nil
}

// Terminates execution of javascript
func (w *Worker) TerminateExecution() {

}
