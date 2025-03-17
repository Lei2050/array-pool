package arraypool

import (
	"math/rand"
	"testing"
)

type TestPoolStruct struct {
	Val int
}

func TestPool(t *testing.T) {
	sttAp := NewPool[TestArrayPoolStruct](3)
	idx, ptr := sttAp.Alloc()
	ptr.Val = 1
	t.Logf("idx: %d, ptr: %+v\n", idx, ptr)
	idx, ptr = sttAp.Alloc()
	ptr.Val = 2
	t.Logf("idx: %d, ptr: %+v\n", idx, ptr)
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)

	for i := range 32 {
		_, ptr = sttAp.Alloc()
		ptr.Val = i + 3
	}
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)

	sttAp.Free(0)
	sttAp.Free(4)
	sttAp.Free(34)
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)
	sttAp.Free(27)
	sttAp.Free(3)
	sttAp.Free(16)
	sttAp.Free(11)
	sttAp.Free(12)
	sttAp.Free(8)
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)

	t.Logf("====================================\n")
	idx, ptr = sttAp.Alloc()
	ptr.Val = rand.Int() % 10000
	t.Logf("idx: %d, ptr: %+v\n", idx, ptr)
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)

	t.Logf("====================================\n")
	idx, ptr = sttAp.Alloc()
	ptr.Val = rand.Int() % 10000
	t.Logf("idx: %d, ptr: %+v\n", idx, ptr)
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)

	for range 13 {
		_, ptr = sttAp.Alloc()
		ptr.Val = rand.Int() % 10000
	}
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)

	sttAp.Clear()
	t.Logf("sttAl: %+v, %+v\n", sttAp, sttAp.items)
}
