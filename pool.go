package arraypool

// Pool is a generic pool for managing objects of type T.
// It maintains a list of free indices and a list of actual objects.
// When an object is removed, it is reset to the default value.
type Pool[T any] struct {
	//// frees stores the indices of the free objects in the pool.
	//frees *ArrayList[int]
	//// items stores the actual objects in the pool.
	//// items is a grow-only ArrayList, with freed indices maintained in frees.
	//// This guarantees that items only appends elements and frees only performs
	//// additions or removals at the tail, leading to high speed.
	//items *ArrayList[T]

	// frees stores the indices of the free objects in the pool.
	// Considering data security, we switched to using a map here.
	// Although we sacrificed a little performance, it makes it easier to
	// prevent data corruption caused by repeated releases when freeing objects.
	frees map[int]struct{}
	// items stores the actual objects in the pool.
	// items is a grow-only ArrayList, with freed indices maintained in frees.
	// This guarantees that items only appends elements and frees only performs
	// additions or removals at the tail, leading to high speed.
	items *ArrayList[T]
}

func NewPool[T any](segmentSize int) *Pool[T] {
	return &Pool[T]{
		//frees: NewArrayList[int](segmentSize),
		frees: make(map[int]struct{}),
		items: NewArrayList[T](segmentSize),
	}
}

// Alloc allocates an object from the pool.
// Make sure not to retain the pointer returned by this function for a prolonged time,
// for example, by saving it within a heap-allocated object;
// use it only temporarily on the current stack.
func (p *Pool[T]) Alloc() (int, *T) {
	//if p.frees.Count() > 0 {
	//	id := p.frees.Get(p.frees.Count() - 1)
	//	p.frees.RemoveLast()
	//	return id, p.items.GetRef(id)
	//}
	if len(p.frees) > 0 {
		for k := range p.frees {
			delete(p.frees, k)
			return k, p.items.GetRef(k)
		}
	}

	return p.items.Count(), p.items.Alloc()
}

// Free releases an object back to the pool.
// It resets the object to the default value and adds its index to the list of free indices.
func (p *Pool[T]) Free(id int) {
	if id < 0 || id >= p.items.Count() {
		return
	}
	if _, ok := p.frees[id]; ok {
		return
	}

	*p.items.GetRef(id) = p.items.defaultValue
	// Add the index to the list of free indices.
	//p.frees.Add(id)
	p.frees[id] = struct{}{}
}

func (p *Pool[T]) Clear() {
	clear(p.frees)
	p.items.Clear()
}

// Count returns the number of allocated objects in the pool.
func (p *Pool[T]) Count() int {
	return p.items.Count() - len(p.frees)
}
