<template>
  <div class="home">
    <!-- Hero Banner -->
    <section class="hero" :style="{ backgroundImage: 'linear-gradient(to right, rgba(20,20,20,0.9) 0%, rgba(20,20,20,0.4) 100%), url(' + heroMovie.image_url + ')' }" v-if="heroMovie">
      <div class="container hero-content">
        <span class="badge">Featured</span>
        <h1 class="hero-title">{{ heroMovie.title }}</h1>
        <p class="hero-desc">{{ heroMovie.description }}</p>
        <div class="hero-meta">
          <span class="meta-item">{{ heroMovie.release_year }}</span>
          <span class="meta-item badge-outline">{{ heroMovie.genre }}</span>
          <span class="meta-item rating">⭐ {{ heroMovie.average_rating.toFixed(1) }}</span>
        </div>
        <router-link :to="`/movie/${heroMovie.id}`" class="btn btn-primary mt-4">
          View Details
        </router-link>
      </div>
    </section>

    <!-- Movie Carousel / Grid -->
    <section class="catalog container">
      <h2 class="section-title">Explore Catalog</h2>
      
      <div class="controls-bar glass-card">
        <input type="text" v-model="search" placeholder="Search by title..." @input="debouncedFetch" class="form-control search-input" />
        <select v-model="genre" @change="fetchMovies" class="form-control">
          <option value="">All Genres</option>
          <option value="Action">Action</option>
          <option value="Comedy">Comedy</option>
          <option value="Drama">Drama</option>
          <option value="Sci-Fi">Sci-Fi</option>
          <option value="Horror">Horror</option>
        </select>
        <select v-model="rating" @change="fetchMovies" class="form-control">
          <option value="0">Any Rating</option>
          <option value="3">3+ Stars</option>
          <option value="4">4+ Stars</option>
        </select>
        <select v-model="sort" @change="fetchMovies" class="form-control">
          <option value="">Newest</option>
          <option value="rating_desc">Highest Rated</option>
          <option value="rating_asc">Lowest Rated</option>
        </select>
      </div>
      
      <div v-if="loading" class="loader">Loading movies...</div>
      <div v-else-if="error" class="error-msg">{{ error }}</div>
      
      <div v-else class="movie-grid">
        <router-link 
          v-for="movie in displayMovies" 
          :key="movie.id" 
          :to="`/movie/${movie.id}`"
          class="movie-card"
        >
          <div class="card-img-wrapper">
            <img :src="movie.image_url" :alt="movie.title" class="card-img" />
            <div class="card-overlay">
              <span class="overlay-rating">⭐ {{ movie.average_rating.toFixed(1) }}</span>
            </div>
          </div>
          <div class="card-body">
            <h3 class="card-title">{{ movie.title }}</h3>
            <div class="card-footer">
              <span class="card-year">{{ movie.release_year }}</span>
              <span class="card-genre">{{ movie.genre }}</span>
            </div>
          </div>
        </router-link>
      </div>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPages > 1">
        <button class="btn btn-secondary" :disabled="page === 1" @click="changePage(page - 1)">Previous</button>
        <span class="page-info">Page {{ page }} of {{ totalPages }}</span>
        <button class="btn btn-secondary" :disabled="page === totalPages" @click="changePage(page + 1)">Next</button>
      </div>
    </section>
  </div>
</template>

<script>
import api from '../api'

export default {
  name: 'Home',
  data() {
    return {
      movies: [],
      loading: true,
      error: null,
      search: '',
      genre: '',
      rating: 0,
      sort: '',
      page: 1,
      totalPages: 1,
      debounceTimer: null
    }
  },
  computed: {
    heroMovie() {
      // Return the top rated movie or first movie for the hero banner
      if (this.movies.length === 0) return null;
      return [...this.movies].sort((a, b) => b.average_rating - a.average_rating)[0];
    },
    displayMovies() {
      return this.movies;
    }
  },
  methods: {
    async fetchMovies() {
      this.loading = true;
      try {
        const response = await api.get('/movies', {
          params: {
            search: this.search,
            genre: this.genre,
            rating: this.rating,
            sort: this.sort,
            page: this.page
          }
        });
        // The new API format returns an object: { movies: [...], page: X, total_pages: Y }
        if (response.data && response.data.movies !== undefined) {
           this.movies = response.data.movies;
           this.page = response.data.page;
           this.totalPages = response.data.total_pages;
        } else {
           // Fallback if older API
           this.movies = response.data;
        }
      } catch (err) {
        this.error = 'Failed to load movie catalog. Please try again later.';
        console.error(err);
      } finally {
        this.loading = false;
      }
    },
    debouncedFetch() {
      clearTimeout(this.debounceTimer);
      this.debounceTimer = setTimeout(() => {
        this.page = 1;
        this.fetchMovies();
      }, 500);
    },
    changePage(newPage) {
      if (newPage >= 1 && newPage <= this.totalPages) {
        this.page = newPage;
        this.fetchMovies();
      }
    }
  },
  mounted() {
    this.fetchMovies();
  }
}
</script>

<style scoped>
.hero {
  height: 60vh;
  min-height: 400px;
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: center;
  margin-bottom: 3rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.hero-content {
  max-width: 600px;
}

.badge {
  display: inline-block;
  background-color: var(--primary-color);
  color: white;
  padding: 0.2rem 0.6rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 1rem;
}

.badge-outline {
  border: 1px solid rgba(255,255,255,0.3);
  padding: 0.1rem 0.5rem;
  border-radius: 4px;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  line-height: 1.1;
  margin-bottom: 1rem;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
}

.hero-desc {
  font-size: 1.1rem;
  color: #ddd;
  margin-bottom: 1.5rem;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.8);
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.hero-meta {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-bottom: 1.5rem;
  font-size: 0.95rem;
  color: #ccc;
}

.rating {
  color: #f5c518; /* IMDB gold */
  font-weight: 600;
}

.mt-4 {
  margin-top: 1rem;
}

.section-title {
  font-size: 1.8rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  border-left: 4px solid var(--primary-color);
  padding-left: 1rem;
}

.controls-bar {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  background: rgba(255,255,255,0.05);
  padding: 1rem;
  border-radius: var(--border-radius);
}

.form-control {
  padding: 0.6rem 1rem;
  border-radius: 4px;
  border: 1px solid rgba(255,255,255,0.1);
  background: rgba(0,0,0,0.5);
  color: white;
}

.search-input {
  flex-grow: 1;
  min-width: 200px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 3rem;
}

.page-info {
  font-weight: 600;
}

/* Grid */
.movie-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 2rem;
  padding-bottom: 4rem;
}

.movie-card {
  background-color: var(--bg-color-light);
  border-radius: var(--border-radius);
  overflow: hidden;
  transition: transform var(--transition), box-shadow var(--transition);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.movie-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.5);
}

.card-img-wrapper {
  position: relative;
  aspect-ratio: 2 / 3;
  overflow: hidden;
}

.card-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.movie-card:hover .card-img {
  transform: scale(1.05);
}

.card-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  padding: 1rem;
  background: linear-gradient(to bottom, rgba(0,0,0,0.8) 0%, transparent 100%);
  opacity: 0;
  transition: opacity var(--transition);
}

.movie-card:hover .card-overlay {
  opacity: 1;
}

.overlay-rating {
  background: rgba(0,0,0,0.7);
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #f5c518;
  font-weight: bold;
}

.card-body {
  padding: 1rem;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.loader, .error-msg {
  text-align: center;
  padding: 4rem 0;
  font-size: 1.2rem;
  color: var(--text-secondary);
}
.error-msg { color: var(--primary-color); }
</style>
