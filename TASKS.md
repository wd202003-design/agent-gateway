# TASKS — Agent Gateway

> 每次任务完成后同步实际交付、验证证据、阻塞和后续工作；不得把计划项当作已实现功能。

## 已完成
- [x] T-001：初始化仓库文档（README.md、AGENTS.md、docs/architecture.md、docs/security.md、docs/roadmap.md、docs/api.md）
- [x] T-002：修正定位为 LLM API Gateway，明确与 Agent Infra 的职责边界
- [x] T-003：提交竞品分析 docs/competitive-analysis.md
- [x] T-004：提交领域模型 v0.1 docs/domain-model.md
- [x] T-005：确定内部 Agent Infra 与外部商业租户并重的产品方向
- [x] T-006：建立 PLANS.md / TASKS.md 任务治理约定

## 进行中
- [x] P-001：设计交付完成，待架构评审（不代表实现验收）

## 未开始
- [ ] P-002：工程 MVP（首轮代码已提交，尚未通过运行验证）

## P-001 交付记录（2026-11-17）
- [x] 领域模型 v0.2 与 Mermaid ERD：docs/domain-model.md
- [x] 路由/质量/熔断/故障转移：docs/routing-design.md
- [x] 采购成本、售价、计量、账务与对账：docs/pricing-billing.md
- [x] 控制面/数据面与配置分发：docs/control-data-plane.md
- [x] 协议兼容矩阵与版本跟踪：docs/protocol-compatibility.md
- 验证：GitHub 文件创建/更新操作返回 commit SHA；未执行运行时、性能或协议端到端测试。
- 遗留：严格预算一致性、跨模型 fallback 授权、事件总线/数据库选型、账单规则及协议实测；需后续 ADR 和评审。
- 相关提交：e02682a、b71dc8c、6ac75c9、e6b0a46、9f5ca94。

## P-001A 架构决策固化（2026-11-17）
- [x] 记录 PostgreSQL / ClickHouse / Redis / NATS JetStream 职责及一致性：docs/adr/0001-platform-stack.md
- [x] 记录严格预算并发预留与故障降级：docs/adr/0002-budget-control.md
- [x] 记录双协议 MVP 和跨模型 fallback 默认禁止：docs/adr/0003-api-and-fallback.md
- [x] 同步 docs/roadmap.md、docs/protocol-compatibility.md、docs/control-data-plane.md、docs/pricing-billing.md、docs/routing-design.md
- 验证：GitHub API 返回文件提交 SHA；尚未运行并发、故障、协议兼容或端到端测试。
- 待验证：Redis/PG 跨系统恢复、JetStream 重放、预算严格性、Responses SSE 与上游适配。
- 关键提交：1a611aa、d5cd9dd、05daed6。

## P-002 首轮工程交付（2026-11-17）
- [x] Go 工程入口与 HTTP Server：go.mod、cmd/gateway/main.go
- [x] 双协议开发切片：internal/gateway/gateway.go；静态 SK、模型别名、Chat Completions/Responses 的非流与 SSE 透传
- [x] Mock 上游测试代码：internal/gateway/gateway_test.go
- [x] 本地启动说明：docs/mvp-running.md、.env.example
- [x] PostgreSQL 初始结构：db/migrations/0001_core.sql（未执行）
- [x] JetStream 事件契约：docs/events.md（未接入）
- [ ] 执行 go test ./... 并修复失败
- [ ] 真实上游协议兼容与 SSE 取消/错误测试
- [ ] PostgreSQL/Redis/NATS/ClickHouse 实际连接、预算预留、用量与账务
- [ ] 多渠道路由、熔断、配置分发与商业租户认证
- 验证结果：GitHub API 已返回提交 SHA；尝试通过 git clone 进行本地测试时因执行环境无法解析 github.com 而失败，因此没有 Go 测试通过的证据。
- 风险：当前是单 SK/单渠道、无计费的开发演示代码，不可用于生产或商业租户流量。
- 关键提交：f5ee67e、d1882a8、cf33823、26bd141、5027c6a、b5f4dd7。

## 当前状态
架构与需求设计阶段；无可运行服务、未完成端到端验证。以上“完成”仅表示相应文档或设计决策已落地。

## 更新规则
1. 开工前在 PLANS.md 新增/更新计划，明确验收标准。
2. 完工后在 TASKS.md 记录实际结果、提交、测试与遗留问题。
3. 计划变更须回写 PLANS.md；失败或阻塞也必须如实记录。
