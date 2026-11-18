# Agent Gateway 架构设计草案

## 目标—要素—结构—环境
**目标**：统一 LLM 接入，保障可用性与 API 兼容，合理分配上游渠道并实现租户商业化运营。

**要素**：
- Control Plane：Provider、Channel、Upstream Model、Public Model Alias、Model Group、Cost、Price Plan、Tenant、SK、Policy、Billing、API Version Registry。
- Data Plane：Protocol Adapter、Authentication、Quota、Router、Channel Pool、Upstream Adapter、Streaming Proxy、Usage Emitter。
- Feedback：Quality Monitor、Health Scoring、Cost Analysis、Audit、Billing Reconciliation。

**结构**：
1. 客户端使用 SK 和公开模型名调用统一 API。
2. 数据面鉴权、授权、配额检查和协议标准化。
3. 模型别名与模型组策略解析为候选上游模型/渠道。
4. 路由器按能力、健康度、权重、成本、限额和区域选择渠道。
5. 协议适配器调用上游，流式返回并采集 TTFT、TPOT、tokens、finish reason、错误分类。
6. 用量事件异步进入账单与成本核算；监控反馈路由权重。
7. 控制面通过版本化快照发布策略，支持审计与回滚。

## 核心领域区分
Provider：供应商（如某模型厂商）。
Channel：供应商账号、凭据、区域、endpoint、额度及限流等具体可调用入口。
Upstream Model：渠道实际提供的模型标识、能力、版本与采购成本。
Public Model：对外兼容的稳定模型名/别名。
Model Group：对租户可见的模型集合及销售策略。
Price Plan：租户或分组适用的输入、输出、缓存、推理等价格维度与生效版本。
Usage Record：每次请求可核对的 tokens、价格版本、上游成本、收费金额、请求状态与来源。

## 路由与可靠性
- 健康检查 + 被动错误率 + TTFT/TPOT + 429/5xx 分类。
- 断路器、退避、限流、配额与权重调整；策略配置可回滚。
- 在响应尚未开始时可按策略故障转移；流式响应已输出后不得无条件跨渠道重试，以免重复计费或语义破坏。
- 区分模型能力（tools、vision、JSON、reasoning、streaming）及协议差异。
- 不承诺无法保证的精确一次上游执行。

## 与 Agent Infra 的关系
Agent Infra 管理任务/沙箱/Worker；Agent Gateway 为其提供统一模型 API、SK、用量与可用性保障。两者可共享身份和可观测标准，但各自拥有业务数据与生命周期。

## 待决策
协议支持矩阵、成本价格精度、计费规则、配额一致性、存储/事件总线、控制面部署模式、跨区域容灾和 SLA。
