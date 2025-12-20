<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" class="text-center">
        <v-icon
          icon="mdi-folder-off-outline"
          size="64"
          color="medium-emphasis"
          class="mb-4"
        ></v-icon>
        <h2 class="text-h4 font-weight-bold mb-2">功能已下架</h2>
        <p class="text-body-1 text-medium-emphasis mb-6">
          抱歉，文件浏览功能目前已停止提供服务。
        </p>
        <v-btn
          color="primary"
          variant="flat"
          @click="$emit('go-home')"
        >
          返回首页
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
defineEmits(['go-home']);
</script>