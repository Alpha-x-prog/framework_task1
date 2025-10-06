<template>
  <div id="app">
    <AppHeader v-if="isAuthenticated" />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useAuthStore } from './stores/auth'
import AppHeader from './components/AppHeader.vue'

export default {
  name: 'App',
  components: {
    AppHeader
  },
  setup() {
    const authStore = useAuthStore()
    
    const isAuthenticated = computed(() => authStore.isAuthed)
    
    return {
      isAuthenticated
    }
  }
}
</script>

<style scoped>
.main-content {
  flex: 1;
  padding: 2rem 0;
}
</style>
