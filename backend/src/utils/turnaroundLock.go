package utils

import "sync"

// TurnaroundLocker 按航班 ID 互斥：同一航班的放行/签收/关闭串行化，并发提交只生效一次。
type TurnaroundLocker struct {
	mu    sync.Mutex
	locks map[int64]*sync.Mutex
}

func NewTurnaroundLocker() *TurnaroundLocker {
	return &TurnaroundLocker{locks: make(map[int64]*sync.Mutex)}
}

func (l *TurnaroundLocker) get(id int64) *sync.Mutex {
	l.mu.Lock()
	defer l.mu.Unlock()
	lock, ok := l.locks[id]
	if !ok {
		lock = &sync.Mutex{}
		l.locks[id] = lock
	}
	return lock
}

// Lock 锁定指定航班的联动写操作。
func (l *TurnaroundLocker) Lock(id int64) func() {
	lock := l.get(id)
	lock.Lock()
	return lock.Unlock
}
