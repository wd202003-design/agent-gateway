# API 草案（尚未实现）

Base path: `/v1`。租户身份由服务端认证上下文确定，不信任请求体中的 tenant_id。

## 创建任务
`POST /v1/tasks`，请求头 `Idempotency-Key`。

```json
{"runtime":"pi","input":{"prompt":"检查仓库并运行测试"},"limits":{"timeout_seconds":600}}
```

响应 `202 Accepted`，包含 task_id、status、created_at。

## 查询任务
`GET /v1/tasks/{task_id}`：状态、时间戳、错误摘要及结果引用。

## 取消任务
`POST /v1/tasks/{task_id}/cancel`：幂等取消；接受取消请求不代表底层进程立即停止。

## 事件
`GET /v1/tasks/{task_id}/events`：待决定 SSE 或分页事件、游标和重放语义。

## 约束
任务 ID 不是授权凭证；错误包含稳定 code 与 request_id；同一幂等键对应不同请求体返回冲突；终态和取消竞态需显式定义。
