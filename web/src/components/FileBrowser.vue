<template>
  <v-container fluid class="pa-2">
    <v-row dense>
      <v-col cols="12">
        <v-text-field
          v-model="searchQuery"
          prepend-inner-icon="mdi-magnify"
          label="搜索文件名、版本或启动器"
          variant="outlined"
          density="compact"
          hide-details
          clearable
          class="mb-2"
        />
      </v-col>
    </v-row>

    <v-snackbar v-model="showSnackbar" :color="snackbarColor" timeout="2000" min-height="40">
      {{ snackbarMessage }}
    </v-snackbar>

    <v-row v-if="loading" dense class="mt-4">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate size="32" width="3" />
        <p class="text-caption mt-2">加载中...</p>
      </v-col>
    </v-row>

    <v-tabs v-model="activeTab" v-if="!loading && filteredLaunchers.length" show-arrows class="mb-2">
      <v-tab v-for="launcher in filteredLaunchers" :key="launcher.name" :value="launcher.name" class="text-none">
        <v-icon start>mdi-application</v-icon>
        {{ launcher.name }}
        <v-chip size="x-small" class="ml-2">{{ launcher.totalFiles }}</v-chip>
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" v-if="!loading">
      <v-window-item v-for="launcher in filteredLaunchers" :key="launcher.name" :value="launcher.name">
        <v-row dense>
          <v-col
            v-for="file in launcher.files"
            :key="file.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
            xl="2"
          >
            <v-card
              class="file-card"
              :class="{ 'latest-border': file.isLatest }"
              hover
              density="compact"
            >
              <v-card-item class="py-2">
                <template v-slot:prepend>
                  <v-icon :icon="getFileIcon(file.name)" :color="file.isLatest ? 'primary' : ''" />
                </template>
                <v-card-title class="text-body-2 text-truncate">
                  {{ file.name }}
                </v-card-title>
                <v-card-subtitle class="text-caption">
                  {{ launcher.name }} · {{ file.version }}
                </v-card-subtitle>
              </v-card-item>

              <v-card-text class="py-1">
                <v-chip
                  v-if="file.isLatest"
                  size="x-small"
                  color="success"
                  class="mr-1"
                >
                  Latest
                </v-chip>
                <span class="text-caption text-medium-emphasis">
                  {{ formatDate(file.published_at) }}
                </span>
              </v-card-text>

              <v-card-actions class="pa-2">
                <v-btn
                  :href="file.downloadUrl"
                  target="_blank"
                  variant="flat"
                  size="small"
                  color="primary"
                  prepend-icon="mdi-download"
                  class="flex-grow-1"
                >
                  下载
                </v-btn>
                <v-btn
                  icon="mdi-content-copy"
                  variant="text"
                  size="small"
                  @click="copyUrl(file.downloadUrl)"
                />
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>
    </v-window>

    <v-row v-if="!loading && !filteredLaunchers.length" dense class="mt-8">
      <v-col cols="12" class="text-center">
        <v-icon icon="mdi-folder-search-outline" size="48" color="medium-emphasis" />
        <h3 class="text-h6 mt-2">{{ searchQuery ? '无匹配结果' : '暂无数据' }}</h3>
        <p class="text-caption text-medium-emphasis">
          {{ searchQuery ? '尝试其他关键词' : '检查API连接或刷新重试' }}
        </p>
        <v-btn color="primary" variant="flat" size="small" class="mt-2" @click="loadData">
          刷新
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { getStatus, getLatest } from '../services/api';

const loading = ref(false);
const showSnackbar = ref(false);
const snackbarMessage = ref('');
const snackbarColor = ref('info');
const searchQuery = ref('');
const launchers = ref([]);
const latestData = ref({});
const activeTab = ref('');

const loadData = async () => {
  loading.value = true;
  showSnackbar.value = false;
  
  try {
    const [statusRes, latestRes] = await Promise.all([getStatus(), getLatest()]);
    
    launchers.value = Object.entries(statusRes.data).map(([name, versions]) => ({
      name,
      versions: versions.sort((a, b) => 
        String(b.tag_name || b.name).localeCompare(String(a.tag_name || a.name))
      )
    }));
    
    latestData.value = latestRes.data;
    
    if (launchers.value.length) {
      activeTab.value = launchers.value[0].name;
    }
  } catch (error) {
    snackbarMessage.value = '加载失败，请重试';
    snackbarColor.value = 'error';
    showSnackbar.value = true;
  } finally {
    loading.value = false;
  }
};

const getFileIcon = (filename) => {
  const ext = filename.split('.').pop()?.toLowerCase();
  const icons = {
    exe: 'mdi-application',
    msi: 'msi:mdi-application',
    dmg: 'mdi-apple',
    zip: 'mdi-folder-zip-outline',
    tar: 'mdi-folder-zip-outline',
    gz: 'mdi-folder-zip-outline',
    deb: 'mdi-debian',
    rpm: 'mdi-redhat',
    appimage: 'mdi-application',
    jar: 'mdi-language-java'
  };
  return icons[ext] || 'mdi-file-outline';
};

const formatDate = (dateString) => {
  if (!dateString) return '未知';
  try {
    return new Date(dateString).toLocaleDateString('zh-CN', {
      month: 'short',
      day: 'numeric'
    });
  } catch {
    return dateString;
  }
};

const copyUrl = async (url) => {
  try {
    await navigator.clipboard.writeText(url);
    snackbarMessage.value = '链接已复制';
    snackbarColor.value = 'success';
    showSnackbar.value = true;
  } catch {
    snackbarMessage.value = '复制失败';
    snackbarColor.value = 'error';
    showSnackbar.value = true;
  }
};

const filteredLaunchers = computed(() => {
  const query = searchQuery.value.toLowerCase().trim();
  const result = [];

  launchers.value.forEach(launcher => {
    const files = [];
    
    launcher.versions.forEach(version => {
      const versionName = version.tag_name || version.name;
      const isLatest = latestData.value[launcher.name] === versionName;
      
      version.assets?.forEach(asset => {
        const matchesSearch = !query || 
          launcher.name.toLowerCase().includes(query) ||
          versionName.toLowerCase().includes(query) ||
          asset.name.toLowerCase().includes(query);
        
        if (matchesSearch) {
          files.push({
            id: `${launcher.name}-${versionName}-${asset.name}`,
            name: asset.name,
            version: versionName,
            published_at: version.published_at,
            isLatest,
            launcher: launcher.name,
            downloadUrl: asset.url && asset.url.startsWith('http') 
              ? asset.url 
              : `${window.location.origin}/download/${launcher.name}/${versionName}/${asset.name}`
          });
        }
      });
    });

    if (files.length) {
      result.push({
        name: launcher.name,
        files: files.sort((a, b) => b.isLatest - a.isLatest || b.published_at.localeCompare(a.published_at)),
        totalFiles: files.length
      });
    }
  });

  return result.sort((a, b) => a.name.localeCompare(b.name));
});

onMounted(() => {
  loadData();
});

defineExpose({ loadData, refresh: loadData });
</script>

<style scoped>
.file-card.latest-border {
  border-left: 3px solid #4CAF50;
}

.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.v-card-actions {
  gap: 4px;
}
</style>
