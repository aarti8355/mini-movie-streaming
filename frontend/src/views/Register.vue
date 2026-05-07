<template>
  <div class="auth-page container">
    <div class="glass-card auth-card">
      <div class="auth-header">
        <h2>Create an Account</h2>
        <p>Join MiniMovie to save and review movies</p>
      </div>
      
      <div v-if="error" class="alert-error">
        {{ error }}
      </div>
      
      <div v-if="success" class="alert-success">
        {{ success }} You can now <router-link to="/login">login here</router-link>.
      </div>

      <form @submit.prevent="handleRegister" v-if="!success">
        <div class="form-group">
          <label for="name">Full Name</label>
          <input 
            type="text" 
            id="name" 
            v-model="name" 
            class="form-control" 
            placeholder="John Doe"
            required
          />
        </div>

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
            placeholder="Create a secure password"
            required
            minlength="6"
          />
        </div>
        
        <button type="submit" class="btn btn-primary w-100" :disabled="loading">
          <span v-if="loading">Creating Account...</span>
          <span v-else>Sign Up</span>
        </button>
      </form>
      
      <div class="auth-footer">
        <p>Already have an account? <router-link to="/login">Sign in</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'

export default {
  name: 'Register',
  data() {
    return {
      name: '',
      email: '',
      password: '',
      loading: false,
      error: null,
      success: null
    }
  },
  methods: {
    async handleRegister() {
      this.loading = true;
      this.error = null;
      this.success = null;
      
      if (this.password.length < 8) {
        this.error = "Password must be at least 8 characters long.";
        this.loading = false;
        return;
      }

      try {
        await api.post('/register', {
          name: this.name,
          email: this.email,
          password: this.password
        });
        
        this.success = "Registration successful!";
        // Clear forms
        this.name = '';
        this.email = '';
        this.password = '';
      } catch (err) {
        this.error = err.response?.data?.error || 'Failed to register. Email might already be in use.';
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

.alert-success {
  background-color: rgba(40, 167, 69, 0.1);
  border-left: 4px solid #28a745;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border-radius: 4px;
  font-size: 0.9rem;
}

.alert-success a {
  color: #28a745;
  font-weight: 600;
  text-decoration: underline;
}

.auth-footer {
  margin-top: 2rem;
  text-align: center;
  font-size: 0.9rem;
  color: var(--text-secondary);
}
</style>
