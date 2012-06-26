package sequence

import (
	"testing"
)

func TestUint16SequenceNew(t *testing.T) {
	var s Uint16Sequence
	if s.Peek() != 0 {
		t.Fail()
	}
}

func TestUint16SequenceNext(t *testing.T) {
	var s Uint16Sequence

	currentId := s.Peek()
	nextId := s.Next()

	if nextId != currentId {
		t.Fail()
	}
}

func TestUint16SequencePeek(t *testing.T) {
	var s Uint16Sequence

	prevId := s.Next()

	if (prevId + 1) != s.Peek() {
		t.Fail()
	}
}

func TestUint32SequenceNew(t *testing.T) {
	var s Uint32Sequence
	if s.Peek() != 0 {
		t.Fail()
	}
}

func TestUint32SequenceNext(t *testing.T) {
	var s Uint32Sequence

	currentId := s.Peek()
	nextId := s.Next()

	if nextId != currentId {
		t.Fail()
	}
}

func TestUint32SequencePeek(t *testing.T) {
	var s Uint32Sequence

	prevId := s.Next()

	if (prevId + 1) != s.Peek() {
		t.Fail()
	}
}

func TestUint64SequenceNew(t *testing.T) {
	var s Uint64Sequence
	if s.Peek() != 0 {
		t.Fail()
	}
}

func TestUint64SequenceNext(t *testing.T) {
	var s Uint64Sequence

	currentId := s.Peek()
	nextId := s.Next()

	if nextId != currentId {
		t.Fail()
	}
}

func TestUint64SequencePeek(t *testing.T) {
	var s Uint64Sequence

	prevId := s.Next()

	if (prevId + 1) != s.Peek() {
		t.Fail()
	}
}
