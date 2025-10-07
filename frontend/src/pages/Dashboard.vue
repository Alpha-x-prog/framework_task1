<template>
  <div class="container">
    <h1>Проекты</h1>
    
    <!-- Форма добавления проекта -->
    <div class="form">
      <h3>Добавить проект</h3>
      <form @submit.prevent="handleCreateProject">
        <div class="form-row">
          <div class="form-group">
            <label for="project-name">Название *</label>
            <input
              id="project-name"
              v-model="newProject.name"
              type="text"
              class="form-control"
              required
              :disabled="loading"
              placeholder="Название проекта"
            />
          </div>
          <div class="form-group">
            <label for="project-customer">Заказчик</label>
            <input
              id="project-customer"
              v-model="newProject.customer"
              type="text"
              class="form-control"
              :disabled="loading"
              placeholder="Название заказчика"
            />
          </div>
        </div>
        <button 
          type="submit" 
          class="btn btn-primary"
          :disabled="loading"
        >
          {{ loading ? 'Создание...' : 'Создать проект' }}
        </button>
      </form>
    </div>
    
    <!-- Список проектов -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">Список проектов</h3>
        <button @click="loadProjects" class="btn btn-secondary" :disabled="loading">
          Обновить
        </button>
      </div>
      
      <div v-if="loading && projects.length === 0" class="loading">
        Загрузка проектов...
      </div>
      
      <div v-else-if="projects.length === 0" class="loading">
        Проекты не найдены
      </div>
      
      <table v-else class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Заказчик</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="project in projects" :key="project.id">
            <td>{{ project.id }}</td>
            <td>{{ project.name }}</td>
            <td>{{ project.customer || '-' }}</td>
            <td>
              <router-link :to="`/defects?project_id=${project.id}`" class="btn btn-primary">
                Дефекты
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Компонент скачивания отчёта -->
    <DownloadReport />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { projectsApi } from '../api'
import DownloadReport from '../components/DownloadReport.vue'

export default {
  name: 'Dashboard',
  components: {
    DownloadReport
  },
  setup() {
    const projects = ref([])
    const loading = ref(false)
    const error = ref('')
    
    const newProject = ref({
      name: '',
      customer: ''
    })
    
    const loadProjects = async () => {
      loading.value = true
      error.value = ''
      
      try {
        projects.value = await projectsApi.getProjects()
      } catch (err) {
        error.value = err.response?.data?.message || 'Ошибка загрузки проектов'
        console.error('Ошибка загрузки проектов:', err)
      } finally {
        loading.value = false
      }
    }
    
    const handleCreateProject = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const projectData = {
          name: newProject.value.name
        }
        
        if (newProject.value.customer) {
          projectData.customer = newProject.value.customer
        }
        
        await projectsApi.createProject(projectData)
        
        // Очищаем форму
        newProject.value = {
          name: '',
          customer: ''
        }
        
        // Перезагружаем список проектов
        await loadProjects()
        
        alert('Проект успешно создан!')
      } catch (err) {
        error.value = err.response?.data?.message || 'Ошибка создания проекта'
        alert(error.value)
      } finally {
        loading.value = false
      }
    }
    
    onMounted(() => {
      loadProjects()
    })
    
    return {
      projects,
      loading,
      error,
      newProject,
      loadProjects,
      handleCreateProject
    }
  }
}
</script>
