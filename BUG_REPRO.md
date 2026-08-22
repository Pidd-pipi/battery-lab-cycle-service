# BUG_REPRO

## Bug 是什么
`buildReportRows` 让所有报告行共享同一个 `IDs` 底层数组，并跨调用复用 `opsReportBuffer` 切片；`OpsService.Report` 重复追加行并截断末行。

## 如何触发
- 生成含多个站点/状态的报告；
- 连续生成两份报告并比对。

## 真实错误信息
报告行 ID 列表互相串扰、旧报告被后一份覆盖、行重复或丢失（行数与分组键不匹配）。
