<template>
  <AppHeader />
  <el-container class="chat-container">
    <!-- 左侧聊天会话列表 -->
    <el-aside width="320px" class="chat-sidebar">
      <div class="sidebar-content">
        <div class="sidebar-header">
          <div class="sidebar-brand">
            <el-icon :size="28" color="#409EFF"><ChatDotRound /></el-icon>
            <h3 class="brand-title">AI 对话</h3>
          </div>
          <div class="connection-status">
            <el-tag 
              :type="wsConnected ? 'success' : 'danger'"
              effect="light"
              size="small"
              round
              class="status-tag"
            >
              <el-icon :size="12" style="margin-right: 4px;">
                <CircleCheckFilled v-if="wsConnected" />
                <CircleCloseFilled v-else />
              </el-icon>
              {{ wsConnected ? '已连接' : '未连接' }}
            </el-tag>
          </div>
        </div>

        <!-- 新建会话按钮 -->
        <div class="new-chat-container">
          <el-button 
            type="primary" 
            @click="createNewChat" 
            :loading="creatingChat" 
            class="new-chat-btn"
          >
            <el-icon><Plus /></el-icon>
            新建对话
          </el-button>
        </div>

        <!-- 会话列表区域 -->
        <div class="session-list-wrapper">
          <div v-if="sessionsLoading" class="loading-container">
            <el-skeleton :rows="5" animated />
          </div>
          <el-scrollbar v-else class="session-scrollbar">
            <el-menu 
              :default-active="selectedChatID?.toString()" 
              class="chat-session-menu"
              @select="handleChatSelect"
            >
              <el-menu-item 
                v-for="chat in chatSessions" 
                :key="chat.id" 
                :index="chat.id.toString()"
                class="chat-session-item"
              >
                <template #title>
                  <div class="chat-item-content">
                    <el-icon class="chat-icon"><ChatLineRound /></el-icon>
                    <div class="chat-info">
                      <span class="chat-title">{{ formatChatTitle(chat) }}</span>
                    </div>
                  </div>
                </template>
              </el-menu-item>
            </el-menu>
          </el-scrollbar>
        </div>
      </div>
    </el-aside>

    <!-- 右侧聊天窗口 -->
    <el-main class="chat-main">
      <div class="chat-window-container">
        <!-- 聊天窗口头部 -->
        <div class="chat-header">
          <div class="chat-header-left">
            <el-avatar 
              :size="36"
              style="background: linear-gradient(135deg, #409EFF, #66B3FF); margin-right: 12px;"
            >
              <el-icon><Robot /></el-icon>
            </el-avatar>
            <div>
              <h4 class="chat-title">AI 智能助手</h4>
              <p class="chat-subtitle">基于多模态能力的智能对话系统</p>
            </div>
          </div>
          <div class="chat-header-right">
            <el-tag v-if="!selectedChatID" type="info" size="small" effect="light">新对话</el-tag>
            <el-tag v-else type="success" size="small" effect="light">对话 {{ selectedChatID }}</el-tag>
          </div>
        </div>

        <!-- 聊天消息列表 -->
        <div class="message-list-container" ref="messageContainer">
          <div v-if="messagesLoading" class="loading-container">
            <el-skeleton :rows="5" animated />
          </div>
          <div v-else-if="messages.length === 0" class="empty-chat-window">
            <div class="empty-content">
              <el-icon :size="64" color="#409EFF" class="empty-icon"><ChatLineSquare /></el-icon>
              <h3 class="empty-title">开始新的对话</h3>
              <p class="empty-subtitle">输入消息开始与AI助手聊天，支持发送图片进行多模态对话</p>
            </div>
          </div>
          <div v-else class="message-list-inner">
            <div 
              v-for="(msg, index) in messages" 
              :key="index" 
              :class="['message-item', { 'sent': !msg.is_reply, 'received': msg.is_reply }]"
            >
              <div class="avatar-wrapper" v-if="msg.is_reply">
                <el-avatar 
                  :size="42"
                  style="background: linear-gradient(135deg, #409EFF, #66B3FF); box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);"
                >
                  <el-icon><Robot /></el-icon>
                </el-avatar>
              </div>
              <div class="message-bubble-wrapper">
                <div class="message-content-wrapper">
                  <div class="message-content">
                    <div v-if="msg.content" v-html="formatMessage(msg.content)" class="message-text"></div>
                    <!-- 显示图片 -->
                    <div v-if="msg.images && msg.images.length > 0" class="message-images">
                      <el-image
                        v-for="(img, imgIndex) in msg.images"
                        :key="imgIndex"
                        :src="getImageURL(img.image_path)"
                        :preview-src-list="[getOriginImageURL(img.image_path)]"
                        fit="cover"
                        class="message-image"
                        :preview-teleported="true"
                      />
                    </div>
                  </div>
                </div>
                <div class="message-meta">
                  <span class="time-stamp">{{ formatTime(msg.created_time) }}</span>
                </div>
              </div>
              <div class="avatar-wrapper" v-if="!msg.is_reply">
                <el-avatar 
                  :size="42"
                  style="background: linear-gradient(135deg, #67C23A, #85CE61); box-shadow: 0 2px 8px rgba(103, 194, 58, 0.3);"
                >
                  {{ user.username.charAt(0).toUpperCase() }}
                </el-avatar>
              </div>
            </div>
          </div>
        </div>

        <!-- 消息输入区域 -->
        <div class="input-area-container" ref="inputContainer">
          <div class="resize-handle" @mousedown="startResize"></div>
          <div class="input-wrapper">
            <el-input
              v-model="newMessageContent"
              type="textarea"
              placeholder="请输入消息..."
              :autosize="{ minRows: 1, maxRows: 10 }"
              @keydown.enter.exact.prevent="handleEnterKey"
              @keydown.enter.shift.exact="handleShiftEnter"
              class="message-textarea"
            />
            <div class="send-button-container">
              <el-button 
                type="primary" 
                @click="sendMessage" 
                :loading="sendingMessage" 
                :disabled="!wsConnected || !newMessageContent.trim()"
                class="send-button"
              >
                <el-icon v-if="!sendingMessage"><Promotion /></el-icon>
                <span>{{ sendingMessage ? '发送中...' : '发送' }}</span>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onBeforeUnmount } from 'vue';
