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
	value interface{}
	next  *ListItem
	prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	count int
}

func (l list) Len() int { return l.count }

func (l list) Front() *ListItem { return l.first }

func (l list) Back() *ListItem { return l.last }

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{value: v}
	if l.count == 0 {
		l.first = item
		l.last = item
	} else {
		item.next = l.first
		l.first.prev = item
		l.first = item
	}
	l.count++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{value: v}
	if l.count == 0 {
		l.first = item
		l.last = item
	} else {
		item.prev = l.last
		l.last.next = item
		l.last = item
	}
	l.count++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.prev != nil {
		i.prev.next = i.next
	}
	if i.next != nil {
		i.next.prev = i.prev
	}
	if i == l.first {
		l.first = i.next
	}
	if i == l.last {
		l.last = i.prev
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.first {
		return
	}
	if i.prev != nil {
		i.prev.next = i.next
	}
	if i.next != nil {
		i.next.prev = i.prev
	}
	if i == l.last {
		l.last = i.prev
	}
	i.next = l.first
	i.prev = nil
	if l.first != nil {
		l.first.prev = i
	}
	l.first = i
}

func (l *list) Clear() {
	l.first = nil
	l.last = nil
	l.count = 0
}

func NewList() List {
	return new(list)
}
