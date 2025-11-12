<template>
  <section class="home">
    <header class="hero card">
      <h1>Discover fresh stories</h1>
      <p>
        Read the latest articles, insights, and tutorials from our community. Use the
        filters below to find content that inspires you.
      </p>
      <div class="hero__filters">
        <input
          v-model="filters.search"
          class="search"
          placeholder="Search posts..."
          type="search"
        />

        <select v-model="filters.category">
          <option value="">All categories</option>
          <option v-for="category in categories" :key="category.id" :value="category.slug">
            {{ category.name }}
          </option>
        </select>

        <select v-model="filters.tag">
          <option value="">All tags</option>
          <option v-for="tag in tags" :key="tag.id" :value="tag.slug">
            {{ tag.name }}
          </option>
        </select>
      </div>
    </header>

    <div class="chips" v-if="filters.category || filters.tag || filters.search">
      <span class="muted">Active filters:</span>
      <span v-if="filters.search" class="chip">
        Search: {{ filters.search }}
        <button type="button" @click="filters.search = ''">×</button>
      </span>
      <span v-if="filters.category" class="chip">
        Category: {{ filters.category }}
        <button type="button" @click="filters.category = ''">×</button>
      </span>
      <span v-if="filters.tag" class="chip">
        Tag: {{ filters.tag }}
        <button type="button" @click="filters.tag = ''">×</button>
      </span>
      <button class="btn btn-secondary btn-clear" type="button" @click="resetFilters">
        Clear filters
      </button>
    </div>

    <div class="list">
      <div v-if="isLoading" class="card muted">Loading posts...</div>
      <div v-else-if="error" class="card error">{{ error }}</div>
      <div v-else-if="displayPosts.length === 0" class="card muted">
        No posts found. Try adjusting the filters.
      </div>
      <div v-else class="card-grid">
        <PostCard
          v-for="post in displayPosts"
          :key="post.id"
          :post="post"
          @select="handleTagSelect"
        />
      </div>
    </div>

    <PaginationControls
      :current-page="filters.page"
      :page-size="pageSize"
      :total="total"
      @update:page="handlePageChange"
    />
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PostCard from '@/components/posts/PostCard.vue'
import PaginationControls from '@/components/common/PaginationControls.vue'
import * as postService from '@/services/posts'
import * as categoryService from '@/services/categories'
import * as tagService from '@/services/tags'
import type { Category, Post, Tag } from '@/types'

const router = useRouter()
const route = useRoute()

const pageSize = 6
const posts = ref<Post[]>([])
const total = ref(0)
const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const initialized = ref(false)

const filters = reactive({
  search: (route.query.search as string) ?? '',
  category: (route.query.category as string) ?? '',
  tag: (route.query.tag as string) ?? '',
  page: Number(route.query.page) || 1,
})

const statusLabels: Record<Post['status'], string> = {
  draft: 'Draft',
  published: 'Published',
  archived: 'Archived',
}

const displayPosts = computed(() =>
  posts.value.map((post) => ({
    ...post,
    statusLabel: statusLabels[post.status] ?? post.status,
  })),
)

const updateRouteQuery = () => {
  const query: Record<string, string> = {}
  if (filters.search) query.search = filters.search
  if (filters.category) query.category = filters.category
  if (filters.tag) query.tag = filters.tag
  if (filters.page > 1) query.page = String(filters.page)
  router.replace({ query })
}

const fetchPosts = async () => {
  isLoading.value = true
  error.value = null
  try {
    const response = await postService.fetchPosts({
      page: filters.page,
      pageSize,
      search: filters.search || undefined,
      category: filters.category || undefined,
      tag: filters.tag || undefined,
    })
    posts.value = response.data
    total.value = response.total
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load posts. Please try again later.'
  } finally {
    isLoading.value = false
  }
}

const loadTaxonomies = async () => {
  try {
    const [cats, tgs] = await Promise.all([
      categoryService.fetchCategories(),
      tagService.fetchTags(),
    ])
    categories.value = cats
    tags.value = tgs
  } catch (err) {
    console.error('Failed to load taxonomies', err)
  }
}

const resetFilters = () => {
  filters.search = ''
  filters.category = ''
  filters.tag = ''
}

const handleTagSelect = (slug: string) => {
  filters.tag = filters.tag === slug ? '' : slug
}

const handlePageChange = (page: number) => {
  filters.page = page
}

watch(
  () => [filters.search, filters.category, filters.tag],
  () => {
    if (!initialized.value) return
    if (filters.page !== 1) {
      filters.page = 1
      return
    }
    updateRouteQuery()
    fetchPosts()
  },
)

watch(
  () => filters.page,
  () => {
    if (!initialized.value) return
    updateRouteQuery()
    fetchPosts()
  },
)

onMounted(async () => {
  await loadTaxonomies()
  await fetchPosts()
  initialized.value = true
})
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.hero {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  background: linear-gradient(145deg, rgba(59, 130, 246, 0.08), transparent);
}

.hero h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
}

.hero p {
  margin: 0;
  color: var(--color-text-secondary);
  max-width: 640px;
}

.hero__filters {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  gap: 1rem;
}

.hero__filters .search {
  padding-left: 3rem;
  background-image: url('data:image/svg+xml,%3Csvg width="20" height="20" fill="none" xmlns="http://www.w3.org/2000/svg"%3E%3Ccircle cx="9" cy="9" r="7" stroke="%236B7280" stroke-width="2"/%3E%3Cpath d="m15 15 4 4" stroke="%236B7280" stroke-width="2" stroke-linecap="round"/%3E%3C/svg%3E');
  background-repeat: no-repeat;
  background-position: 1rem center;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.75rem;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.35rem 0.8rem;
  border-radius: 999px;
  background-color: rgba(59, 130, 246, 0.14);
  color: var(--color-primary-dark);
  font-weight: 600;
}

.chip button {
  border: none;
  background: transparent;
  cursor: pointer;
  font-weight: 700;
}

.btn-clear {
  padding: 0.35rem 0.75rem;
}

.card-grid {
  display: grid;
  gap: 1.5rem;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
}

.error {
  color: var(--color-danger);
  font-weight: 600;
}

@media (max-width: 900px) {
  .hero__filters {
    grid-template-columns: 1fr;
  }
}
</style>

