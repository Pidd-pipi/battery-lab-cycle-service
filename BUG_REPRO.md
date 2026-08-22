# BUG_REPRO

## Bug 是什么
`OpsError.Unwrap` 缺失且 `wrapOps` 用 `%v` 断链，`opsCode`/`opsIs*` 无法通过 `errors.Is` 识别领域哨兵错误，分类误判为 internal；`RunWithRetry` 丢失领域错误停判，对 not_found/conflict 也会重试。

## 如何触发
- 任一领域操作失败后调用 `errors.Is` / `opsCode` 分类；
- 用 `RunWithRetry` 重试返回领域错误的函数。

## 真实错误信息
`opsCode(wrapOps(...ErrOpsNotFound))` 返回 `internal`；`errors.Is(err, ErrOpsNotFound)` 为 false；重试对 `ErrOpsConflict` 连续执行 3 次。
