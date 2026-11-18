# Agent Gateway 架构草案

## 目标—要素—结构—环境
**目标**：统一多租户 Agent 的任务提交、计算调度、安全执行和审计。

**要素**：API Gateway、Identity/Policy、Task Store、Scheduler、Worker、Runtime Adapter、Sandbox、Tool Proxy、Artifact Store、Event/Audit。

**结构**：
1. API 认证并提交任务，绑定身份、租户、配额与幂等键。
2. Scheduler 持久化状态、排队、分配 Worker 租约。
3. Worker 在隔离环境中启动可插拔 Runtime（首选验证 Pi）。
4. Tool Proxy 依据策略执行工具授权，收集事件和制品。
5. Worker 终结任务，Scheduler 协调超时、取消与故障恢复。

**环境**：优先支持单机/虚拟机及受限内网，不强制依赖 Kubernetes。容器、gVisor、Kata、microVM、WASM 均为候选执行后端，以探针选型。

## 状态机
CREATED → QUEUED → LEASED → STARTING → RUNNING → SUCCEEDED / FAILED / CANCELLED / TIMED_OUT。

Worker 失联不能视为进程已经终止；租约与 fencing token 防止过期 Worker 提交结果。

## 反馈回路
- 负载：队列延迟与节点资源 → 调整调度并发。
- 安全：异常调用与拒绝记录 → 隔离、阻断、告警。
- 可靠性：心跳、超时、失败率 → 重试、熔断与恢复。
- 成本：模型消耗、沙箱时长 → 限额与路由优化。

## 待定架构决策
任务存储与队列、执行器通信、Runtime 适配接口、网络与凭据隔离、pause/resume 语义。
