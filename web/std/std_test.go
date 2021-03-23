package std

import (
	"testing"
)

func TestStrings(t *testing.T) {

	var st = NewT("[\"11\",\"22\"]")
	t.Log(st.Strings())
	t.Log(st.Array().Index(0).Int())
	t.Log(st.Index(0).Int())

}
