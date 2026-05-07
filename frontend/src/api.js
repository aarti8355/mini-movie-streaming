import axios from 'axios';
import { useAuthStore } from './stores/auth';

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
});

// Intercept requests to attach the JWT token
api.interceptors.request.use(config => {
  const authStore = useAuthStore();
  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`;
  }
  return config;
});

export default api;
