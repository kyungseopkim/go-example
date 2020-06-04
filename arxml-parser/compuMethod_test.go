package arxml

import (
    "fmt"
    "testing"
)

func TestComputeNum_String(t *testing.T) {
    cn := NewCompuNum(0, 0.25)
    fmt.Println(cn)
}
