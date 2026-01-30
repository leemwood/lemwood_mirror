<script setup>
import { ref, onMounted } from 'vue';
import { Cookie } from 'lucide-vue-next';
import Button from '@/components/ui/Button.vue';

const isOpen = ref(false);

onMounted(() => {
  const consented = localStorage.getItem('cookies-consented');
  if (!consented) {
    isOpen.value = true;
  }
});

const accept = () => {
  localStorage.setItem('cookies-consented', 'true');
  isOpen.value = false;
};
</script>

<template>
  <div v-if="isOpen" class="fixed bottom-0 left-0 right-0 z-50 p-4 animate-in slide-in-from-bottom-full duration-500">
    <div class="mx-auto max-w-4xl rounded-xl border bg-background/80 backdrop-blur-md p-4 shadow-lg flex flex-col sm:flex-row items-center gap-4 justify-between">
      <div class="flex items-center gap-3">
        <div class="rounded-full bg-primary/10 p-2 text-primary">
          <Cookie class="h-6 w-6" />
        </div>
        <div class="text-sm text-muted-foreground">
          <p class="font-medium text-foreground">Cookie 使用提示</p>
          <p>本网站使用 Cookies 以优化您的体验。继续浏览即表示您同意我们使用 Cookies。</p>
        </div>
      </div>
      <div class="flex gap-2">
        <Button size="sm" @click="accept">
          我知道了
        </Button>
      </div>
    </div>
  </div>
</template>
