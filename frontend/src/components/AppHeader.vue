<template>
  <header class="header">
    <div class="container">
      <div class="header-content">
        <div class="logo">Defects App</div>
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

export default {
  name: 'AppHeader',
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()
    
    const user = computed(() => authStore.user)
    
    const handleLogout = () => {
      authStore.logout()
      router.push('/login')
    }
    
    return {
      user,
      handleLogout
    }
  }
}
</script>
