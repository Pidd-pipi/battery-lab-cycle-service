# BUG_REPRO

## Bug 是什么
`CellStore.List` 读路径未加 `RLock`，`UpdateStatus` 读改写未持锁，并发状态写入产生 data race 且丢失更新；完成态每次调用都重复累加 `Cycle`，服务层返回值与存储不一致。

## 如何触发
- 多个请求并发点样品的「完成」并刷新列表；
- 重复对已完成样品提交 `complete`。

## 真实错误信息
`go test -race` 报 `WARNING: DATA RACE`（`CellStore.UpdateStatus` 未持锁读取与写操作竞争）；循环次数出现多算/少算。
