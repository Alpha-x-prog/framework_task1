<template>
  <div class="container">
    <div class="form" style="max-width: 400px; margin: 2rem auto;">
      <h2 style="margin-bottom: 2rem; text-align: center;">Регистрация</h2>
      
      <div v-if="error" class="alert alert-error">
        {{ error }}
      </div>
      
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="email">Email</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            class="form-control"
            required
            :disabled="loading"
          />
        </div>
        
        <div class="form-group">
          <label for="password">Пароль</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            class="form-control"
            required
            :disabled="loading"
          />
        </div>
        
        <div class="form-group">
          <label for="role">Роль</label>
          <select
            id="role"
            v-model="form.role"
            class="form-control"
            :disabled="loading"
          >
            <option value="">Выберите роль</option>
            <option 
              v-for="role in roles" 
              :key="role.id" 
              :value="role.name"
            >
              {{ role.label }}
            </option>
          </select>
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary" 
          style="width: 100%;"
          :disabled="loading"
        >
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
      </form>
      
      <div style="text-align: center; margin-top: 1rem;">
        <p>Уже есть аккаунт? <router-link to="/login">Войти</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { refsApi } from '../api'

export default {
  name: 'Register',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = ref({
      email: '',
      password: '',
      role: ''
    })
    
    const roles = ref([])
    const loading = ref(false)
    const error = ref('')
    
    const loadRoles = async () => {
      try {
        // Загружаем роли из API
        const apiRoles = await refsApi.getRoles()
        
        // Маппинг ролей для отображения
        const roleLabels = {
          'manager': 'Руководитель / менеджер проекта',
          'engineer': 'Инженер-исполнитель',
          'viewer': 'Заказчик / наблюдатель',
          'lead': 'Руководитель отдела'
        }
        
        // Преобразуем роли из API с добавлением лейблов
        roles.value = apiRoles.map(role => ({
          ...role,
          label: roleLabels[role.name] || role.name
        }))
        
        // Устанавливаем роль по умолчанию "engineer"
        const engineerRole = roles.value.find(role => role.name === 'engineer')
        if (engineerRole) {
          form.value.role = engineerRole.name
        }
      } catch (err) {
        console.error('Ошибка загрузки ролей:', err)
        // Fallback к статическому списку в случае ошибки API
        roles.value = [
          { id: 1, name: 'engineer', label: 'Инженер-исполнитель' },
          { id: 2, name: 'manager', label: 'Руководитель / менеджер проекта' },
          { id: 3, name: 'viewer', label: 'Заказчик / наблюдатель' },
          { id: 4, name: 'lead', label: 'Руководитель отдела' }
        ]
        form.value.role = 'engineer'
      }
    }
    
    const handleRegister = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const registerData = {
          email: form.value.email,
          password: form.value.password
        }
        
        // Добавляем роль только если она выбрана
        if (form.value.role) {
          // form.value.role содержит name роли из API
          registerData.role = form.value.role
        }
        
        console.log('Отправляем данные регистрации:', registerData)
        await authStore.register(registerData)
        router.push('/projects')
      } catch (err) {
        error.value = err.response?.data?.message || 'Ошибка регистрации'
      } finally {
        loading.value = false
      }
    }
    
    onMounted(() => {
      loadRoles()
    })
    
    return {
      form,
      roles,
      loading,
      error,
      handleRegister
    }
  }
}
</script>
