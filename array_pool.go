package arraypool

import "fmt"

// ArrayPool is a generic array pool that manages the allocation and deallocation of T.
// The first element of the array is a sentinel and is not allocated.
type ArrayPool[T any] struct {
	// arr stores the actual array elements. arr[0] is a sentinel and won't be allocated.
	arr []T
	// alloc indicates the next position to allocate.
	alloc int
	// free keeps track of the freed indices.
	free map[int]struct{}
}

// NewArrayPool creates a new ArrayPool with the specified capacity.
// If the capacity is less than zero, it panics.
// The actual capacity is increased by one to accommodate the sentinel.
func NewArrayPool[T any](cap int) *ArrayPool[T] {
	if cap < 0 {
		panic("cap is less than zero")
	}

	cap++
	return &ArrayPool[T]{
		arr:   make([]T, cap),
		alloc: 1,
		free:  make(map[int]struct{}),
	}
}

// nextCap calculates the new capacity when the array pool needs to grow.
// If the old capacity is less than 256, it doubles the capacity.
// Otherwise, it increases the capacity by a certain factor.
func (ap *ArrayPool[T]) nextCap(oldCap int) int {
	doubleCap := oldCap + oldCap
	const threshold = 256
	if oldCap < threshold {
		return doubleCap
	}
	return oldCap + ((oldCap + 3*threshold) >> 2)
}

// grow increases the capacity of the array pool.
// It creates a new array with the calculated capacity and copies the old elements.
func (ap *ArrayPool[T]) grow() {
	newCap := ap.nextCap(len(ap.arr))
	newArray := make([]T, newCap)
	copy(newArray, ap.arr)
	ap.arr = newArray
}

// Alloc allocates an index from the array pool.
// It returns an index greater than or equal to 1.
// If there is a free index, it uses that; otherwise, it grows the pool if necessary.
func (ap *ArrayPool[T]) Alloc() int {
	if ap.alloc < len(ap.arr) {
		id := ap.alloc
		ap.alloc++
		return id
	}

	// Check if there are any freed indices
	if len(ap.free) > 0 {
		for k := range ap.free {
			delete(ap.free, k)
			return k
		}
	}

	ap.grow()

	return ap.Alloc()
}

// Free marks an index as free for future allocation.
// If the index is invalid, it panics.
// If the index is the last allocated one, it simply decrements the allocation pointer.
// Otherwise, it adds the index to the free list.
func (ap *ArrayPool[T]) Free(id int) {
	if id <= 0 || id >= ap.alloc {
		panic(fmt.Errorf("free invalid id:%d, next alloc pos:%d", id, ap.alloc))
	}

	// Reset the element to the zero value to prevent memory leaks
	ap.arr[id] = ap.arr[0]

	if id == ap.alloc-1 {
		ap.alloc--
		return
	}

	// ap.free = append(ap.free, id)

	_, ok := ap.free[id]
	if ok {
		return
	}
	ap.free[id] = struct{}{}
}

func (ap *ArrayPool[T]) Get(id int) T {
	return ap.arr[id]
}

// Make sure not to retain the pointer returned by this function for a prolonged time,
// for example, by saving it within a heap-allocated object;
// use it only temporarily on the current stack.
func (ap *ArrayPool[T]) GetRef(id int) *T {
	return &ap.arr[id]
}
