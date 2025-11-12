<template>
  <article class="post-card">
    <RouterLink :to="`/posts/${post.slug}`" class="post-card__title">
      {{ post.title }}
    </RouterLink>

    <p class="post-card__summary" v-if="post.summary">
      {{ post.summary }}
    </p>

    <div class="post-card__meta">
      <span class="badge">{{ post.statusLabel }}</span>
      <span class="muted">
        {{ formattedDate }}
      </span>
      <TagChip
        v-for="tag in post.tags"
        :key="tag.id"
        :tag="tag"
        @select="handleTagSelect"
      />
    </div>

    <RouterLink class="post-card__link" :to="`/posts/${post.slug}`">
      Read more ->
    </RouterLink>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import TagChip from '@/components/posts/TagChip.vue'
import type { Post } from '@/types'

const props = defineProps<{
  post: Post & { statusLabel?: string }
}>()

const emit = defineEmits<{
  (e: 'select', slug: string): void
}>()

const formattedDate = computed(() => {
  const date = props.post.publishedAt ?? props.post.createdAt
  return new Date(date).toLocaleDateString()
})

const handleTagSelect = (slug: string) => {
  emit('select', slug)
}
</script>

<style scoped>
.post-card {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.75rem;
  border-radius: 1rem;
  background-color: var(--color-surface);
  border: 1px solid rgba(148, 163, 184, 0.16);
  box-shadow: var(--shadow-sm);
  transition: transform 0.15s ease, box-shadow 0.2s ease;
}

.post-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.post-card__title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
}

.post-card__title:hover {
  color: var(--color-primary);
}

.post-card__summary {
  margin: 0;
  color: var(--color-text-secondary);
  line-height: 1.6;
}

.post-card__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  align-items: center;
  font-size: 0.9rem;
}

.post-card__link {
  font-weight: 600;
  color: var(--color-primary);
  margin-top: auto;
}
</style>
