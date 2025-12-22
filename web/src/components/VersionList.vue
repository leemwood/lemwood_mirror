<template>
  <div class="mb-6">
    <div class="d-flex align-center justify-space-between mb-6">
       <div>
         <h2 class="text-h4 font-weight-bold">版本探索</h2>
         <p class="text-medium-emphasis mt-1">发现并下载最新的启动器组件</p>
       </div>
       <v-btn icon="mdi-refresh" variant="text" @click="loadData" :loading="loading"></v-btn>
    </div>
    
    <div v-if="loading" class="d-flex justify-center pa-12">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <div v-else-if="!launcherList.length" class="text-center text-medium-emphasis pa-12">
      <v-icon size="64" class="mb-4">mdi-package-variant-closed</v-icon>
      <div>暂无数据</div>
    </div>

    <v-row v-else>
      <v-col 
        v-for="item in launcherList" 
        :key="item.name" 
        cols="12" 
        sm="6"
      >
        <v-hover v-slot="{ isHovering, props }">
          <v-card 
            v-bind="props"
            :elevation="isHovering ? 2 : 0"
            class="h-100 transition-swing rounded d-flex flex-column"
            border
          >
            <div class="logo-background-container pa-6 d-flex justify-center align-center position-relative" style="height: 160px; overflow: hidden;">
              <img 
                :src="item.logoUrl" 
                class="logo-background"
                alt=""
              />
              <img 
                :src="item.logoUrl" 
                class="logo-foreground"
                :alt="item.displayName"
              />
              <v-chip 
                v-if="item.latest" 
                color="success" 
                size="small" 
                variant="flat" 
                class="position-absolute top-0 right-0 ma-4 font-weight-bold"
              >
                {{ item.latest }}
              </v-chip>
            </div>

            <v-card-item class="pt-4">
              <v-card-title class="text-h6 font-weight-bold text-center">
                {{ item.displayName }}
              </v-card-title>
              <v-card-subtitle class="text-center mt-1">
                最近更新: {{ formatDate(item.lastUpdated) }}
              </v-card-subtitle>
            </v-card-item>

            <v-spacer></v-spacer>
            
            <v-card-actions class="pa-4 pt-0 ma-auto">
              <v-btn-group class="w-100" density="comfortable">
                <v-btn 
                  v-if="item.hasAssets"
                  color="primary" 
                  variant="flat"
                  size="large"
                  class="flex-1"
                  :href="item.latestDownloadUrl"
                  prepend-icon="mdi-download"
                >
                  下载最新版
                </v-btn>
                
                <v-btn 
                  variant="tonal"
                  size="large"
                  class="flex-1"
                  @click="openHistory(item)"
                  prepend-icon="mdi-history"
                >
                  历史版本
                </v-btn>
                
                <v-btn 
                  v-if="item.hasAssets"
                  color="primary"
                  variant="outlined"
                  size="large"
                  class="flex-1"
                  @click="copyLink(item.latestDownloadUrl)"
                  prepend-icon="mdi-content-copy"
                >
                  复制链接
                </v-btn>
              </v-btn-group>
            </v-card-actions>
          </v-card>
        </v-hover>
      </v-col>
    </v-row>

    <v-dialog v-model="historyDialog" max-width="900" scrollable transition="dialog-bottom-transition">
      <v-card class="rounded" v-if="selectedLauncher">
        <v-toolbar color="surface" class="px-2 border-b">
           <v-toolbar-title class="font-weight-bold">
             {{ selectedLauncher.displayName }} 版本历史
           </v-toolbar-title>
           <v-spacer></v-spacer>
           <v-btn icon="mdi-close" variant="text" @click="historyDialog = false"></v-btn>
        </v-toolbar>

        <v-card-text class="pa-4 bg-surface-light">
           <v-row>
             <v-col v-for="v in selectedLauncher.versions" :key="v.tag_name || v.name" cols="12" md="6">
               <v-card variant="flat" class="border rounded h-100">
                 <v-card-item>
                   <template v-slot:title>
                     <div class="d-flex align-center justify-space-between">
                       <span class="text-subtitle-1 font-weight-bold">{{ v.tag_name || v.name }}</span>
                       <v-chip v-if="selectedLauncher.latest === (v.tag_name || v.name)" color="success" size="x-small">LATEST</v-chip>
                     </div>
                   </template>
                   <template v-slot:subtitle>
                     {{ formatDate(v.published_at) }}
                   </template>
                 </v-card-item>
                 
                 <v-divider class="mx-4"></v-divider>
                 
                 <v-list density="compact" class="bg-transparent py-2">
                    <v-list-item v-for="asset in v.assets" :key="asset.name" active-color="primary">
                       <template v-slot:prepend>
                          <v-icon icon="mdi-file-outline" size="small" class="text-medium-emphasis"></v-icon>
                       </template>
                       <v-list-item-title class="text-body-2 font-weight-medium">
                          {{ asset.name }}
                       </v-list-item-title>
                       <template v-slot:append>
                          <v-btn 
                            icon="mdi-download" 
                            size="small" 
                            variant="text" 
                            color="primary" 
                            :href="getAssetUrl(selectedLauncher.name, v, asset)"
                            :download="asset.name"
                          ></v-btn>
                          <v-btn 
                            icon="mdi-content-copy" 
                            size="small" 
                            variant="text" 
                            color="medium-emphasis"
                            @click="copyLink(getAssetUrl(selectedLauncher.name, v, asset))"
                          ></v-btn>
                       </template>
                    </v-list-item>
                 </v-list>
               </v-card>
             </v-col>
           </v-row>
        </v-card-text>
      </v-card>
    </v-dialog>
    
    <v-snackbar v-model="snackbar" :timeout="2000" color="success" rounded>
       链接已复制
    </v-snackbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { getStatus, getLatest } from '../services/api';

