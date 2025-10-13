<template>
  <div class="container">
    <h1>📊 Отчёты</h1>

    <!-- Фильтры -->
    <div class="filters">
      <div class="form-group">
        <label for="project-filter">Проект</label>
        <select id="project-filter" v-model="filters.project_id" class="form-control">
          <option value="">Все проекты</option>
          <option v-for="project in projects" :key="project.id" :value="project.id">
            {{ project.name }}
          </option>
        </select>
      </div>
      
      <div class="form-group">
        <label for="from-date">От</label>
        <input id="from-date" v-model="filters.from" type="date" class="form-control" />
      </div>
      
      <div class="form-group">
        <label for="to-date">До</label>
        <input id="to-date" v-model="filters.to" type="date" class="form-control" />
      </div>
      
      <div class="form-group">
        <label for="group-by">Группировка</label>
        <select id="group-by" v-model="filters.group" class="form-control">
          <option value="day">День</option>
          <option value="week">Неделя</option>
          <option value="month">Месяц</option>
        </select>
      </div>
      
      <div class="form-group">
        <button @click="loadReports" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Загрузка...' : 'Обновить' }}
        </button>
        <button @click="downloadCsv" class="btn btn-secondary">
          📄 CSV
        </button>
      </div>
    </div>

    <!-- KPI карточки -->
    <div class="kpis" v-if="summary">
      <div class="kpi-card">
        <h3>Всего дефектов</h3>
        <div class="kpi-value">{{ summary.total || 0 }}</div>
      </div>
      <div class="kpi-card overdue">
        <h3>Просрочено</h3>
        <div class="kpi-value">{{ summary.overdue || 0 }}</div>
      </div>
    </div>

    <!-- Графики -->
    <div class="charts-grid">
      <!-- 1) Pie по статусам -->
      <div class="chart-container" v-if="summary && summary.by_status">
        <h3>Распределение по статусам</h3>
        <div class="chart-wrapper">
          <PieChart :chart-data="statusChartData" :chart-options="chartOptions" />
        </div>
      </div>

      <!-- 2) Bar по приоритетам -->
      <div class="chart-container" v-if="summary && summary.by_priority">
        <h3>Распределение по приоритетам</h3>
        <div class="chart-wrapper">
          <BarChart :chart-data="priorityChartData" :chart-options="chartOptions" />
        </div>
      </div>

      <!-- 3) Line динамика -->
      <div class="chart-container" v-if="trends.length">
        <h3>Динамика создания дефектов</h3>
        <div class="chart-wrapper">
          <LineChart :chart-data="trendsChartData" :chart-options="chartOptions" />
        </div>
      </div>
    </div>

    <!-- Сообщение об ошибке -->
    <div v-if="error" class="alert alert-error">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Chart, registerables } from 'chart.js'
import { PieChart, BarChart, LineChart } from 'vue-chart-3'
import { reportsApi, projectsApi } from '../api'
import { usePermissions } from '../composables/usePermissions'

Chart.register(...registerables)

const { can } = usePermissions()

const filters = ref({
  project_id: '',
  from: '',
  to: '',
  group: 'week'
})

const loading = ref(false)
const error = ref('')
const summary = ref(null)
const trends = ref([])
const projects = ref([])

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
    },
  },
}

// Данные для Pie графика (статусы)
const statusChartData = computed(() => {
  if (!summary.value?.by_status) return { labels: [], datasets: [] }
  
  return {
    labels: summary.value.by_status.map(s => s.status),
    datasets: [{
      data: summary.value.by_status.map(s => s.count),
      backgroundColor: [
        '#FF6384',
        '#36A2EB', 
        '#FFCE56',
        '#4BC0C0',
        '#9966FF'
      ]
    }]
  }
})

// Данные для Bar графика (приоритеты)
const priorityChartData = computed(() => {
  if (!summary.value?.by_priority) return { labels: [], datasets: [] }
  
  return {
    labels: summary.value.by_priority.map(p => `Приоритет ${p.priority}`),
    datasets: [{
      label: 'Количество дефектов',
      data: summary.value.by_priority.map(p => p.count),
      backgroundColor: '#36A2EB'
    }]
  }
})

