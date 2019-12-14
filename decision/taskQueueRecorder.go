package decision

import "github.com/ipfs/go-peertaskqueue/peertask"

type taskQueueRecorder struct{
	blockNumber int64
	blockSize	int64
	maxNumber	int64
	maxSize		int64
}

func(t *taskQueueRecorder) BlocksAdd(){

}

func(t *taskQueueRecorder) BlocksPop(task *peertask.TaskBlock){
	//task.Tasks
}

/*
 * Whether the taskQueue is full.
 * return true if the number of blocks or size of blocks extends limit.
 */
func(t *taskQueueRecorder) IsFull() bool{
	if t.maxNumber > 0 && t.blockNumber >= t.maxNumber {
		return true
	}
	if t.maxSize > 0 && t.blockSize >= t.maxSize {
		return true
	}
	return false
}