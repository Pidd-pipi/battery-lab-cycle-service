package main

import "context"

// opsBatchWorker 消费任务队列并执行状态迁移，结果写入 results。
func opsBatchWorker(ctx context.Context, jobs <-chan string, results chan<- BatchOutcome, svc *OpsService, target OpsStatus, actor string) {
	for id := range jobs {
		_, err := svc.Transition(ctx, id, 0, target, actor)
		results <- BatchOutcome{ID: id, OK: err == nil, Err: err}
	}
}
