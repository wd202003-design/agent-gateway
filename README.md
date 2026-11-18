# Agent Gateway

面向 Agent Infra、AI 应用和外部租户的多租户 LLM API 接入与运营网关。定位类似 OpenRouter、LiteLLM、New API、Bifrost 等，而非 Agent Runtime / 沙箱任务调度系统。

## 目标
- LLM API 协议兼容、版本追踪和流式透传
- 上游 Provider / Channel / Model 管理及高可用路由
- 成本管理、模型分组与租户销售价格管理
- API Key (SK) 权限、配额、限流、调用记录和租户账单
- TTFT（首 Token 延迟）、TPOT、成功率、错误率、吞吐和成本监控
- 控制面与数据面解耦，支持配置版本发布和快速回滚

## 架构
Control Plane（渠道、模型、价格、租户、策略、账单） → 版本化配置快照 → Data Plane（鉴权、协议适配、路由、上游请求、流式响应） → 异步 Usage / Telemetry / Billing。

Agent Infra 负责 Agent 生命周期、计算调度与沙箱执行；Agent Gateway 只负责 LLM 接入、路由与商业运营。

## 文档
- [架构](docs/architecture.md)
- [安全](docs/security.md)
- [路线图](docs/roadmap.md)
- [API 草案](docs/api.md)
- [Agent 开发约定](AGENTS.md)

当前处于设计基线阶段，尚无可运行服务。
