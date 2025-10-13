import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'

export function usePermissions() {
  const auth = useAuthStore()

  const can = computed(() => ({
    // Проекты
    createProject: auth.role === 'manager',
    viewProjects: true, // Все роли могут просматривать проекты
    
    // Дефекты
    createDefect: ['engineer', 'manager'].includes(auth.role),
    editDefect: ['engineer', 'manager'].includes(auth.role),
    updateDefectStatus: ['engineer', 'manager'].includes(auth.role),
    viewDefects: true, // Все роли могут просматривать дефекты
    
    // Комментарии и вложения
    addComment: ['engineer', 'manager'].includes(auth.role),
    uploadAttachment: ['engineer', 'manager'].includes(auth.role),
    
    // Отчеты
    downloadReport: ['manager', 'lead', 'viewer'].includes(auth.role),
    viewReports: ['manager', 'lead', 'viewer'].includes(auth.role),
    
    // Назначение исполнителей (только менеджер)
    assignDefect: auth.role === 'manager',
    
    // Управление пользователями (только менеджер)
    manageUsers: auth.role === 'manager'
  }))

  const isReadOnly = computed(() => auth.role === 'viewer' || auth.role === 'lead')

  return { 
    can,
    isReadOnly,
    role: auth.role,
    isManager: auth.isManager,
    isEngineer: auth.isEngineer,
    isLead: auth.isLead,
    isViewer: auth.isViewer
  }
}
