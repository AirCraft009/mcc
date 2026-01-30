package helper

type Queue[V any] struct {
	values []V
}

func NewQueue[V any]() *Queue[V] {
	return &Queue[V]{
		values: make([]V, 0),
	}
}

func (q *Queue[V]) Enqueue(value V) {
	q.values = append(q.values, value)
}

func (q *Queue[V]) Dequeue() V {

	// Return the zero value for the type if empty
	if len(q.values) == 0 {
		var zero V
		return zero
	}

	// Get leading value
	value := q.values[0]

	// Remove the first value from the queue
	q.values = q.values[1:]

	return value
}

func (q *Queue[V]) Peek() V {

	// Return the zero value for the type if empty
	if len(q.values) == 0 {
		var zero V
		return zero
	}

	return q.values[0]
}

func (q *Queue[V]) Length(value V) int {
	return len(q.values)
}
