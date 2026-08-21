package ut

import (
	"game/duck/exit"
	"time"
)

// 任务列表 管理器

type Task func()

type TaskMgr struct {
	ch    chan Task
	close chan struct{}
}

func NewTaskMgr() *TaskMgr {
	t := &TaskMgr{
		ch:    make(chan Task, 1024),
		close: make(chan struct{}),
	}

	exit.Close("等待TaskMgr执行完成", t)

	go t.run()

	return t
}

func (t *TaskMgr) Add(task Task) {
	t.ch <- task
}

func (t *TaskMgr) run() {
	for {
		select {
		case one := <-t.ch:
			t.exec(one)
		case <-t.close:
			return
		}
	}
}

func (t *TaskMgr) exec(one Task) {
	one()
}

func (t *TaskMgr) Close() {

	for {
		if len(t.ch) == 0 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	close(t.close)
}
