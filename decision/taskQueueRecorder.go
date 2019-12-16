package decision

import "github.com/SJTU-OpenNetwork/go-peertaskqueue/peertask"


type taskQueueRecorder struct{
	blockNumber int64
	blockSize	int64
	maxNumber	int64
	maxSize		int64
}

func(t *taskQueueRecorder) BlockAdd(taskNumber int, blocksize int){
	//t.blockSize += task.Tasks
	t.blockNumber += int64(taskNumber)
	t.blockSize += int64(blocksize)
}

func(t *taskQueueRecorder) BlockAddNumberOnly(taskNumber int){
	t.blockNumber += int64(taskNumber)
}

func(t *taskQueueRecorder) BlockPop(task *peertask.TaskBlock){
	//task.Tasks
	//blockNumber +=
    if task == nil {
        return
    }

	if t.blockNumber <= 0{
		log.Error("pop when record counts 0 blocks")
		return
	}
	t.blockNumber -= int64(len(task.Tasks))
}

func (t *taskQueueRecorder) RemoveTask(){
	if t.blockNumber <= 0{
		log.Error("remove when record counts 0 blocks")
		return
	}
	t.blockNumber --
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

/*
 * The predicted complete time for the tasks in task queue.
 * It was used to estimate the timestamp of ticket.
 * return: the predicted time in millsecond
 */
func (t *taskQueueRecorder) PredictTime() int64 {
	//return (t.blockNumber * defaultBlockSize)/
	return t.blockNumber * 100
}

func (t *taskQueueRecorder) PredictTimeByNumOfTasks(num int) int64 {
	//return (t.blockNumber * defaultBlockSize)/
	return int64(num) * 100
}
