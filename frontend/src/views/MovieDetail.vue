<template>
  <div class="movie-detail container">
    <div v-if="loading" class="loader">Loading movie details...</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>
    
    <template v-else-if="movie">
      <!-- Top Section -->
      <div class="detail-header">
        <div class="poster-wrapper">
          <img :src="movie.image_url" :alt="movie.title" class="poster-img" />
        </div>
        
        <div class="info-wrapper">
          <h1 class="movie-title">{{ movie.title }}</h1>
          
          <div class="meta-row">
            <span class="meta-item year">{{ movie.release_year }}</span>
            <span class="meta-item badge-outline">{{ movie.genre }}</span>
            <span class="meta-item rating">⭐ {{ movie.average_rating.toFixed(1) }}</span>
          </div>
          
          <div class="description glass-card">
            <p>{{ movie.description }}</p>
          </div>
          
          <div class="actions" v-if="isAuthenticated">
            <button @click="addToWatchlist" class="btn btn-primary" :disabled="watchlistLoading">
              {{ watchlistLoading ? 'Adding...' : '+ Add to Watchlist' }}
            </button>
            <span v-if="watchlistMessage" class="ml-3 text-success">{{ watchlistMessage }}</span>
          </div>
          <div class="actions" v-else>
            <p class="text-secondary"><router-link to="/login">Sign in</router-link> to add this to your watchlist or leave a review.</p>
          </div>
        </div>
      </div>
      
      <!-- Review Section -->
      <div class="review-section">
        <h2 class="section-title">User Reviews</h2>
        
        <div v-if="isAuthenticated" class="glass-card add-review-card">
          <h3>Write a Review</h3>
          <form @submit.prevent="submitReview">
            <div class="form-group row-flex">
              <label for="rating">Rating (1-5)</label>
              <select id="rating" v-model.number="newReview.rating" class="form-control short-input">
                <option v-for="n in 5" :key="n" :value="n">{{ n }} Stars</option>
              </select>
            </div>
            
            <div class="form-group">
              <label for="comment">Your thoughts</label>
              <textarea 
                id="comment" 
                v-model="newReview.comment" 
                class="form-control" 
                rows="4"
                placeholder="What did you think of the movie?"
              ></textarea>
            </div>
            
            <button type="submit" class="btn btn-secondary" :disabled="reviewSubmitting">
              {{ reviewSubmitting ? 'Posting...' : 'Post Review' }}
            </button>
            <p v-if="reviewError" class="text-error mt-2">{{ reviewError }}</p>
          </form>
        </div>

        <!-- Review List -->
        <div class="reviews-list">
          <p v-if="reviews.length === 0" class="text-secondary">No reviews yet. Be the first to review!</p>
          
          <div v-for="review in reviews" :key="review.id" class="review-item glass-card">
            <div class="review-header">
              <div class="rating">⭐ {{ review.rating }} / 5</div>
              <div class="review-date">{{ new Date(review.created_at).toLocaleDateString() }}</div>
            </div>
            <p class="review-comment">{{ review.comment }}</p>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import api from '../api'
import { mapState } from 'pinia'
import { useAuthStore } from '../stores/auth'

export default {
  name: 'MovieDetail',
  data() {
    return {
      movie: null,
      reviews: [],
      loading: true,
      error: null,
      
      watchlistLoading: false,
      watchlistMessage: '',
      
      newReview: {
        rating: 5,
        comment: ''
      },
      reviewSubmitting: false,
      reviewError: null
    }
  },
  computed: {
    ...mapState(useAuthStore, ['isAuthenticated'])
  },
  methods: {
    async fetchMovie() {
      const id = this.$route.params.id;
      this.loading = true;
      try {
        const [movieRes, reviewsRes] = await Promise.all([
          api.get(`/movies/${id}`),
          api.get(`/movies/${id}/reviews`)
        ]);
        this.movie = movieRes.data;
        this.reviews = reviewsRes.data;
      } catch (err) {
        this.error = "Failed to load movie. It might not exist.";
      } finally {
        this.loading = false;
      }
    },
    async addToWatchlist() {
      this.watchlistLoading = true;
      this.watchlistMessage = '';
      try {
        await api.post('/watchlist', {
          movie_id: this.movie.id
        });
        this.watchlistMessage = "Added to watchlist!";
        setTimeout(() => { this.watchlistMessage = '' }, 3000);
      } catch (err) {
        this.watchlistMessage = "Error adding to watchlist.";
      } finally {
        this.watchlistLoading = false;
      }
    },
    async submitReview() {
      if (this.newReview.comment && this.newReview.comment.length < 10) {
        this.reviewError = "Review comment must be at least 10 characters long.";
        return;
      }
      
      this.reviewSubmitting = true;
      this.reviewError = null;
      try {
        const payload = {
          movie_id: this.movie.id,
          rating: this.newReview.rating,
          comment: this.newReview.comment
        };
        const res = await api.post('/reviews', payload);
        
        // Push new review to array 
        this.reviews.push(res.data);
        
        // Reset form
        this.newReview.comment = '';
        this.newReview.rating = 5;
        
        // Slight hack: we don't recalculate the average rating blindly here
        // but fetching again would update the overall score smoothly
        this.fetchMovie();
      } catch (err) {
        this.reviewError = err.response?.data?.error || "Error posting review. Note: You can only review a movie once.";
      } finally {
        this.reviewSubmitting = false;
      }
    }
  },
  mounted() {
    this.fetchMovie();
  }
}
</script>

<style scoped>
.movie-detail {
  padding: 3rem 1rem;
}

.detail-header {
  display: flex;
  gap: 3rem;
  margin-bottom: 4rem;
}

@media (max-width: 768px) {
  .detail-header {
    flex-direction: column;
  }
}

.poster-wrapper {
  flex: 0 0 300px;
  border-radius: var(--border-radius);
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0,0,0,0.5);
  align-self: flex-start;
}

.poster-img {
  width: 100%;
  height: auto;
  display: block;
}

.info-wrapper {
  flex: 1;
}

.movie-title {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 1rem;
  line-height: 1.1;
}

.meta-row {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-bottom: 2rem;
  font-size: 1.1rem;
}

.badge-outline {
  border: 1px solid rgba(255,255,255,0.3);
  padding: 0.2rem 0.6rem;
  border-radius: 4px;
}

.rating {
  color: #f5c518;
  font-weight: 600;
}

.description {
  font-size: 1.1rem;
  line-height: 1.8;
  margin-bottom: 2rem;
  padding: 1.5rem;
}

.actions {
  display: flex;
  align-items: center;
}

.actions a {
  text-decoration: underline;
  color: var(--primary-color);
}

.ml-3 { margin-left: 1rem; }
.text-success { color: #28a745; font-weight: 600; }
.text-error { color: var(--primary-color); }
.text-secondary { color: var(--text-secondary); }
.mt-2 { margin-top: 1rem; }

/* Reviews */
.section-title {
  font-size: 1.8rem;
  font-weight: 600;
  margin-bottom: 2rem;
  border-left: 4px solid var(--primary-color);
  padding-left: 1rem;
}

.add-review-card {
  margin-bottom: 3rem;
  padding: 2rem;
}

.add-review-card h3 {
  margin-bottom: 1.5rem;
}

.row-flex {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.short-input {
  width: auto;
  min-width: 120px;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.review-item {
  padding: 1.5rem;
  border-left: 4px solid #333;
}

.review-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1rem;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.review-header .rating {
  font-size: 1rem;
}

.review-comment {
  line-height: 1.6;
}
</style>
