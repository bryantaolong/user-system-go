<template>
  <div class="not-found-container">
    <h1>404</h1>
    <p>抱歉，您访问的页面不存在。</p>
    <el-button type="primary" @click="goHome">返回首页</el-button>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

const goHome = () => {
  router.push('/');
};
</script>

<style scoped>
.not-found-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 80vh; /* 让内容垂直居中 */
  text-align: center;
  font-family: 'Arial', sans-serif;
  color: #333;
}

.not-found-container h1 {
  font-size: 10em;
  margin-bottom: 0.1em;
  color: #409EFF; /* Element Plus primary color */
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.not-found-container p {
  font-size: 1.5em;
  margin-bottom: 2em;
}
</style>