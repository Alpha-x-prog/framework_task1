import http from './http'

// Auth API
export const authApi = {
  register: (credentials) => http.post('/auth/register', credentials).then(res => res.data),
  login: (credentials) => http.post('/auth/login', credentials).then(res => res.data),
  me: () => http.get('/api/me').then(res => res.data)
}

// Projects API
export const projectsApi = {
  getProjects: () => http.get('/api/projects').then(res => res.data),
  createProject: (data) => http.post('/api/projects', data).then(res => res.data)
}

// Defects API
export const defectsApi = {
  getDefects: (params = {}) => {
    const queryParams = new URLSearchParams()
    Object.keys(params).forEach(key => {
      if (params[key] !== null && params[key] !== undefined && params[key] !== '') {
        queryParams.append(key, params[key])
      }
    })
    return http.get(`/api/defects?${queryParams}`).then(res => res.data)
  },
  createDefect: (data) => http.post('/api/defects', data).then(res => res.data),
  getComments: (defectId) => http.get(`/api/defects/${defectId}/comments`).then(res => res.data),
  addComment: (defectId, data) => http.post(`/api/defects/${defectId}/comments`, data).then(res => res.data),
  uploadAttachment: (defectId, file) => {
    const formData = new FormData()
    formData.append('file', file)
    return http.post(`/api/defects/${defectId}/attachments`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    }).then(res => res.data)
  },
  updateStatus: (defectId, statusId) => http.patch(`/api/defects/${defectId}/status`, { status_id: statusId })
}

// References API
export const refsApi = {
  getStatuses: () => http.get('/api/refs/statuses').then(res => res.data),
  getRoles: () => http.get('/api/refs/roles').then(res => res.data)
}

// Reports API
export const reportsApi = {
  // JSON-сводка: total, by_status, by_priority, overdue
  getSummary: (params = {}) => {
    const queryParams = new URLSearchParams()
    Object.keys(params).forEach(key => {
      if (params[key] !== null && params[key] !== undefined && params[key] !== '') {
        queryParams.append(key, params[key])
      }
    })
    return http.get(`/api/reports/summary?${queryParams}`).then(res => res.data)
  },
  
  // Тренды по времени: series с bucket + счётчики по статусам
  getTrends: (params = {}) => {
    const queryParams = new URLSearchParams()
    Object.keys(params).forEach(key => {
      if (params[key] !== null && params[key] !== undefined && params[key] !== '') {
        queryParams.append(key, params[key])
      }
    })
    return http.get(`/api/reports/trends?${queryParams}`).then(res => res.data)
  },
  
  // CSV-отчёт
  downloadSummaryCsv: async () => {
    const response = await http.get('/api/reports/summary.csv', {
      responseType: 'blob'
    })
    
    const blob = new Blob([response.data], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'summary.csv'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  }
}

// Lightweight helpers (flat exports) requested by UI pieces
// Users list with optional params, e.g. role=engineer
export const getUsers = (params = {}) => http.get('/api/users', { params })

// Defects helpers mirroring above but returning raw axios response
export const getDefects = (params = {}) => http.get('/api/defects', { params })
export const createDefect = (payload) => http.post('/api/defects', payload)