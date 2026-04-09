import request from '@/utils/request'

export interface FormField {
  ID?: number
  label: string
  field_key: string
  field_type: 'input' | 'select' | 'date' | 'number' | 'checkbox'
  is_required: boolean
  options?: string // 逗号分隔的选项（select/checkbox 类型用）
  order: number
  show_in_home: boolean   // 是否在首页展示（同时决定是否可作为筛选项）
  show_in_filter: boolean // 兼容旧字段，保持与 show_in_home 联动
  show_in_list: boolean   // legacy，不再使用
  show_in_admin: boolean  // legacy，不再使用
  default_value: string
}

export interface Form {
  ID: number
  name: string
  description: string
  fields: FormField[]
}

export function getForms() {
  return request<Form[]>({
    url: '/admin/forms',
    method: 'get'
  })
}

/** 获取全局唯一属性表单（不存在时自动创建） */
export function getGlobalForm() {
  return request<Form>({
    url: '/admin/forms/global',
    method: 'get'
  })
}

export function createForm(data: { name: string; description: string }) {
  return request<Form>({
    url: '/admin/forms',
    method: 'post',
    data
  })
}

export function updateForm(id: number, data: { name: string; description: string }) {
  return request({
    url: `/admin/forms/${id}`,
    method: 'put',
    data
  })
}

export function deleteForm(id: number) {
  return request({
    url: `/admin/forms/${id}`,
    method: 'delete'
  })
}

export function saveFormFields(id: number, fields: FormField[]) {
  return request({
    url: `/admin/forms/${id}/fields`,
    method: 'post',
    data: { fields }
  })
}

export function bindCategoriesToForm(id: number, category_ids: number[]) {
  return request({
    url: `/admin/forms/${id}/categories`,
    method: 'put',
    data: { category_ids }
  })
}
