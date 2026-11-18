# LLM Gateway 竞品分析（2026-11-17）

## 项目定位
Agent Gateway 是为 Agent Infra 和外部客户提供的 LLM API 接入、渠道调度、协议兼容、成本定价与租户计费平台，不负责 Agent Runtime 或沙箱调度。

## 对标
| 项目 | 重点能力 | 借鉴方向 |
|---|---|---|
| LiteLLM | 多协议适配、模型别名、路由、fallback、预算与费用追踪 | Adapter / Router / Spend |
| New API | 渠道、分组、Token 权限、模型倍率、运营计费 | 商业化控制面 |
| Bifrost | Go 高性能代理、负载均衡、failover、虚拟 Key | 数据面实现 |
| Agent Router（原 Envoy AI Gateway） | Envoy 数据面、策略与配额控制 | 控制/数据面分离 |
| MLflow AI Gateway | Endpoint、流量拆分、fallback、Tracing | 可观测与模型治理 |
| OpenRouter | 成本、TTFT、吞吐、限流与可靠性联合路由 | 动态评分与分发 |

## 关键差异
- 不能仅用一个 provider 字段表示渠道、模型和部署；路由单元应是 ChannelDeployment。
- 上游成本与租户售价必须独立版本化。
- 单次请求可能有多个上游 attempt；客户收费与供应商成本不能简单一一对应。
- 流式输出开始后不允许无条件 fallback。
- 模型兼容必须按 API 族、版本、tools/vision/reasoning 等能力逐项验证。
- 预算控制应有预留和结算语义，避免并发超支。

## 参考
- https://docs.litellm.ai/docs/routing
- https://docs.litellm.ai/docs/proxy/cost_tracking
- https://github.com/QuantumNous/new-api
- https://github.com/maximhq/bifrost
- https://github.com/theagentrouter/agent-router
- https://mlflow.org/genai/ai-gateway/
- https://github.com/OpenRouterTeam/docs/blob/main/guides/routing/provider-selection.mdx

## 后续验证
以具体 release/tag 固定版本，对协议、流式、异常、重试、计量、成本和授权做自动化探针，以上为文档层面的定性比较而非性能实测。
