<template>
  <div class="mb-6">
    <div class="d-flex align-center justify-space-between mb-4">
        <h2 class="text-h5">文件浏览</h2>
        <v-btn
          icon="mdi-code-json"
          variant="text"
          size="small"
          :color="showRaw ? 'primary' : ''"
          @click="showRaw = !showRaw"
          title="切换原始数据视图"
        ></v-btn>
    </div>
    
    <v-card variant="outlined">
      <v-card-text>
        <div class="d-flex gap-2 align-center mb-4">
          <v-btn icon="mdi-arrow-up" variant="tonal" @click="goUp" :disabled="isRoot" color="secondary"></v-btn>
          
          <v-text-field
            v-model="path"
            label="当前路径"
            placeholder="."
            variant="outlined"
            density="compact"
            hide-details
            @keyup.enter="loadFiles"
            prepend-inner-icon="mdi-folder-open"
          ></v-text-field>
          
          <v-btn color="primary" @click="loadFiles" :loading="loading" icon="mdi-arrow-right"></v-btn>
        </div>

        <div v-if="loading">
            <v-skeleton-loader type="list-item-avatar@3"></v-skeleton-loader>
        </div>

        <div v-else-if="error" class="text-center text-error pa-4">
            <v-icon icon="mdi-alert-circle" color="error" class="mb-2"></v-icon>
            <div>加载失败</div>
            <div class="text-caption">{{ error }}</div>
        </div>
        
        <div v-else-if="showRaw">
             <v-code class="d-block pa-4 rounded bg-grey-lighten-4" style="white-space: pre-wrap; max-height: 500px; overflow: auto; color: black;">{{ JSON.stringify(rawData, null, 2) }}</v-code>
        </div>

        <v-list v-else density="compact" lines="one">
          <v-list-item
            v-for="(item, index) in items"
            :key="index"
            @click="onItemClick(item)"
            :value="item"
            rounded
          >
            <template v-slot:prepend>
              <v-icon :color="item.type === 'dir' ? 'amber' : 'primary'">
                {{ item.type === 'dir' ? 'mdi-folder' : 'mdi-file' }}
              </v-icon>
            </template>

            <v-list-item-title class="font-weight-medium">
                {{ item.name }}
            </v-list-item-title>
            
            <template v-slot:append>
                <span class="text-caption text-medium-emphasis mr-2" v-if="item.size">
                    {{ formatSize(item.size) }}
                </span>
                <span class="text-caption text-medium-emphasis" v-if="item.mtime">
                    {{ formatDate(item.mtime) }}
                </span>
            </template>
          </v-list-item>
          
          <v-list-item v-if="items.length === 0" class="text-center text-medium-emphasis">
            <v-list-item-title>空目录</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { getFiles } from '../services/api';

const path = ref('.');
const rawData = ref(null);
const loading = ref(false);
const error = ref(null);
const showRaw = ref(false);

const isRoot = computed(() => {
    return !path.value || path.value === '.' || path.value === '/' || path.value === '';
});

// Adapter to handle unknown API response shape
const items = computed(() => {
  const data = rawData.value;
  if (!data) return [];
  
  let list = [];
  
  if (Array.isArray(data)) {
    list = data.map(item => {
      if (typeof item === 'string') {
        // Assume string ending in / is dir
        const isDir = item.endsWith('/');
        return { 
            name: isDir ? item.slice(0, -1) : item, 
            type: isDir ? 'dir' : 'file',
            raw: item
        };
      } else if (typeof item === 'object') {
        // Try to detect common properties
        const name = item.name || item.filename || item.key || 'Unknown';
        const isDir = item.type === 'directory' || item.isDirectory === true || item.isDir === true || item.mode?.startsWith('d');
        return {
            name: name,
            type: isDir ? 'dir' : 'file',
            size: item.size || item.length,
            mtime: item.mtime || item.lastModified,
            raw: item
        };
      }
      return { name: String(item), type: 'file' };
    });
  } else if (typeof data === 'object') {
      // Maybe it's a map? key -> details
      list = Object.keys(data).map(key => ({
          name: key,
          type: 'file', // Default
          ...data[key]
      }));
  }

  // Sort: Directories first, then alphabetical
  return list.sort((a, b) => {
      if (a.type === b.type) return a.name.localeCompare(b.name);
      return a.type === 'dir' ? -1 : 1;
  });
});

const loadFiles = async () => {
    loading.value = true;
    error.value = null;
    try {
        // Normalize path
        if (!path.value) path.value = '.';
        
        const res = await getFiles(path.value);
        rawData.value = res.data;
    } catch (e) {
        error.value = e.message || '加载失败';
        console.error(e);
    } finally {
        loading.value = false;
    }
};

const goUp = () => {
    if (isRoot.value) return;
    const parts = path.value.split('/').filter(p => p && p !== '.');
    parts.pop();
    path.value = parts.length ? parts.join('/') : '.';
    loadFiles();
};

const onItemClick = (item) => {
    if (item.type === 'dir') {
        const current = path.value === '.' ? '' : path.value;
        const separator = current.endsWith('/') || !current ? '' : '/';
        path.value = current + separator + item.name;
        loadFiles();
    } else {
        // For files, ideally we would download or preview.
        // For now, let's just log it or maybe copy path?
        console.log('File clicked:', item);
    }
};

const formatSize = (bytes) => {
    if (bytes === 0) return '0 B';
    if (!bytes) return '';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const formatDate = (date) => {
    if (!date) return '';
    try {
        return new Date(date).toLocaleDateString();
    } catch (e) {
        return date;
    }
};

onMounted(() => {
    loadFiles();
});
</script>
