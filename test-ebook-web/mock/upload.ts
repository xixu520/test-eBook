import { MockMethod } from 'vite-plugin-mock'

export default [
  // 模拟文件上传
  {
    url: '/api/v1/upload',
    method: 'post',
    timeout: 2000,
    response: () => {
      return {
        code: 200,
        message: '上传成功',
        data: {
          id: Math.floor(Math.random() * 10000),
          filename: 'uploaded_file.pdf',
        },
      }
    },
  },
  // 获取 OCR 任务列表
  {
    url: '/api/v1/ocr/tasks',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: [
          { id: 1, name: 'GB 50007-2011 地基基础设计规范', status: 'completed', progress: 100, time: '2024-03-22 10:00' },
          { id: 2, name: 'GB 50010-2010 混凝土结构设计规范', status: 'processing', progress: 45, time: '2024-03-22 12:45' },
          { id: 3, name: 'JGJ 94-2008 建筑桩基技术规范', status: 'pending', progress: 0, time: '2024-03-22 12:50' },
          { id: 4, name: '损坏的文档.pdf', status: 'failed', progress: 15, error: '文件格式损坏', time: '2024-03-22 09:30' },
        ],
      }
    },
  },
] as MockMethod[]
