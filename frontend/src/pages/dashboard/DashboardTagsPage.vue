<template>
  <section class="taxonomies">
    <div class="card">
      <header class="header">
        <div>
          <h2>{{ editingId ? 'Edit tag' : 'Create tag' }}</h2>
          <p class="muted">
            Tags let readers discover related content across categories. Keep them short and descriptive.
          </p>
        </div>
        <button v-if="editingId" class="btn btn-secondary" type="button" @click="resetForm">
          Cancel edit
        </button>
      </header>

      <form class="form" @submit.prevent="submit">
        <div class="form-group">
          <label class="form-label" for="tag-name">Name</label>
          <input
            id="tag-name"
            v-model="form.name"
            placeholder="Vue, Go, Architecture..."
            required
            type="text"
          />
        </div>
        <div class="form-group">
          <label class="form-label" for="tag-slug">Slug</label>
          <input
            id="tag-slug"
            v-model="form.slug"
            placeholder="vue, go, architecture"
            type="text"
          />
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" :disabled="isSubmitting" type="submit">
            <span v-if="isSubmitting">{{ editingId ? 'Saving...' : 'Creating...' }}</span>
            <span v-else>{{ editingId ? 'Save changes' : 'Create tag' }}</span>
          </button>
        </div>
        <p v-if="error" class="form-error">{{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
      </form>
    </div>

    <div class="card">
      <header class="header">
        <div>
          <h2>All tags</h2>
          <p class="muted">Use consistent naming for a better navigation experience.</p>
        </div>
        <button class="btn btn-secondary" type="button" @click="loadTags">Refresh</button>
      </header>

      <div v-if="isLoading" class="muted">Loading tags...</div>
      <div v-else-if="tags.length === 0" class="muted">No tags yet.</div>
      <table v-else class="table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Slug</th>
            <th>Created</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="tag in tags" :key="tag.id">
            <td>{{ tag.name }}</td>
            <td>{{ tag.slug }}</td>
            <td>{{ new Date(tag.createdAt).toLocaleDateString() }}</td>
            <td class="actions">
              <button class="btn btn-secondary" type="button" @click="startEdit(tag)">
                Edit
              </button>
              <button class="btn btn-danger" type="button" @click="remove(tag.id)">
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { Tag } from '@/types'
import * as tagService from '@/services/tags'

const tags = ref<Tag[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)
const editingId = ref<number | null>(null)

const form = reactive({
  name: '',
  slug: '',
})

const resetForm = () => {
  form.name = ''
  form.slug = ''
  editingId.value = null
  error.value = null
  success.value = null
}

const loadTags = async () => {
  isLoading.value = true
  try {
    tags.value = await tagService.fetchTags()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load tags.'
  } finally {
    isLoading.value = false
  }
}

const startEdit = (tag: Tag) => {
  editingId.value = tag.id
  form.name = tag.name
  form.slug = tag.slug
  error.value = null
  success.value = null
}

const submit = async () => {
  if (!form.name) {
    error.value = 'Name is required.'
    return
  }
  isSubmitting.value = true
  error.value = null
  success.value = null
  try {
    if (editingId.value) {
      await tagService.updateTag(editingId.value, {
        name: form.name,
        slug: form.slug || undefined,
      })
      success.value = 'Tag updated.'
    } else {
      await tagService.createTag({
        name: form.name,
        slug: form.slug || undefined,
      })
      success.value = 'Tag created.'
    }
    await loadTags()
    resetForm()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to save tag.'
  } finally {
    isSubmitting.value = false
  }
}

const remove = async (id: number) => {
  if (!confirm('Delete this tag?')) return
  try {
    await tagService.deleteTag(id)
    await loadTags()
  } catch (err) {
    console.error(err)
    alert('Failed to delete tag.')
  }
}

onMounted(loadTags)
</script>

<style scoped>
.taxonomies {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.success {
  color: var(--color-success);
  font-weight: 600;
}

.table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.table th,
.table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid rgba(148, 163, 184, 0.2);
}

.table .actions {
  display: flex;
  gap: 0.5rem;
}

@media (max-width: 768px) {
  .header {
    flex-direction: column;
  }

  .table {
    display: block;
    overflow-x: auto;
  }
}
</style>
