package main

// #cgo CFLAGS: -Wall
// #cgo LDFLAGS: -L${SRCDIR}/lib -lsum
// #include "sum.h"
import "C"

import "fmt"

func main() {
    x, y := 1, 2
    s := C.sum(C.int(x), C.int(y))
    fmt.Printf("sum(%d, %d) = %d\n", x, y, s)
}
