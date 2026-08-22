# BUG_REPRO

## Bug 是什么
`normalizeOpsRecord` 不初始化 `Labels`/`Revision`，`OpsRecord.Clone` 浅拷贝保留 nil map，零值 `OpsClock` 无兜底；向 nil map 写标签直接 panic，新建记录首次更新报版本冲突。

## 如何触发
- 对未带标签的记录克隆后写标签；
- 新建记录后以 revision=1 做首次更新；
- 使用零值 `OpsClock`。

## 真实错误信息
```
panic: assignment to entry in nil map [recovered]
	panic: assignment to entry in nil map
goroutine 35 [running]:
panic({0x103396620?, 0x1033f8e20?})
```
首次更新返回 `operations revision conflict`。
