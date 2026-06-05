import client from './client'

export interface User {
  id: string
  email: string
  is_temp: boolean
  created_at: string
}

export interface HistoryItem {
  id: string
  task_id: string
  text_key: string
  category: string
  duration_seconds: number
  outcome: 'completed' | 'skipped'
  created_at: string
}

export interface Stats {
  total_completed: number
  total_skipped: number
}

export async function getMe(): Promise<User> {
  const { data } = await client.get('/me')
  return data
}

export async function getHistory(): Promise<HistoryItem[]> {
  const { data } = await client.get('/history')
  return data
}

export async function getStats(): Promise<Stats> {
  const { data } = await client.get('/history/stats')
  return data
}

export interface BindResp {
  access_token: string
  refresh_token: string
  expires_in: number
}

export async function bindEmail(email: string, password: string): Promise<BindResp> {
  const { data } = await client.post('/auth/bind', { email, password })
  return data
}
