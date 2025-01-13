package main

import "syscall/js"

type Widget interface {
	Update(svg js.Value)
}
