<template>
  <AppHeader />
  <el-container style="height: calc(100vh - 60px); background-color: #f0f2f5;">
    <!-- 左侧聊天会话列表 -->
    <el-aside width="300px" style="background-color: #fff; border-right: 1px solid #e8e8e8;">
      <div class="logo-container">
        <h3 style="color: #409eff; margin: 0;">AI Chat</h3>
        <div class="connection-status">
          <el-badge 
            :value="wsConnected ? '已连接' : '未连接'" 
            :type="wsConnected ? 'success' : 'danger'"
            class="status-badge"
          />
        </div>
      </div>

      <!-- 新建会话按钮 -->
      <div class="new-chat-container">
        <el-button type="primary" @click="createNewChat" :loading="creatingChat" style="width: 100%;">
          <el-icon><Plus /></el-icon>
          新建对话
        </el-button>
      </div>

      <div v-if="sessionsLoading" class="loading-container">
        <el-skeleton :rows="5" animated />
      </div>
      <el-scrollbar v-else height="calc(100vh - 180px)">
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
                <span class="chat-title">对话 {{ chat.id }}</span>
                <span class="chat-time">{{ formatChatTime(chat.created_at) }}</span>
              </div>
            </template>
          </el-menu-item>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <!-- 右侧聊天窗口 -->
    <el-main style="position: relative; overflow: hidden; padding: 0;">
      <div class="chat-window-container">
        <!-- 聊天窗口头部 -->
        <div class="chat-header">
          <h4 style="color: #333; margin: 0;">AI 助手</h4>
          <el-tag v-if="!selectedChatID" type="info" size="small">新对话</el-tag>
        </div>

        <!-- 聊天消息列表 -->
        <div class="message-list-container" ref="messageContainer">
          <div v-if="messagesLoading" class="loading-container">
            <el-skeleton :rows="5" animated />
          </div>
          <div v-else class="message-list-inner">
            <div 
              v-for="(msg, index) in messages" 
              :key="index" 
              :class="['message-item', { 'sent': !msg.is_reply, 'received': msg.is_reply }]"
            >
              <div class="avatar-wrapper" v-if="msg.is_reply">
                <el-avatar 
                  :size="40"
                  style="background-color: #409eff;"
                >
                  AI
                </el-avatar>
              </div>
              <div class="message-bubble-wrapper">
                <div class="message-content-wrapper">
                  <div class="message-content">
                    <div v-if="msg.content" v-html="formatMessage(msg.content)"></div>
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
                <div class="time-stamp">{{ formatTime(msg.created_time) }}</div>
              </div>
              <div class="avatar-wrapper" v-if="!msg.is_reply">
                <el-avatar 
                  :size="40"
                  :src="user.avatarURL"
                  style="background-color: #67c23a;"
                >
                  我
                </el-avatar>
              </div>
            </div>
          </div>
        </div>

        <!-- 消息输入区域 -->
        <div class="input-area-container" ref="inputContainer">
          <div class="resize-handle" @mousedown="startResize"></div>
          <el-input
            v-model="newMessageContent"
            type="textarea"
            placeholder="输入消息..."
            :autosize="{ minRows: 1, maxRows: 10 }"
            @keydown.enter.exact.prevent="handleEnterKey"
            @keydown.enter.shift.exact="handleShiftEnter"
            class="message-textarea"
          />
          <div class="send-button-container">
            <el-button type="primary" @click="sendMessage" :loading="sendingMessage" class="send-button">
              发送
            </el-button>
          </div>
        </div>
      </div>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onBeforeUnmount } from 'vue';
import { ElMessage } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
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
        username: parsedUserInfo.username || user.value.username,
        avatarURL: parsedUserInfo.avatarURL || user.value.avatarURL,
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

// 格式化消息内容（支持换行）
const formatMessage = (content: string): string => {
  return content.replace(/\n/g, '<br>');
};

// 使用工具函数获取图片URL（已在上面导入）

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
.logo-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 64px;
  padding: 0 16px;
  background-color: #fff;
  border-bottom: 1px solid #e8e8e8;
}

.status-badge {
  font-size: 12px;
}

.new-chat-container {
  padding: 16px;
  background-color: #fff;
  border-bottom: 1px solid #e8e8e8;
}

.loading-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px 20px;
  color: #666;
}

.chat-session-menu {
  border-right: none;
}

.chat-session-item {
  height: 60px;
  line-height: 60px;
}

.chat-item-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  line-height: 1.5;
}

.chat-title {
  font-weight: 500;
  color: #333;
  font-size: 14px;
}

.chat-time {
  font-size: 12px;
  color: #8c8c8c;
}

.chat-window-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  background-color: #fff;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e8e8e8;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.message-list-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background-color: #f0f2f5;
  display: flex;
  flex-direction: column;
}

.message-list-inner {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.message-item {
  display: flex;
  align-items: flex-end;
  margin-bottom: 15px;
  max-width: 70%;
}

.message-item.sent {
  align-self: flex-end;
  flex-direction: row;
}

.message-item.received {
  align-self: flex-start;
  flex-direction: row-reverse;
}

.message-bubble-wrapper {
  display: flex;
  flex-direction: column;
  max-width: 100%;
}

.message-content-wrapper {
  display: flex;
  max-width: 100%;
}

.message-content {
  padding: 12px 16px;
  border-radius: 20px;
  line-height: 1.5;
  font-size: 14px;
  word-break: break-word;
}

.message-item.received .message-content {
  background-color: #fff;
  color: #333;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  border-bottom-left-radius: 2px;
}

.message-item.sent .message-content {
  background-color: #409eff;
  color: #fff;
  border-bottom-right-radius: 2px;
}

.message-images {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.message-image {
  width: 100px;
  height: 100px;
  border-radius: 8px;
  cursor: pointer;
}

.avatar-wrapper {
  margin: 0 10px;
}

.time-stamp {
  font-size: 10px;
  color: #999;
  margin-top: 4px;
  margin-left: 10px;
}

.message-item.sent .time-stamp {
  margin-right: 10px;
  margin-left: 0;
  text-align: right;
}

.input-area-container {
  padding: 16px 24px;
  background-color: #fff;
  border-top: 1px solid #e8e8e8;
  position: relative;
  min-height: 60px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.resize-handle {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background-color: #d9d9d9;
  cursor: ns-resize;
  transition: background-color 0.2s;
}

.resize-handle:hover {
  background-color: #409eff;
}

.message-textarea {
  flex: 1;
  resize: none;
}

.send-button-container {
  display: flex;
  justify-content: flex-end;
}

.send-button {
  background-color: #409eff;
  border-color: #409eff;
}

.send-button:hover {
  background-color: #66b1ff;
  border-color: #66b1ff;
}

.empty-chat-window {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>

