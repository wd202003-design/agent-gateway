# P-001 控制面与数据面边界

## Control Plane
模块：Tenant/Project/SK 管理、Provider/Channel/Model Catalog、Model Group、Cost/Price、Routing Policy、Config Release、Billing Admin、Audit Admin。
持久化业务事实、审核配置变更，生成带 schema_version、revision、checksum、effective_at、signature 的不可变快照。

## Data Plane
模块：Ingress、Protocol Parser、Key Verifier、Entitlement/Quota、Model Resolver、Route Selector、Provider Adapter、SSE Proxy、Usage Emitter、Metrics。
只使用本地已验证配置快照和运行时健康状态处理请求；不得每次调用同步查询控制面数据库。密钥通过安全凭据通道单独分发，不能明文嵌入可下载快照。

## 配置分发
DRAFT → VALIDATED → APPROVED → PUBLISHED → ACKNOWLEDGED / ROLLED_BACK。
数据面校验签名、schema 兼容性和单调版本；原子切换快照，保留 last-known-good；版本不兼容拒绝激活并告警。控制面失联时可按策略继续服务已授权请求，但新密钥/撤销/预算变更需要明确失联时效和安全优先级；敏感权限撤销必须有短 TTL 或强制失效通道。

## 运行事件
数据面产生 invocation/attempt/usage/quality/audit 事件，经持久化缓冲或 outbox 传给计量服务。事件总线故障时限流或停止收费型请求，不能无限丢失用量。监控高基数指标与账务事件分开处理。

## 双场景
INTERNAL：Agent Infra 使用租户级 SK、模型权限、预算和成本中心归集。
EXTERNAL：商业租户使用销售价格、账单、余额与发票。
共同使用认证、协议、路由、用量事件和健康评分；内部租户也不得绕过审计。

## 高可用边界
- 控制面故障不直接中断已有数据面；快照过期或密钥撤销同步失败须按风险策略拒绝。
- 数据面水平扩容，路由健康状态与全局配额需协调；单节点缓存不能代表全局严格预算。
- 配置回滚不改变已绑定历史请求的价格和成本版本。

## 待定技术选型
已决定 PostgreSQL、ClickHouse、Redis、NATS JetStream + PostgreSQL Outbox；配置分发、KMS/Secrets、OpenTelemetry 具体方案仍需评审。参见 docs/adr/0001-platform-stack.md。
