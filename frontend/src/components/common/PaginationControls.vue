<template>
  <div class="pagination" v-if="totalPages > 1">
    <button
      class="btn btn-secondary"
      :disabled="currentPage === 1"
      type="button"
      @click="goToPage(currentPage - 1)"
    >
      Previous
    </button>
    <div class="pagination__info">
      Page {{ currentPage }} / {{ totalPages }}
    </div>
    <button
      class="btn btn-secondary"
      :disabled="currentPage === totalPages"
      type="button"
      @click="goToPage(currentPage + 1)"
    >
      Next
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPage: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{
  (e: 'update:page', value: number): void
}>()

const totalPages = computed(() =>
  Math.max(1, Math.ceil(props.total / props.pageSize)),
)

const goToPage = (page: number) => {
  const next = Math.min(Math.max(1, page), totalPages.value)
  emit('update:page', next)
}
</script>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-top: 2rem;
}

.pagination__info {
  font-weight: 600;
  color: var(--color-text-secondary);
}
</style>
