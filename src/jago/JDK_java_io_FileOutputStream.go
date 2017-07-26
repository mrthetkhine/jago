package jago

import (
	"os"
	"unsafe"
)

func register_java_io_FieOutputStream() {
	register("java/io/FileOutputStream.initIDs()V", Java_java_io_FileOutputStream_initIDs)
	register("java/io/FileOutputStream.writeBytes([BIIZ)V", Java_java_io_FileOutputStream_writeBytes)
}

func Java_java_io_FileOutputStream_initIDs() {
	// TODO
}

func Java_java_io_FileOutputStream_writeBytes(this Reference, byteArr ArrayRef, off Int, len Int, append Boolean) {
	bytes := make([]int8, byteArr.Length())
	for i := 0; i < int(byteArr.Length()); i++ {
		bytes[i] = int8(byteArr.GetElement(Int(i)).(Byte))
	}

	bytes = bytes[off : off+len]
	ptr := unsafe.Pointer(&bytes)

	var file *os.File

	fileDescriptor := this.GetInstanceVariableByName("fd", "Ljava/io/FileDescriptor;").(Reference)
	if !fileDescriptor.IsNull() {
		fd := fileDescriptor.GetInstanceVariableByName("fd", "I").(Int)
		switch fd {
		case 0: file = os.Stdin
		case 1: file = os.Stdout
		case 2: file = os.Stderr
		}
	} else {
		path := this.GetInstanceVariableByName("path", "Ljava/lang/String;").(JavaLangString)
		f, err := os.Open(path.toNativeString())
		if err != nil {
			Fatal("Write file %s failed", path)
		}
		file = f
	}

	if append.IsTrue() {
		file.Chmod(os.ModeAppend)
	}

	file.Write(*((*[]byte)(ptr)))
}