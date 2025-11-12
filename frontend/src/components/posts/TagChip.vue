<template>
  <span class="tag-chip" @click="$emit('select', tag.slug)">
    #{{ tag.name }}
  </span>
</template>

<script setup lang="ts">
import type { Tag } from '@/types'

defineProps<{
  tag: Tag
}>()

defineEmits<{
  (e: 'select', slug: string): void
}>()
</script>

<style scoped>
.tag-chip {
  cursor: pointer;
  background-color: rgba(59, 130, 246, 0.12);
  color: var(--color-primary-dark);
  border-radius: 999px;
  padding: 0.25rem 0.75rem;
  font-size: 0.8rem;
  font-weight: 600;
}

.tag-chip:hover {
  background-color: rgba(59, 130, 246, 0.18);
}
</style>

