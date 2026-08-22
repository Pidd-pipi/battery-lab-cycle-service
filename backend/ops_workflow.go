package main

import (
	"context"
	"sync"
)

// BatchOutcome 记录批量工作流中单个记录的处理结果。
type BatchOutcome struct {
	ID  string
	OK  bool
	Err error
}

// batchWorkerCount 是 BatchTransition 使用的并发 worker 数量。
const batchWorkerCount = 3

// BatchTransition 用固定数量的 worker 并发把一批记录迁移到目标状态。
// 每个 worker 独立从任务队列取记录；所有结果收集完成后返回。
func (s *OpsService) BatchTransition(ctx context.Context, ids []string, target OpsStatus, actor string) ([]BatchOutcome, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	// jobs 缓冲到全部入队，避免生产者被 worker 阻塞；
	// results 同样缓冲到 ids 数量，保证 worker 写结果永不阻塞。
	jobs := make(chan string, len(ids))
	results := make(chan BatchOutcome, len(ids))
	var wg sync.WaitGroup
	// 必须在启动 worker 之前累加 WaitGroup，否则下面的 Wait 可能在
	// worker 还未 Add 时提前返回，导致 close(results) 与 worker 的写操作竞争。
	wg.Add(batchWorkerCount)
	for i := 0; i < batchWorkerCount; i++ {
		go func() {
			defer wg.Done()
			opsBatchWorker(ctx, jobs, results, s, target, actor)
		}()
	}
	go func() {
		defer close(jobs)
		for _, id := range ids {
			// 上下文取消后停止投递，剩余 worker 也会尽快退出。
			select {
			case <-ctx.Done():
				return
			case jobs <- id:
			}
		}
	}()
	// 等 worker 全部退出后再关闭 results，绝不能在 worker 仍在写时关闭，
	// 否则 worker 向已关闭的 channel 发送会 panic。
	go func() {
		wg.Wait()
		close(results)
	}()
	out := make([]BatchOutcome, 0, len(ids))
	for r := range results {
		out = append(out, r)
	}
	return out, nil
}
