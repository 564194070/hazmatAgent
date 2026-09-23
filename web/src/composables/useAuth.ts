import { computed, ref } from 'vue'

// 和业务 key 分开，避免和其他前端项目撞 localStorage
const TOKEN_KEY = 'xatu.token'
const USER_KEY = 'xatu.userId'

// 刷新页面时从 localStorage 还原。没有就保持 null，守卫会把人送去登录页
const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
const userId = ref<string | null>(localStorage.getItem(USER_KEY))

// 与 Gin 返回的 user 字段一致（Go 默认 JSON 名，没有 json tag）
export type ApiUser = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  Username: string
  Password: string
  Nickname: string
  Email: string
  Status: number
  Role: string
}

type MsgResp = {
  msg?: string
}

type LoginResp = {
  token: string
  user: ApiUser
}

async function postJson<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  const data = (await res.json().catch(() => ({}))) as T & MsgResp
  if (!res.ok) {
    throw new Error(data.msg || '请求失败')
  }
  return data
}

export function useAuth() {
  const isLoggedIn = computed(() => Boolean(token.value && userId.value))

  async function register(username: string, password: string, nickname: string) {
    const name = username.trim()
    if (!name || !password) {
      throw new Error('请输入用户名和密码')
    }
    await postJson<MsgResp>('/api/user/register', {
      username: name,
      password,
      nickname: nickname.trim(),
    })
  }

  async function login(username: string, password: string) {
    const name = username.trim()
    if (!name || !password) {
      throw new Error('请输入用户名和密码')
    }

    const data = await postJson<LoginResp>('/api/user/login', {
      username: name,
      password,
    })
    token.value = data.token
    userId.value = String(data.user.ID)
    localStorage.setItem(TOKEN_KEY, token.value)
    localStorage.setItem(USER_KEY, userId.value)
  }

  function logout() {
    token.value = null
    userId.value = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return { token, userId, isLoggedIn, register, login, logout }
}
