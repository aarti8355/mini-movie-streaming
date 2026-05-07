<template>
  <div class="auth-page container">
    <div class="glass-card auth-card">
      <div class="auth-header">
        <h2>Welcome Back</h2>
        <p>Sign in to continue your movie journey</p>
      </div>
      
      <div v-if="error" class="alert-error">
        {{ error }}
      </div>

      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">Email</label>
          <input 
            type="email" 
            id="email" 
            v-model="email" 
            class="form-control" 
            placeholder="you@example.com"
            required
          />
        </div>
        
        <div class="form-group">
          <label for="password">Password</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            class="form-control" 
            placeholder="Enter your password"
            required
          />
        </div>
        
        <button type="submit" class="btn btn-primary w-100" :disabled="loading">
          <span v-if="loading">Signing in...</span>
          <span v-else>Sign In</span>
        </button>
      </form>
      
      <div class="auth-footer">
        <p>Don't have an account? <router-link to="/register">Sign up</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'
import { mapActions } from 'pinia'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'Login',
  data() {
    return {
      email: '',
      password: '',
      loading: false,
      error: null
    }
  },
  methods: {
    ...mapActions(useAuthStore, ['setAuth']),
    async handleLogin() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await api.post('/login', {
          email: this.email,
          password: this.password
        });
        
        // Save to Pinia
        this.setAuth(response.data.user, response.data.token);
        
        // Redirect
        this.$router.push('/');
      } catch (err) {
        this.error = err.response?.data?.error || 'Failed to login. Please check your credentials.';
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 100px);
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-header h2 {
  font-size: 1.8rem;
  color: white;
  margin-bottom: 0.5rem;
}

.auth-header p {
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.w-100 {
  width: 100%;
}

.alert-error {
  background-color: rgba(229, 9, 20, 0.1);
  border-left: 4px solid var(--primary-color);
  padding: 1rem;
  margin-bottom: 1.5rem;
  border-radius: 4px;
  font-size: 0.9rem;
}

.auth-footer {
  margin-top: 2rem;
  text-align: center;
  font-size: 0.9rem;
  color: var(--text-secondary);
}
</style>
