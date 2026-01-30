<script setup>
import { ref, computed, onMounted } from 'vue';
import { getStatus, getLatest } from '@/services/api';
import { Search, File, Folder, Download, Copy, Loader2, FileArchive, HardDrive, ChevronRight, Home, ArrowLeft } from 'lucide-vue-next';
import Input from '@/components/ui/Input.vue';
import Button from '@/components/ui/Button.vue';
import Badge from '@/components/ui/Badge.vue';
import Skeleton from '@/components/ui/Skeleton.vue';
import { cn } from '@/lib/utils';
import { useClipboard } from '@vueuse/core';

const loading = ref(true);
const searchQuery = ref('');
const launchers = ref({});
const latestData = ref({});
const { copy, copied } = useClipboard();

// Navigation State
const currentPath = ref([]); // Array of { name: string, id: string, type: 'root'|'launcher'|'version' }

const loadData = async () => {
  loading.value = true;
  try {
    const [statusRes, latestRes] = await Promise.all([getStatus(), getLatest()]);
    // Sort launchers
    const sortedLaunchers = {};
    Object.keys(statusRes.data).sort().forEach(key => {
        sortedLaunchers[key] = statusRes.data[key].sort((a, b) => 
            String(b.tag_name || b.name).localeCompare(String(a.tag_name || a.name))
        );
    });
    launchers.value = sortedLaunchers;
    latestData.value = latestRes.data;
  } catch (error) {
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const getFileIcon = (filename) => {
  const ext = filename.split('.').pop()?.toLowerCase();
  if (['zip', 'tar', 'gz', '7z', 'rar'].includes(ext)) return FileArchive;
  if (['exe', 'msi', 'apk', 'dmg'].includes(ext)) return HardDrive;
  return File;
};

const formatDate = (dateString) => {
  if (!dateString) return '未知';
  try {
    return new Date(dateString).toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
  } catch {
    return dateString;
  }
};

const copyUrl = (url) => {
  copy(url);
};

// --- Navigation Logic ---

const navigateTo = (item, type) => {
    if (type === 'launcher') {
        currentPath.value = [{ name: item, id: item, type: 'launcher' }];
    } else if (type === 'version') {
        currentPath.value.push({ name: item.name, id: item.id, type: 'version', data: item.data });
    }
};

const navigateUp = () => {
    currentPath.value.pop();
};

const navigateToBreadcrumb = (index) => {
    if (index === -1) {
        currentPath.value = [];
    } else {
        currentPath.value = currentPath.value.slice(0, index + 1);
    }
};

// --- Computed Data for Current View ---

const currentItems = computed(() => {
    const query = searchQuery.value.toLowerCase().trim();
    
    // If searching, we might want to show a flat list of matching files across all launchers?
    // For now, let's keep navigation but filter current view. 
    // If search is active at root, maybe filter launchers.
    
    const depth = currentPath.value.length;
    
    if (depth === 0) {
        // Root: Show Launchers
        return Object.keys(launchers.value).map(name => ({
            id: name,
            name: name,
            type: 'launcher',
            count: launchers.value[name].length,
            latest: latestData.value[name]
        })).filter(l => !query || l.name.toLowerCase().includes(query));
    } else if (depth === 1) {
        // Launcher: Show Versions
        const launcherName = currentPath.value[0].id;
        const versions = launchers.value[launcherName] || [];
        
        return versions.map(v => ({
            id: v.tag_name || v.name,
            name: v.tag_name || v.name,
            type: 'version',
            date: v.published_at,
            isLatest: latestData.value[launcherName] === (v.tag_name || v.name),
            data: v,
            fileCount: v.assets?.length || 0
        })).filter(v => !query || v.name.toLowerCase().includes(query));
    } else if (depth === 2) {
         // Version: Show Files
         const versionData = currentPath.value[1].data;
         const launcherName = currentPath.value[0].id;
         const versionName = currentPath.value[1].id;
         
         return (versionData.assets || []).map(asset => ({
             id: asset.name,
             name: asset.name,
             type: 'file',
             size: asset.size, // API doesn't always give size, but if it did...
             downloadUrl: asset.url && asset.url.startsWith('http') 
              ? asset.url 
              : `${window.location.origin}/download/${launcherName}/${versionName}/${asset.name}`
         })).filter(f => !query || f.name.toLowerCase().includes(query));
    }
    return [];
});

onMounted(() => {
  loadData();
});
</script>

<template>
  <div class="flex flex-col h-full space-y-4 max-w-full">
    <!-- Header & Breadcrumbs -->
    <div class="flex flex-col space-y-4 shrink-0">
        <div class="flex items-center justify-between gap-4">
             <div class="flex items-center gap-2 overflow-hidden text-sm font-medium text-muted-foreground">
                 <Button variant="ghost" size="icon" class="h-8 w-8 shrink-0" @click="navigateToBreadcrumb(-1)" :disabled="!currentPath.length">
                     <Home class="h-4 w-4" />
                 </Button>
                 <template v-for="(crumb, index) in currentPath" :key="crumb.id">
                     <ChevronRight class="h-4 w-4 shrink-0 opacity-50" />
                     <Button variant="ghost" size="sm" class="h-8 px-2 truncate max-w-[120px]" @click="navigateToBreadcrumb(index)">
                         {{ crumb.name }}
                     </Button>
                 </template>
             </div>
             <div class="relative w-40 md:w-64 shrink-0">
                 <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                 <Input
                   v-model="searchQuery"
                   type="search"
                   placeholder="筛选..."
                   class="pl-8 h-9 bg-background/50 backdrop-blur border-white/10"
                 />
             </div>
        </div>
    </div>

    <!-- Content Area -->
    <div v-if="loading" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
        <Skeleton class="aspect-square rounded-xl" v-for="i in 10" :key="i" />
    </div>

    <div v-else-if="!currentItems.length" class="flex flex-col items-center justify-center py-20 text-muted-foreground">
        <Folder class="h-16 w-16 mb-4 opacity-20" />
        <p>空文件夹</p>
    </div>

    <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-2 pb-20">
        <!-- Back Button (if not root) -->
        <div 
            v-if="currentPath.length > 0" 
            @click="navigateUp"
            class="group relative flex flex-col items-center justify-center gap-2 p-3 rounded-xl border border-dashed border-muted-foreground/20 bg-muted/5 hover:bg-muted/10 cursor-pointer transition-all hover:scale-[1.02] active:scale-[0.98]"
        >
             <ArrowLeft class="h-6 w-6 text-muted-foreground group-hover:text-foreground transition-colors" />
             <span class="text-xs font-medium text-muted-foreground">返回上一级</span>
        </div>

        <!-- Items -->
        <div 
            v-for="item in currentItems" 
            :key="item.id"
            @click="item.type !== 'file' ? navigateTo(item.type === 'launcher' ? item.name : item, item.type) : null"
            :class="cn(
                'group relative flex flex-col justify-between p-3 rounded-xl border border-white/5 bg-background/40 backdrop-blur-md shadow-sm transition-all duration-200',
                item.type !== 'file' ? 'cursor-pointer hover:bg-background/60 hover:border-white/20 hover:shadow-md hover:scale-[1.02] active:scale-[0.98]' : ''
            )"
        >
            <!-- Icon Area -->
            <div class="flex items-start justify-between mb-2">
                 <div class="p-2 rounded-lg bg-gradient-to-br from-primary/10 to-primary/5 text-primary group-hover:from-primary/20 group-hover:to-primary/10 transition-colors">
                     <Folder v-if="item.type === 'launcher'" class="h-5 w-5" />
                     <Folder v-else-if="item.type === 'version'" class="h-5 w-5" />
                     <component v-else :is="getFileIcon(item.name)" class="h-5 w-5" />
                 </div>
                 <div v-if="item.isLatest" class="px-1.5 py-0.5 rounded-full bg-green-500/10 text-green-600 text-[9px] font-bold uppercase tracking-wider border border-green-500/20">
                     Latest
                 </div>
                 <div v-if="item.type === 'file'" class="flex gap-0.5">
                      <Button size="icon" variant="ghost" class="h-6 w-6" @click.stop="copyUrl(item.downloadUrl)">
                          <Copy class="h-3 w-3" />
                      </Button>
                      <Button size="icon" variant="ghost" class="h-6 w-6" as="a" :href="item.downloadUrl">
                          <Download class="h-3 w-3" />
                      </Button>
                 </div>
            </div>

            <!-- Text Area -->
            <div class="flex flex-col gap-0.5 min-w-0">
                <h3 class="font-medium truncate text-sm leading-none" :title="item.name">{{ item.name }}</h3>
                <div class="flex items-center justify-between text-[10px] text-muted-foreground mt-1">
                    <span v-if="item.type === 'launcher'">{{ item.count }} 版本</span>
                    <span v-else-if="item.type === 'version'">{{ formatDate(item.date) }}</span>
                    <span v-else>文件</span>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>


<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
