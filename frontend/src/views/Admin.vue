<template>
  <div class="admin-page container">
    <div class="admin-header">
      <h1 class="page-title">Admin Dashboard</h1>
    </div>

    <div class="tabs-container">
      <button class="tab-btn" :class="{ active: activeTab === 'movies' }" @click="activeTab = 'movies'">Manage Movies</button>
      <button class="tab-btn" :class="{ active: activeTab === 'reviews' }" @click="activeTab = 'reviews'; fetchReviews()">Moderate Reviews</button>
    </div>

    <!-- MOVIES TAB -->
    <div v-if="activeTab === 'movies'">
      <div class="admin-actions">
        <button class="btn btn-primary" @click="openAddForm">+ Add New Movie</button>
      </div>

    <div v-if="loading" class="loader">Loading catalog...</div>
    <div v-else-if="error" class="error-msg">{{ error }}</div>

    <div v-else>
      <div class="table-responsive glass-card">
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Image</th>
              <th>Title</th>
              <th>Genre</th>
              <th>Year</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="movie in movies" :key="movie.id">
              <td>{{ movie.id }}</td>
              <td>
                <img :src="movie.image_url" class="thumbnail" alt="thumbnail" />
              </td>
              <td><strong>{{ movie.title }}</strong></td>
              <td>{{ movie.genre }}</td>
              <td>{{ movie.release_year }}</td>
              <td class="actions">
                <button class="btn-action edit" @click="openEditForm(movie)">Edit</button>
                <button class="btn-action delete" @click="deleteMovie(movie.id)">Delete</button>
              </td>
            </tr>
            <tr v-if="movies.length === 0">
              <td colspan="6" class="text-center py-4">No movies currently in the catalog.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    </div>

    <!-- REVIEWS TAB -->
    <div v-if="activeTab === 'reviews'">
      <div class="table-responsive glass-card">
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>User</th>
              <th>Movie</th>
              <th>Rating</th>
              <th>Comment</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="review in reviews" :key="review.id">
              <td>{{ review.id }}</td>
              <td><strong>{{ review.user_name }}</strong></td>
              <td>{{ review.movie_title }}</td>
              <td>⭐ {{ review.rating }}</td>
              <td class="comment-cell">{{ review.comment || '-' }}</td>
              <td class="actions">
                <button class="btn-action delete" @click="deleteReview(review.id)">Delete</button>
              </td>
            </tr>
            <tr v-if="reviews.length === 0">
              <td colspan="6" class="text-center py-4">No reviews found.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add / Edit Modal Overlay -->
    <div v-if="showForm" class="modal-overlay">
      <div class="modal-content glass-card">
        <h2>{{ editingMovie ? 'Edit Movie' : 'Add New Movie' }}</h2>
        
        <form @submit.prevent="saveMovie">
          <div class="form-group">
            <label>Title</label>
            <input type="text" v-model="form.title" class="form-control" required />
          </div>
          
          <div class="form-row">
            <div class="form-group half">
              <label>Genre</label>
              <input type="text" v-model="form.genre" class="form-control" required />
            </div>
            
            <div class="form-group half">
              <label>Release Year</label>
              <input type="number" v-model.number="form.release_year" class="form-control" required />
            </div>
          </div>
          
          <div class="form-group">
            <label>Image URL</label>
            <input type="url" v-model="form.image_url" class="form-control" />
          </div>
          
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" class="form-control" rows="4"></textarea>
          </div>
          
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeForm">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? 'Saving...' : 'Save Movie' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import api from '../api'

