# P-001 协议兼容与 API 版本管理

## 初始协议矩阵
| Protocol | Phase | Critical semantics |
|---|---|---|
| OpenAI Chat Completions | M1 | roles, tools, tool_choice, JSON mode, SSE chunks, finish_reason, usage |
| OpenAI Models | M1 | 按 SK 过滤公开模型、能力声明 |
| OpenAI Embeddings | M2 | input batching、dimensions、usage |
| OpenAI Responses | M1 | event types、tool calls、reasoning、response state、stream |
| Anthropic Messages | M3 | content blocks、tool_use、cache、thinking、SSE |
| Gemini GenerateContent | M3 | parts、function calling、safety、stream |

阶段为建议，非功能已实现声明。

## 兼容性层
Client Protocol Parser → Canonical Request（保留语义与 opaque extensions）→ Provider Adapter → Canonical Response/Event → Client Serializer。
不能保证所有协议无损互转；每个适配器声明 supported_features、lossy_transformations、api_versions、streaming_semantics。对于不支持的能力默认返回明确错误，不静默丢字段。

## 版本跟踪
- client_protocol_family/version、gateway_schema_version、adapter_version。
- provider_api_version、upstream_model_name/revision（如供应商可提供）。
- compatibility_test_suite_version、config_revision、route_policy_version。
- 每次 Invocation/Attempt 绑定当时版本，以支持追踪和回归。

## 流式边界
区分 HTTP TTFB、首个 SSE 事件、首个有效输出 token（TTFT）。流事件需要有序、取消传播、超时、错误和 usage 终结标记。开始输出后不得自动跨渠道拼接响应。

## 测试矩阵
非流/流、tools、parallel tool calls、JSON schema、reasoning、vision、长上下文、429、5xx、超时、断连、空 usage、未知字段、供应商版本漂移。每种映射标记 PASS/PARTIAL/UNSUPPORTED，且绑定 adapter 和上游版本。

## API 生命周期
EXPERIMENTAL → STABLE → DEPRECATED → RETIRED；定义变更公告、兼容窗口、灰度、回滚和回归测试。公开模型别名与底层模型 revision 的变化分别审计。

MVP 同时实现 Chat Completions 与 Responses，详见 docs/adr/0003-api-and-fallback.md。
