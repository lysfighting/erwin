# Mission Contract: Phase 2 Animation Enhancements
## Created: 2026-03-02

## Success Criteria (must ALL be true for MISSION_COMPLETE):
□ ExecutionSimulator 并行侧上方新增小型 agent 拓扑图（3-4 节点 + 连线）
□ 拓扑图在模拟运行时显示流动光点动画（dispatch = 绿色, return = 蓝色）
□ SwarmTopologySection 的 SVG 连线上添加流动粒子动画
□ Swarm 动画滚动进入视口时自动播放一轮消息流序列
□ 点击 Swarm 任意节点高亮该节点的消息路径
□ 两组件共享一致的颜色编码（dispatch 绿 / return 蓝 / broadcast 橙）
□ 使用 /visual-test 截图验证视觉效果
□ 零 TypeScript 编译错误，dev server 正常运行

## Scope Boundaries:
- IN: ExecutionSimulator 动画增强, SwarmTopologySection 消息流动画, visual-test 验证
- OUT: ExecutionSimulator V2 完整规格, i18n 补全, 其他 section 改动

## Key Constraints:
- 纯 Framer Motion + SVG，不引入新依赖
- 遵循现有动画模式（whileInView, useInView, staggered delays）
- 不破坏现有内容和功能

## Key Decisions Made:
- ExecutionSimulator: 在并行侧添加小型拓扑 + 流动光点，不重写为 V2 规格
- SwarmTopology: 自动播放 + 点击触发模式
- 统一颜色编码: dispatch=#22c55e, return=#3b82f6, broadcast=#f97316
