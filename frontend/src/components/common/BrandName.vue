<template>
  <span :class="brandClasses">{{ brandStore.brandName }}</span>
</template>

<script setup>
import { computed } from 'vue'
import { useBrandStore } from '../../stores/brand'

const brandStore = useBrandStore()

const props = defineProps({
  size: {
    type: String,
    default: 'base',
    validator: (value) => ['xs', 'sm', 'base', 'lg', 'xl', '2xl', '3xl'].includes(value)
  },
  weight: {
    type: String,
    default: 'bold',
    validator: (value) => ['normal', 'medium', 'semibold', 'bold', 'extrabold'].includes(value)
  },
  color: {
    type: String,
    default: 'gray-800'
  }
})

// 动态计算样式类
const brandClasses = computed(() => {
  const sizeClasses = {
    'xs': 'text-xs',
    'sm': 'text-sm', 
    'base': 'text-base',
    'lg': 'text-lg',
    'xl': 'text-xl',
    '2xl': 'text-2xl',
    '3xl': 'text-3xl'
  }
  
  const weightClasses = {
    'normal': 'font-normal',
    'medium': 'font-medium',
    'semibold': 'font-semibold',
    'bold': 'font-bold',
    'extrabold': 'font-extrabold'
  }
  
  return [
    sizeClasses[props.size],
    weightClasses[props.weight],
    `text-${props.color}`
  ].join(' ')
})
</script> 