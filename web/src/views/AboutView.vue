<script setup>
import Card from '@/components/ui/Card.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardTitle from '@/components/ui/CardTitle.vue'
import CardContent from '@/components/ui/CardContent.vue'
</script>

<template>
  <div class="container max-w-3xl py-10 space-y-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-bold tracking-tight">关于本项目</h1>
      <p class="text-muted-foreground">
        了解更多关于 Lemwood Mirror 的信息。
      </p>
    </div>

    <Card class="bg-background/50 backdrop-blur-sm border-white/10">
      <CardHeader>
        <CardTitle>Lemwood Mirror</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <p>
          Lemwood Mirror 是一个自托管的镜像服务，旨在为 Minecraft 社区提供稳定、快速的启动器和相关工具下载。
          本项目自动追踪并下载 GitHub 上的最新发布版本，确保用户能够第一时间获取到最新的软件更新。
        </p>
        <p>
          主要功能包括：
        </p>
        <ul class="list-disc list-inside space-y-1 text-muted-foreground ml-4">
          <li>自动同步 GitHub Release</li>
          <li>提供直观的版本历史浏览</li>
          <li>详细的下载统计数据</li>
          <li>RESTful API 支持</li>
        </ul>
      </CardContent>
    </Card>

    <Card class="bg-background/50 backdrop-blur-sm border-white/10">
        <CardHeader>
            <CardTitle>关于作者</CardTitle>
        </CardHeader>
        <CardContent>
             <p>本项目由开源社区成员维护。如果您对本项目感兴趣，或者想提供反馈，欢迎通过 GitHub 联系。</p>
        </CardContent>
    </Card>
  </div>
</template>
