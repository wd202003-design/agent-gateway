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
- [ ] P-002：工程技术栈与可运行 MVP

## P-001 交付记录（2026-11-17）
- [x] 领域模型 v0.2 与 Mermaid ERD：docs/domain-model.md
- [x] 路由/质量/熔断/故障转移：docs/routing-design.md
- [x] 采购成本、售价、计量、账务与对账：docs/pricing-billing.md
- [x] 控制面/数据面与配置分发：docs/control-data-plane.md
- [x] 协议兼容矩阵与版本跟踪：docs/protocol-compatibility.md
- 验证：GitHub 文件创建/更新操作返回 commit SHA；未执行运行时、性能或协议端到端测试。
- 遗留：严格预算一致性、跨模型 fallback 授权、事件总线/数据库选型、账单规则及协议实测；需后续 ADR 和评审。
- 相关提交：e02682a、b71dc8c、6ac75c9、e6b0a46、9f5ca94。

## 当前状态
架构与需求设计阶段；无可运行服务、未完成端到端验证。以上“完成”仅表示相应文档或设计决策已落地。

## 更新规则
1. 开工前在 PLANS.md 新增/更新计划，明确验收标准。
2. 完工后在 TASKS.md 记录实际结果、提交、测试与遗留问题。
3. 计划变更须回写 PLANS.md；失败或阻塞也必须如实记录。
