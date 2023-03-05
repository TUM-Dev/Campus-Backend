package utils

import "sync"

func RunTasksInRoutines[T any](tasks *[]T, routine func(T), maxRoutines int) {
	taskLen := len(*tasks)
	partSize := taskLen / maxRoutines
	var wg sync.WaitGroup

	for i := 0; i < maxRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 0; j < partSize; j++ {
				routine((*tasks)[i*partSize+j])
			}
		}(i)
	}

	restTasks := taskLen % maxRoutines
	if restTasks > 0 {
		for j := 0; j < restTasks; j++ {
			routine((*tasks)[taskLen-restTasks+j])
		}
	}

	wg.Wait()
}

func RunXTasksInRoutines(count int, routine func(int), maxRoutines int) {
	partSize := count / maxRoutines
	var wg sync.WaitGroup

	for i := 0; i < maxRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 0; j < partSize; j++ {
				routine(i*partSize + j)
			}
		}(i)
	}

	restTasks := count % maxRoutines
	if restTasks > 0 {
		for j := 0; j < restTasks; j++ {
			routine(count - restTasks + j)
		}
	}

	wg.Wait()
}
