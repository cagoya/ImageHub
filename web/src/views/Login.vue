<template>
  <div>
    <div class="screen">
      <h2>用户登录</h2>
      <div class="name">
        <label for="name">用户名</label>
        <div class="item">
          <input type="name" name="name" placeholder="请输入用户名" v-model="loginForm.name" />
        </div>
      </div>
  
      <div class="password">
        <label for="password">密码</label>
        <div class="item">
          <input class="pas" :type="passwordFieldType" name="password" placeholder="请输入密码" v-model="loginForm.password" />
          <ion-icon class="show-hide" @click="togglePasswordVisibility"></ion-icon>
        </div>
      </div>
  
      <button class="login" @click="handleLogin">登录</button>

      <div class="footer">
        <span @click="handleRegister">注册</span>
        <span @click="handleForgotPassword">忘记密码?</span>
      </div>
    </div>
    
    <el-dialog
      v-model="dialogVisible"
      title="修改密码"
      width="500"
    >
      <el-form>
        <el-form-item label="用户名">
          <el-input v-model="forgetPasswordForm.username" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="forgetPasswordForm.email" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="forgetPasswordForm.newPassword" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmitForgotPassword">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { UserAPI } from '../api';
  import { ElDialog, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'
  
  const router = useRouter();

  const dialogVisible = ref(false);
  const forgetPasswordForm = ref({
    username: '',
    newPassword: '',
    email: ''
  });

  // 登录表单数据
  const loginForm = ref({
    name: '',
    password: ''
  });
  
  // 密码显示/隐藏状态
  const showPassword = ref(false);
  const passwordFieldType = ref<'password' | 'text'>('password');
  
  // 切换密码的可见性
  const togglePasswordVisibility = () => {
    showPassword.value = !showPassword.value;
    passwordFieldType.value = showPassword.value ? 'text' : 'password';
  };
  
// 处理登录逻辑
const handleLogin = async () => {
  try {
    if (!loginForm.value.name || !loginForm.value.password) {
      ElMessage.error('请输入用户名和密码');
      return;
    }
    const response = await UserAPI.login(loginForm.value.name, loginForm.value.password);
    if (response.status === 200) {
      localStorage.setItem('userToken', response.data.data.token);
      localStorage.setItem('userInfo', JSON.stringify(response.data.data.user));
      localStorage.setItem('userID', response.data.data.user.id);
      ElMessage.success('登录成功');
      router.push('/dashboard');
    }
  } catch (error) {
    ElMessage.error('登录失败: ' + (error instanceof Error ? error.message : '请检查用户名和密码'));
  }
};
  
  // 处理注册链接点击
  const handleRegister = () => {
    router.push('/register');
  };
  
  // 处理忘记密码链接点击
  const handleForgotPassword = () => {
    dialogVisible.value = true;
  };

  // 处理提交修改密码表单
  const handleSubmitForgotPassword = async () => {
    try {
      if (!forgetPasswordForm.value.username || !forgetPasswordForm.value.newPassword || !forgetPasswordForm.value.email) {
        ElMessage.error('请填写完整信息');
        return;
      }
      if (forgetPasswordForm.value.newPassword.length < 6) {
        ElMessage.error('密码长度不能小于6位');
        return;
      }
      const response = await UserAPI.forgetPassword(forgetPasswordForm.value.username, forgetPasswordForm.value.email, forgetPasswordForm.value.newPassword);
      if (response.status === 200) {
        ElMessage.success('修改密码成功');
      }
      dialogVisible.value = false;
    } catch (error) {
      ElMessage.error('修改密码失败: ' + (error instanceof Error ? error.message : '请检查用户名和邮箱是否匹配'));
    }
  };
  
</script>
  
<style scoped>
.screen {
  user-select: none;
  background: #f1f7fe;
  padding: 3em;
  display: flex;
  flex-direction: column;
  border-radius: 30px;
  box-shadow: 0 0 2em #e6e9f9;
  gap: 1.5em;
  min-width: 250px;
  max-width: 350px;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

h2 {
  text-align: center;
  color: #3e4684;
  margin: 0 0 0.5em;
}

.name, .password {
  background: white;
  box-shadow: 0 0 2em #e6e9f9;
  padding: 1em;
  display: flex;
  flex-direction: column;
  gap: 0.5em;
  border-radius: 20px;
  color: #4d4d4d;
}

.item {
  display: flex;
  align-items: center;
}

input {
  outline: none;
  border: none;
  flex-grow: 1;
  padding-left: 0.5em;
}

input::placeholder {
  color: black;
  font-size: 0.9em;
}

ion-icon {
  color: #4d4d4d;
  margin-bottom: -0.2em;
}

.show-hide {
  cursor: pointer;
}

.login {
  padding: 1em;
  background: #3e4684;
  color: white;
  border: none;
  border-radius: 30px;
  font-weight: 600;
  cursor: pointer;
}

.footer {
  display: flex;
  justify-content: space-between;
  font-size: 0.7em;
  color: #5e5e5e;
  gap: 14em;
  padding-bottom: 10em;
}

.footer span {
  cursor: pointer;
}
</style>