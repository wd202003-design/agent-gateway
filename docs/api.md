# Agent Gateway API 草案（尚未实现）

## 数据面
- `POST /v1/chat/completions`：OpenAI 风格请求，支持 stream 与非 stream。
- `GET /v1/models`：返回当前 SK 有权访问的公开模型名与能力（待确定扩展字段）。
- `POST /v1/embeddings`、`POST /v1/responses`：后续阶段按协议兼容矩阵实现。

Authorization: Bearer <SK>；服务端根据 SK 绑定租户、模型权限、额度和计费策略。

## 控制面（路径待 ADR 确定）
- Providers / Channels：endpoint、凭据引用、区域、限额、健康度。
- Models / Groups：上游模型、公开别名、能力、版本、分组。
- Costs / Prices：上游采购成本与租户销售价格，均保留生效版本。
- Tenants / Keys：SK 生命周期、模型授权、预算与配额。
- Usage / Billing：请求调用记录、token、TTFT、上游成本、销售金额、账单及对账。
- Config Releases：配置快照发布、灰度和回滚。

## 兼容性与记录
记录请求 API 协议版本、目标上游 API 版本（如可获知）、模型标识与版本、渠道、首 token 时间、总时长、usage 来源、错误类型、流式中断状态和 trace_id。

请求重试须遵循是否已向客户端发送输出、幂等性及计费风险，不允许对已开始的流式响应进行无条件自动重试。
