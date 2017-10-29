package utils

var FreeList = make(chan [1024]byte, 100)
