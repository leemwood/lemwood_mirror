<script setup>
import { Megaphone, ExternalLink, Sparkles } from 'lucide-vue-next'
import Card from '@/components/ui/Card.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardTitle from '@/components/ui/CardTitle.vue'
import CardContent from '@/components/ui/CardContent.vue'

const commonLinks = [
  { label: '赞助网站', url: 'https://zs.lemwood.cn/', color: 'text-pink-500', desc: '支持我们的发展' },
  { label: 'MT管理器下载', url: 'https://mt2.cn', color: 'text-green-500', desc: '强大的文件管理工具' },
  { label: 'Shizuku 官网', url: 'https://shizuku.rikka.app/zh-hans/', color: 'text-blue-500', desc: '系统API调试工具' },
];

const zeinkLink = {
    label: 'ZeinkLab 泽客来宾',
    url: 'https://zeinklab.com',
    desc: '泽客来宾博客网是林梦泽（化名）搭建的一个个人分享日常、技术和游戏的综合内容网站，投身科技，用代码丰富生活，用数据填满热心。',
    isAd: true
}
</script>

<template>
  <div class="grid gap-6 md:grid-cols-2">
      <!-- Site Announcements -->
      <Card class="h-full bg-background/50 backdrop-blur-sm border-white/10 shadow-sm hover:shadow-md transition-all">
        <CardHeader class="pb-3">
          <CardTitle class="flex items-center gap-2 text-lg">
            <Megaphone class="h-5 w-5 text-primary" />
            <span>常用链接</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div class="grid gap-3">
            <a
              v-for="(link, index) in commonLinks"
              :key="index"
              :href="link.url"
              target="_blank"
              class="group flex items-center gap-3 rounded-lg border border-transparent p-2 hover:border-border hover:bg-muted/50 transition-all"
            >
              <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted group-hover:bg-background transition-colors">
                  <ExternalLink class="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="font-medium text-sm truncate group-hover:text-primary transition-colors">{{ link.label }}</span>
                <span class="text-xs text-muted-foreground truncate">{{ link.desc }}</span>
              </div>
            </a>
          </div>
        </CardContent>
      </Card>

      <!-- Featured / Ad -->
      <Card class="h-full bg-gradient-to-br from-background/50 to-primary/5 backdrop-blur-sm border-primary/20 shadow-sm hover:shadow-md transition-all relative overflow-hidden">
        <div class="absolute top-0 right-0 p-2">
            <span class="text-[10px] uppercase tracking-wider text-muted-foreground/50 font-bold border border-border px-1.5 py-0.5 rounded">Featured</span>
        </div>
        <CardHeader class="pb-3">
          <CardTitle class="flex items-center gap-2 text-lg">
            <Sparkles class="h-5 w-5 text-amber-500" />
            <span>特别推荐</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
           <a :href="zeinkLink.url" target="_blank" class="block group h-full">
               <div class="flex flex-col h-full gap-2">
                   <h3 class="font-semibold text-primary group-hover:underline underline-offset-4 decoration-primary/50">{{ zeinkLink.label }}</h3>
                   <p class="text-sm text-muted-foreground leading-relaxed line-clamp-4">
                       {{ zeinkLink.desc }}
                   </p>
                   <div class="mt-auto pt-2 flex items-center text-xs text-muted-foreground group-hover:text-primary transition-colors">
                       访问网站 <ExternalLink class="ml-1 h-3 w-3" />
                   </div>
               </div>
           </a>
        </CardContent>
      </Card>
  </div>
</template>