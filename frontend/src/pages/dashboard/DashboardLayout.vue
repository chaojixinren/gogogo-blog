<template>
  <div class="dashboard">
    <aside class="sidebar card">
      <h2>Dashboard</h2>
      <p class="muted">Manage your posts, categories and tags.</p>

      <nav>
        <RouterLink :class="{ active: route.name === 'dashboard-posts' }" to="/dashboard/posts">
          Posts
        </RouterLink>
        <RouterLink
          :class="{ active: route.name === 'dashboard-categories' }"
          to="/dashboard/categories"
        >
          Categories
        </RouterLink>
        <RouterLink :class="{ active: route.name === 'dashboard-tags' }" to="/dashboard/tags">
          Tags
        </RouterLink>
      </nav>
    </aside>

    <main class="content">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'

const route = useRoute()
</script>

<style scoped>
.dashboard {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 1.5rem;
  align-items: start;
}

.sidebar {
  position: sticky;
  top: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  height: fit-content;
}

.sidebar nav {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.sidebar a {
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  transition: background-color 0.2s ease, color 0.2s ease;
}

.sidebar a:hover,
.sidebar a.active {
  background-color: rgba(59, 130, 246, 0.15);
  color: var(--color-primary);
}

.content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

@media (max-width: 900px) {
  .dashboard {
    grid-template-columns: 1fr;
  }

  .sidebar {
    position: static;
  }

  .sidebar nav {
    flex-direction: row;
    flex-wrap: wrap;
  }
}
</style>
