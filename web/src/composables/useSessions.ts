import type { Conversation } from 'ant-design-x-vue'
import type { MessageInfo } from 'ant-design-x-vue'
import { computed, ref } from 'vue'

export interface ChatSession {
  key: string
  label: string
  /** 毫秒时间戳，用来算「今天 / 昨天 / 更早」 */
  timestamp: number
  messages: MessageInfo<string>[]
}

const GROUP_ORDER = ['今天', '昨天', '更早']

/** 每个新会话都带同一句开场白，和原来 App.vue 的 defaultMessages 一致 */
export function createGreeting(): MessageInfo<string>[] {
  return [
    {
      id: 'welcome',
      message: '你好，我是 xatu-agent。',
      status: 'success',
    },
  ]
}

function createSession(): ChatSession {
  return {
    // 这个 key 以后原样作为 session_id 传给后端
    key: crypto.randomUUID(),
    label: '新会话',
    timestamp: Date.now(),
    messages: createGreeting(),
  }
}

function groupOf(timestamp: number): string {
  const start = new Date()
  start.setHours(0, 0, 0, 0)
  const day = 24 * 60 * 60 * 1000
  if (timestamp >= start.getTime()) return '今天'
  if (timestamp >= start.getTime() - day) return '昨天'
  return '更早'
}

// 和 useAuth 一样放在模块级：退出登录时要能清掉，不能做成每个组件各一份
const sessions = ref<ChatSession[]>([createSession()])
const activeKey = ref(sessions.value[0]?.key ?? '')

export function useSessions() {
  // 只把列表需要的字段交给 Conversations，消息留在 sessions 里
  const items = computed<Conversation[]>(() =>
    sessions.value.map((session) => ({
      key: session.key,
      label: session.label,
      timestamp: session.timestamp,
      group: groupOf(session.timestamp),
    })),
  )

  function get(key: string) {
    return sessions.value.find((session) => session.key === key)
  }

  function saveMessages(key: string, messages: MessageInfo<string>[]) {
    const session = get(key)
    if (session) session.messages = messages
  }

  function create() {
    const session = createSession()
    // 新会话插到顶部，侧栏第一眼能看到
    sessions.value = [session, ...sessions.value]
    activeKey.value = session.key
    return session
  }

  /**
   * 删除后保证至少还剩一个会话，聊天区不会出现「没有 activeKey」。
   * 返回删除之后应该显示的 key。
   */
  function remove(key: string) {
    const rest = sessions.value.filter((session) => session.key !== key)
    sessions.value = rest.length > 0 ? rest : [createSession()]
    if (activeKey.value === key || !get(activeKey.value)) {
      activeKey.value = sessions.value[0].key
    }
    return activeKey.value
  }

  function rename(key: string, label: string) {
    const session = get(key)
    const next = label.trim()
    if (session && next) session.label = next
  }

  /** 用户发出第一句话时，用这句话当标题，避免侧栏全是「新会话」 */
  function touchTitle(key: string, text: string) {
    const session = get(key)
    if (session && session.label === '新会话') {
      session.label = text.trim().slice(0, 20) || '新会话'
    }
  }

  function reset() {
    const session = createSession()
    sessions.value = [session]
    activeKey.value = session.key
  }

  return {
    sessions,
    activeKey,
    items,
    /** Conversations 的 groupable.sort：按今天、昨天、更早，而不是按分组名字母序 */
    groupSort: (a: string, b: string) => GROUP_ORDER.indexOf(a) - GROUP_ORDER.indexOf(b),
    get,
    saveMessages,
    create,
    remove,
    rename,
    touchTitle,
    reset,
  }
}