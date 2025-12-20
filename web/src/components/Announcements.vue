<template>
  <v-card class="mb-6 rounded-lg" variant="flat" color="surface-light" border>
    <v-card-item>
      <template v-slot:prepend>
         <v-icon color="primary" class="mr-2">mdi-bullhorn-variant</v-icon>
      </template>
      <v-card-title class="text-subtitle-1 font-weight-bold">
        站点公告
      </v-card-title>
    </v-card-item>

    <v-divider class="mx-4"></v-divider>

    <v-card-text class="pt-2">
      <v-row dense>
        <v-col v-for="(link, index) in links" :key="index" cols="12" sm="4">
          <v-btn
            block
            variant="text"
            :href="link.url"
            target="_blank"
            class="justify-start text-none px-2"
            color="on-surface"
          >
            <template v-slot:prepend>
               <v-icon :color="link.color || 'primary'" size="small">mdi-link-variant</v-icon>
            </template>
            <div class="d-flex flex-column align-start ml-1">
               <span class="font-weight-medium">{{ link.label }}</span>
               <span class="text-caption text-medium-emphasis text-truncate" style="max-width: 150px;">
                  {{ link.text || link.url }}
               </span>
            </div>
          </v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
const links = [
  { label: '赞助网站', url: 'https://zs.lemwood.cn/', color: 'pink' },
  { label: 'MT管理器下载', url: 'https://mt2.cn', text: 'mt2.cn', color: 'green' },
  { label: 'Shizuku 官网', url: 'https://shizuku.rikka.app/zh-hans/', color: 'blue' },
];
</script>
