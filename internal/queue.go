package internal

//! not tested

type BuilderQueue struct {
	elems []Job
	seen  map[Job]struct{}
}

func NewBuilderQueue() *BuilderQueue {
	return &BuilderQueue{seen: make(map[Job]struct{})}
}

// to poll?
func (q *BuilderQueue) Len() int {
	return len(q.elems)
}

func (q *BuilderQueue) Enqueue(j Job) {
	if _, exists := q.seen[j]; !exists {
		// enqueue
		q.elems = append(q.elems, j)
		q.seen[j] = struct{}{}
	}
}

func (q *BuilderQueue) Dequeue() Job {

	// return front Job & Dequeue item
	return nil
}

/*
func (q *UniqueQueue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    item := q.items[0]
    q.items = q.items[1:]
    delete(q.seen, item)
    return item, true
}

func main() {
    q := NewUniqueQueue[int]()

    q.Enqueue(1)
    q.Enqueue(2)
    q.Enqueue(1) // wird ignoriert
    q.Enqueue(3)

    for q.Len() > 0 {
        v, _ := q.Dequeue()
        fmt.Println(v)
    }
}
*/
