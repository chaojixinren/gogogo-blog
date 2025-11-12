<template>
  <section class="auth card">
    <header>
      <h1>Create an account</h1>
      <p class="muted">Join the community and start sharing your own articles.</p>
    </header>

    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label class="form-label" for="register-username">Username</label>
        <input
          id="register-username"
          v-model="form.username"
          autocomplete="username"
          placeholder="Pick a unique username"
          required
          type="text"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="register-displayName">Display name</label>
        <input
          id="register-displayName"
          v-model="form.displayName"
          placeholder="How should we call you?"
          type="text"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="register-email">Email</label>
        <input
          id="register-email"
          v-model="form.email"
          autocomplete="email"
          placeholder="Optional email address"
          type="email"
        />
      </div>

      <div class="form-group">
        <label class="form-label" for="register-password">Password</label>
        <input
          id="register-password"
          v-model="form.password"
          autocomplete="new-password"
          placeholder="At least 6 characters"
          required
          minlength="6"
          type="password"
        />
      </div>

      <p v-if="error" class="form-error">{{ error }}</p>

      <button class="btn btn-primary" :disabled="auth.isProcessing" type="submit">
        <span v-if="auth.isProcessing">Creating account...</span>
        <span v-else>Sign up</span>
      </button>
    </form>

    <p class="muted footer">
      Already have an account?
      <RouterLink to="/login">Sign in</RouterLink>
    </p>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const auth = useAuthStore()
const error = ref<string | null>(null)

const form = reactive({
  username: '',
  email: '',
  password: '',
  displayName: '',
})

const handleRegister = async () => {
  error.value = null
  try {
    await auth.register(form)
    router.replace('/dashboard')
  } catch (err) {
    console.error(err)
    error.value = 'Failed to create account. Please try again.'
  }
}
</script>

<style scoped>
.auth {
  max-width: 520px;
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