import { ElMessage } from 'element-plus';
import { 
  Plus, 
  ChatDotRound, 
  ChatLineRound, 
  ChatLineSquare, 
  Promotion,
  CircleCheckFilled,
  CircleCloseFilled
} from '@element-plus/icons-vue';
import { MCPAPI, backendURL } from '../api/index';
import { getImageUrl as getImageURL, getOriginImageUrl as getOriginImageURL } from '../utils/url';
import AppHeader from '../components/AppHeader.vue';

// 接口定义
interface ChatSession {
  id: number;
  created_at: string;
}

interface Message {
  content: string;
  is_reply: boolean;
  created_time: string;
  images?: Array<{
    image_path: string;
    [key: string]: any;
  }>;
}

interface UserInfo {
  username: string;
  avatarURL: string;
  id: number;
}

interface ChatRequest {
  is_new_chat: boolean;
  chat_id: number;
  content: string;
}

interface ChatResponse {
  chat_id: number;
  reply: string;
  images: Array<{
    image_path: string;
    [key: string]: any;
  }>;
}

// 状态定义
const user = ref<UserInfo>({
  username: '',
  avatarURL: '',
  id: 0
});

const chatSessions = ref<ChatSession[]>([]);
const selectedChatID = ref<number | null>(null);
const messages = ref<Message[]>([]);
const newMessageContent = ref<string>('');

const ws = ref<WebSocket | null>(null);
const messageContainer = ref<HTMLDivElement | null>(null);
const inputContainer = ref<HTMLDivElement | null>(null);

