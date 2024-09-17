package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
	Clear()
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First *ListItem
	Last  *ListItem
	Count int
}

func (l list) Len() int { return l.Count }

func (l list) Front() *ListItem { return l.First }

func (l list) Back() *ListItem { return l.Last }

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.Count == 0 {
		l.First = item
		l.Last = item
	} else {
		item.Next = l.First
		l.First.Prev = item
		l.First = item
	}
	l.Count++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.Count == 0 {
		l.First = item
		l.Last = item
	} else {
		item.Prev = l.Last
		l.Last.Next = item
		l.Last = item
	}
	l.Count++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.First {
		l.First = i.Next
	}
	if i == l.Last {
		l.Last = i.Prev
	}
	l.Count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.First {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.Last {
		l.Last = i.Prev
	}
	i.Next = l.First
	i.Prev = nil
	if l.First != nil {
		l.First.Prev = i
	}
	l.First = i
}

func (l *list) Clear() {
	l.First = nil
	l.Last = nil
	l.Count = 0
}

func NewList() List {
	return new(list)
}
