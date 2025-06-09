package queue

import "media_transcoder/dto"

type Queue struct {
	node  QueueNode
	front *Queue
	rear  *Queue
}

type QueueNode struct {
	filename string
	data     dto.Format
}

func (queue *Queue) Enqueue() {}
func (queue *Queue) Dequeue() {}
func (queue *Queue) Status()  {}

// HasChanged should be a bool that must be toggled on any task completion?
// Required for increasing efficiency of queue status
// func (queue *Queue) HasChanged() {}
