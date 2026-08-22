# Battery Lab Cycle Service

标准库电池循环实验记录服务。`GET /health`、`GET /api/cells` 和 `POST /api/cells/status` 分别用于健康检查、读取样品集合和修改状态。POST 示例：`{"id":"cell-a14","status":"complete"}`；状态支持 `cycling`、`paused`、`complete`。默认端口为 8080，可用 `PORT` 覆盖。

## Verification

验证日期：2026-08-21。`gofmt -w *.go`、`go build ./...`、`go test ./...` 均成功。运行服务后真实调用 health、cells 集合、合法状态 POST 均为 200；非法状态和无效 JSON 为 400；未知样品为 404；`/` 和 `/app.js` 均为 200。验证完成后服务已关闭。

## Engineering Notes

电池循环流程代码按领域模型、校验、状态转换、并发安全存储、审计事件和 HTTP 生命周期分层。请求会保留请求标识并经过恢复与超时保护；状态写入使用版本校验，错误通过可识别的领域错误返回。

除现有接口回归测试外，项目还保留可复用的分页、过滤、策略、工作流和运行健康能力，便于后续扩展而不把业务规则堆积到处理器中。

## Enterprise Layout

```text
.
├── backend/                 # Go module, source code, embedded web assets, Dockerfile
├── database/                # Database extension documentation
├── output/                  # Verification record
├── prompt.txt               # Bugfix task prompt
└── runtime_smoke.json       # Startup contract
```

Health check: `GET /health`. API endpoints: `GET /api/cells` and `POST /api/cells/status`.
