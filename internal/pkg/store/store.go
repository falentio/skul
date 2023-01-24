// this package inpsired with svelte/store
// some websocket connection handled with this store
package store

import (
	"reflect"
)

type StoreSubscriber[T any] func(value T)

type Store[T any] struct {
	value       T
	subscribers []StoreSubscriber[T]
}

type RStore[T any] interface {
	Get() T
	Subscribe(s StoreSubscriber[T])
}

type WStore[T any] interface {
	Set(t T)
	Update(func(t T) T)
}

type RWStore[T any] interface {
	RStore[T]
	WStore[T]
}

func (s *Store[T]) Get() T {
	return s.value
}

func (s *Store[T]) Subscribe(subscriber StoreSubscriber[T]) (unsubscribe func()) {
	s.subscribers = append(s.subscribers, subscriber)
	subscriber(s.value)
	return func() {
		for i := range s.subscribers {
			if isSameSubscriber(subscriber, s.subscribers[i]) {
				s.subscribers[i] = s.subscribers[len(s.subscribers)-1]
				s.subscribers = s.subscribers[:len(s.subscribers)-1]
			}
		}
	}
}

func (s *Store[T]) Set(v T) {
	s.value = v
	for _, subs := range s.subscribers {
		subs(v)
	}
}

func (s *Store[T]) Update(u func(prev T) T) {
	v := u(s.value)
	s.Set(v)
}

func isSameSubscriber[T any](a, b StoreSubscriber[T]) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
