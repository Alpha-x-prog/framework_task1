<template>
  <div class="form-group">
    <label for="project-select">Проект</label>
    <select 
      id="project-select" 
      v-model="selectedProjectId" 
      class="form-control"
      @change="handleChange"
    >
      <option value="">Все проекты</option>
      <option 
        v-for="project in projects" 
        :key="project.id" 
        :value="project.id"
      >
        {{ project.name }} {{ project.customer ? `(${project.customer})` : '' }}
      </option>
    </select>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { projectsApi } from '../api'

export default {
  name: 'ProjectSelect',
  props: {
    modelValue: {
      type: [String, Number],
      default: ''
    }
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const projects = ref([])
    const selectedProjectId = ref(props.modelValue)
    
    // Следим за изменением modelValue и обновляем selectedProjectId
    watch(() => props.modelValue, (newValue) => {
      selectedProjectId.value = newValue
    })
    
    const loadProjects = async () => {
      try {
        projects.value = await projectsApi.getProjects()
      } catch (error) {
        console.error('Ошибка загрузки проектов:', error)
      }
    }
    
    const handleChange = () => {
      emit('update:modelValue', selectedProjectId.value)
    }
    
    onMounted(() => {
      loadProjects()
    })
    
    return {
      projects,
      selectedProjectId,
      handleChange
    }
  }
}
</script>
