# Agent Gateway 核心领域模型 v0.1

## Bounded Contexts
1. Catalog：ModelFamily、ModelRevision、ModelCapability、UpstreamModel、PublicModel、ModelGroup、ModelGroupMember。
2. Supply：Provider、Channel、ChannelCredential、ChannelDeployment、ChannelQuota、CostRateCard、CostRateItem。
3. Routing：RoutePolicy、RouteCandidate、RouteConstraint、ChannelHealth、CircuitBreaker、RouteDecision、RouteAttempt。
4. Access：Tenant、Project、APIKey、ModelEntitlement、RateLimitPolicy、BudgetAccount、BudgetReservation。
5. Pricing：PricePlan、PriceRule、PriceAssignment、PriceSnapshot、PromotionPolicy。
6. Metering：Invocation、InvocationAttempt、UsageRecord、CostRecord、ChargeRecord、BillingLedger、Invoice、Reconciliation。

## 核心关系
Provider 1:N Channel；Channel N:M UpstreamModel，通过 ChannelDeployment 表达。
PublicModel 1:N RoutePolicy；RoutePolicy 1:N RouteCandidate；RouteCandidate 指向 ChannelDeployment。
ModelGroup N:M PublicModel；Tenant N:M ModelGroup（经 ModelEntitlement）。
Tenant 1:N Project 1:N APIKey；PriceAssignment 关联 Tenant / ModelGroup / PricePlan。
Invocation 1:N InvocationAttempt；InvocationAttempt 1:N UsageRecord；UsageRecord 经版本化 CostRateCard/PricePlan 生成成本与收费记录；BillingLedger 记录结算和调整。

## 必须保持的不变量
- APIKey 只保存校验值或受保护引用，不保存明文。
- 每个 Invocation 绑定可信 tenant_id、key_id、public_model_id、route_policy_version、price_version。
- 每个 Attempt 绑定 channel_deployment_id、upstream_model_id、protocol_version。
- 使用不可变的计量与价格快照重算历史费用；金额用 decimal/整数最小单位，禁止二进制浮点。
- UsageRecord 有 source（provider/gateway_estimated/reconciled）、status 和去重键。
- 账务调整使用冲正/补记，不直接改写已入账流水。
- 跨模型 fallback 必须满足能力和客户授权，并有明确价格规则。
- 控制面配置经版本化快照下发；数据面热路径不依赖控制面 DB。
- RouteDecision 与 RouteAttempt 分开；响应已开始后禁止无条件重试。

## 调用生命周期
AUTHENTICATE → AUTHORIZE → RESERVE_BUDGET → RESOLVE_PUBLIC_MODEL → FILTER_CAPABILITIES → SELECT_DEPLOYMENT → ADAPT_REQUEST → UPSTREAM_ATTEMPT → STREAM/RESPOND → COLLECT_USAGE → SETTLE → EMIT_LEDGER → RECONCILE。

## 质量指标
TTFT、TPOT、E2E latency、成功率、429、5xx、超时、tokens/s、可用配额、单位成本；按 channel_deployment + api_protocol + region + 时间窗聚合。
首 Token 时间应以首个有效模型输出 token 定义，不能简单用 HTTP 首字节时间替代。

## 待确认
协议族支持范围、是否允许跨模型 fallback、成本币种与汇率、账单计费粒度、用量估算规则、预算超额策略、数据保留期限和审计策略。