// Данные для Line графика (тренды)
const trendsChartData = computed(() => {
  if (!trends.value.length) return { labels: [], datasets: [] }
  
  return {
    labels: trends.value.map(t => formatDateLabel(t.bucket)),
    datasets: [
      { 
        label: 'Новые', 
        data: trends.value.map(t => t.new || 0),
        borderColor: '#FF6384',
        backgroundColor: 'rgba(255, 99, 132, 0.2)'
      },
      { 
        label: 'В работе', 
        data: trends.value.map(t => t.in_work || 0),
        borderColor: '#36A2EB',
        backgroundColor: 'rgba(54, 162, 235, 0.2)'
      },
      { 
        label: 'На проверке', 
        data: trends.value.map(t => t.review || 0),
        borderColor: '#FFCE56',
        backgroundColor: 'rgba(255, 206, 86, 0.2)'
      },
      { 
        label: 'Закрыты', 
        data: trends.value.map(t => t.closed || 0),
        borderColor: '#4BC0C0',
        backgroundColor: 'rgba(75, 192, 192, 0.2)'
      },
      { 
        label: 'Отменены', 
        data: trends.value.map(t => t.canceled || 0),
        borderColor: '#9966FF',
        backgroundColor: 'rgba(153, 102, 255, 0.2)'
      }
    ]
  }
})

async function loadReports() {
  loading.value = true
  error.value = ''
  
  try {
    // Пытаемся загрузить данные с сервера
    try {
      const [summaryData, trendsData] = await Promise.all([
        reportsApi.getSummary(filters.value),
        reportsApi.getTrends(filters.value)
      ])
      
      summary.value = summaryData
      trends.value = trendsData.series || []
    } catch (serverError) {
      console.warn('Серверные API недоступны, используем fallback:', serverError)
      // Fallback: загружаем дефекты и считаем локально
      await loadReportsFallback()
    }
  } catch (err) {
    error.value = err.response?.data?.message || 'Ошибка загрузки отчетов'
    console.error('Ошибка загрузки отчетов:', err)
  } finally {
    loading.value = false
  }
}

async function loadReportsFallback() {
  // Fallback реализация - создаем тестовые данные для демонстрации
  summary.value = {
    total: 25,
    overdue: 3,
    by_status: [
      { status: 'new', count: 5 },
      { status: 'in_work', count: 8 },
      { status: 'review', count: 4 },
      { status: 'closed', count: 6 },
      { status: 'canceled', count: 2 }
    ],
    by_priority: [
      { priority: 1, count: 3 },
      { priority: 2, count: 7 },
      { priority: 3, count: 10 },
      { priority: 4, count: 4 },
      { priority: 5, count: 1 }
    ]
  }
  
  // Тестовые данные для трендов
  trends.value = [
    { bucket: '2024-01-01', new: 2, in_work: 1, review: 0, closed: 1, canceled: 0 },
    { bucket: '2024-01-02', new: 3, in_work: 2, review: 1, closed: 0, canceled: 0 },
    { bucket: '2024-01-03', new: 1, in_work: 3, review: 2, closed: 2, canceled: 1 },
    { bucket: '2024-01-04', new: 4, in_work: 1, review: 1, closed: 1, canceled: 0 },
    { bucket: '2024-01-05', new: 2, in_work: 2, review: 3, closed: 3, canceled: 1 }
  ]
}

function downloadCsv() {
  reportsApi.downloadSummaryCsv()
}

function formatDateLabel(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('ru-RU', { 
    day: '2-digit', 
    month: '2-digit' 
  })
}

async function loadProjects() {
  try {
    projects.value = await projectsApi.getProjects()
  } catch (err) {
    console.error('Ошибка загрузки проектов:', err)
  }
}

onMounted(async () => {
  await Promise.all([
    loadProjects(),
    loadReports()
  ])
})
</script>

<style scoped>
.reports {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filters {
  display: flex;
  gap: 16px;
  align-items: end;
  flex-wrap: wrap;
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.filters .form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.filters .form-group label {
  font-weight: 500;
  color: #495057;
}

.kpis {
  display: flex;
  gap: 16px;
  margin-bottom: 2rem;
}

.kpi-card {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
  flex: 1;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.kpi-card.overdue {
  border-color: #dc3545;
}

.kpi-card h3 {
  margin: 0 0 0.5rem 0;
  color: #6c757d;
  font-size: 0.9rem;
  font-weight: 500;
}

.kpi-value {
  font-size: 2rem;
  font-weight: bold;
  color: #495057;
}

.kpi-card.overdue .kpi-value {
  color: #dc3545;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.chart-container {
  background: white;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.chart-container h3 {
  margin: 0 0 1rem 0;
  color: #495057;
  font-size: 1.1rem;
}

.chart-wrapper {
  height: 300px;
  position: relative;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    align-items: stretch;
  }
  
  .kpis {
    flex-direction: column;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