export default {
  name: 'Admin',
  data() {
    return {
      activeTab: 'movies',
      movies: [],
      reviews: [],
      loading: true,
      error: null,
      
      showForm: false,
      saving: false,
      editingMovie: null,
      form: {
        title: '',
        genre: '',
        description: '',
        release_year: new Date().getFullYear(),
        image_url: ''
      }
    }
  },
  methods: {
    async fetchMovies() {
      this.loading = true;
      try {
        const res = await api.get('/movies');
        if (res.data && res.data.movies !== undefined) {
          this.movies = res.data.movies;
        } else {
          this.movies = res.data;
        }
      } catch (err) {
        this.error = "Failed to load movies.";
      } finally {
        this.loading = false;
      }
    },
    async fetchReviews() {
      this.loading = true;
      try {
        const res = await api.get('/reviews');
        this.reviews = res.data;
      } catch (err) {
        this.error = "Failed to load reviews.";
      } finally {
        this.loading = false;
      }
    },
    openAddForm() {
      this.editingMovie = null;
      this.form = {
        title: '',
        genre: '',
        description: '',
        release_year: new Date().getFullYear(),
        image_url: ''
      };
      this.showForm = true;
    },
    openEditForm(movie) {
      this.editingMovie = movie;
      this.form = { ...movie };
      this.showForm = true;
    },
    closeForm() {
      this.showForm = false;
    },
    async saveMovie() {
      this.saving = true;
      try {
        if (this.editingMovie) {
          // Update via PUT request
          await api.put(`/movies/${this.editingMovie.id}`, this.form);
        } else {
          // Create via POST request
          await api.post('/movies', this.form);
        }
        this.closeForm();
        this.fetchMovies(); // reload catalog
      } catch (err) {
        alert("Failed to save movie: " + (err.response?.data?.error || err.message));
      } finally {
        this.saving = false;
      }
    },
    async deleteMovie(id) {
      try {
        await api.delete(`/movies/${id}`);
        this.fetchMovies(); // reload catalog after delete
      } catch (err) {
        alert("Failed to delete movie: " + (err.response?.data?.error || err.message));
      }
    },
    async deleteReview(id) {
      try {
        await api.delete(`/reviews/${id}`);
        this.fetchReviews();
      } catch (err) {
        alert("Failed to delete review.");
      }
    }
  },
  mounted() {
    this.fetchMovies();
  }
}
</script>

<style scoped>
.admin-page {
  padding: 3rem 1rem;
}

.admin-header {
  margin-bottom: 2rem;
}

.tabs-container {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  border-bottom: 2px solid rgba(255,255,255,0.1);
}

.tab-btn {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 1.1rem;
  font-weight: 600;
  padding: 0.5rem 1rem;
  cursor: pointer;
  position: relative;
}

.tab-btn.active {
  color: var(--primary-color);
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  right: 0;
  height: 2px;
  background-color: var(--primary-color);
}

.admin-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}

.page-title {
  font-size: 2rem;
  border-left: 4px solid var(--primary-color);
  padding-left: 1rem;
}

.table-responsive {
  overflow-x: auto;
  padding: 0;
}

.table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.table th, .table td {
  padding: 1.2rem 1.5rem;
  border-bottom: 1px solid rgba(255,255,255,0.05);
}

.table th {
  background-color: rgba(0,0,0,0.3);
  font-weight: 600;
  color: var(--text-secondary);
}

.table tr:hover td {
  background-color: rgba(255,255,255,0.03);
}

.thumbnail {
  width: 50px;
  height: 70px;
  object-fit: cover;
  border-radius: 4px;
}

.btn-action {
  background: transparent;
  border: none;
  font-size: 0.9rem;
  padding: 0.4rem 0.8rem;
  cursor: pointer;
  border-radius: 4px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.btn-action.edit {
  color: #3f83f8;
  margin-right: 0.5rem;
}

.btn-action.edit:hover {
  background: rgba(63, 131, 248, 0.1);
}

.btn-action.delete {
  color: #f05252;
}

.btn-action.delete:hover {
  background: rgba(240, 82, 82, 0.1);
}

/* Modal styling */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 2rem;
}

.modal-content h2 {
  margin-bottom: 1.5rem;
  border-bottom: 1px solid rgba(255,255,255,0.1);
  padding-bottom: 1rem;
}

.form-row {
  display: flex;
  gap: 1.5rem;
}

.half {
  flex: 1;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255,255,255,0.1);
}

.text-center { text-align: center; }
.py-4 { padding: 2rem 0; }
</style>
