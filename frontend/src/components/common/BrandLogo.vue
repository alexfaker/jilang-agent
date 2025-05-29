<template>
  <div :class="containerClasses">
    <CpuChipIcon :class="iconClasses" />
    <BrandName 
      :size="textSize" 
      :weight="textWeight" 
      :color="textColor"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { CpuChipIcon } from '@heroicons/vue/24/outline'
import BrandName from './BrandName.vue'

const props = defineProps({
  size: {
    type: String,
    default: 'base',
    validator: (value) => ['sm', 'base', 'lg', 'xl'].includes(value)
  },
  layout: {
    type: String,
    default: 'horizontal',
    validator: (value) => ['horizontal', 'vertical'].includes(value)
  },
  iconColor: {
    type: String,
    default: 'indigo-600'
  },
  textColor: {
    type: String,
    default: 'gray-800'
  },
  textWeight: {
    type: String,
    default: 'bold'
  },
  spacing: {
    type: String,
    default: 'normal',
    validator: (value) => ['tight', 'normal', 'loose'].includes(value)
  }
})

// 容器样式类
const containerClasses = computed(() => {
  const layoutClasses = {
    'horizontal': 'flex items-center',
    'vertical': 'flex flex-col items-center'
  }
  
  const spacingClasses = {
    'tight': props.layout === 'horizontal' ? 'space-x-1' : 'space-y-1',
    'normal': props.layout === 'horizontal' ? 'space-x-2' : 'space-y-2', 
    'loose': props.layout === 'horizontal' ? 'space-x-3' : 'space-y-3'
  }
  
  return [
    layoutClasses[props.layout],
    spacingClasses[props.spacing]
  ].join(' ')
})

// 图标样式类
const iconClasses = computed(() => {
  const sizeClasses = {
    'sm': 'h-6 w-6',
    'base': 'h-8 w-8',
    'lg': 'h-10 w-10',
    'xl': 'h-12 w-12'
  }
  
  return [
    sizeClasses[props.size],
    `text-${props.iconColor}`
  ].join(' ')
})

// 文本尺寸映射
const textSize = computed(() => {
  const sizeMap = {
    'sm': 'base',
    'base': 'xl',
    'lg': '2xl',
    'xl': '3xl'
  }
  return sizeMap[props.size]
})
</script> 