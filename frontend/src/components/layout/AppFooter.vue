<template>
  <footer class="footer">
    <div class="footer__inner">
      <p>© {{ year }} GoGo 博客 · 分享创造，记录灵感</p>
      <div class="footer__links">
        <RouterLink to="/">首页</RouterLink>
        <RouterLink to="/dashboard">控制台</RouterLink>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

const year = new Date().getFullYear()
</script>

<style scoped>
.footer {
  background-color: var(--color-surface);
  border-top: 1px solid rgba(148, 163, 184, 0.18);
}

.footer__inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.footer__links {
  display: flex;
  gap: 1rem;
}

.footer__links a {
  color: inherit;
  transition: color 0.2s ease;
}

.footer__links a:hover {
  color: var(--color-primary);
}

@media (max-width: 768px) {
  .footer__inner {
    flex-direction: column;
    gap: 0.75rem;
    text-align: center;
  }
}
</style>

