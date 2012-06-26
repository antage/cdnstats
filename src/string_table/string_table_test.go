package string_table

import (
	"testing"
)

func TestStringTableNew(t *testing.T) {
	st := New()

	if st.Len() != 0 {
		t.Fail()
	}
}

func TestStringTableStore(t *testing.T) {
	st := New()
	id := st.Store("abcdefgh")
	id2 := st.Store("abcdefgh")

	if id != id2 {
		t.Fail()
	}

	if st.Len() != 1 {
		t.Fail()
	}
}

func TestStringTableLookup(t *testing.T) {
	st := New()

	s := "abcdefgh"
	id := st.Store(s)

	s2, ok := st.Lookup(id)
	if !ok {
		t.Fail()
	}
	if s2 != s {
		t.Fail()
	}
}

func TestStringTableLookupNotExistedString(t *testing.T) {
	st := New()

	if _, ok := st.Lookup(9999); ok {
		t.Fail()
	}
}
