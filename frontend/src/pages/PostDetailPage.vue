<template>
  <section class="post-detail" v-if="!isLoading && post">
    <article class="card post">
      <header>
        <RouterLink class="muted back" to="/"> <- Back to posts </RouterLink>
        <h1>{{ post.title }}</h1>
        <div class="meta">
          <span>By {{ post.author.displayName }}</span>
          <span>{{ formattedDate }}</span>
          <span class="badge">{{ statusLabel }}</span>
          <RouterLink
            v-if="post.category"
            class="badge"
            :to="{ name: 'home', query: { category: post.category.slug } }"
          >
            {{ post.category.name }}
          </RouterLink>
        </div>
        <div class="tags" v-if="post.tags.length">
          <RouterLink
            v-for="tag in post.tags"
            :key="tag.id"
            class="tag"
            :to="{ name: 'home', query: { tag: tag.slug } }"
          >
            #{{ tag.name }}
          </RouterLink>
        </div>
      </header>

      <div class="content" v-html="compiledContent" />
    </article>

    <aside class="card comments">
      <h2>Comments ({{ comments.length }})</h2>
      <ul v-if="comments.length" class="comment-list">
        <li v-for="comment in comments" :key="comment.id">
          <div class="comment-header">
            <strong>{{ comment.authorName }}</strong>
            <span class="muted">{{ new Date(comment.createdAt).toLocaleString() }}</span>
          </div>
          <p>{{ comment.body }}</p>
        </li>
      </ul>
      <p v-else class="muted">No comments yet. Be the first to share your thoughts.</p>

      <form class="comment-form" @submit.prevent="submitComment">
        <div class="form-group" v-if="!auth.isAuthenticated">
          <label class="form-label" for="authorName">Your name</label>
          <input
            id="authorName"
            v-model="commentForm.authorName"
            placeholder="Enter your display name"
            required
            type="text"
          />
        </div>

        <div class="form-group">
          <label class="form-label" for="commentBody">Leave a comment</label>
          <textarea
            id="commentBody"
            v-model="commentForm.body"
            placeholder="Share your thoughts..."
            required
          />
        </div>

        <div class="form-actions">
          <button class="btn btn-primary" :disabled="isSubmitting" type="submit">
            <span v-if="isSubmitting">Submitting...</span>
            <span v-else>Submit comment</span>
          </button>
        </div>
        <p v-if="commentError" class="form-error">{{ commentError }}</p>
      </form>
    </aside>
  </section>

  <section v-else class="card">
    <p v-if="isLoading">Loading post...</p>
    <p v-else-if="error" class="form-error">{{ error }}</p>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import DOMPurify from 'dompurify'
import { marked } from 'marked'
import { useAuthStore } from '@/store/auth'
import * as postService from '@/services/posts'
import * as commentService from '@/services/comments'
import type { Comment, Post } from '@/types'

const route = useRoute()
const auth = useAuthStore()

const post = ref<Post | null>(null)
const comments = ref<Comment[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)
const commentError = ref<string | null>(null)

const commentForm = reactive({
  authorName: '',
  body: '',
})

const statusLabels: Record<Post['status'], string> = {
  draft: 'Draft',
  published: 'Published',
  archived: 'Archived',
}

const formattedDate = computed(() => {
  if (!post.value) return ''
  const date = post.value.publishedAt ?? post.value.createdAt
  return new Date(date).toLocaleString()
})

const statusLabel = computed(() => {
  if (!post.value) return ''
  return statusLabels[post.value.status] ?? post.value.status
})

const compiledContent = computed(() => {
  if (!post.value) return ''
  const raw = marked.parse(post.value.content)
  return DOMPurify.sanitize(raw)
})

const loadPost = async () => {
  const slug = route.params.slug as string
  if (!slug) return

  isLoading.value = true
  error.value = null
  try {
    const result = await postService.fetchPostBySlug(slug)
    post.value = result
    comments.value = result.comments ?? []
    if (auth.isAuthenticated && auth.user) {
      commentForm.authorName = auth.user.displayName
    }
  } catch (err) {
    console.error(err)
    error.value = 'Post not found.'
  } finally {
    isLoading.value = false
  }
}

const submitComment = async () => {
  if (!post.value) return
  commentError.value = null
  isSubmitting.value = true
  try {
    const payload = {
      body: commentForm.body.trim(),
      authorName: auth.isAuthenticated ? auth.user?.displayName : commentForm.authorName.trim(),
    }
    if (!payload.body) {
      commentError.value = 'Comment body is required.'
      return
    }
    if (!auth.isAuthenticated && !payload.authorName) {
      commentError.value = 'Name is required.'
      return
    }

    const newComment = await commentService.createComment(post.value.id, payload)
    comments.value = [...comments.value, newComment]
    commentForm.body = ''
    if (!auth.isAuthenticated) {
      commentForm.authorName = ''
    }
  } catch (err) {
    console.error(err)
    commentError.value = 'Failed to submit comment.'
  } finally {
    isSubmitting.value = false
  }
}

watch(
  () => route.params.slug,
  () => {
    loadPost()
  },
)

onMounted(loadPost)
</script>

<style scoped>
.post-detail {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
  align-items: start;
}

.post h1 {
  margin-top: 0.5rem;
  font-size: 2.2rem;
}

.meta {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  align-items: center;
  color: var(--color-text-secondary);
}

.tags {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-top: 0.5rem;
}

.tag {
  padding: 0.3rem 0.7rem;
  border-radius: 999px;
  background-color: rgba(59, 130, 246, 0.12);
  color: var(--color-primary-dark);
  font-weight: 600;
}

.content {
  margin-top: 1.5rem;
  line-height: 1.7;
  color: var(--color-text-primary);
}

.content :deep(pre) {
  background-color: #1f2937;
  color: #f9fafb;
  padding: 1rem;
  border-radius: 0.75rem;
  overflow-x: auto;
}

.content :deep(code) {
  font-family: 'Fira Code', 'JetBrains Mono', monospace;
}

.comments h2 {
  margin-top: 0;
}

.comment-list {
  list-style: none;
  padding: 0;
  margin: 1rem 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.35rem;
}

.comment-form {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.back {
  font-size: 0.9rem;
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
}

@media (max-width: 900px) {
  .post-detail {
    grid-template-columns: 1fr;
  }
}
</style>
