# PriorityQueue

基于数组的对象池，目标是极大减少gc扫描时间。非线程安全。

## Feature & Example usage

```golang
import (
	arraypool "github.com/Lei2050/array-pool"
)

type TestStruct struct {
	Val int
}

func main() {
	pool := arraypool.New[TestStruct](8)
	id := pool.Alloc()
	ptr := pool.Get(id)
	//do something on ptr
	//...
	pool.Free(id)
}
```