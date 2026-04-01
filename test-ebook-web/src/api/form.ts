import request from '@/utils/request'

export interface FormField {
  ID?: number
  label: string
  field_key: string
  field_type: 'input' | 'select' | 'date' | 'number'
  is_required: boolean
  options?: string // JSON 数组字符串
  order: number
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
