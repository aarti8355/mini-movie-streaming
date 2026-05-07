<template>
  <div class="watchlist-page container">
    <h1 class="page-title">My Watchlist</h1>
    
    <div v-if="loading" class="loader">Fetching your list...</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>
    
    <div v-else-if="items.length === 0" class="empty-state glass-card">
      <h2>Your watchlist is empty</h2>
      <p>Start exploring our catalog and add movies you want to watch later!</p>
      <router-link to="/" class="btn btn-primary mt-4">Explore Movies</router-link>
    </div>

    <div v-else class="watchlist-grid">
      <div v-for="item in items" :key="item.id" class="watchlist-item glass-card">
        <router-link :to="`/movie/${item.movie_id}`" class="item-link">
          <div class="img-wrapper">
            <img :src="item.movie_image_url" :alt="item.movie_title" class="item-img" />
          </div>
          <div class="item-content">
            <h3 class="item-title">{{ item.movie_title }}</h3>
            <p class="item-date">Added: {{ new Date(item.created_at).toLocaleDateString() }}</p>
          </div>
        </router-link>
        
        <button @click="removeItem(item.id)" class="btn btn-remove" :disabled="removing === item.id">
          {{ removing === item.id ? '...' : 'Remove' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'

export default {
  name: 'Watchlist',
  data() {
    return {
      items: [],
      loading: true,
      error: null,
      removing: null
    }
  },
  methods: {
    async fetchWatchlist() {
      this.loading = true;
      try {
        const res = await api.get('/watchlist');
        this.items = res.data;
      } catch (err) {
        this.error = "Failed to load watchlist.";
      } finally {
        this.loading = false;
      }
    },
    async removeItem(watchlistId) {
      this.removing = watchlistId;
      try {
        await api.delete(`/watchlist/${watchlistId}`);
        // Remove from local array
        this.items = this.items.filter(i => i.id !== watchlistId);
      } catch (err) {
        console.error("Failed to remove item", err);
      } finally {
        this.removing = null;
      }
    }
  },
  mounted() {
    this.fetchWatchlist();
  }
}
</script>

<style scoped>
.watchlist-page {
  padding: 3rem 1rem;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 3rem;
  border-left: 4px solid var(--primary-color);
  padding-left: 1rem;
}

.empty-state {
  text-align: center;
  padding: 5rem 2rem;
}

.empty-state h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.empty-state p {
  color: var(--text-secondary);
  font-size: 1.1rem;
}

.mt-4 {
  margin-top: 1.5rem;
}

/* Grid */
.watchlist-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 2rem;
}

.watchlist-item {
  padding: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: transform var(--transition);
}

.watchlist-item:hover {
  transform: translateY(-5px);
}

.item-link {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.img-wrapper {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.item-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.watchlist-item:hover .item-img {
  transform: scale(1.05);
}

.item-content {
  padding: 1.5rem;
  flex-grow: 1;
}

.item-title {
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-date {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.btn-remove {
  background-color: transparent;
  color: #fff;
  border-top: 1px solid rgba(255,255,255,0.1);
  padding: 1rem;
  width: 100%;
  text-transform: uppercase;
  font-size: 0.85rem;
  font-weight: bold;
  letter-spacing: 1px;
  border-radius: 0;
}

.btn-remove:hover {
  background-color: rgba(229, 9, 20, 0.2);
  color: var(--primary-color);
}
</style>