const LAUNCHER_INFO_MAP = {
  'zl': { displayName: 'ZalithLauncher', logoUrl: 'https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/34c1ec9e07f826df.webp' },
  'zl2': { displayName: 'ZalithLauncher2', logoUrl: 'https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/ee0028bd82493eb3.webp' },
  'hmcl': { displayName: 'HMCL', logoUrl: 'https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/3835841e4b9b7abf.jpeg' },
  'MG': { displayName: 'MobileGlues', logoUrl: 'https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/3625548d2639a024.png' },
  'fcl': { displayName: 'FoldCraftLauncher', logoUrl: 'https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/dc5e0ee14d8f54f0.png' },
  'shizuku': { displayName: 'Shizuku', logoUrl: 'https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/f7067665f073b4cc.png' }
};

const rawLaunchers = ref({});
const latestMap = ref({});
const loading = ref(true);
const historyDialog = ref(false);
const selectedLauncher = ref(null);
const snackbar = ref(false);

const launcherList = computed(() => {
  return Object.keys(rawLaunchers.value).map(name => {
    const versions = rawLaunchers.value[name];
    const latestVersion = latestMap.value[name];
    const latestObj = versions.find(v => (v.tag_name || v.name) === latestVersion) || versions[0];
    const info = LAUNCHER_INFO_MAP[name] || { displayName: name, logoUrl: 'https://cdn.zeinklab.com/gh/JanePHPDev/Blog-Static-Resource@main/images/b4ee27d31312bdb9.svg' };
    
    const latestDownloadUrl = latestObj && latestObj.assets && latestObj.assets.length > 0
      ? getAssetUrl(name, latestObj, latestObj.assets[0])
      : '#';
    
    return {
      name,
      displayName: info.displayName,
      logoUrl: info.logoUrl,
      versions,
      latest: latestVersion,
      lastUpdated: versions.length ? versions[0].published_at : null,
      hasAssets: latestObj && latestObj.assets && latestObj.assets.length > 0,
      latestObj,
      latestDownloadUrl
    };
  });
});

const loadData = async () => {
  loading.value = true;
  try {
    const [statusRes, latestRes] = await Promise.all([getStatus(), getLatest()]);
    
    const data = statusRes.data;
    for (const key in data) {
        data[key].sort((a, b) => String(b.tag_name || b.name).localeCompare(String(a.tag_name || b.name)));
    }
    rawLaunchers.value = data;
    latestMap.value = latestRes.data;
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const openHistory = (item) => {
  selectedLauncher.value = item;
  historyDialog.value = true;
};

const formatDate = (dateStr) => {
    if (!dateStr) return '未知时间';
    return new Date(dateStr).toLocaleDateString();
};

const getAssetUrl = (launcherName, version, asset) => {
     if (asset.url && (asset.url.startsWith('http://') || asset.url.startsWith('https://'))) {
        return asset.url;
    }
    return `/download/${launcherName}/${version.tag_name || version.name}/${asset.name}`;
};

const copyLink = (url) => {
    const fullUrl = url.startsWith('http') ? url : window.location.origin + url;
    navigator.clipboard.writeText(fullUrl).then(() => {
        snackbar.value = true;
    });
};

onMounted(() => {
    loadData();
});

defineExpose({ refresh: loadData });
</script>

<style scoped>
.gap-2 { 
  gap: 8px; 
}

.logo-background-container {
  position: relative;
}

.logo-background {
  position: absolute;
  width: 150%;
  height: 150%;
  object-fit: cover;
  filter: blur(20px);
  opacity: 0.3;
  transform: scale(1.2);
}

.logo-foreground {
  position: relative;
  width: 80px;
  height: 80px;
  object-fit: contain;
  z-index: 1;
  border-radius: 12px;
}

@media (max-width: 600px) {
  .logo-foreground {
    width: 60px;
    height: 60px;
  }
  
  .v-btn {
    font-size: 14px;
    padding: 0 8px !important;
  }
  
  .v-btn--density-comfortable {
    height: 40px;
  }
}
</style>
