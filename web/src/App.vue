<template>
  <v-app>
    <v-app-bar :elevation="0" border color="surface-light" class="px-2">
      <template v-slot:prepend>
        <v-app-bar-nav-icon v-if="!mobile" variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
        <v-avatar size="32" class="ml-2 mr-2 cursor-pointer" @click="showAbout = true">
           <v-img src="https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/b4ee27d31312bdb9.svg" alt="Logo"></v-img>
        </v-avatar>
      </template>

      <v-app-bar-title 
        class="font-weight-bold cursor-pointer" 
        @click="showAbout = true"
        style="user-select: none;"
      >
        柠枺镜像
        <v-tooltip activator="parent" location="bottom">关于本站</v-tooltip>
      </v-app-bar-title>

      <template v-slot:append>
        <v-btn
          :variant="mobile ? 'text' : 'flat'"
          :icon="mobile"
          color="primary"
          class="mr-2"
          @click="manualRefresh"
          :loading="refreshing"
        >
          <v-icon>mdi-refresh</v-icon>
          <span v-if="!mobile" class="ml-2">刷新</span>
        </v-btn>

        <v-btn icon @click="toggleTheme" variant="text">
          <v-icon>{{ theme.global.current.value.dark ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
        </v-btn>
      </template>
    </v-app-bar>

    <v-navigation-drawer v-if="!mobile" v-model="drawer" elevation="1">
      <v-list nav class="pa-2">
        <v-list-item prepend-icon="mdi-home-variant" title="首页" value="home" @click="tab = 'home'" :active="tab === 'home'" rounded></v-list-item>
        <v-list-item prepend-icon="mdi-folder-multiple" title="文件浏览" value="files" @click="tab = 'files'" :active="tab === 'files'" rounded></v-list-item>
        <v-list-item prepend-icon="mdi-chart-timeline-variant" title="数据统计" value="stats" @click="tab = 'stats'" :active="tab === 'stats'" rounded></v-list-item>
        <v-list-item prepend-icon="mdi-api" title="API 文档" value="api" @click="tab = 'api'" :active="tab === 'api'" rounded></v-list-item>
      </v-list>
      
      <template v-slot:append>
        <div class="pa-4 text-caption text-center text-medium-emphasis">
          v1.0.0
        </div>
      </template>
    </v-navigation-drawer>

    <v-main class="bg-background">
      <v-container class="pb-16 px-4 pt-6" fluid style="max-width: 1200px">
        <v-window v-model="tab" disabled transition="scroll-x-transition" reverse-transition="scroll-x-reverse-transition">
          <v-window-item value="home">
            <Announcements />
            <VersionList ref="versionListRef" />
          </v-window-item>

          <v-window-item value="files">
            <FileBrowser @go-home="tab = 'home'" />
          </v-window-item>

          <v-window-item value="stats">
            <Statistics ref="statsRef" />
          </v-window-item>

          <v-window-item value="api">
             <ApiDocs />
          </v-window-item>
        </v-window>
        
        <v-footer class="text-center d-flex flex-column py-6 text-medium-emphasis mt-8 bg-transparent">
            <div class="d-flex align-center gap-2">
               <span>&copy; {{ new Date().getFullYear() }} Lemwood Mirror</span>
            </div>
            <div class="text-caption mt-2">
                <a href="https://beian.miit.gov.cn" target="_blank" class="text-decoration-none text-medium-emphasis hover-link">
                   新ICP备2024015133号-5
                </a>
            </div>
        </v-footer>
      </v-container>
    </v-main>

    <v-bottom-navigation v-if="mobile" v-model="tab" grow color="primary" elevation="1">
      <v-btn value="home">
        <v-icon>mdi-home-variant</v-icon>
        <span>首页</span>
      </v-btn>

      <v-btn value="files">
        <v-icon>mdi-folder-multiple</v-icon>
        <span>文件</span>
      </v-btn>

      <v-btn value="stats">
        <v-icon>mdi-chart-timeline-variant</v-icon>
        <span>统计</span>
      </v-btn>

      <v-btn value="api">
        <v-icon>mdi-api</v-icon>
        <span>API</span>
      </v-btn>
    </v-bottom-navigation>

    <!-- About Dialog -->
    <v-dialog v-model="showAbout" max-width="400">
      <v-card class="text-center pa-4 rounded">
        <v-card-text>
          <v-avatar size="64" class="mb-4">
            <v-img src="https://cdn.mengze.vip/gh/JanePHPDev/Blog-Static-Resource@main/images/b4ee27d31312bdb9.svg" alt="Logo"></v-img>
          </v-avatar>
          <h3 class="text-h5 font-weight-bold mb-1">柠枺镜像</h3>
          <div class="text-subtitle-2 text-medium-emphasis mb-4">Lemwood Mirror Service</div>
          
          <v-divider class="mb-4"></v-divider>
          
          <p class="text-body-2 mb-4">
            这是一个提供各类Minecraft启动器及核心组件下载的高速镜像站，致力于为开发者与玩家打造稳定、快速的资源获取体验。基于腾讯云200Mbps高速服务器与自研Go语言Mirror架构，我们确保主流启动器如HMCL、PCL2、BakaXL及其依赖库的分发效率，无论是日常更新还是开发调试，都能享受流畅可靠的下载服务。
          </p>
          
          <div class="d-flex justify-center gap-4 mb-2">
             <v-btn variant="text" size="small" prepend-icon="mdi-web" href="https://zs.lemwood.cn/" target="_blank">赞赏我</v-btn>
             <v-btn variant="text" size="small" prepend-icon="mdi-github" href="https://github.com" target="_blank">GitHub</v-btn>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" variant="tonal" block @click="showAbout = false">关闭</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useTheme, useDisplay } from 'vuetify';
import { scan } from './services/api';

import ApiDocs from './components/ApiDocs.vue';
import Announcements from './components/Announcements.vue';
import VersionList from './components/VersionList.vue';
import Statistics from './components/Statistics.vue';
import FileBrowser from './components/FileBrowser.vue';

const theme = useTheme();
const { mobile } = useDisplay();
const refreshing = ref(false);
const tab = ref('home');
const drawer = ref(true);
const showAbout = ref(false);

const versionListRef = ref(null);
const statsRef = ref(null);

const toggleTheme = () => {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark';
};

const manualRefresh = async () => {
  refreshing.value = true;
  try {
    await scan();
    if (versionListRef.value) await versionListRef.value.refresh();
    if (statsRef.value) await statsRef.value.refresh();
  } catch (e) {
    console.error('Manual refresh failed:', e);
  } finally {
    refreshing.value = false;
  }
};

onMounted(() => {
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme) theme.global.name.value = savedTheme;
  
  const savedTab = localStorage.getItem('currentTab');
  if (savedTab) tab.value = savedTab;
});

watch(tab, (val) => localStorage.setItem('currentTab', val));
watch(() => theme.global.name.value, (val) => localStorage.setItem('theme', val));
</script>

<style>
.bg-background {
  background-color: rgb(var(--v-theme-background));
}
.cursor-pointer {
  cursor: pointer;
}
.hover-link:hover {
  color: rgb(var(--v-theme-primary)) !important;
  text-decoration: underline !important;
}
</style>

<style>
/* Global background override for Vuetify 3 app to ensure full coverage */
.bg-background {
  background-color: rgb(var(--v-theme-background));
}
</style>