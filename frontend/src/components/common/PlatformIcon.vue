<template>
  <!-- Raw mode: multi-color SVG (e.g. MuleRun pixel art) -->
  <svg
    v-if="icon.mode === 'raw'"
    :class="sizeClass"
    viewBox="0 0 24 24"
    v-html="icon.rawSvg"
  />
  <!-- Fill mode: standard SVG paths with fill="currentColor" -->
  <svg
    v-else
    :class="sizeClass"
    :viewBox="icon.viewBox ?? '0 0 24 24'"
    fill="currentColor"
    :fill-rule="icon.fillRule"
  >
    <path v-for="(d, i) in icon.paths" :key="i" :d="d" />
  </svg>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { GroupPlatform } from '@/types'
import { iconData, type IconData } from '@/data/icons'

interface Props {
  platform?: GroupPlatform
  size?: 'xs' | 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'sm'
})

/** Map platform name to icon key in iconData */
const PLATFORM_ICON_KEYS: Record<string, string> = {
  anthropic: 'claude',
  openai: 'openai',
  gemini: 'gemini',
  antigravity: 'antigravity',
  mulerun: 'mulerun',
}

const FALLBACK: IconData = iconData['fallback']

const icon = computed((): IconData => {
  const key = PLATFORM_ICON_KEYS[props.platform ?? '']
  return (key ? iconData[key] : undefined) ?? FALLBACK
})

const sizeClass = computed(() => {
  const sizes = {
    xs: 'w-3 h-3',
    sm: 'w-3.5 h-3.5',
    md: 'w-4 h-4',
    lg: 'w-5 h-5'
  }
  return sizes[props.size] + ' flex-shrink-0'
})
</script>
