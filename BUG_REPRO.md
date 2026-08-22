# BUG_REPRO

## Bug 是什么
`OpsAudit.For`/`Since` 用 `a.events[:0]` 复用底层数组，过滤写回覆盖历史事件；`opsClonePage` 返回共享切片；`opsBounds` 缺少越界钳制导致分页切片越界。

## 如何触发
- 查询某记录的审计历史后再查完整历史；
- 用超大页码（如 999）分页查询。

## 真实错误信息
审计历史首条记录被后续查询覆盖；分页查询 panic：`slice bounds out of range`。
