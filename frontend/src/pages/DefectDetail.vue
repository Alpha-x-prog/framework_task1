<template>
  <div class="container">
    <div v-if="loading && !defect" class="loading">
      Загрузка дефекта...
    </div>
    
    <div v-else-if="!defect" class="loading">
      Дефект не найден
    </div>
    
    <div v-else>
      <!-- Заголовок дефекта -->
      <div class="card">
        <div class="card-header">
          <h1 class="card-title">
            #{{ defect.id }} - {{ defect.title }}
          </h1>
          <span :class="`status status-${getStatusClass(defect.status_id)}`">
            {{ getStatusName(defect.status_id) }}
          </span>
        </div>
        
        <div v-if="defect.description" style="margin-bottom: 1rem;">
          <h4>Описание:</h4>
          <p>{{ defect.description }}</p>
        </div>
        
        <div class="form-row">
          <div>
            <strong>Проект:</strong> {{ getProjectName(defect.project_id) }}
          </div>
          <div>
            <strong>Приоритет:</strong> 
            <span :class="`priority priority-${defect.priority}`">
              {{ defect.priority }}
            </span>
          </div>
          <div v-if="defect.due_date">
            <strong>Срок:</strong> {{ formatDate(defect.due_date) }}
          </div>
          <div>
            <strong>Создан:</strong> {{ formatDate(defect.created_at) }}
          </div>
        </div>
      </div>
      
      <!-- Продвижение статуса (только для инженеров и менеджеров) -->
      <div v-if="can.updateDefectStatus" class="card">
        <h3>Сменить статус</h3>
        <div class="form-row">
          <div class="form-group" style="display: flex; align-items: center; gap: 12px;">
            <div>
              <label>Текущий статус</label>
              <div>
                <span :class="`status status-${getStatusClass(defect.status_id)}`">
                  {{ getStatusName(defect.status_id) }}
                </span>
              </div>
            </div>
            <button
              type="button"
              class="btn btn-secondary"
              :disabled="statusLoading || !nextStatusId"
              @click="handleAdvance"
              title="Перевести на следующий этап"
            >
              {{ nextStatusLabel }}
            </button>
          </div>
        </div>

        <div v-if="statusMessage" :class="`alert ${statusMessage.type}`">
          {{ statusMessage.text }}
        </div>
      </div>
      
      <!-- Комментарии -->
      <div class="card">
        <h3>Комментарии</h3>
        
        <!-- Форма добавления комментария (только для инженеров и менеджеров) -->
        <form v-if="can.addComment" @submit.prevent="handleAddComment" style="margin-bottom: 2rem;">
          <div class="form-group">
            <label for="comment-body">Добавить комментарий</label>
            <textarea
              id="comment-body"
              v-model="newComment.body"
              class="form-control"
              rows="3"
              required
              :disabled="commentLoading"
              placeholder="Введите комментарий..."
            ></textarea>
          </div>
          <button 
            type="submit" 
            class="btn btn-primary"
            :disabled="commentLoading"
          >
            {{ commentLoading ? 'Добавление...' : 'Добавить комментарий' }}
          </button>
        </form>
        
        <!-- Список комментариев -->
        <div v-if="commentsLoading" class="loading">
          Загрузка комментариев...
        </div>
        
        <div v-else-if="comments.length === 0" class="loading">
          Комментариев пока нет
        </div>
        
        <div v-else>
          <div v-for="comment in comments" :key="comment.id" class="comment">
            <div class="comment-header">
              <span>Автор: {{ commentLogin(comment) }}</span>
              <span>{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="comment-body">{{ comment.body }}</div>
          </div>
        </div>
      </div>
      
      <!-- Вложения -->
      <div class="card">
        <h3>Вложения</h3>
        
        <!-- Форма загрузки файла (только для инженеров и менеджеров) -->
        <form v-if="can.uploadAttachment" @submit.prevent="handleFileUpload" style="margin-bottom: 2rem;">
          <div class="form-group">
            <label for="file-upload">Загрузить файл</label>
            <input
              id="file-upload"
              ref="fileInput"
              type="file"
              class="form-control"
              :disabled="fileLoading"
              @change="handleFileSelect"
            />
          </div>
          <button 
            type="submit" 
            class="btn btn-primary"
            :disabled="fileLoading || !selectedFile"
          >
            {{ fileLoading ? 'Загрузка...' : 'Загрузить файл' }}
          </button>
        </form>
        
        <div v-if="fileMessage" :class="`alert ${fileMessage.type}`">
          {{ fileMessage.text }}
        </div>
        
        <div v-if="attachments.length > 0" class="attachments-grid">
          <div v-for="a in attachments" :key="a.id" class="attachment-item">
            <div class="attachment-meta">
              <div class="attachment-name" :title="a.file_path">{{ shortFileName(a.file_path) }}</div>
              <div class="attachment-sub">ID: {{ a.id }} · {{ a.mime || 'file' }} · {{ formatDate(a.created_at) }}</div>
            </div>
            <div class="attachment-actions">
              <button class="btn btn-secondary" @click="download(a)">Скачать</button>
            </div>
          </div>
        </div>
        <div v-else class="loading">Файлов пока нет</div>
      </div>
      
      <!-- Скачать отчет (для менеджеров, руководителей и наблюдателей) -->
      <div v-if="can.downloadReport" class="card">
        <h3>Отчеты</h3>
        <button @click="handleDownloadCsv" class="btn btn-secondary">
          Скачать отчет (CSV)
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { defectsApi, projectsApi, refsApi, reportsApi } from '../api'
import { loginFromEmail } from '../utils/login'
import { usePermissions } from '../composables/usePermissions'

export default {
  name: 'DefectDetail',
  setup() {
    const route = useRoute()
    const { can } = usePermissions()
    
    const defect = ref(null)
    const projects = ref([])
    const statuses = ref([])
    const comments = ref([])
    const attachments = ref([])
    
    const loading = ref(false)
    const commentsLoading = ref(false)
    const commentLoading = ref(false)
    const statusLoading = ref(false)
    const fileLoading = ref(false)
    
    // ручной выбор статуса убран
    
    const newComment = ref({
      body: ''
    })
    
    const selectedFile = ref(null)
    const fileInput = ref(null)
    
    const statusMessage = ref(null)
    const fileMessage = ref(null)
    
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
      if (name === 'new' || name.includes('new')) return 'new'
      if (name === 'in_work' || name.includes('in_work') || name.includes('work') || name.includes('progress')) return 'in-progress'
      if (name === 'review' || name.includes('review')) return 'review'
      if (name === 'closed' || name.includes('closed')) return 'closed'
      if (name === 'canceled' || name.includes('cancel')) return 'canceled'
      return 'unknown'
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      return new Date(dateString).toLocaleString('ru-RU')
    }

    // Порядок продвижения статусов
    const statusOrder = ['new', 'in_work', 'review', 'closed']

    const currentStatusName = computed(() => {
      if (!defect.value) return ''
      const s = statuses.value.find(x => x.id === defect.value.status_id)
      return s ? s.name : ''
    })

    const nextStatusId = computed(() => {
      const curName = currentStatusName.value
      if (!curName) return null
      const idx = statusOrder.indexOf(curName)
      if (idx === -1 || idx >= statusOrder.length - 1) return null
      const nextName = statusOrder[idx + 1]
      const next = statuses.value.find(x => x.name === nextName)
      return next ? next.id : null
    })

    const nextStatusLabel = computed(() => {
      if (!nextStatusId.value) return 'Дальше'
      const s = statuses.value.find(x => x.id === nextStatusId.value)
      return s ? `Дальше → ${s.name}` : 'Дальше'
    })
    
    const loadDefect = async () => {
      loading.value = true
      
      try {
        // Загружаем дефект через список дефектов с фильтром по ID
        const defects = await defectsApi.getDefects({ id: route.params.id })
        if (defects.length > 0) {
          defect.value = defects[0]
        }
      } catch (err) {
        console.error('Ошибка загрузки дефекта:', err)
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
    
    const loadComments = async () => {
      if (!defect.value) return
      
      commentsLoading.value = true
      
      try {
        comments.value = await defectsApi.getComments(defect.value.id)
      } catch (err) {
        console.error('Ошибка загрузки комментариев:', err)
      } finally {
        commentsLoading.value = false
      }
    }

    const loadAttachments = async () => {
      if (!defect.value) return
      try {
        attachments.value = await defectsApi.listAttachments(defect.value.id)
      } catch (err) {
        console.error('Ошибка загрузки вложений:', err)
      }
    }
    
    // ручное обновление статуса удалено

    const handleAdvance = async () => {
      if (!defect.value || !nextStatusId.value) return
      statusLoading.value = true
      statusMessage.value = null
      try {
        await defectsApi.updateStatus(defect.value.id, nextStatusId.value)
        // Обновляем локально для мгновенного UI
        defect.value.status_id = nextStatusId.value
        statusMessage.value = { type: 'alert-success', text: 'Статус продвинут' }
        // Перезагружаем страницу, чтобы гарантированно синхронизировать всё состояние
        // (это также устранит временные ошибки рендера)
        window.location.reload()
      } catch (err) {
        statusMessage.value = {
          type: 'alert-error',
          text: err.response?.data?.message || 'Не удалось обновить статус'
        }
      } finally {
        statusLoading.value = false
      }
    }
    
    const handleAddComment = async () => {
      if (!defect.value || !newComment.value.body.trim()) return
      
      commentLoading.value = true
      
      try {
        await defectsApi.addComment(defect.value.id, { body: newComment.value.body })
        
        newComment.value.body = ''
        await loadComments()
      } catch (err) {
        alert(err.response?.data?.message || 'Ошибка добавления комментария')
      } finally {
        commentLoading.value = false
      }
    }
    
    const handleFileSelect = (event) => {
      selectedFile.value = event.target.files[0]
    }
    
    const handleFileUpload = async () => {
      if (!defect.value || !selectedFile.value) return
      
      fileLoading.value = true
      fileMessage.value = null
      
      try {
        const response = await defectsApi.uploadAttachment(defect.value.id, selectedFile.value)
        
        attachments.value.push(response)
        
        fileMessage.value = {
          type: 'alert-success',
          text: `Файл успешно загружен: ${response.file_path}`
        }
        
        selectedFile.value = null
        if (fileInput.value) {
          fileInput.value.value = ''
        }
      } catch (err) {
        fileMessage.value = {
          type: 'alert-error',
          text: err.response?.data?.message || 'Ошибка загрузки файла'
        }
      } finally {
        fileLoading.value = false
      }
    }
    
    const handleDownloadCsv = async () => {
      try {
        await reportsApi.downloadSummaryCsv()
      } catch (err) {
        alert(err.response?.data?.message || 'Ошибка скачивания отчета')
      }
    }
    
    onMounted(async () => {
      await Promise.all([
        loadProjects(),
        loadStatuses(),
        loadDefect()
      ])
      
      if (defect.value) {
        await Promise.all([
          loadComments(),
          loadAttachments()
        ])
      }
    })
    
    return {
      defect,
      projects,
      statuses,
      comments,
      attachments,
      loading,
      commentsLoading,
      commentLoading,
      statusLoading,
      fileLoading,
      newComment,
      selectedFile,
      fileInput,
      statusMessage,
      fileMessage,
      nextStatusId,
      nextStatusLabel,
      getProjectName,
      getStatusName,
      getStatusClass,
      formatDate,
      commentLogin: (c) => c.author_email ? loginFromEmail(c.author_email) : `ID: ${c.author_id}`,
      // handleStatusUpdate removed
      handleAdvance,
      handleAddComment,
      handleFileSelect,
      handleFileUpload,
      handleDownloadCsv,
      loadAttachments,
      shortFileName: (p) => {
        if (!p) return ''
        const idx = p.lastIndexOf('/')
        const idx2 = p.lastIndexOf('\\')
        const cut = Math.max(idx, idx2)
        return cut >= 0 ? p.slice(cut + 1) : p
      },
      download: async (a) => {
        try {
          const resp = await defectsApi.downloadAttachment(a.id)
          const blob = new Blob([resp.data], { type: a.mime || 'application/octet-stream' })
          const url = URL.createObjectURL(blob)
          const link = document.createElement('a')
          link.href = url
          link.download = (a.file_path && (a.file_path.split('/').pop() || a.file_path.split('\\').pop())) || `attachment-${a.id}`
          document.body.appendChild(link)
          link.click()
          document.body.removeChild(link)
          URL.revokeObjectURL(url)
        } catch (err) {
          alert(err.response?.data?.message || 'Ошибка скачивания файла')
        }
      },
      can
    }
  }
}
</script>
