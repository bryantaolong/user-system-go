<template>
  <div class="not-found">
    <h1>404</h1>
    <p>抱歉，您访问的页面不存在。</p>
    <a-button type="primary" @click="$router.push('/')">返回首页</a-button>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
}

.not-found h1 {
  font-size: 120px;
  color: #409eff;
  margin-bottom: 16px;
}

.not-found p {
  font-size: 18px;
  color: #606266;
  margin-bottom: 24px;
}
</style>
