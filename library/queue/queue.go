package queue

import "sync"

// 支持 O(1) 时间复杂度的队列
// 当元素数量不超过 maxLen 时，插入立即成功
// 当达到 maxLen 时，每次插入都会将最先入队的元素挤出去
type Queue struct {
	q      []string
	m      map[string]bool
	maxLen int
	lock   *sync.Mutex
}

func New(maxLen int) *Queue {
	return &Queue{
		q:      make([]string, 0),
		m:      make(map[string]bool),
		maxLen: maxLen,
		lock:   &sync.Mutex{},
	}
}

func (q *Queue) Insert(val string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.q) < q.maxLen {
		// 插入立即成功
		q.q = append(q.q, val)
		q.m[val] = true
	} else {
		// 先出队队首元素，再插入
		delete(q.m, q.q[0])
		copy(q.q, q.q[1:])
		q.q[q.maxLen-1] = val
		q.m[val] = true
	}
}

func (q *Queue) Search(val string) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.m[val]
}

func (q *Queue) Count() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.q)
}
