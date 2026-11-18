# P-001 路由与服务质量设计

## 两阶段选择
**硬过滤**：租户授权 → 公开模型映射 → API 协议及版本 → tools/vision/JSON/reasoning/context 能力 → 地域/数据驻留 → 渠道启用 → 配额与并发 → 熔断状态 → 合同成本/价格约束。无候选时返回稳定错误码（model_unavailable / capability_unsupported / quota_exceeded）。

**软评分**：对候选 deployment d 计算 score(d) = w1*health(d) + w2*latency(d) + w3*throughput(d) + w4*cost_efficiency(d) + w5*configured_weight(d)。各维度先标准化并限制上下界；权重按路由策略版本固定。数据不足时回退静态权重，不将缺样本视为健康。

## 健康度
主动探针检测可达性/鉴权/最小响应；被动统计按 deployment + protocol + region + 1m/5m/1h 窗口计算请求成功率、429、5xx、timeout、TTFT p50/p95、TPOT、tokens/s、容量与限额。TTFT 定义为请求开始至首个有效模型输出 token，不等同 TCP TTFB。区分客户端取消与上游失败。

## 熔断和恢复
CLOSED → OPEN → HALF_OPEN → CLOSED；按最小样本量和错误率触发，OPEN 后退避探测；HALF_OPEN 限制探测并发。429 优先触发限流退避，避免与 5xx 混同。

## 重试/fallback
- 仅在尚未向客户端输出不可撤销内容且请求语义允许时进行重试。
- SSE 已发出内容后不透明切换上游；应终止流并记录 partial_output。
- 每个 request 设置最大尝试次数、时间预算、指数退避和重试预算，避免重试风暴。
- 同模型跨渠道 fallback 默认允许（满足协议、区域、授权和价格约束）。
- 跨模型 fallback 默认关闭；开启需租户显式授权、能力等价验证及计费规则。
- 上游已执行但响应丢失时可能产生成本；不得宣称 exactly-once。

## 可解释性
RouteDecision 保存候选列表摘要、过滤原因、策略版本、健康快照版本、选中 deployment 与随机种子/权重信息。高基数事件异步落地，热路径仅保留必要字段。

## 反馈稳定性
采用 EWMA + 最小样本量 + hysteresis + 权重调整速率限制，避免渠道在短窗口内振荡。灰度发布评分策略并比较错误率、成本和 p95 TTFT。

## 验收场景
正常权重分流、渠道 429、连续 5xx、断路恢复、无候选、tools 不兼容、SSE 首 token 前失败、SSE 中断、双渠道成本记录、控制面失联。

## 跨模型 Fallback 决策
默认禁止；租户管理员可启用并配置替代模型白名单，项目/SK 只能收紧权限。每次替代需验证协议、能力、区域、成本上限并记录审计。详见 docs/adr/0003-api-and-fallback.md。
