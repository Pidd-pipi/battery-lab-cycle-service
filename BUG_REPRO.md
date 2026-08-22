# BUG_REPRO

## Bug 是什么
`OpsService.Transition`/`Search` 把入参 ctx 替换为 `context.Background()`，取消传播断点；`newEnterpriseServer` 请求级超时全部清零；`requestIDMiddleware` 不递增请求标识。

## 如何触发
- 请求取消后继续调用 `Transition`/`Search`；
- 检查 `ReadHeaderTimeout`/`ReadTimeout`/`WriteTimeout`/`IdleTimeout`；
- 连续发起多个无 `X-Request-ID` 的请求。

## 真实错误信息
取消的 ctx 仍能成功执行状态迁移；服务各超时字段为 0；多个请求的 `X-Request-ID` 均为 `req-0`。
