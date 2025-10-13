<template>
  <div class="container">
    <div class="forbidden-page">
      <div class="forbidden-content">
        <h1>🚫 Доступ запрещен</h1>
        <p>У вас нет прав для доступа к этой странице.</p>
        <p>Ваша роль: <strong>{{ role }}</strong></p>
        <div class="actions">
          <router-link to="/projects" class="btn btn-primary">
            ← Вернуться к проектам
          </router-link>
          <button @click="logout" class="btn btn-secondary">
            Выйти из системы
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  name: 'Forbidden',
  setup() {
    const authStore = useAuthStore()
    
    const logout = () => {
      authStore.logout()
      window.location.href = '/login'
    }
    
    return {
      role: authStore.role,
      logout
    }
  }
}
</script>

<style scoped>
.forbidden-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
  text-align: center;
}

.forbidden-content {
  max-width: 500px;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.forbidden-content h1 {
  color: #dc3545;
  margin-bottom: 1rem;
  font-size: 2rem;
}

.forbidden-content p {
  margin-bottom: 1rem;
  color: #666;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 2rem;
}

.actions .btn {
  padding: 0.75rem 1.5rem;
  text-decoration: none;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn:hover {
  opacity: 0.9;
}
</style>
