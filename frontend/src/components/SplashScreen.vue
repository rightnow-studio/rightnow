<template>
  <Transition name="splash-fade">
    <div v-if="visible" class="splash-screen" :style="{ background: bgColor }">
      <div
        class="bubble-wrap"
        :class="{ 'is-popping': isPopping, 'is-popped': isPopped }"
        :style="wrapStyle"
      >
        <!-- 气泡主环（渐变边框） -->
        <div class="bubble-ring"></div>

        <!-- 内部玻璃高光 -->
        <div class="bubble-glare"></div>

        <!-- 底部珊瑚红光晕 -->
        <div class="bubble-glow"></div>

        <!-- 4个细小碎片，静静漂散 -->
        <div class="frag frag-0"></div>
        <div class="frag frag-1"></div>
        <div class="frag frag-2"></div>
        <div class="frag frag-3"></div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const emit = defineEmits<{ done: [] }>()

const DURATION = Number(import.meta.env.VITE_SPLASH_DURATION ?? 2500)
const bgColor = String(import.meta.env.VITE_SPLASH_BG_COLOR ?? '#1abcb0')
const outlineColor = String(import.meta.env.VITE_BUBBLE_OUTLINE_COLOR ?? '#2dd4bf')

const visible = ref(true)
const isPopping = ref(false)
const isPopped = ref(false)

// 将动态时长和颜色注入为 CSS 变量
const wrapStyle = computed(() => ({
  '--outline-color': outlineColor,
  '--breathe-dur': `${Math.round(DURATION * 0.45)}ms`,
  '--pop-dur': `${Math.round(DURATION * 0.18)}ms`,
}))

const timers: ReturnType<typeof setTimeout>[] = []

onMounted(() => {
  // 60% 处开始戳破
  timers.push(setTimeout(() => { isPopping.value = true }, DURATION * 0.60))
  // 75% 处环消失，碎片飘散
  timers.push(setTimeout(() => { isPopped.value = true }, DURATION * 0.75))
  // 结束：隐藏整个开屏层
  timers.push(setTimeout(() => {
    visible.value = false
    emit('done')
  }, DURATION))
})

onUnmounted(() => timers.forEach(clearTimeout))
</script>

<style scoped>
/* ─── 全屏覆盖层 ─── */
.splash-screen {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ─── Vue Transition：整体淡出 ─── */
.splash-fade-leave-active {
  transition: opacity 0.4s ease;
}
.splash-fade-leave-to {
  opacity: 0;
}

/* ─── 气泡容器 ─── */
.bubble-wrap {
  position: relative;
  width: 120px;
  height: 120px;
  animation: bubbleIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both,
             bubbleBreathe var(--breathe-dur, 1100ms) ease-in-out 0.5s infinite;
}

/* 戳破中：停止呼吸，缓慢膨胀 + 淡出 */
.bubble-wrap.is-popping {
  animation: bubbleIn 0s both, bubblePop var(--pop-dur, 450ms) ease-out forwards;
}

/* 已破：彻底隐藏主环，触发碎片 */
.bubble-wrap.is-popped .bubble-ring,
.bubble-wrap.is-popped .bubble-glare,
.bubble-wrap.is-popped .bubble-glow {
  opacity: 0;
}
.bubble-wrap.is-popped .frag {
  opacity: 0;
  transition: transform 0.6s ease-out, opacity 0.6s ease-out;
}
.bubble-wrap.is-popped .frag-0 { transform: translate(-22px, -18px); }
.bubble-wrap.is-popped .frag-1 { transform: translate(22px, -18px); }
.bubble-wrap.is-popped .frag-2 { transform: translate(22px,  18px); }
.bubble-wrap.is-popped .frag-3 { transform: translate(-22px,  18px); }

/* ─── 气泡主环（渐变边框技巧） ─── */
.bubble-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2.5px solid transparent;
  background:
    linear-gradient(transparent, transparent) padding-box,
    conic-gradient(
      from 210deg,
      rgba(255, 107, 107, 0.95),
      rgba(255, 160, 130, 0.6),
      rgba(255, 255, 255, 0.85),
      rgba(45, 212, 191, 0.5),
      rgba(255, 107, 107, 0.95)
    ) border-box;
  transition: opacity var(--pop-dur, 450ms) ease-out;
}

/* ─── 内部玻璃高光（右上角弧光） ─── */
.bubble-glare {
  position: absolute;
  top: 10px;
  right: 14px;
  width: 30px;
  height: 18px;
  border-radius: 50%;
  background: radial-gradient(
    ellipse at center,
    rgba(255, 255, 255, 0.45) 0%,
    transparent 70%
  );
  transform: rotate(-30deg);
  transition: opacity var(--pop-dur, 450ms) ease-out;
}

/* ─── 底部珊瑚红光晕 ─── */
.bubble-glow {
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  width: 80px;
  height: 40px;
  border-radius: 50%;
  background: radial-gradient(
    ellipse at center,
    rgba(255, 107, 107, 0.35) 0%,
    transparent 70%
  );
  filter: blur(6px);
  transition: opacity var(--pop-dur, 450ms) ease-out;
}

/* ─── 碎片（4个细小点） ─── */
.frag {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 5px;
  height: 5px;
  margin: -2.5px 0 0 -2.5px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.7);
  opacity: 0;
  transition: none;
}
/* 初始散布位置（正常状态下贴近圆环） */
.frag-0 { transform: translate(-38px, -28px); opacity: 0.6; }
.frag-1 { transform: translate( 38px, -28px); opacity: 0.6; }
.frag-2 { transform: translate( 38px,  28px); opacity: 0.6; }
.frag-3 { transform: translate(-38px,  28px); opacity: 0.6; }

/* ─── Keyframes ─── */
@keyframes bubbleIn {
  from {
    transform: scale(0.7);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes bubbleBreathe {
  0%, 100% { transform: scale(1); }
  50%       { transform: scale(1.04); }
}

@keyframes bubblePop {
  0%   { transform: scale(1);    opacity: 1; }
  40%  { transform: scale(1.08); opacity: 0.8; }
  100% { transform: scale(1.18); opacity: 0; }
}
</style>
