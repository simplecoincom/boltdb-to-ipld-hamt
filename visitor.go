package boltdbtoipldhamt

import "sync"

type OnVisitFunc = func(target interface{})

type Visitor struct {
	mutex   sync.RWMutex
	visited map[interface{}]bool
}

func NewVisitor() Visitor {
	return Visitor{
		visited: make(map[interface{}]bool),
	}
}

func (v *Visitor) Visit(target interface{}, onVisitFunc ...OnVisitFunc) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.visited[target] = true
	for _, visitorFunc := range onVisitFunc {
		if visitorFunc != nil {
			visitorFunc(target)
		}
	}
}

func (v *Visitor) Visited(target interface{}) bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	_, ok := v.visited[target]
	return !ok
}
