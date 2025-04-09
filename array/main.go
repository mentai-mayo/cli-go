package array

type Array[T any] struct {
	inner []T
}

// ----- info -----

// array length
func (a *Array[T]) Len() int {
	return len(a.inner)
}

// array capacity
func (a *Array[T]) Cap() int {
	return cap(a.inner)
}

// ----- iterator -----

// get iterator
func (a *Array[T]) Iter() func(yield func(i uint, v T) bool) {
	return func(yield func(i uint, v T) bool) {
		for i := 0; i < len(a.inner); i++ {
			v := a.inner[i]
			if !yield(uint(i), v) {
				return
			}
		}
	}
}

// ----- mutable methods -----
func (a *Array[T]) Push(elem T) {
	a.inner = append(a.inner, elem)
}

func (a *Array[T]) Pop() (*T, bool) {
	if len(a.inner) == 0 {
		return nil, false
	}
	elem := a.inner[len(a.inner)-1]
	a.inner = a.inner[:len(a.inner)-1]
	return &elem, true
}

func (a *Array[T]) Dequeue() (*T, bool) {
	if len(a.inner) == 0 {
		return nil, false
	}
	elem := a.inner[0]
	a.inner = a.inner[1:]
	return &elem, true
}

// ----- convert -----
func (a *Array[T]) Into() []T {
	slice := make([]T, len(a.inner), len(a.inner))
	copy(slice, a.inner)
	return slice
}

func (a *Array[T]) IntoInverse() []T {
	slice := a.Into()
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[len(slice)-1] = slice[len(slice)-1], slice[i]
	}
	return slice
}

// ----- static -----
func New[T any](cap uint) Array[T] {
	inner := make([]T, 0, cap)
	return Array[T]{inner}
}

func FromSlice[T any](slice []T) Array[T] {
	inner := make([]T, len(slice), cap(slice))
	copy(inner, slice)
	return Array[T]{inner}
}
