<template>
  <header class="header">
    <div class="header__inner">
      <RouterLink to="/" class="brand">
        <span class="brand__dot" />
        GoGo 博客
      </RouterLink>

      <nav class="nav">
        <RouterLink :class="{ active: route.name === 'home' }" to="/">首页</RouterLink>
        <RouterLink
          v-if="auth.isAuthenticated"
          :class="{ active: route.path.startsWith('/dashboard') }"
          to="/dashboard"
        >
          控制台
        </RouterLink>
      </nav>

      <div class="actions">
        <template v-if="auth.isAuthenticated">
          <RouterLink class="user-badge" to="/dashboard">
            <span class="avatar">{{ userInitials }}</span>
            <span class="name">{{ auth.user?.displayName }}</span>
          </RouterLink>
          <button class="btn btn-secondary" type="button" @click="handleLogout">
            退出
          </button>
        </template>
        <template v-else>
          <RouterLink class="btn btn-secondary" to="/login"> 登录 </RouterLink>
          <RouterLink class="btn btn-primary" to="/register"> 注册 </RouterLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const userInitials = computed(() => {
  if (!auth.user?.displayName) return '你'
  const parts = auth.user.displayName.trim().split(/\s+/)
  return parts.map((part) => part[0]?.toUpperCase()).join('').slice(0, 2)
})

const handleLogout = async () => {
  await auth.logout()
  router.push({ name: 'home' })
}
</script>

<style scoped>
.header {
  background-color: var(--color-surface);
  border-bottom: 1px solid rgba(148, 163, 184, 0.2);
  box-shadow: var(--shadow-sm);
}

.header__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
}

.brand {
  font-weight: 700;
  font-size: 1.2rem;
  display: inline-flex;
  align-items: center;
  gap: 0.75rem;
  color: var(--color-primary);
}

.brand__dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-dark));
  box-shadow: 0 0 20px rgba(37, 99, 235, 0.45);
}

.nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-weight: 600;
}

.nav a {
  color: var(--color-text-secondary);
  padding: 0.35rem 0.75rem;
  border-radius: 0.5rem;
  transition: background-color 0.2s ease, color 0.2s ease;
}

.nav a:hover,
.nav a.active {
  background-color: rgba(59, 130, 246, 0.1);
  color: var(--color-primary);
}

.actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.35rem 0.75rem;
  border-radius: 999px;
  background-color: rgba(59, 130, 246, 0.12);
  color: var(--color-primary-dark);
  font-weight: 600;
}

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: var(--color-primary);
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
  font-weight: 700;
}

@media (max-width: 768px) {
  .header__inner {
    padding: 0.9rem 1rem;
    flex-wrap: wrap;
  }

  .nav {
    order: 3;
    width: 100%;
    justify-content: flex-start;
  }

  .actions {
    margin-left: auto;
  }
}
</style>

