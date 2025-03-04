package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now().UnixMilli()
	for {
		newTime := time.Now().UnixMilli()

		fmt.Println(newTime - currentTime)

		currentTime = newTime

		time.Sleep(100 * time.Millisecond)
	}
}
