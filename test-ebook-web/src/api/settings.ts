import request from '@/utils/request'

export interface SystemSetting {
  id?: number
  key: string
  value: string
  description?: string
}

export function getSettings() {
  return request<SystemSetting[]>({
    url: '/settings',
    method: 'get'
  })
}

export function saveSettings(data: any) {
  // Convert nested object to flattened Key-Value list for backend storage
  const settingsList: SystemSetting[] = []
  
  const flatten = (obj: any, prefix = '') => {
    for (const key in obj) {
      const fullKey = prefix ? `${prefix}.${key}` : key
      if (typeof obj[key] === 'object' && !Array.isArray(obj[key]) && obj[key] !== null) {
        flatten(obj[key], fullKey)
      } else {
        settingsList.push({
          key: fullKey,
          value: typeof obj[key] === 'object' ? JSON.stringify(obj[key]) : String(obj[key])
        })
      }
    }
  }
  
  flatten(data)
  
  return request({
    url: '/settings',
    method: 'put',
    data: settingsList
  })
}

export function testOcrConnection(data: any) {
  return request({
    url: '/settings/ocr-test',
    method: 'post',
    data: {
      api_key: data.baidu_config?.api_key || data.paddle_config?.token || '',
      secret_key: data.baidu_config?.secret_key || ''
    }
  })
}
