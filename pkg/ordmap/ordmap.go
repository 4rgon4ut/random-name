package ordmap

import (
	"sync"
)

type Node struct {
	key   string
	value string
	next  *Node
}

// OrderedMap combines a map and a linked list to provide O(1) access
// while maintaining insertion order
type OrderedMap struct {
	data map[string]*Node
	head *Node
	tail *Node
	mu   sync.RWMutex // allows concurrent reads but exclusive writes
}

func New() *OrderedMap {
	return &OrderedMap{
		data: make(map[string]*Node),
	}
}

func (om *OrderedMap) Get(key string) (string, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if node, ok := om.data[key]; ok {
		return node.value, true
	}
	return "", false
}

func (om *OrderedMap) Set(key, value string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if node, ok := om.data[key]; ok {
		node.value = value
		return
	}

	// new keys are always added to the end, maintaining insertion order
	newNode := &Node{key: key, value: value}
	om.data[key] = newNode

	if om.tail == nil {
		om.head = newNode
		om.tail = newNode
	} else {
		om.tail.next = newNode
		om.tail = newNode
	}
}

func (om *OrderedMap) Delete(key string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	node, ok := om.data[key]
	if !ok {
		return
	}

	delete(om.data, key)

	if om.head == node {
		om.head = node.next
		if om.head == nil {
			om.tail = nil
		}
	} else {
		current := om.head
		for current != nil && current.next != node {
			current = current.next
		}
		if current != nil {
			current.next = node.next
			if current.next == nil {
				om.tail = current
			}
		}
	}
}

func (om *OrderedMap) GetAll() map[string]string {
	// RLock allows multiple readers to access the map
    om.mu.RLock()
	defer om.mu.RUnlock()

	result := make(map[string]string, len(om.data))
	current := om.head
	for current != nil {
		result[current.key] = current.value
		current = current.next
	}
	return result
}