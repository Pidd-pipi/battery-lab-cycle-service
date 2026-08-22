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

// BatchTransition 用固定数量的 worker 并发把一批记录迁移到目标状态。
// 每个 worker 独立从任务队列取记录；所有结果收集完成后返回。
func (s *OpsService) BatchTransition(ctx context.Context, ids []string, target OpsStatus, actor string) ([]BatchOutcome, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	jobs := make(chan string)
	results := make(chan BatchOutcome, len(ids))
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opsBatchWorker(ctx, jobs, results, s, target, actor)
		}()
	}
	go func() {
		defer close(jobs)
		for _, id := range ids {
			if err := ctx.Err(); err != nil {
				return
			}
			select {
			case jobs <- id:
			case <-ctx.Done():
				return
			}
		}
	}()
	wg.Wait()
	close(results)
	out := make([]BatchOutcome, 0, len(ids))
	for r := range results {
		out = append(out, r)
	}
	return out, nil
}
