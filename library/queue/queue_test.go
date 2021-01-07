package queue

import (
	"log"
	"testing"
)

func TestQueue_Insert(t *testing.T) {

	q := New(3)

	q.Insert("a")
	if !check(q, []string{"a"}, map[string]bool{"a": true}) {
		t.Error("error")
	}

	q.Insert("b")
	if !check(q, []string{"a", "b"}, map[string]bool{
		"a": true,
		"b": true,
	}) {
		t.Error("error")
	}

	q.Insert("c")
	if !check(q, []string{"a", "b", "c"}, map[string]bool{
		"a": true,
		"b": true,
		"c": true,
	}) {
		t.Error("error")
	}

	q.Insert("d")
	if !check(q, []string{"b", "c", "d"}, map[string]bool{
		"b": true,
		"c": true,
		"d": true,
	}) {
		t.Error("error")
	}

	q.Insert("e")
	if !check(q, []string{"c", "d", "e"}, map[string]bool{
		"c": true,
		"d": true,
		"e": true,
	}) {
		t.Error("error")
	}

}

func check(q *Queue, rq []string, rm map[string]bool) bool {
	for k, v := range q.q {
		if rq[k] != v {
			log.Printf("q[%d]=%s, rq[%d]=%s", k, v, k, rq[k])
			return false
		}
	}

	for k, v := range q.m {
		if rm[k] != v {
			log.Printf("q[%s]=%v, rm[%s]=%v", k, v, k, rm[k])
			return false
		}
	}

	return true
}
