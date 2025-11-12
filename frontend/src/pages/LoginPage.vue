<template>
  <section class="auth card">
    <header>
      <h1>Welcome back</h1>
      <p class="muted">Sign in to manage your posts and share new stories.</p>
    </header>

    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label class="form-label" for="login-username">Username</label>
        <input
          id="login-username"
          v-model="form.username"
          autocomplete="username"
          placeholder="Enter your username"
          required
          type="text"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="login-password">Password</label>
        <input
          id="login-password"
          v-model="form.password"
          autocomplete="current-password"
          placeholder="Enter your password"
          required
          type="password"
        />
      </div>

      <p v-if="error" class="form-error">{{ error }}</p>

      <button class="btn btn-primary" :disabled="auth.isProcessing" type="submit">
        <span v-if="auth.isProcessing">Signing in...</span>
        <span v-else>Sign in</span>
      </button>
    </form>

    <p class="muted footer">
      Don't have an account?
      <RouterLink to="/register">Create one</RouterLink>
    </p>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const form = reactive({
  username: '',
  password: '',
})

const error = ref<string | null>(null)

const handleLogin = async () => {
  error.value = null
  try {
    await auth.login(form)
    const redirect = (route.query.redirect as string) ?? '/dashboard'
    router.replace(redirect)
  } catch (err) {
    console.error(err)
    error.value = 'Invalid username or password.'
  }
}
</script>

<style scoped>
.auth {
  max-width: 480px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 2.5rem;
}

.footer {
  text-align: center;
}

.footer a {
  color: var(--color-primary);
  font-weight: 600;
}
</style>
