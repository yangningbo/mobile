// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

// Android redirects stdout and stderr to /dev/null.
// As these are common debugging utilities in Go,
// we redirect them to logcat.
//
// Unfortunately, logcat is line oriented, so we must buffer.

/*
#cgo LDFLAGS: -llog
#include <android/log.h>
#include <string.h>
*/
import "C"

import (
	"bufio"
	"os"
	"unsafe"
)

var ctag = C.CString("Go")

func lineLog(f *os.File, priority C.int) {
	const logSize = 1024 // matches android/log.h.
	r := bufio.NewReaderSize(f, logSize)
	for {
		line, _, err := r.ReadLine()
		str := string(line)
		if err != nil {
			str += " " + err.Error()
		}
		cstr := C.CString(str)
		C.__android_log_write(priority, ctag, cstr)
		C.free(unsafe.Pointer(cstr))
		if err != nil {
			break
		}
	}
}

func init() {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stderr = w
	go lineLog(r, C.ANDROID_LOG_ERROR)

	r, w, err = os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	go lineLog(r, C.ANDROID_LOG_INFO)
}