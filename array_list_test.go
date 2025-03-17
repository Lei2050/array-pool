package arraypool

import (
	"testing"
)

type TestArrayListStruct struct {
	Val int
}

func TestArrayList(t *testing.T) {
	segmentSize := 4
	sttAl := NewArrayList[TestArrayListStruct](segmentSize)
	v := 1
	sttAl.Add(TestArrayListStruct{Val: v})
	v++
	t.Logf("sttAl: %+v\n", sttAl)

	for range segmentSize {
		sttAl.Add(TestArrayListStruct{Val: v})
		v++
	}
	t.Logf("sttAl: %+v\n", sttAl)

	for range segmentSize * 2 {
		sttAl.Add(TestArrayListStruct{Val: v})
		v++
	}
	t.Logf("sttAl: %+v\n", sttAl)

	for range segmentSize*4 - 2 {
		sttAl.Add(TestArrayListStruct{Val: v})
		v++
	}
	t.Logf("sttAl: %+v\n", sttAl)

	idx := 0
	t.Logf("at %d = %+v\n", idx, sttAl.GetRef(idx))
	idx = 5
	t.Logf("at %d = %+v\n", idx, sttAl.GetRef(idx))
	idx = 16
	t.Logf("at %d = %+v\n", idx, sttAl.GetRef(idx))
	idx = 26
	t.Logf("at %d = %+v\n", idx, sttAl.GetRef(idx))
	idx = 100
	// t.Logf("at %d = %+v\n", idx, sttAl.GetRef(idx)) //panic

	idx = 0
	sttAl.FastRemoveAt(idx)
	t.Logf("sttAl: %+v\n", sttAl)
	idx = 5
	sttAl.FastRemoveAt(idx)
	t.Logf("sttAl: %+v\n", sttAl)
	idx = 16
	sttAl.FastRemoveAt(idx)
	t.Logf("sttAl: %+v\n", sttAl)
	idx = 23
	sttAl.FastRemoveAt(idx)
	t.Logf("sttAl: %+v\n", sttAl)

	sttAl.FastRemoveAt(22)
	sttAl.FastRemoveAt(21)
	sttAl.FastRemoveAt(20)
	t.Logf("sttAl: %+v\n", sttAl)

	sttAl.FastRemoveAt(19)
	sttAl.FastRemoveAt(18)
	sttAl.FastRemoveAt(17)
	sttAl.FastRemoveAt(16)
	t.Logf("sttAl: %+v\n", sttAl)

	for i := range 16 {
		sttAl.FastRemoveAt(15 - i)
	}
	t.Logf("sttAl: %+v\n", sttAl)
}
