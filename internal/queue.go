package internal

/*
type UniqueQueue[T comparable] struct {
    items []T
    seen  map[T]struct{}
}

func NewUniqueQueue[T comparable]() *UniqueQueue[T] {
    return &UniqueQueue[T]{seen: make(map[T]struct{})}
}

func (q *UniqueQueue[T]) Enqueue(item T) {
    if _, exists := q.seen[item]; !exists {
        q.items = append(q.items, item)
        q.seen[item] = struct{}{}
    }
}

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

func (q *UniqueQueue[T]) Len() int {
    return len(q.items)
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
