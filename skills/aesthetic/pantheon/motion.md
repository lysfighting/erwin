# Pantheon — Motion & Interaction

> 动效是设计的时间维度。每一帧都要有理由存在。

---

## Saul Bass（1920–1996）— 动效叙事的发明者

**核心思想**：发明了现代电影片头序列作为叙事装置，而非仅仅是演职员表。证明了几秒钟的动态设计可以建立整部电影的情感架构。

**代表作**：*Vertigo* 片头序列（1958）；*Anatomy of a Murder* 海报与片头（1959）；AT&T"地球"标志（1969）

**认知方法 — Bass Essence Extraction（本质提取）**：
> 把整个体验提炼成一个视觉动作。什么一个形状、一个运动、一种色彩关系捕捉了*整体*？如果你的加载动效、英雄区域或引导流程无法被简化为一个核心视觉意图——还没有设计好。找到那个意图，其他的是修饰。

---

## Ollie Johnston & Frank Thomas（迪士尼九大元老）— 12 大动画原则

**核心思想**：在 *The Illusion of Life*（1981）中整理了动画的 12 大原则：挤压拉伸、预备动作、舞台布置、跟进运动、缓入缓出、弧线运动、次要动作、时机掌握、夸张、立体感、吸引力、直线/姿势过渡。所有动态设计的基础。

**代表作**：*The Illusion of Life: Disney Animation*（1981）；《小鹿斑比》、《木偶奇遇记》动画

**认知方法 — 12-Principle Motion Audit（12原则动效审查）**：
> 当一个动效感觉"不对"时，诊断它违反了哪些原则：
>
> **UI 动效的四大高频问题**：
> 1. 缺乏缓入缓出（ease-in/ease-out）——线性运动是机械的
> 2. 没有预备动作（anticipation）——重要动作之前缺少微小预示
> 3. 所有元素相同速度运动（缺乏时机变化）
> 4. 动作停止后没有跟进（follow-through）——感觉被硬切
>
> 修复这四点，80% 的动效问题会消失。

---

## Norman McLaren（1914–1987）— 帧间动效的哲学家

**核心思想**："动画不是运动中的画的艺术，而是被画下来的运动的艺术。"最重要的是帧**之间**发生的事，而不是帧本身。

**代表作**：*Neighbours*（1952）；*Pas de deux*（1968）

**认知方法 — McLaren In-Between Test（中间帧测试）**：
> 过渡的质量由状态**之间**发生的事决定，而非由状态本身决定。相同的起止状态配以不同的插值曲线，创造完全不同的情感体验。专注于旅程，不只是目的地。
>
> **实践**：`ease-in-out` 是基础。`spring` 弹性曲线比贝塞尔曲线更有生命感。Framer Motion 的 `type: "spring"` 比 `duration` 更接近自然物理。

---

## Oskar Fischinger（1900–1967）— 视觉音乐的先驱

**核心思想**：视觉音乐——抽象视觉形式与音乐结构的同步。颜色、形状和运动可以创造一种与音乐平行的纯视觉语言。

**认知方法 — Fischinger Rhythm Audit（节奏审查）**：
> 你的动效有节奏感吗？就像音乐有节拍、小节和乐句一样，动态设计应该有时间结构。同时发生的动效应该共享一个公共时序网格。随机时序 = 视觉噪音。
>
> **实践**：让相关元素共享同一个 `delay` 步长（比如每 50ms 一个元素进场），而不是随机延迟。

---

## UI 动效实践参考

```
时间基准（尊重用户的时间）：
  微交互（hover, focus, active）: 100–150ms
  状态变化（modal open, tab switch）: 200–300ms
  页面过渡: 300–400ms
  重要提示/信息：最多 500ms
  超过 500ms 的 UI 动效 = 对用户时间不敬

缓动曲线选择：
  进入（ease-out）:  cubic-bezier(0, 0, 0.2, 1)   ← 快进慢停，感觉轻盈
  离开（ease-in）:   cubic-bezier(0.4, 0, 1, 1)   ← 慢出快消，感觉决断
  标准（ease）:      cubic-bezier(0.4, 0, 0.2, 1)  ← 综合，通用
  弹性（spring）:    Framer Motion spring config    ← 最有生命感

Tailwind 常用动效组合：
  fade in:    opacity-0 → opacity-100, transition-opacity duration-200
  slide up:   translate-y-2 opacity-0 → translate-y-0 opacity-100
  scale in:   scale-95 opacity-0 → scale-100 opacity-100

什么时候用动效 vs. 不用：
  ✅ 状态变化（loading → done）
  ✅ 空间变化（drawer open, list expand）
  ✅ 注意力引导（where to look next）
  ❌ 纯粹装饰（"这看起来很酷"）
  ❌ 重复/循环动效（用户会疲劳）
  ❌ 背景移动（分散注意力）
```

---

## 动效与情感

| 情感目标 | 动效语言 |
|---|---|
| 快速、效率 | 短时长（100–150ms）、直接移动、无过多缓动 |
| 精致、高端 | 轻微放慢（200–300ms）、优雅缓出、subtle scale |
| 友好、亲切 | 弹性 spring 效果、轻微弹跳 |
| 专业、严肃 | 最少动效、fade 为主、no bounce |
| 创意、活力 | 可以夸张（scale 1.05→1）、交错进场 |
