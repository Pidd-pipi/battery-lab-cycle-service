# BUG_REPRO

## Bug 是什么
`POST /api/cells/status` handler 丢失 `errors.Is(err, ErrCellNotFound)` 分支、所有错误统一映射 500、坏请求体 500；`ValidateCellStatus` 校验被旁路，任意状态被接受。

## 如何触发
- POST 未知样品 id；
- POST 非法状态或非 JSON 请求体。

## 真实错误信息
未知样品返回 500（应为 404）；非法状态/坏请求体返回 500（应为 400）；`{"status":"failed"}` 被接受并落库。
