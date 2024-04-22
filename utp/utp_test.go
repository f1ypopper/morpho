package utp

import (
	"fmt"
	"testing"
)

func TestUTP(t *testing.T) {
	conn, err := Dial("localhost:1111")
	if err != nil {
		fmt.Printf("CONNECTION ERR: %e\n", err)
	}
	conn.Write([]byte("Hello World\n"))
	buf := make([]byte, 10)
	conn.Read(buf)
}
