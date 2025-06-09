package queue

import "media_transcoder/dto"

type Queue struct {
	front *QueueNode
	rear  *QueueNode
}

type QueueNode struct {
	filename string
	data     dto.Format
	next     *QueueNode
}

// QueueNode Functions
func initQueueNode(filename string, data dto.Format) *QueueNode {
	return &QueueNode{filename: filename, data: data, next: nil}
}

// Queue Functions
func (queue *Queue) Enqueue(filename string, data dto.Format) {
	newNode := initQueueNode(filename, data)

	// For underflow condition
	if queue.front == nil {
		queue.front = newNode
		queue.rear = newNode
	}

	// Normal condition
	queue.rear.next = newNode
	queue.rear = newNode
}

// This returns the data of the file to be processed
func (queue *Queue) Dequeue() *QueueNode {
	returnNode := queue.front
	queue.front = queue.front.next

	return returnNode
}

func (queue *Queue) Status() {
}

// HasChanged should be a bool that must be toggled on any task completion?
// Required for increasing efficiency of queue status
// func (queue *Queue) HasChanged() {}
