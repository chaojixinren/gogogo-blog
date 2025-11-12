<template>
  <section class="posts">
    <div class="card editor">
      <header class="editor__header">
        <div>
          <h2>{{ editingPostId ? 'Edit post' : 'Create new post' }}</h2>
          <p class="muted">
            Write your article and publish when you are ready. Drafts stay private until they are published.
          </p>
        </div>
        <button class="btn btn-secondary" type="button" @click="resetForm" v-if="editingPostId">
          Cancel edit
        </button>
      </header>

      <form class="editor__form" @submit.prevent="submitPost">
        <div class="form-row">
          <div class="form-group">
            <label class="form-label" for="post-title">Title</label>
            <input
              id="post-title"
              v-model="postForm.title"
              placeholder="Post title"
              required
              type="text"
            />
          </div>
          <div class="form-group">
            <label class="form-label" for="post-status">Status</label>
            <select id="post-status" v-model="postForm.status">
              <option value="draft">Draft</option>
              <option value="published">Published</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label" for="post-summary">Summary</label>
          <textarea
            id="post-summary"
            v-model="postForm.summary"
            placeholder="A short teaser for the post"
            rows="3"
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="post-content">Content</label>
          <textarea
            id="post-content"
            v-model="postForm.content"
            placeholder="Write using Markdown. Code blocks, lists and quotes are supported."
            required
            rows="12"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label" for="post-category">Category</label>
            <select id="post-category" v-model.number="postForm.categoryId">
              <option :value="null">No category</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label" for="post-tags">Tags</label>
            <input
              id="post-tags"
              v-model="tagInput"
              placeholder="Comma separated tags"
              type="text"
            />
            <small class="muted">
              Suggestions:
              <button
                v-for="tag in tags"
                :key="tag.id"
                class="tag-suggestion"
                type="button"
                @click="addTag(tag.name)"
              >
                #{{ tag.name }}
              </button>
            </small>
          </div>
        </div>

        <div class="form-actions">
          <button class="btn btn-primary" :disabled="isSubmitting" type="submit">
            <span v-if="isSubmitting">{{ editingPostId ? 'Saving...' : 'Publishing...' }}</span>
            <span v-else>{{ editingPostId ? 'Save changes' : 'Publish post' }}</span>
          </button>
        </div>
        <p v-if="formError" class="form-error">{{ formError }}</p>
        <p v-if="formSuccess" class="success">{{ formSuccess }}</p>
      </form>
    </div>

    <div class="card list">
      <header class="list__header">
        <div>
          <h2>My posts</h2>
          <p class="muted">You have {{ total }} posts.</p>
        </div>
        <button class="btn btn-secondary" type="button" @click="refreshPosts">
          Refresh
        </button>
      </header>

      <div v-if="isLoading" class="muted">Loading posts...</div>
      <div v-else-if="myPosts.length === 0" class="muted">
        You have not published any posts yet.
      </div>
      <table v-else class="post-table">
        <thead>
          <tr>
            <th>Title</th>
            <th>Status</th>
            <th>Updated</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in myPosts" :key="item.id">
            <td>
              <RouterLink :to="`/posts/${item.slug}`">
                {{ item.title }}
              </RouterLink>
              <p class="muted small">{{ item.summary }}</p>
            </td>
            <td>
              <span class="badge">{{ statusLabels[item.status] ?? item.status }}</span>
            </td>
            <td>{{ new Date(item.updatedAt).toLocaleString() }}</td>
            <td class="actions">
              <button class="btn btn-secondary" type="button" @click="startEdit(item)">
                Edit
              </button>
              <button class="btn btn-danger" type="button" @click="deletePost(item.id)">
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <PaginationControls
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="handlePageChange"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import PaginationControls from '@/components/common/PaginationControls.vue'
import type { Category, Post, Tag } from '@/types'
import * as postService from '@/services/posts'
import * as categoryService from '@/services/categories'
import * as tagService from '@/services/tags'

const pageSize = 5
const myPosts = ref<Post[]>([])
const total = ref(0)
const page = ref(1)

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])

