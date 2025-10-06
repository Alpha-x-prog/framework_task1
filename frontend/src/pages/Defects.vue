<template>
  <div class="container">
    <h1>Дефекты</h1>
    
    <!-- Фильтры -->
    <div class="filters">
      <h3>Фильтры</h3>
      <form @submit.prevent="applyFilters">
        <div class="filters-row">
          <ProjectSelect v-model="filters.project_id" />
          
          <div class="form-group">
            <label for="status-filter">Статус</label>
            <select id="status-filter" v-model="filters.status_id" class="form-control">
              <option value="">Все статусы</option>
              <option v-for="status in statuses" :key="status.id" :value="status.id">
                {{ status.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="priority-filter">Приоритет</label>
            <select id="priority-filter" v-model="filters.priority" class="form-control">
              <option value="">Все приоритеты</option>
              <option value="1">1 - Критический</option>
              <option value="2">2 - Высокий</option>
              <option value="3">3 - Средний</option>
              <option value="4">4 - Низкий</option>
              <option value="5">5 - Очень низкий</option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="search">Поиск</label>
            <input
              id="search"
              v-model="filters.q"
              type="text"
              class="form-control"
              placeholder="Поиск по названию..."
            />
          </div>
          
          <div class="form-group">
            <label for="due-from">Дата с</label>
            <input
              id="due-from"
              v-model="filters.due_from"
              type="date"
              class="form-control"
            />
          </div>
          
          <div class="form-group">
            <label for="due-to">Дата по</label>
            <input
              id="due-to"
              v-model="filters.due_to"
              type="date"
              class="form-control"
            />
          </div>
          
          <div class="form-group">
            <button type="submit" class="btn btn-primary">Применить</button>
            <button type="button" @click="clearFilters" class="btn btn-secondary">Очистить</button>
          </div>
        </div>
      </form>
    </div>
    
    <!-- Список дефектов -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">Список дефектов</h3>
        <button @click="loadDefects" class="btn btn-secondary" :disabled="loading">
          Обновить
        </button>
      </div>
      
      <div v-if="loading && defects.length === 0" class="loading">
        Загрузка дефектов...
      </div>
      
      <div v-else-if="defects.length === 0" class="loading">
        Дефекты не найдены
      </div>
      
      <table v-else class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Проект</th>
            <th>Статус</th>
            <th>Приоритет</th>
            <th>Срок</th>
            <th>Обновлен</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="defect in defects" :key="defect.id">
            <td>{{ defect.id }}</td>
            <td>
              <router-link :to="`/defects/${defect.id}`">
                {{ defect.title }}
              </router-link>
            </td>
            <td>{{ getProjectName(defect.project_id) }}</td>
            <td>
              <span :class="`status status-${getStatusClass(defect.status_id)}`">
                {{ getStatusName(defect.status_id) }}
              </span>
            </td>
            <td>
              <span :class="`priority priority-${defect.priority}`">
                {{ defect.priority }}
              </span>
            </td>
            <td>{{ defect.due_date ? formatDate(defect.due_date) : '-' }}</td>
            <td>{{ formatDate(defect.updated_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Форма создания дефекта -->
    <div class="form">
      <h3>Создать дефект</h3>
      <form @submit.prevent="handleCreateDefect">
        <div class="form-row">
          <div class="form-group">
            <label for="defect-project">Проект *</label>
            <select
              id="defect-project"
              v-model="newDefect.project_id"
              class="form-control"
              required
              :disabled="loading"
            >
              <option value="">Выберите проект</option>
              <option v-for="project in projects" :key="project.id" :value="project.id">
                {{ project.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="defect-title">Название *</label>
            <input
              id="defect-title"
              v-model="newDefect.title"
              type="text"
              class="form-control"
              required
              :disabled="loading"
              placeholder="Название дефекта"
            />
          </div>
        </div>
        
        <div class="form-group">
          <label for="defect-description">Описание</label>
          <textarea
            id="defect-description"
            v-model="newDefect.description"
            class="form-control"
            rows="3"
            :disabled="loading"
            placeholder="Описание дефекта"
          ></textarea>
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label for="defect-priority">Приоритет</label>
            <select
              id="defect-priority"
              v-model="newDefect.priority"
              class="form-control"
              :disabled="loading"
            >
              <option value="3">3 - Средний</option>
              <option value="1">1 - Критический</option>
              <option value="2">2 - Высокий</option>
              <option value="4">4 - Низкий</option>
              <option value="5">5 - Очень низкий</option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="defect-status">Статус</label>
            <select
              id="defect-status"
              v-model="newDefect.status_id"
              class="form-control"
              :disabled="loading"
            >
              <option value="">Выберите статус</option>
              <option v-for="status in statuses" :key="status.id" :value="status.id">
                {{ status.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="defect-due">Срок</label>
            <input
              id="defect-due"
              v-model="newDefect.due_date"
              type="date"
              class="form-control"
              :disabled="loading"
            />
          </div>
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary"
          :disabled="loading"
        >
          {{ loading ? 'Создание...' : 'Создать дефект' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { defectsApi, projectsApi, refsApi } from '../api'
import ProjectSelect from '../components/ProjectSelect.vue'

export default {
  name: 'Defects',
  components: {
    ProjectSelect
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    const defects = ref([])
    const projects = ref([])
    const statuses = ref([])
    const loading = ref(false)
    const error = ref('')
    
    const filters = ref({
      project_id: route.query.project_id || '',
      status_id: '',
      priority: '',
      q: '',
      due_from: '',
      due_to: '',
      assignee_id: ''
    })
    
    const newDefect = ref({
      project_id: route.query.project_id || '',
      title: '',
      description: '',
      priority: '3',
      status_id: '',
      due_date: '',
      assignee_id: ''
    })
    
    const getProjectName = (projectId) => {
      const project = projects.value.find(p => p.id === projectId)
      return project ? project.name : `ID: ${projectId}`
    }
    
    const getStatusName = (statusId) => {
      const status = statuses.value.find(s => s.id === statusId)
      return status ? status.name : `ID: ${statusId}`
    }
    
    const getStatusClass = (statusId) => {
      const status = statuses.value.find(s => s.id === statusId)
      if (!status) return 'unknown'
      
      const name = status.name.toLowerCase()
      if (name.includes('new')) return 'new'
      if (name.includes('progress')) return 'in-progress'
      if (name.includes('resolved')) return 'resolved'
      if (name.includes('closed')) return 'closed'
      return 'unknown'
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      return new Date(dateString).toLocaleDateString('ru-RU')
    }
    
    const loadDefects = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const params = {}
        Object.keys(filters.value).forEach(key => {
          if (filters.value[key]) {
            params[key] = filters.value[key]
          }
        })
        
        defects.value = await defectsApi.getDefects(params)
      } catch (err) {
        error.value = err.response?.data?.message || 'Ошибка загрузки дефектов'
        console.error('Ошибка загрузки дефектов:', err)
      } finally {
        loading.value = false
      }
    }
    
    const loadProjects = async () => {
      try {
        projects.value = await projectsApi.getProjects()
      } catch (err) {
        console.error('Ошибка загрузки проектов:', err)
      }
    }
    
    const loadStatuses = async () => {
      try {
        statuses.value = await refsApi.getStatuses()
      } catch (err) {
        console.error('Ошибка загрузки статусов:', err)
      }
    }
    
    const applyFilters = () => {
      loadDefects()
    }
    
    const clearFilters = () => {
      filters.value = {
        project_id: '',
        status_id: '',
        priority: '',
        q: '',
        due_from: '',
        due_to: '',
        assignee_id: ''
      }
      loadDefects()
    }
    
    const handleCreateDefect = async () => {
      loading.value = true
      error.value = ''
      
      try {
        const defectData = {
          project_id: newDefect.value.project_id,
          title: newDefect.value.title
        }
        
        if (newDefect.value.description) {
          defectData.description = newDefect.value.description
        }
        if (newDefect.value.priority) {
          defectData.priority = parseInt(newDefect.value.priority)
        }
        if (newDefect.value.status_id) {
          defectData.status_id = newDefect.value.status_id
        }
        if (newDefect.value.due_date) {
          defectData.due_date = newDefect.value.due_date
        }
        if (newDefect.value.assignee_id) {
          defectData.assignee_id = parseInt(newDefect.value.assignee_id)
        }
        
        const response = await defectsApi.createDefect(defectData)
        
        // Очищаем форму
        newDefect.value = {
          project_id: '',
          title: '',
          description: '',
          priority: '3',
          status_id: '',
          due_date: '',
          assignee_id: ''
        }
        
        // Переходим к созданному дефекту
        router.push(`/defects/${response.id}`)
      } catch (err) {
        error.value = err.response?.data?.message || 'Ошибка создания дефекта'
        alert(error.value)
      } finally {
        loading.value = false
      }
    }
    
    onMounted(async () => {
      await Promise.all([
        loadProjects(),
        loadStatuses(),
        loadDefects()
      ])
    })
    
    return {
      defects,
      projects,
      statuses,
      loading,
      error,
      filters,
      newDefect,
      getProjectName,
      getStatusName,
      getStatusClass,
      formatDate,
      loadDefects,
      applyFilters,
      clearFilters,
      handleCreateDefect
    }
  }
}
</script>
