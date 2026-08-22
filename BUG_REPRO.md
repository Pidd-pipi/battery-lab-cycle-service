# BUG_REPRO

## Bug 是什么
状态机转换表缺失 `paused->active` 与 `active->closed` 两条边，`opsStatusValid` 漏掉 `paused`，`opsStatusTerminal` 把 `paused` 当终态；`OpsGuard` 的恢复/关闭/迁移判定与状态机不一致。

## 如何触发
- 对 paused 记录执行恢复（Resume）；
- 对 active 记录执行关闭（Close）；
- 统计/终态判断 paused 记录。

## 真实错误信息
`paused -> active`、`active -> closed` 迁移被 `ErrOpsTransition` 拒绝；`opsStatusTerminal(paused)` 为 true。
