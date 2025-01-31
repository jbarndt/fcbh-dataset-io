package search

// https://reintech.io/blog/breadth-first-search-algorithm-in-go

type LanguageNode struct {
	Lang  *Language
	Depth int
}

type LanguageQueue struct {
	items []*LanguageNode
}

func (q *LanguageQueue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *LanguageQueue) Size() int {
	return len(q.items)
}

func (q *LanguageQueue) Enqueue(lang *Language, depth int) {
	q.items = append(q.items, &LanguageNode{Lang: lang, Depth: depth})
}

func (q *LanguageQueue) Dequeue() *LanguageNode {
	if q.IsEmpty() {
		return nil
	} else {
		lang := q.items[0]
		q.items = q.items[1:]
		return lang
	}
}