// 加载状态
const sessionsLoading = ref<boolean>(false);
const messagesLoading = ref<boolean>(false);
const wsConnected = ref<boolean>(false);
const sendingMessage = ref<boolean>(false);
const creatingChat = ref<boolean>(false);

// 输入框调整高度相关状态
const isResizing = ref<boolean>(false);
const startY = ref<number>(0);
const startHeight = ref<number>(0);

// 从localStorage中加载用户信息
const loadUserInfo = () => {
  const userInfoString = localStorage.getItem('userInfo');
  if (userInfoString) {
    try {
      const parsedUserInfo = JSON.parse(userInfoString);
      user.value = {
        username: parsedUserInfo.username || '用户',
        avatarURL: parsedUserInfo.avatarURL || '',
        id: parsedUserInfo.id || 0
      };
    } catch (e) {
      console.error('解析用户信息失败', e);
    }
  }
};

// 格式化时间
const formatTime = (time: string): string => {
  const date = new Date(time);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

// 格式化聊天时间
const formatChatTime = (time: string): string => {
  const date = new Date(time);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  
  if (days === 0) {
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  } else if (days === 1) {
    return '昨天';
  } else if (days < 7) {
    return `${days}天前`;
  } else {
    return date.toLocaleDateString();
  }
};

// 格式化聊天标题
const formatChatTitle = (chat: ChatSession): string => {
  return `对话 ${chat.id}`;
};

// 格式化消息内容（支持换行）
const formatMessage = (content: string): string => {
  return content.replace(/\n/g, '<br>');
};

// 滚动到底部
const scrollToBottom = () => {
  if (messageContainer.value) {
    nextTick(() => {
      messageContainer.value!.scrollTop = messageContainer.value!.scrollHeight;
    });
  }
};

// 开始调整输入框高度
const startResize = (e: MouseEvent) => {
  e.preventDefault();
  isResizing.value = true;
  startY.value = e.clientY;
  if (inputContainer.value) {
    startHeight.value = inputContainer.value.offsetHeight;
  }
  
  document.addEventListener('mousemove', handleResize);
  document.addEventListener('mouseup', stopResize);
  document.body.style.cursor = 'ns-resize';
  document.body.style.userSelect = 'none';
};

// 处理拖动调整
const handleResize = (e: MouseEvent) => {
  if (!isResizing.value || !inputContainer.value) return;
  
  const deltaY = startY.value - e.clientY;
  const newHeight = startHeight.value + deltaY;
  const maxHeight = window.innerHeight * 0.25;
  const minHeight = 60;
  
  const clampedHeight = Math.max(minHeight, Math.min(newHeight, maxHeight));
  inputContainer.value.style.height = `${clampedHeight}px`;
};

// 停止调整
const stopResize = () => {
  isResizing.value = false;
  document.removeEventListener('mousemove', handleResize);
  document.removeEventListener('mouseup', stopResize);
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
};

// 获取聊天会话列表
const fetchChatSessions = async () => {
  sessionsLoading.value = true;
  try {
    const response = await MCPAPI.getChats();
    if (response.data.code === 200) {
      chatSessions.value = response.data.data?.chats || [];
      // 按创建时间倒序排列
      chatSessions.value.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    } else {
      ElMessage.error('获取会话列表失败');
    }
  } catch (error) {
    ElMessage.error('获取会话列表失败');
    console.error('Error fetching chat sessions:', error);
  } finally {
    sessionsLoading.value = false;
  }
};

// 获取消息列表
const fetchMessages = async (chatID: number) => {
  messagesLoading.value = true;
  try {
    const response = await MCPAPI.getChatMessages(chatID);
    if (response.data.code === 200) {
      // 后端返回的是按时间倒序（最新的在前），需要反转以便显示时最新的在底部
      const messagesData = response.data.data || [];
      messages.value = messagesData.reverse();
      nextTick(() => {
        scrollToBottom();
      });
    } else {
      ElMessage.error('获取消息失败');
    }
  } catch (error) {
    ElMessage.error('获取消息失败');
    console.error('Error fetching messages:', error);
  } finally {
    messagesLoading.value = false;
  }
};

// 创建新对话
const createNewChat = async () => {
  creatingChat.value = true;
  try {
    // 先清空当前选中的会话和消息
    selectedChatID.value = null;
    messages.value = [];
    
    // 刷新会话列表，新创建的会话会在用户发送第一条消息时自动创建
    await fetchChatSessions();
    
    // 提示用户可以开始输入消息
    ElMessage.success('新对话已准备就绪，请输入消息开始聊天');
  } catch (error) {
    ElMessage.error('创建对话失败');
    console.error('Error creating chat:', error);
  } finally {
    creatingChat.value = false;
  }
};

// 选择聊天会话
const handleChatSelect = (chatID: string) => {
  const id = parseInt(chatID);
  if (id && id !== selectedChatID.value) {
    selectedChatID.value = id;
    fetchMessages(id);
  }
};

// 发送消息
const sendMessage = async () => {
  if (!newMessageContent.value.trim()) {
    ElMessage.warning('请输入消息内容');
    return;
  }
  
  if (!ws.value || ws.value.readyState !== WebSocket.OPEN) {
    ElMessage.error('WebSocket连接未建立，无法发送消息');
    return;
  }
  
  const content = newMessageContent.value.trim();
  newMessageContent.value = '';
  
  // 如果当前没有选中的会话，创建新会话
  const isNewChat = !selectedChatID.value;
  
  // 立即在本地界面显示用户消息
  const userMessage: Message = {
    content: content,
    is_reply: false,
    created_time: new Date().toISOString()
  };
  messages.value.push(userMessage);
  nextTick(() => {
    scrollToBottom();
  });
  
  sendingMessage.value = true;
  
  try {
    const request: ChatRequest = {
      is_new_chat: isNewChat,
      chat_id: selectedChatID.value || 0,
      content: content
    };
    
    ws.value.send(JSON.stringify(request));
  } catch (error) {
    ElMessage.error('发送消息失败，请重试');
    console.error('Error sending message:', error);
    sendingMessage.value = false;
  }
};

// 处理 Enter 键（发送消息）
const handleEnterKey = () => {
  sendMessage();
};

// 处理 Shift+Enter（换行）
const handleShiftEnter = () => {
  // 不做任何处理，让默认行为发生（换行）
};

// 连接WebSocket
const connectWebSocket = () => {
  if (!user.value.id) {
    console.error('用户ID不存在，无法建立WebSocket连接');
    return;
  }
  
  const wsUrl = `${backendURL.replace('http', 'ws')}/api/mcp/chat?userID=${user.value.id}`;
  ws.value = new WebSocket(wsUrl);
  
  ws.value.onopen = () => {
    console.log('WebSocket连接已建立');
    wsConnected.value = true;
    ElMessage.success('已连接到AI助手');
  };
  
  ws.value.onmessage = (event) => {
    try {
      const data: ChatResponse | { error: string } = JSON.parse(event.data);
      
      // 检查是否是错误消息
      if ('error' in data) {
        ElMessage.error(data.error);
        sendingMessage.value = false;
        return;
      }
      
      // 处理正常响应
      const response = data as ChatResponse;
      
      // 更新选中的会话ID（如果是新对话或收到第一条回复）
      if (response.chat_id) {
        // 如果是新对话（selectedChatID 为 null），或者 chat_id 与当前不同，则更新
        if (!selectedChatID.value || response.chat_id !== selectedChatID.value) {
          selectedChatID.value = response.chat_id;
          // 刷新会话列表
          fetchChatSessions();
        }
      }
      
      // 添加AI回复到消息列表
      const aiMessage: Message = {
        content: response.reply,
        is_reply: true,
        created_time: new Date().toISOString(),
        images: response.images || []
      };
      
      messages.value.push(aiMessage);
      sendingMessage.value = false;
      
      nextTick(() => {
        scrollToBottom();
      });
      
    } catch (error) {
      console.error('解析WebSocket消息失败:', error);
      sendingMessage.value = false;
    }
  };
  
  ws.value.onerror = (error) => {
    console.error('WebSocket连接错误:', error);
    ElMessage.error('WebSocket连接失败');
    wsConnected.value = false;
    sendingMessage.value = false;
  };
  
  ws.value.onclose = (event) => {
    console.log('WebSocket连接已关闭:', event.code, event.reason);
    wsConnected.value = false;
    sendingMessage.value = false;
    
    // 如果不是正常关闭，尝试重连
    if (event.code !== 1000) {
      console.log('WebSocket连接异常关闭，3秒后尝试重连...');
      setTimeout(() => {
        if (user.value.id) {
          connectWebSocket();
        }
      }, 3000);
    }
  };
};

onMounted(async () => {
  loadUserInfo();
  await fetchChatSessions();
  connectWebSocket();
});

onBeforeUnmount(() => {
  if (ws.value) {
    ws.value.close();
  }
  // 清理事件监听器
  document.removeEventListener('mousemove', handleResize);
  document.removeEventListener('mouseup', stopResize);
});

// 监听 selectedChatID 变化，如果变化了，滚动到底部
watch(selectedChatID, () => {
  nextTick(() => {
    scrollToBottom();
  });
});

</script>

<style scoped>
.chat-container {
  height: calc(100vh - 70px);
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

/* 左侧边栏样式 */
.chat-sidebar {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(64, 158, 255, 0.1);
  box-shadow: 2px 0 12px rgba(64, 158, 255, 0.08);
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 80px;
  min-height: 80px;
  padding: 0 24px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  flex-shrink: 0;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-title {
  font-size: 18px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.status-tag {
  font-weight: 500;
  border: none;
}

.new-chat-container {
  padding: 16px 20px;
  background: rgba(255, 255, 255, 0.95);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  flex-shrink: 0;
}

.new-chat-btn {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.new-chat-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

.session-list-wrapper {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.session-scrollbar {
  flex: 1;
  padding: 10px;
  overflow: hidden;
}

.chat-session-menu {
  border-right: none;
  background: transparent;
  width: 100%;
}

.chat-session-item {
  height: auto;
  min-height: 56px;
  margin: 1px 0;
  border-radius: 12px;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  padding: 10px 16px !important;
}

:deep(.chat-session-item .el-menu-item-content) {
  width: 100%;
  display: flex;
  align-items: center;
}

:deep(.chat-session-item .el-menu-item) {
  padding: 0 !important;
  height: auto;
}

.chat-item-content {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.chat-icon {
  color: #409EFF;
  font-size: 20px;
  flex-shrink: 0;
}

.chat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
  flex: 1;
  min-width: 0;
}

.chat-title {
  font-size: 14px;
  font-weight: 600;
  color: #2c3e50;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
}

.chat-time {
  font-size: 12px;
  color: #7f8c8d;
  white-space: nowrap;
}

:deep(.chat-session-item:hover),
:deep(.chat-session-item.is-active) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-color: rgba(64, 158, 255, 0.2);
  transform: translateX(4px);
}

:deep(.chat-session-item.is-active) {
  background: linear-gradient(135deg, #e6f3ff, #d9edff);
  border-color: #409EFF;
}

/* 右侧聊天窗口样式 */
.chat-main {
  position: relative;
  overflow: hidden;
  padding: 0;
  background: rgba(255, 255, 255, 0.95);
}

.chat-window-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 聊天头部 */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30px;
  height: 80px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.08);
}

.chat-header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.chat-title {
  font-size: 20px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 4px 0;
}

.chat-subtitle {
  font-size: 12px;
  color: #7f8c8d;
  margin: 0;
}

.chat-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* 消息列表容器 */
.message-list-container {
  flex: 1;
  overflow-y: auto;
  padding: 30px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
}

.message-list-inner {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  gap: 20px;
}

/* 空状态样式 */
.empty-chat-window {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.empty-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  text-align: center;
  max-width: 400px;
}

.empty-icon {
  opacity: 0.8;
}

.empty-title {
  font-size: 22px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0;
}

.empty-subtitle {
  font-size: 14px;
  color: #7f8c8d;
  margin: 0;
  line-height: 1.6;
}

/* 消息项样式 */
.message-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-item.sent {
  flex-direction: row-reverse;
}

.message-bubble-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 70%;
}

.message-content-wrapper {
  display: flex;
  max-width: 100%;
}

.message-content {
  padding: 16px 20px;
  border-radius: 20px;
  line-height: 1.6;
  font-size: 15px;
  word-break: break-word;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.message-item.received .message-content {
  background: white;
  color: #333;
  border-bottom-left-radius: 4px;
  border: 1px solid rgba(64, 158, 255, 0.1);
}

.message-item.sent .message-content {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  color: white;
  border-bottom-right-radius: 4px;
}

.message-text {
  line-height: 1.6;
}

/* 图片消息样式 */
.message-images {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 12px;
}

.message-image {
  width: 120px;
  height: 120px;
  border-radius: 12px;
  cursor: pointer;
  transition: transform 0.3s ease;
  border: 2px solid rgba(64, 158, 255, 0.1);
}

.message-image:hover {
  transform: scale(1.05);
  border-color: #409EFF;
}

/* 消息元信息 */
.message-meta {
  display: flex;
  justify-content: flex-end;
  margin-top: 6px;
  padding: 0 8px;
}

.time-stamp {
  font-size: 11px;
  color: #7f8c8d;
  font-weight: 400;
}

/* 输入区域样式 */
.input-area-container {
  padding: 24px 30px;
  background: white;
  border-top: 1px solid rgba(64, 158, 255, 0.1);
  position: relative;
  min-height: 120px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: 0 -2px 12px rgba(64, 158, 255, 0.08);
}

.resize-handle {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 6px;
  background: linear-gradient(135deg, #e6f3ff, #d9edff);
  cursor: ns-resize;
  transition: all 0.3s ease;
  border-radius: 3px;
  opacity: 0.6;
}

.resize-handle:hover {
  height: 8px;
  opacity: 1;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
}

:deep(.message-textarea .el-textarea__inner) {
  background: #f8fbff;
  border: 2px solid #e6f3ff;
  border-radius: 16px;
  padding: 16px 20px;
  font-size: 15px;
  line-height: 1.6;
  transition: all 0.3s ease;
  resize: none;
  box-shadow: none;
}

:deep(.message-textarea .el-textarea__inner:hover),
:deep(.message-textarea .el-textarea__inner:focus) {
  border-color: #409EFF;
  background: white;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

:deep(.message-textarea .el-textarea__inner::placeholder) {
  color: #a0b4c0;
}

.send-button-container {
  display: flex;
  justify-content: flex-end;
}

.send-button {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 12px;
  padding: 12px 32px;
  font-weight: 600;
  font-size: 15px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}

.send-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

.send-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* 加载状态样式 */
.loading-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-container {
    flex-direction: column;
  }
  
  .chat-sidebar {
    width: 100% !important;
    height: 300px;
    border-right: none;
    border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  }
  
  .sidebar-content {
    height: 100%;
  }
  
  .session-list-wrapper {
    flex: 1;
    min-height: 0;
  }
  
  .chat-main {
    height: calc(100% - 300px);
  }
  
  .message-list-container {
    padding: 20px;
  }
  
  .chat-header,
  .sidebar-header {
    padding: 0 20px;
  }
  
  .input-area-container {
    padding: 20px;
  }
}

@media (max-width: 480px) {
  .message-content {
    max-width: 85%;
  }
  
  .chat-header-left {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .chat-title {
    font-size: 18px;
  }
  
  .chat-subtitle {
    display: none;
  }
}
</style>