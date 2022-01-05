package main

import (
	"os"
	"sync"
)
var wg sync.WaitGroup
func main() {

	file,_:=os.OpenFile()
	file.Read()
}