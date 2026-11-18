# Agent Gateway 实施路线图

## M0 需求与架构
- [x] 修正项目定位、领域模型、控制面/数据面边界与 API 草案
- [ ] 竞品能力矩阵：LiteLLM / New API / Bifrost / Agent Router / MLflow AI Gateway
- [ ] ADR：协议矩阵、存储、事件总线、配置分发、计费精度与技术栈

## M1 最小可用 LLM 网关
- [ ] OpenAI Chat Completions（含 SSE）兼容接入
- [ ] Provider/Channel/Model 目录和凭据管理
- [ ] SK 认证、模型别名、静态权重路由、故障转移
- [ ] 请求日志、TTFT、token 用量、错误分类
- [ ] 最小部署、端到端协议兼容测试

## M2 高可用与商业化
- [ ] 多渠道健康评分、熔断、限流、配额和智能路由
- [ ] 上游成本、模型分组与租户价格版本
- [ ] 幂等用量事件、账单流水、成本核对和租户账单
- [ ] SK 权限、额度、预算与调用分析

## M3 协议与运营扩展
- [ ] Responses / Embeddings / 其他协议适配与版本跟踪
- [ ] 多区域部署、灰度配置、容量预测、SLO 告警
- [ ] 对账、账单修正、审计与运营控制台

文档完成不代表服务已实现。
