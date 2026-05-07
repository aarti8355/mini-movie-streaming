<template>
  <div id="app">
    <nav class="navbar">
      <div class="container nav-content">
        <router-link to="/" class="brand">
          <span class="brand-text">MiniMovie</span>
        </router-link>
        
        <div class="nav-links">
          <router-link to="/">Home</router-link>
          
          <template v-if="isAuthenticated">
            <router-link v-if="isAdmin" to="/admin" class="nav-admin">Dashboard (Admin)</router-link>
            <router-link to="/watchlist">My Watchlist</router-link>
            <div class="user-menu">
              <span class="user-greeting">Hi, {{ user?.name }}</span>
              <button @click="handleLogout" class="btn btn-secondary btn-sm">Logout</button>
            </div>
          </template>
          
          <template v-else>
            <router-link to="/login" class="nav-btn">Login</router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">Sign Up</router-link>
          </template>
        </div>
      </div>
    </nav>
    
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script>
import { mapState, mapActions } from 'pinia'
import { useAuthStore } from './stores/auth'

export default {
  name: 'App',
  computed: {
    ...mapState(useAuthStore, ['isAuthenticated', 'user', 'isAdmin'])
  },
  methods: {
    ...mapActions(useAuthStore, ['logout']),
    handleLogout() {
      this.logout()
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.navbar {
  background-color: var(--nav-bg);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 100;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 1rem 0;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  display: flex;
  align-items: center;
}

.brand-text {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--primary-color);
  letter-spacing: 1px;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.nav-links a {
  font-weight: 500;
  font-size: 0.95rem;
}

.nav-links a.router-link-active:not(.brand) {
  color: var(--primary-color);
}

.nav-admin {
  background-color: var(--bg-color-light);
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  border: 1px solid rgba(255,255,255,0.1);
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-left: 1rem;
  padding-left: 1rem;
  border-left: 1px solid rgba(255, 255, 255, 0.2);
}

.user-greeting {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.btn-sm {
  padding: 0.4rem 0.8rem;
  font-size: 0.85rem;
}

.main-content {
  min-height: calc(100vh - 70px);
}
</style>