const isLoading = ref(false)
const isSubmitting = ref(false)
const formError = ref<string | null>(null)
const formSuccess = ref<string | null>(null)
const editingPostId = ref<number | null>(null)
const tagInput = ref('')

const statusLabels: Record<Post['status'], string> = {
  draft: 'Draft',
  published: 'Published',
  archived: 'Archived',
}

const postForm = reactive({
  title: '',
  summary: '',
  content: '',
  status: 'draft',
  categoryId: null as number | null,
})

const loadTaxonomies = async () => {
  try {
    const [cats, tgs] = await Promise.all([
      categoryService.fetchCategories(),
      tagService.fetchTags(),
    ])
    categories.value = cats
    tags.value = tgs
  } catch (err) {
    console.error('Failed to load categories or tags', err)
  }
}

const refreshPosts = async () => {
  isLoading.value = true
  try {
    const response = await postService.fetchMyPosts({
      page: page.value,
      pageSize,
      includeContent: false,
    })
    myPosts.value = response.data
    total.value = response.total
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const resetForm = () => {
  editingPostId.value = null
  postForm.title = ''
  postForm.summary = ''
  postForm.content = ''
  postForm.status = 'draft'
  postForm.categoryId = null
  tagInput.value = ''
  formError.value = null
  formSuccess.value = null
}

const startEdit = (post: Post) => {
  editingPostId.value = post.id
  postForm.title = post.title
  postForm.summary = post.summary ?? ''
  postForm.content = post.content
  postForm.status = post.status
  postForm.categoryId = post.category?.id ?? null
  tagInput.value = post.tags.map((tag) => tag.name).join(', ')
  formError.value = null
  formSuccess.value = null
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const addTag = (tagName: string) => {
  const tagsSet = new Set(
    tagInput.value
      .split(',')
      .map((tag) => tag.trim())
      .filter(Boolean),
  )
  tagsSet.add(tagName)
  tagInput.value = Array.from(tagsSet).join(', ')
}

const parseTags = () =>
  tagInput.value
    .split(',')
    .map((tag) => tag.trim())
    .filter(Boolean)

const submitPost = async () => {
  formError.value = null
  formSuccess.value = null

  if (!postForm.title || !postForm.content) {
    formError.value = 'Title and content are required.'
    return
  }

  isSubmitting.value = true
  try {
    const payload = {
      title: postForm.title,
      summary: postForm.summary,
      content: postForm.content,
      status: postForm.status,
      categoryId: postForm.categoryId ?? undefined,
      tags: parseTags(),
    }

    if (editingPostId.value) {
      await postService.updatePost(editingPostId.value, payload)
      formSuccess.value = 'Post updated successfully.'
    } else {
      await postService.createPost(payload)
      formSuccess.value = 'Post created successfully.'
    }

    await refreshPosts()
    resetForm()
  } catch (err) {
    console.error(err)
    formError.value = 'Failed to save post.'
  } finally {
    isSubmitting.value = false
  }
}

const deletePost = async (id: number) => {
  if (!confirm('Are you sure you want to delete this post?')) return
  try {
    await postService.deletePost(id)
    await refreshPosts()
  } catch (err) {
    console.error(err)
    alert('Failed to delete post.')
  }
}

onMounted(async () => {
  await loadTaxonomies()
  await refreshPosts()
})

const handlePageChange = async (value: number) => {
  page.value = value
  await refreshPosts()
}
</script>

<style scoped>
.posts {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.editor {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.editor__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.editor__form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.success {
  color: var(--color-success);
  font-weight: 600;
}

.tag-suggestion {
  background: transparent;
  border: none;
  color: var(--color-primary);
  font-weight: 600;
  cursor: pointer;
  margin-right: 0.5rem;
}

.list__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.post-table th,
.post-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid rgba(148, 163, 184, 0.2);
}

.post-table .small {
  font-size: 0.85rem;
}

.post-table .actions {
  display: flex;
  gap: 0.5rem;
}

@media (max-width: 768px) {
  .editor__header {
    flex-direction: column;
    gap: 1rem;
  }

  .post-table {
    display: block;
    overflow-x: auto;
  }
}
</style>
