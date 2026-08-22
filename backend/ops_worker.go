package main

import "context"

// opsBatchWorker 消费任务队列并执行状态迁移，结果写入 results。
// 每个 id 都会产生一条 BatchOutcome（成功或失败都写）；
// 上下文取消时停止取任务，剩余任务不再处理也不再上报，
// 由调用方在 worker 全部退出后再 close(results)。
func opsBatchWorker(ctx context.Context, jobs <-chan string, results chan<- BatchOutcome, svc *OpsService, target OpsStatus, actor string) {
	for {
		// 先看上下文：已取消则立刻退出，避免取消后还在跑。
		select {
		case <-ctx.Done():
			return
		default:
		}
		select {
		case id, ok := <-jobs:
			if !ok {
				return
			}
			_, err := svc.Transition(ctx, id, 0, target, actor)
			results <- BatchOutcome{ID: id, OK: err == nil, Err: err}
		case <-ctx.Done():
			return
		}
	}
}
