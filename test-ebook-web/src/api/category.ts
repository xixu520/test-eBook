import request from '@/utils/request'

/** 分类数据结构（字段名与后端 GORM gorm.Model 对齐） */
export interface Category {
  ID: number
  name: string
  parent_id: number
  order: number
  doc_count: number
  form_id?: number
  children: Category[]
  CreatedAt?: string
  UpdatedAt?: string
}

/** 新增/编辑分类的请求体 */
export interface CategoryForm {
  name: string
  parent_id: number
  order: number
  form_id?: number
}

/** 获取分类树 */
export function getCategories(): Promise<Category[]> {
  return request({
    url: '/categories',
    method: 'get',
  })
}

/** 新增分类 */
export function addCategory(data: CategoryForm) {
  return request({
    url: '/categories',
    method: 'post',
    data
  })
}

/** 修改分类 */
export function updateCategory(id: number, data: CategoryForm) {
  return request({
    url: `/categories/${id}`,
    method: 'put',
    data
  })
}

/** 删除分类 */
export function deleteCategory(id: number) {
  return request({
    url: `/categories/${id}`,
    method: 'delete',
  })
}
