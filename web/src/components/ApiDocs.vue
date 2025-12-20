<template>
  <div class="mb-6">
    <div class="mb-6">
      <h2 class="text-h4 font-weight-bold mb-2">API 文档</h2>
      <p class="text-medium-emphasis">
        Lemwood Mirror 提供了 RESTful API 接口，供开发者程序化获取镜像站数据。
        所有接口的基础路径为 <code class="bg-surface-variant px-1 rounded">/api</code>。
      </p>
    </div>

    <v-row>
      <v-col cols="12" md="4" class="hidden-sm-and-down">
         <v-card position="sticky" style="top: 80px" variant="outlined" class="rounded">
           <v-list density="compact" nav>
             <v-list-subheader>Endpoint Index</v-list-subheader>
             <v-list-item 
               v-for="(endpoint, i) in endpoints" 
               :key="i"
               @click="scrollTo(i)"
               :value="i"
               color="primary"
             >
               <template v-slot:prepend>
                 <span :class="`text-${getMethodColor(endpoint.method)} font-weight-bold mr-2 text-caption`" style="width: 30px">{{ endpoint.method }}</span>
               </template>
               <v-list-item-title class="text-caption">{{ endpoint.path }}</v-list-item-title>
             </v-list-item>
           </v-list>
         </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card
          v-for="(endpoint, index) in endpoints"
          :key="index"
          class="mb-4 rounded border"
          variant="flat"
          :id="`endpoint-${index}`"
        >
          <v-card-item>
            <template v-slot:prepend>
              <v-chip 
                :color="getMethodColor(endpoint.method)" 
                size="small" 
                label 
                class="font-weight-bold mr-2"
              >
                {{ endpoint.method }}
              </v-chip>
            </template>
            <v-card-title class="text-body-1 font-family-monospace">
              {{ endpoint.path }}
            </v-card-title>
          </v-card-item>
          
          <v-divider></v-divider>

          <v-card-text>
            <div class="text-subtitle-1 font-weight-bold mb-1">{{ endpoint.title }}</div>
            <p class="text-body-2 text-medium-emphasis mb-3">{{ endpoint.desc }}</p>
            
            <div class="bg-surface-variant rounded pa-3 font-family-monospace text-caption">
               <span class="text-medium-emphasis"># Example Request</span><br/>
               curl -X {{ endpoint.method }} "https://mirror.lemwood.cn{{ endpoint.path.replace('{launcher}', 'hmcl') }}"
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
const endpoints = [
  { method: 'GET', path: '/api/status', title: '获取所有版本状态', desc: '返回系统中所有已注册启动器及其所有可用版本的详细列表。数据包含版本号、发布时间、下载链接等。' },
  { method: 'GET', path: '/api/status/{launcher}', title: '获取指定启动器状态', desc: '返回特定启动器（如 hmcl, pcl2）的所有历史版本信息。' },
  { method: 'GET', path: '/api/latest', title: '获取所有最新版本', desc: '返回所有启动器的最新版本号及简要信息，适合用于检查更新。' },
  { method: 'GET', path: '/api/latest/{launcher}', title: '获取指定启动器最新版本', desc: '查询特定启动器的最新发布版本。' },
  { method: 'GET', path: '/api/stats', title: '获取统计数据', desc: '获取全站的访问量、下载量趋势、热门下载排行以及访客地域分布数据。' },
  { method: 'POST', path: '/api/scan', title: '触发手动扫描', desc: '强制服务器立即与上游仓库同步，检查是否有新版本发布。此操作有频率限制。' },
];

const getMethodColor = (method) => {
    switch(method) {
        case 'GET': return 'primary';
        case 'POST': return 'success';
        case 'PUT': return 'warning';
        case 'DELETE': return 'error';
        default: return 'grey';
    }
};

const scrollTo = (index) => {
    const el = document.getElementById(`endpoint-${index}`);
    if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
};
</script>

<style scoped>
.font-family-monospace {
    font-family: 'Roboto Mono', monospace, monospace;
}
</style>
