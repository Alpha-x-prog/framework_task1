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
          <button class="btn btn-secondary" @click="toggleTheme" :aria-label="`Toggle theme`">
            {{ theme === 'dark' ? '🌙' : '☀️' }}
          </button>
          <span class="user-email" v-if="user?.email">{{ user.email }}</span>
          <button @click="handleLogout" class="btn btn-secondary">Выйти</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { computed, ref, onMounted } from 'vue'
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
    const theme = ref('light')

    const applyTheme = (t) => {
      const el = document.documentElement
      if (t === 'dark') {
        el.setAttribute('data-theme', 'dark')
      } else {
        el.removeAttribute('data-theme')
      }
    }

    const toggleTheme = () => {
      theme.value = theme.value === 'dark' ? 'light' : 'dark'
      localStorage.setItem('theme', theme.value)
      applyTheme(theme.value)
    }

    onMounted(() => {
      const saved = localStorage.getItem('theme')
      if (saved === 'dark' || saved === 'light') {
        theme.value = saved
      } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        theme.value = 'dark'
      }
      applyTheme(theme.value)
    })
    
    const handleLogout = () => {
      authStore.logout()
      router.push('/login')
    }
    
    return {
      user,
      handleLogout,
      can,
      theme,
      toggleTheme
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
  color: var(--color-text);
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: var(--color-surface-alt);
}

.nav-link.router-link-active {
  background-color: var(--color-primary);
  color: white;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-email {
  opacity: 0.85;
}
</style>
