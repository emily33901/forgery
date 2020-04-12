package main

import (
	"github.com/emily33901/forgery-go/forgery"
)

func main() {
	f := forgery.Get()

	f.Run()
}
