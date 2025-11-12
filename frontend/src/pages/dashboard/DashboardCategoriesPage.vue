<template>
  <section class="taxonomies">
    <div class="card">
      <header class="header">
        <div>
          <h2>{{ editingId ? 'Edit category' : 'Create category' }}</h2>
          <p class="muted">
            Organize posts by assigning them to categories. A clear structure makes it easier for readers to find content.
          </p>
        </div>
        <button type="button" class="btn btn-secondary" v-if="editingId" @click="resetForm">
          Cancel edit
        </button>
      </header>

      <form class="form" @submit.prevent="submit">
        <div class="form-group">
          <label class="form-label" for="category-name">Name</label>
          <input
            id="category-name"
            v-model="form.name"
            placeholder="Frontend, Backend, DevOps..."
            required
            type="text"
          />
        </div>
        <div class="form-group">
          <label class="form-label" for="category-slug">Slug</label>
          <input
            id="category-slug"
            v-model="form.slug"
            placeholder="frontend, backend"
            type="text"
          />
        </div>
        <div class="form-group">
          <label class="form-label" for="category-description">Description</label>
          <textarea
            id="category-description"
            v-model="form.description"
            placeholder="Optional description"
            rows="3"
          />
        </div>
        <div class="form-actions">
          <button class="btn btn-primary" :disabled="isSubmitting" type="submit">
            <span v-if="isSubmitting">{{ editingId ? 'Saving...' : 'Creating...' }}</span>
            <span v-else>{{ editingId ? 'Save changes' : 'Create category' }}</span>
          </button>
        </div>
        <p v-if="error" class="form-error">{{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
      </form>
    </div>

    <div class="card">
      <header class="header">
        <div>
          <h2>All categories</h2>
          <p class="muted">You can edit or delete categories at any time.</p>
        </div>
        <button class="btn btn-secondary" type="button" @click="loadCategories">Refresh</button>
      </header>

      <div v-if="isLoading" class="muted">Loading categories...</div>
      <div v-else-if="categories.length === 0" class="muted">No categories yet.</div>
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
          <tr v-for="category in categories" :key="category.id">
            <td>{{ category.name }}</td>
            <td>{{ category.slug }}</td>
            <td>{{ new Date(category.createdAt).toLocaleDateString() }}</td>
            <td class="actions">
              <button class="btn btn-secondary" type="button" @click="startEdit(category)">
                Edit
              </button>
              <button class="btn btn-danger" type="button" @click="remove(category.id)">
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
import type { Category } from '@/types'
import * as categoryService from '@/services/categories'

const categories = ref<Category[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const error = ref<string | null>(null)
const success = ref<string | null>(null)
const editingId = ref<number | null>(null)

const form = reactive({
  name: '',
  slug: '',
  description: '',
})

const resetForm = () => {
  form.name = ''
  form.slug = ''
  form.description = ''
  editingId.value = null
  error.value = null
  success.value = null
}

const loadCategories = async () => {
  isLoading.value = true
  try {
    categories.value = await categoryService.fetchCategories()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to load categories.'
  } finally {
    isLoading.value = false
  }
}

const startEdit = (category: Category) => {
  editingId.value = category.id
  form.name = category.name
  form.slug = category.slug
  form.description = category.description ?? ''
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
      await categoryService.updateCategory(editingId.value, {
        name: form.name,
        slug: form.slug || undefined,
        description: form.description || undefined,
      })
      success.value = 'Category updated.'
    } else {
      await categoryService.createCategory({
        name: form.name,
        slug: form.slug || undefined,
        description: form.description || undefined,
      })
      success.value = 'Category created.'
    }
    await loadCategories()
    resetForm()
  } catch (err) {
    console.error(err)
    error.value = 'Failed to save category.'
  } finally {
    isSubmitting.value = false
  }
}

const remove = async (id: number) => {
  if (!confirm('Delete this category? Posts will keep their content but lose the category.')) {
    return
  }
  try {
    await categoryService.deleteCategory(id)
    await loadCategories()
  } catch (err) {
    console.error(err)
    alert('Failed to delete category.')
  }
}

onMounted(loadCategories)
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
