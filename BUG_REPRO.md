# BUG_REPRO

## Bug 是什么
`OpsService.BatchTransition` 的 `wg.Add(1)` 写在 goroutine 内部、缺少 `wg.Wait()` 就 `close(results)`，channel 生命周期错位；worker 与收集端吞掉失败结果。

## 如何触发
- 并发调用批量关闭；
- 批次中包含不存在的记录 id。

## 真实错误信息
`panic: send on closed channel`；失败记录的结果缺失（outcomes 数量少于输入）。
