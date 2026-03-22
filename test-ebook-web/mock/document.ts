import { MockMethod } from 'vite-plugin-mock'

interface Document {
  id: number
  standard_no: string
  name: string
  category_id: number
  category_name: string
  publisher: string
  status: 'current' | 'obsolete' | 'upcoming'
  issue_date: string
  implement_date: string
  obsolete_date?: string | null
  version: string
  is_latest: boolean
  ocr_status: 'pending' | 'processing' | 'completed' | 'failed'
  verify_status: 'pending' | 'pass' | 'retry'
  uploader_name: string
  upload_time: string
}

const documentList: Document[] = [
  {
    id: 101,
    standard_no: 'GB 50007-2011',
    name: '建筑地基基础设计规范',
    category_id: 1,
    category_name: '地基基础',
    publisher: '住房和城乡建设部',
    status: 'current',
    issue_date: '2011-07-26',
    implement_date: '2012-08-01',
    version: '2011',
    is_latest: true,
    ocr_status: 'completed',
    verify_status: 'pass',
    uploader_name: 'admin',
    upload_time: '2023-10-01 10:00:00',
  },
  {
    id: 1011,
    standard_no: 'GB 50007-2002',
    name: '建筑地基基础设计规范',
    category_id: 1,
    category_name: '地基基础',
    publisher: '住房和城乡建设部',
    status: 'obsolete',
    issue_date: '2002-02-20',
    implement_date: '2002-09-01',
    version: '2002',
    is_latest: false,
    ocr_status: 'completed',
    verify_status: 'pass',
    uploader_name: 'admin',
    upload_time: '2023-09-15 08:30:00',
  },
  {
    id: 102,
    standard_no: 'GB 50202-2018',
    name: '地基基础工程施工质量验收标准',
    category_id: 1,
    category_name: '地基基础',
    publisher: '住房和城乡建设部',
    status: 'current',
    issue_date: '2018-03-16',
    implement_date: '2018-10-01',
    version: '2018',
    is_latest: true,
    ocr_status: 'completed',
    verify_status: 'pending',
    uploader_name: 'admin',
    upload_time: '2023-10-02 11:30:00',
  },
  {
    id: 1021,
    standard_no: 'GBJ 202-83',
    name: '地基与基础工程施工及验收规范',
    category_id: 1,
    category_name: '地基基础',
    publisher: '中国建筑工业出版社',
    status: 'obsolete',
    issue_date: '1983-12-30',
    implement_date: '1984-06-01',
    version: '1983',
    is_latest: false,
    ocr_status: 'completed',
    verify_status: 'pass',
    uploader_name: 'admin',
    upload_time: '2023-09-20 14:00:00',
  },
  {
    id: 103,
    standard_no: 'GB 50009-2012',
    name: '建筑结构荷载规范',
    category_id: 2,
    category_name: '结构工程',
    publisher: '住房和城乡建设部',
    status: 'current',
    issue_date: '2012-08-01',
    implement_date: '2013-02-01',
    version: '2012',
    is_latest: true,
    ocr_status: 'completed',
    verify_status: 'pass',
    uploader_name: 'admin',
    upload_time: '2023-10-05 09:00:00',
  },
  {
    id: 104,
    standard_no: 'JGJ 130-2011',
    name: '建筑施工扣件式钢管脚手架安全技术规范',
    category_id: 3,
    category_name: '施工技术',
    publisher: '住房和城乡建设部',
    status: 'current',
    issue_date: '2011-06-01',
    implement_date: '2011-12-01',
    version: '2011',
    is_latest: true,
    ocr_status: 'completed',
    verify_status: 'pass',
    uploader_name: 'admin',
    upload_time: '2023-10-06 14:00:00',
  },
]

export default [
  {
    url: '/api/v1/documents',
    method: 'get',
    response: ({ query }: { query: any }) => {
      const { page = 1, size = 10, keyword, category_id, publisher, status } = query

      let filtered = [...documentList]
      if (keyword) {
        const lower = (keyword as string).toLowerCase()
        filtered = filtered.filter(
          (i) =>
            i.standard_no.toLowerCase().includes(lower) ||
            i.name.toLowerCase().includes(lower) ||
            i.publisher.toLowerCase().includes(lower)
        )
      }
      if (category_id) {
        filtered = filtered.filter((i) => i.category_id === Number(category_id))
      }
      if (publisher) {
        filtered = filtered.filter((i) => i.publisher === publisher)
      }
      if (status) {
        filtered = filtered.filter((i) => i.status === status)
      }

      const total = filtered.length
      const list = filtered.slice((Number(page) - 1) * Number(size), Number(page) * Number(size))

      return {
        code: 200,
        message: '获取成功',
        data: {
          total,
          page: Number(page),
          size: Number(size),
          list,
        },
      }
    },
  },
  {
    url: '/api/v1/documents/history',
    method: 'get',
    response: ({ query }: { query: any }) => {
      const { standard_no } = query
      if (!standard_no) {
        return { code: 400, message: '参数错误', data: null }
      }

      const baseNo = (standard_no as string).split('-')[0].trim()
      const history = documentList.filter((d) => d.standard_no.startsWith(baseNo))

      return {
        code: 200,
        message: '获取成功',
        data: history,
      }
    },
  },
  {
    url: '/api/v1/announcements/active',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: {
          content: '系统将于本周五晚进行升级维护，届时将暂停服务约 2 小时。',
          update_time: '2026-03-22 09:00:00',
        },
      }
    },
  },
] as MockMethod[]
