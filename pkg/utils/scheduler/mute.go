package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MuteScheduler struct {
	tasks sync.Map // map[key]cancelFunc
}

func NewMuteScheduler() *MuteScheduler {
	return &MuteScheduler{}
}

func (s *MuteScheduler) makeKey(groupId, userId int64) string {
	return fmt.Sprintf("%d:%d", groupId, userId)
}

func (s *MuteScheduler) Register(groupId, userId int64, delay time.Duration, task func()) {
	key := s.makeKey(groupId, userId)

	if val, ok := s.tasks.Load(key); ok {
		val.(context.CancelFunc)()
		s.tasks.Delete(key)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.tasks.Store(key, cancel)

	go func() {
		select {
		case <-time.After(delay):
			task()
			s.tasks.Delete(key)
		case <-ctx.Done():
		}
	}()
}

func (s *MuteScheduler) Remove(groupId, userId int64) {
	key := s.makeKey(groupId, userId)
	if val, ok := s.tasks.Load(key); ok {
		val.(context.CancelFunc)()
		s.tasks.Delete(key)
	}
}
