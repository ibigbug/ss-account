package utils

var FreeList = make(chan []byte, 100)
