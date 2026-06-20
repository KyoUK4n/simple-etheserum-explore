package utils

import "unsafe"

// StringToBytes 字符串转字节
func StringToBytes(s string) []byte {
	b := make([]byte, 0)
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
	return b
}

// BytesToString 字节转字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}