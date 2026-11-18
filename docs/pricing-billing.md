# P-001 成本、定价、计量与账务

## 三层财务模型
1. CostRateCard：渠道部署的上游采购成本（供应商币种）。
2. PricePlan：公开模型/模型组与租户销售价格（合同币种）。
3. BillingLedger：实际结算与资金/额度变动的不可变流水。

## 价格规则
PriceRule 使用 metric + unit + conditions + decimal_price；支持 input_tokens、output_tokens、cache_read_tokens、cache_write_tokens、reasoning_tokens、image_units、audio_seconds、request_count、tool_calls 及 provider_specific。每条用量应标注是否为其他计量项的子集，避免重复计费。条件支持模型、上下文长度、时间、租户等级、地区、有效期和折扣优先级。版本发布后不可改写。

## 计量与快照
请求进入时固定 tenant、public_model、价格版本与汇率版本；选择上游时固定渠道成本版本。UsageRecord 记录来源 PROVIDER / GATEWAY_ESTIMATED / RECONCILED、单位、精度、原始证据引用。供应商未返回 usage 时使用可解释估算并标记 provisional；后续按供应商账单冲正调整。

## 预算与额度
请求前预算预留 estimated_max_charge，完成后按最终用量结算并释放余额；异常断连设置超时补偿。高并发严格预算需要原子预留存储或等价机制；缓存近似计数只能提供软预算。若无法估算上界，必须设置最大输出 tokens 或选择有界预算策略。

## 示例
渠道 A 输入成本 0.5 USD/1M tokens、输出成本 2 USD/1M；渠道 B 成本 0.7/2.5；销售价格输入 1.0、输出 4.0 USD/1M。请求输入 2000 tokens、输出 1000 tokens，A 成本 0.003 USD，销售收费 0.006 USD。以上为示意价格，非市场报价。

## 账务一致性
- Invocation 1:N Attempt，成本按每次实际尝试核算；客户收费按合同和最终可计费用量计算。
- 账务写入使用唯一 (request_id, charge_type, revision) 幂等约束；流水为 append-only。
- Outbox 在同一数据库事务中保存结算结果与待投递事件；消费者至少一次投递并去重。
- 账单按已过账流水汇总，不能直接 sum 原始 API 日志。
- 币种独立记录，换汇时固定汇率版本、时间和舍入模式。
- 修正使用 reversal + adjustment，保留原始凭证与审计链。

## 异常矩阵
| 场景 | 用量 | 客户收费 | 上游成本 |
|---|---|---|---|
| 正常完成 | 上游 usage 优先 | 按快照结算 | 按成功 attempt |
| 首 token 前失败并重试 | 记录每次 attempt | 依据合同通常按最终结果 | 可能有失败成本 |
| 流中断 | 记录 partial/estimated | 按合同及可证实用量 | 可能有实际成本 |
| 上游缺 usage | provisional 估算 | 暂估或延迟结算 | 待对账 |
| 重复事件 | 去重 | 不重复收费 | 不重复入账 |
| 后续供应商账单差异 | reconciliation | 冲正/补记（按合同） | 调整成本 |

## 验收
同请求多 attempt、重试无重复收费、预算并发、部分输出、缺失 usage、价格变更后历史重算、跨币种对账。
