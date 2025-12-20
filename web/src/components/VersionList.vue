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
        md="4" 
        lg="3"
      >
        <v-hover v-slot="{ isHovering, props }">
          <v-card 
            v-bind="props"
            :elevation="isHovering ? 2 : 0"
            class="h-100 transition-swing rounded d-flex flex-column"
            border
          >
            <div :class="`bg-${item.color}-lighten-5 pa-6 d-flex justify-center align-center position-relative`" style="height: 160px;">
               <v-avatar :color="item.color" size="80">
                  <v-icon size="40" color="white">{{ item.icon }}</v-icon>
               </v-avatar>
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
                {{ item.name }}
              </v-card-title>
              <v-card-subtitle class="text-center mt-1">
                最近更新: {{ formatDate(item.lastUpdated) }}
              </v-card-subtitle>
            </v-card-item>

            <v-spacer></v-spacer>

            <v-card-actions class="pa-4 pt-0 d-flex flex-column gap-2">
              <v-btn 
                block 
                color="primary" 
                variant="flat" 
                size="large" 
                prepend-icon="mdi-download"
                class="rounded"
                :href="getLatestDownloadUrl(item)"
                v-if="item.hasAssets"
              >
                下载最新版
              </v-btn>
              
              <v-btn 
                block 
                variant="tonal" 
                size="large"
                class="rounded ml-0"
                @click="openHistory(item)"
              >
                历史版本
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-hover>
      </v-col>
    </v-row>

    <!-- History Dialog -->
    <v-dialog v-model="historyDialog" max-width="900" scrollable transition="dialog-bottom-transition">
      <v-card class="rounded" v-if="selectedLauncher">
        <v-toolbar color="surface" class="px-2 border-b">
           <v-toolbar-title class="font-weight-bold">
             {{ selectedLauncher.name }} 版本历史
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
    
    return {
      name,
      versions,
      latest: latestVersion,
      lastUpdated: versions.length ? versions[0].published_at : null,
      icon: getIcon(name),
      color: getColor(name),
      hasAssets: latestObj && latestObj.assets && latestObj.assets.length > 0,
      latestObj
    };
  });
});

const loadData = async () => {
  loading.value = true;
  try {
    const [statusRes, latestRes] = await Promise.all([getStatus(), getLatest()]);
    
    const data = statusRes.data;
    // Sort contents
    for (const key in data) {
        data[key].sort((a, b) => String(b.tag_name || b.name).localeCompare(String(a.tag_name || a.name)));
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

const getLatestDownloadUrl = (item) => {
   if (!item.latestObj || !item.latestObj.assets || !item.latestObj.assets.length) return '#';
   // Default to first asset of latest version
   return getAssetUrl(item.name, item.latestObj, item.latestObj.assets[0]);
};

const copyLink = (url) => {
    const fullUrl = url.startsWith('http') ? url : window.location.origin + url;
    navigator.clipboard.writeText(fullUrl).then(() => {
        snackbar.value = true;
    });
};

const getIcon = (name) => {
    const n = name.toLowerCase();
    if (n.includes('hmcl')) return 'mdi-cube-outline';
    if (n.includes('pcl')) return 'mdi-controller';
    if (n.includes('baka')) return 'mdi-ghost';
    if (n.includes('shizuku')) return 'mdi-water';
    if (n.includes('mt')) return 'mdi-file-tree';
    return 'mdi-package-variant';
};

const getColor = (name) => {
    const n = name.toLowerCase();
    if (n.includes('hmcl')) return 'orange';
    if (n.includes('pcl')) return 'light-blue';
    if (n.includes('baka')) return 'deep-purple';
    if (n.includes('shizuku')) return 'cyan';
    if (n.includes('mt')) return 'green';
    return 'blue-grey';
};

onMounted(() => {
    loadData();
});

defineExpose({ refresh: loadData });
</script>

<style scoped>
.gap-2 { gap: 8px; }
</style>