<template>
  <header class="header">
    <div class="container">
      <div class="header-content">
        <div class="logo">Defects App</div>
        
        <nav class="nav">
          <router-link to="/projects" class="nav-link">Проекты</router-link>
          <router-link to="/defects" class="nav-link">Дефекты</router-link>
          <router-link v-if="can.downloadReport" to="/reports" class="nav-link">Отчёты</router-link>
        </nav>
        
        <div class="user-info">
          <span>{{ user?.email }}</span>
          <button @click="handleLogout" class="btn btn-secondary">Выйти</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { usePermissions } from '../composables/usePermissions'

export default {
  name: 'AppHeader',
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()
    const { can } = usePermissions()
    
    const user = computed(() => authStore.user)
    
    const handleLogout = () => {
      authStore.logout()
      router.push('/login')
    }
    
    return {
      user,
      handleLogout,
      can
    }
  }
}
</script>

<style scoped>
.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
}

.nav {
  display: flex;
  gap: 1rem;
}

.nav-link {
  color: #495057;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: #f8f9fa;
}

.nav-link.router-link-active {
  background-color: #007bff;
  color: white;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}
</style>
