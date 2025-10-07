<template>
  <div class="download-report">
    <button 
      @click="handleDownload"
      class="btn btn-primary download-btn"
      :disabled="loading"
      :class="{ 'loading': loading }"
    >
      <span v-if="loading" class="spinner"></span>
      {{ loading ? '⏳ Скачиваю…' : '📊 Скачать отчёт (CSV)' }}
    </button>
  </div>
</template>

<script>
import { ref } from 'vue'
import { reportsApi } from '../api'

export default {
  name: 'DownloadReport',
  setup() {
    const loading = ref(false)
    
    const handleDownload = async () => {
      loading.value = true
      
      try {
        await reportsApi.downloadSummaryCsv()
      } catch (error) {
        console.error('Ошибка скачивания отчёта:', error)
        alert('Не удалось скачать отчёт')
      } finally {
        loading.value = false
      }
    }
    
    return {
      loading,
      handleDownload
    }
  }
}
</script>

<style scoped>
.download-report {
  margin: 20px 0;
}

.download-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  background-color: #007bff;
  color: white;
}

.download-btn:hover:not(:disabled) {
  background-color: #0056b3;
  transform: translateY(-1px);
}

.download-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.download-btn.loading {
  background-color: #6c757d;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
