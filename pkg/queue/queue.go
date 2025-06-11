package queue

import (
	"media_transcoder/dto"
	"sync"
)

// Static hasChanged status & queueStatus slice
// Reasoning: Implemented to reduce multiple status checks and form a cache-like mechanism using
// slice queueStatus to reduce the number of traversals needed through the queue
// Queue traversal for status check works on O(n) time so its inefficient
// to traverse again and again if nothing has changed
var hasChanged bool
var queueStatus []*QueueNode

type Queue struct {
	front *QueueNode
	rear  *QueueNode
	lock  sync.Mutex
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
	queue.lock.Lock()
	defer queue.lock.Unlock()

	newNode := initQueueNode(filename, data)

	// For underflow condition
	if queue.front == nil {
		queue.front = newNode
		queue.rear = newNode
	}

	// Normal condition
	queue.rear.next = newNode
	queue.rear = newNode

	// Update queue status
	queue.changeStatus()
}

// This returns the data of the file to be processed
func (queue *Queue) Dequeue() *QueueNode {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	returnNode := queue.front
	queue.front = queue.front.next

	// Update queue status
	queue.changeStatus()
	return returnNode
}

// Returns the current queueStatus
func (queue *Queue) Status() []*QueueNode {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	if !queue.getStatus() {
		return queueStatus
	}

	// Resetting the queueStatus for making room for new queueStat
	queueStatus = nil

	// Traverse the queue and get queueStatus
	tmp := queue.front
	for tmp != queue.rear {
		queueStatus = append(queueStatus, tmp)
		tmp = tmp.next
	}
	// Append the queue.rear in queueStatus
	queueStatus = append(queueStatus, tmp)

	return queueStatus
}

// Below private functions do not need to have locks and unlocks
// because their calling function implement it already
// Adding locks to below functions keeps them in a perpetual lock state

// Static variable hasChanged(bool) functions
// hasChanged bool toggler function to set to true on queue status change
func (queue *Queue) changeStatus() {
	hasChanged = true
}

// hasChanged bool checker function
func (queue *Queue) getStatus() bool {
	// Save hasChanged in currentStatus to return present state as its being changed in func
	currentStatus := hasChanged

	// sets hasChanged to false to notify checked for the current instance unless queue is updated
	hasChanged = false
	return currentStatus
}
