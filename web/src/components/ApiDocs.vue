<template>
  <v-container fluid class="pa-2">
    <v-row dense>
      <v-col cols="12">
        <div class="d-flex align-center mb-2">
          <h1 class="text-h5 font-weight-bold">API 文档</h1>
          <v-spacer />
          <v-btn
            v-if="mobile"
            icon="mdi-menu"
            variant="text"
            size="small"
            @click="showMobileMenu = !showMobileMenu"
          />
        </div>
        <p class="text-caption text-medium-emphasis mb-3">
          RESTful API · 基础路径: <code class="bg-black px-1 rounded">/api</code>
        </p>
        
        <v-text-field
          v-model="searchQuery"
          prepend-inner-icon="mdi-magnify"
          label="搜索接口"
          variant="outlined"
          density="compact"
          hide-details
          clearable
          class="mb-2"
        />
      </v-col>
    </v-row>

    <v-expand-transition>
      <v-row v-if="mobile && showMobileMenu" dense>
        <v-col cols="12">
          <v-card class="mb-2" variant="outlined">
            <v-list density="compact" nav>
              <v-list-subheader class="text-caption">ENDPOINTS</v-list-subheader>
              <v-list-item
                v-for="(endpoint, i) in filteredEndpoints"
                :key="i"
                @click="scrollTo(i); showMobileMenu = false"
                density="compact"
              >
                <template v-slot:prepend>
                  <v-chip :color="getMethodColor(endpoint.method)" size="x-small" label class="mr-2" style="min-width: 40px; justify-content: center">
                    {{ endpoint.method }}
                  </v-chip>
                </template>
                <v-list-item-title class="text-caption">{{ endpoint.path }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </v-expand-transition>

    <v-snackbar v-model="showSnackbar" :color="snackbarColor" timeout="1500" min-height="36">
      {{ snackbarMessage }}
    </v-snackbar>

    <v-row dense>
      <v-col cols="12" md="4" class="hidden-sm-and-down">
        <v-card position="sticky" style="top: 70px" variant="outlined" class="rounded">
          <v-list density="compact" nav>
            <v-list-subheader class="text-caption">ENDPOINTS</v-list-subheader>
            <v-list-item
              v-for="(endpoint, i) in filteredEndpoints"
              :key="i"
              @click="scrollTo(i)"
              :active="activeIndex === i"
              density="compact"
            >
              <template v-slot:prepend>
                <v-chip :color="getMethodColor(endpoint.method)" size="x-small" label class="mr-2" style="min-width: 40px; justify-content: center">
                  {{ endpoint.method }}
                </v-chip>
              </template>
              <v-list-item-title class="text-caption">{{ endpoint.path }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <div v-if="!filteredEndpoints.length" class="text-center py-8">
          <v-icon icon="mdi-magnify" size="32" color="medium-emphasis" class="mb-2" />
          <p class="text-caption text-medium-emphasis">无匹配结果</p>
        </div>

        <!-- /api/status -->
        <v-card id="endpoint-0" class="mb-3 rounded endpoint-card" variant="outlined" @mouseenter="activeIndex = 0">
          <v-card-item class="py-2">
            <template v-slot:prepend>
              <v-chip color="primary" size="small" label class="font-weight-bold mr-2">GET</v-chip>
            </template>
            <v-card-title class="text-body-2 font-weight-bold font-mono">/api/status</v-card-title>
          </v-card-item>
          <v-divider />
          <v-card-text class="py-3">
            <div class="text-subtitle-2 font-weight-bold mb-2">获取所有版本状态</div>
            <p class="text-caption text-medium-emphasis mb-4">返回所有启动器及完整版本列表，包含版本号、发布时间、下载链接</p>
            
            <v-expansion-panels class="mb-3" variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-code-json</v-icon>
                  响应示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto" style="max-height: 300px">
                    <pre class="ma-0"><span class="text-yellow">{</span>
  <span class="text-green">"hmcl"</span>: [
    <span class="text-yellow">{</span>
      <span class="text-green">"tag_name"</span>: <span class="text-cyan">"v3.5.9"</span>,
      <span class="text-green">"name"</span>: <span class="text-cyan">"HMCL v3.5.9"</span>,
      <span class="text-green">"published_at"</span>: <span class="text-cyan">"2024-01-15T10:30:00Z"</span>,
      <span class="text-green">"assets"</span>: <span class="text-yellow">[</span>
        <span class="text-yellow">{</span>
          <span class="text-green">"name"</span>: <span class="text-cyan">"HMCL-3.5.9.exe"</span>,
          <span class="text-green">"size"</span>: <span class="text-orange">2856128</span>
        <span class="text-yellow">}</span>
      <span class="text-yellow">]</span>
    <span class="text-yellow">}</span>
  <span class="text-yellow">]</span>,
  <span class="text-green">"pcl2"</span>: <span class="text-yellow">[]</span>
<span class="text-yellow">}</span></pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
              
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-language-javascript</v-icon>
                  JavaScript 调用示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-purple">async</span> <span class="text-purple">function</span> <span class="text-green">getAllStatus</span>() {
  <span class="text-purple">const</span> response = <span class="text-purple">await</span> <span class="text-orange">fetch</span>(<span class="text-cyan">'https://mirror.lemwood.icu/api/status'</span>);
  <span class="text-purple">const</span> data = <span class="text-purple">await</span> response.<span class="text-green">json</span>();
  
  <span class="text-purple">return</span> data;
}</pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            
            <div class="bg-black rounded pa-2 font-mono text-caption relative mb-2">
              <span class="text-orange">curl</span>
              <span class="text-yellow ml-1">-X</span>
              <span class="text-green ml-1 font-weight-bold">GET</span>
              <span class="text-cyan ml-2">"https://mirror.lemwood.icu/api/status"</span>
              <v-btn
                icon="mdi-content-copy"
                variant="text"
                size="x-small"
                class="absolute top-1 right-1"
                @click="copyCurl('/api/status', 'GET')"
              />
            </div>
          </v-card-text>
        </v-card>

        <!-- /api/status/{launcher} -->
        <v-card id="endpoint-1" class="mb-3 rounded endpoint-card" variant="outlined" @mouseenter="activeIndex = 1">
          <v-card-item class="py-2">
            <template v-slot:prepend>
              <v-chip color="primary" size="small" label class="font-weight-bold mr-2">GET</v-chip>
            </template>
            <v-card-title class="text-body-2 font-weight-bold font-mono">/api/status/{launcher}</v-card-title>
          </v-card-item>
          <v-divider />
          <v-card-text class="py-3">
            <div class="text-subtitle-2 font-weight-bold mb-2">获取指定启动器状态</div>
            <p class="text-caption text-medium-emphasis mb-4">返回特定启动器的历史版本信息</p>
            
            <v-expansion-panels class="mb-3" variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-tune</v-icon>
                  请求参数
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-table density="compact" class="bg-black rounded">
                    <thead>
                      <tr>
                        <th class="text-caption font-weight-bold">参数</th>
                        <th class="text-caption font-weight-bold">类型</th>
                        <th class="text-caption font-weight-bold">必填</th>
                        <th class="text-caption font-weight-bold">说明</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td class="text-caption font-mono">launcher</td>
                        <td class="text-caption font-mono">string</td>
                        <td class="text-caption">
                          <v-chip size="x-small" color="error">是</v-chip>
                        </td>
                        <td class="text-caption">启动器名称 (hmcl/pcl2/bakaxl)</td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-expansion-panel-text>
              </v-expansion-panel>
              
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-code-json</v-icon>
                  响应示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto" style="max-height: 300px">
                    <pre class="ma-0"><span class="text-yellow">[</span>
  <span class="text-yellow">{</span>
    <span class="text-green">"tag_name"</span>: <span class="text-cyan">"v3.5.9"</span>,
    <span class="text-green">"published_at"</span>: <span class="text-cyan">"2024-01-15T10:30:00Z"</span>,
    <span class="text-green">"assets"</span>: <span class="text-yellow">[]</span>
  <span class="text-yellow">}</span>
<span class="text-yellow">]</span></pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
              
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-language-javascript</v-icon>
                  JavaScript 调用示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-purple">async</span> <span class="text-purple">function</span> <span class="text-green">getLauncherVersions</span>(name) {
  <span class="text-purple">const</span> valid = [<span class="text-cyan">'hmcl'</span>, <span class="text-cyan">'pcl2'</span>, <span class="text-cyan">'bakaxl'</span>];
  <span class="text-purple">if</span> (!valid.<span class="text-green">includes</span>(name)) {
    <span class="text-purple">throw</span> <span class="text-purple">new</span> <span class="text-orange">Error</span>(<span class="text-cyan">'无效的启动器名称'</span>);
  }
  
  <span class="text-purple">const</span> response = <span class="text-purple">await</span> <span class="text-orange">fetch</span>(
    <span class="text-cyan">'https://mirror.lemwood.icu/api/status/'</span> + name
  );
  
  <span class="text-purple">return</span> <span class="text-purple">await</span> response.<span class="text-green">json</span>();
}</pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            
            <div class="bg-black rounded pa-2 font-mono text-caption relative mb-2">
              <span class="text-orange">curl</span>
              <span class="text-yellow ml-1">-X</span>
              <span class="text-green ml-1 font-weight-bold">GET</span>
              <span class="text-cyan ml-2">"https://mirror.lemwood.icu/api/status/hmcl"</span>
              <v-btn
                icon="mdi-content-copy"
                variant="text"
                size="x-small"
                class="absolute top-1 right-1"
                @click="copyCurl('/api/status/hmcl', 'GET')"
              />
            </div>
          </v-card-text>
        </v-card>

        <!-- /api/latest -->
        <v-card id="endpoint-2" class="mb-3 rounded endpoint-card" variant="outlined" @mouseenter="activeIndex = 2">
          <v-card-item class="py-2">
            <template v-slot:prepend>
              <v-chip color="primary" size="small" label class="font-weight-bold mr-2">GET</v-chip>
            </template>
            <v-card-title class="text-body-2 font-weight-bold font-mono">/api/latest</v-card-title>
          </v-card-item>
          <v-divider />
          <v-card-text class="py-3">
            <div class="text-subtitle-2 font-weight-bold mb-2">获取所有最新版本</div>
            <p class="text-caption text-medium-emphasis mb-4">快速检查所有启动器的最新版本号</p>
            
            <v-expansion-panels class="mb-3" variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-code-json</v-icon>
                  响应示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-yellow">{</span>
  <span class="text-green">"hmcl"</span>: <span class="text-cyan">"v3.5.9"</span>,
  <span class="text-green">"pcl2"</span>: <span class="text-cyan">"Snapshot-20240115"</span>,
  <span class="text-green">"bakaxl"</span>: <span class="text-cyan">"v3.5.1"</span>
<span class="text-yellow">}</span></pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
              
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-language-javascript</v-icon>
                  JavaScript 调用示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-purple">async</span> <span class="text-purple">function</span> <span class="text-green">checkUpdates</span>() {
  <span class="text-purple">const</span> latest = <span class="text-purple">await</span> <span class="text-orange">fetch</span>(<span class="text-cyan">'https://mirror.lemwood.icu/api/latest'</span>)
    .<span class="text-green">then</span>(r => r.<span class="text-green">json</span>());
  
  <span class="text-purple">return</span> latest;
}</pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            
            <div class="bg-black rounded pa-2 font-mono text-caption relative mb-2">
              <span class="text-orange">curl</span>
              <span class="text-yellow ml-1">-X</span>
              <span class="text-green ml-1 font-weight-bold">GET</span>
              <span class="text-cyan ml-2">"https://mirror.lemwood.icu/api/latest"</span>
              <v-btn
                icon="mdi-content-copy"
                variant="text"
                size="x-small"
                class="absolute top-1 right-1"
                @click="copyCurl('/api/latest', 'GET')"
              />
            </div>
          </v-card-text>
        </v-card>

        <!-- /api/latest/{launcher} -->
        <v-card id="endpoint-3" class="mb-3 rounded endpoint-card" variant="outlined" @mouseenter="activeIndex = 3">
          <v-card-item class="py-2">
            <template v-slot:prepend>
              <v-chip color="primary" size="small" label class="font-weight-bold mr-2">GET</v-chip>
            </template>
            <v-card-title class="text-body-2 font-weight-bold font-mono">/api/latest/{launcher}</v-card-title>
          </v-card-item>
          <v-divider />
          <v-card-text class="py-3">
            <div class="text-subtitle-2 font-weight-bold mb-2">获取指定启动器最新版本</div>
            <p class="text-caption text-medium-emphasis mb-4">查询单个启动器的最新发布版本</p>
            
            <v-expansion-panels class="mb-3" variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-tune</v-icon>
                  请求参数
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-table density="compact" class="bg-black rounded">
                    <thead>
                      <tr>
                        <th class="text-caption font-weight-bold">参数</th>
                        <th class="text-caption font-weight-bold">类型</th>
                        <th class="text-caption font-weight-bold">必填</th>
                        <th class="text-caption font-weight-bold">说明</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td class="text-caption font-mono">launcher</td>
                        <td class="text-caption font-mono">string</td>
                        <td class="text-caption">
                          <v-chip size="x-small" color="error">是</v-chip>
                        </td>
                        <td class="text-caption">启动器名称</td>
                      </tr>
                    </tbody>
                  </v-table>
                </v-expansion-panel-text>
              </v-expansion-panel>
              
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-code-json</v-icon>
                  响应示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-yellow">{</span>
  <span class="text-green">"tag_name"</span>: <span class="text-cyan">"v3.5.9"</span>,
  <span class="text-green">"name"</span>: <span class="text-cyan">"HMCL v3.5.9"</span>,
  <span class="text-green">"published_at"</span>: <span class="text-cyan">"2024-01-15T10:30:00Z"</span>
<span class="text-yellow">}</span></pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            
            <div class="bg-black rounded pa-2 font-mono text-caption relative mb-2">
              <span class="text-orange">curl</span>
              <span class="text-yellow ml-1">-X</span>
              <span class="text-green ml-1 font-weight-bold">GET</span>
              <span class="text-cyan ml-2">"https://mirror.lemwood.icu/api/latest/hmcl"</span>
              <v-btn
                icon="mdi-content-copy"
                variant="text"
                size="x-small"
                class="absolute top-1 right-1"
                @click="copyCurl('/api/latest/hmcl', 'GET')"
              />
            </div>
          </v-card-text>
        </v-card>

        <!-- /api/stats -->
        <v-card id="endpoint-4" class="mb-3 rounded endpoint-card" variant="outlined" @mouseenter="activeIndex = 4">
          <v-card-item class="py-2">
            <template v-slot:prepend>
              <v-chip color="primary" size="small" label class="font-weight-bold mr-2">GET</v-chip>
            </template>
            <v-card-title class="text-body-2 font-weight-bold font-mono">/api/stats</v-card-title>
          </v-card-item>
          <v-divider />
          <v-card-text class="py-3">
            <div class="text-subtitle-2 font-weight-bold mb-2">获取统计数据</div>
            <p class="text-caption text-medium-emphasis mb-4">访问统计、下载量、热门排行、地域分布</p>
            
            <v-expansion-panels class="mb-3" variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-code-json</v-icon>
                  响应示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-yellow">{</span>
  <span class="text-green">"totalDownloads"</span>: <span class="text-orange">152304</span>,
  <span class="text-green">"totalVisits"</span>: <span class="text-orange">89234</span>,
  <span class="text-green">"topDownloads"</span>: <span class="text-yellow">[</span>
    <span class="text-yellow">{</span>
      <span class="text-green">"launcher"</span>: <span class="text-cyan">"hmcl"</span>,
      <span class="text-green">"version"</span>: <span class="text-cyan">"v3.5.9"</span>,
      <span class="text-green">"count"</span>: <span class="text-orange">5234</span>
    <span class="text-yellow">}</span>
  <span class="text-yellow">]</span>
<span class="text-yellow">}</span></pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            
            <div class="bg-black rounded pa-2 font-mono text-caption relative mb-2">
              <span class="text-orange">curl</span>
              <span class="text-yellow ml-1">-X</span>
              <span class="text-green ml-1 font-weight-bold">GET</span>
              <span class="text-cyan ml-2">"https://mirror.lemwood.icu/api/stats"</span>
              <v-btn
                icon="mdi-content-copy"
                variant="text"
                size="x-small"
                class="absolute top-1 right-1"
                @click="copyCurl('/api/stats', 'GET')"
              />
            </div>
          </v-card-text>
        </v-card>

        <!-- /api/scan -->
        <v-card id="endpoint-5" class="mb-3 rounded endpoint-card" variant="outlined" @mouseenter="activeIndex = 5">
          <v-card-item class="py-2">
            <template v-slot:prepend>
              <v-chip color="success" size="small" label class="font-weight-bold mr-2">POST</v-chip>
            </template>
            <v-card-title class="text-body-2 font-weight-bold font-mono">/api/scan</v-card-title>
          </v-card-item>
          <v-divider />
          <v-card-text class="py-3">
            <div class="text-subtitle-2 font-weight-bold mb-2">触发手动扫描</div>
            <p class="text-caption text-medium-emphasis mb-4">强制同步上游仓库检查新版本（需认证）</p>
            
            <v-expansion-panels class="mb-3" variant="accordion">
              <v-expansion-panel>
                <v-expansion-panel-title class="text-caption font-weight-bold">
                  <v-icon start size="small">mdi-code-json</v-icon>
                  响应示例
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="bg-black rounded pa-2 font-mono text-caption overflow-x-auto">
                    <pre class="ma-0"><span class="text-yellow">{</span>
  <span class="text-green">"success"</span>: <span class="text-purple">true</span>,
  <span class="text-green">"newVersions"</span>: <span class="text-orange">2</span>,
  <span class="text-green">"message"</span>: <span class="text-cyan">"扫描完成"</span>
<span class="text-yellow">}</span></pre>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
            
            <div class="mt-3">
              <div class="text-caption font-weight-bold mb-1">错误码</div>
              <div class="d-flex flex-wrap gap-1">
                <v-chip size="x-small" variant="outlined" class="font-mono">429: 频率限制</v-chip>
                <v-chip size="x-small" variant="outlined" class="font-mono">403: 权限不足</v-chip>
                <v-chip size="x-small" variant="outlined" class="font-mono">500: 服务器错误</v-chip>
              </div>
            </div>
            
            <div class="bg-black rounded pa-2 font-mono text-caption relative mt-3">
              <span class="text-orange">curl</span>
              <span class="text-yellow ml-1">-X</span>
              <span class="text-green ml-1 font-weight-bold">POST</span>
              <span class="text-cyan ml-2">"https://mirror.lemwood.icu/api/scan"</span>
              <v-btn
                icon="mdi-content-copy"
                variant="text"
                size="x-small"
                class="absolute top-1 right-1"
                @click="copyCurl('/api/scan', 'POST')"
              />
            </div>
          </v-card-text>
        </v-card>
        
        <div v-if="!loading" class="text-center text-caption text-medium-emphasis mt-4">
          {{ filteredEndpoints.length }} 个接口
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useDisplay } from 'vuetify';

const { mobile } = useDisplay();
const loading = ref(true);
const showSnackbar = ref(false);
const snackbarMessage = ref('');
const snackbarColor = ref('info');
const searchQuery = ref('');
const activeIndex = ref(0);
const showMobileMenu = ref(false);

const endpoints = [
  { method: 'GET', path: '/api/status', title: '获取所有版本状态', desc: '返回所有启动器及完整版本列表，包含版本号、发布时间、下载链接' },
  { method: 'GET', path: '/api/status/{launcher}', title: '获取指定启动器状态', desc: '返回特定启动器的历史版本信息' },
  { method: 'GET', path: '/api/latest', title: '获取所有最新版本', desc: '快速检查所有启动器的最新版本号' },
  { method: 'GET', path: '/api/latest/{launcher}', title: '获取指定启动器最新版本', desc: '查询单个启动器的最新发布版本' },
  { method: 'GET', path: '/api/stats', title: '获取统计数据', desc: '访问统计、下载量、热门排行、地域分布' },
  { method: 'POST', path: '/api/scan', title: '触发手动扫描', desc: '强制同步上游仓库检查新版本（需认证）' },
];

const getMethodColor = (method) => ({
  GET: 'primary',
  POST: 'success',
  PUT: 'warning',
  DELETE: 'error'
}[method] || 'grey');

const copyCurl = async (path, method) => {
  const cmd = `curl -X ${method} "https://mirror.lemwood.icu${path}"`;
  try {
    await navigator.clipboard.writeText(cmd);
    snackbarMessage.value = '已复制';
    snackbarColor.value = 'success';
  } catch {
    snackbarMessage.value = '复制失败';
    snackbarColor.value = 'error';
  } finally {
    showSnackbar.value = true;
  }
};

const scrollTo = (index) => {
  const el = document.getElementById(`endpoint-${index}`);
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    activeIndex.value = index;
  }
};

const filteredEndpoints = computed(() => {
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return endpoints;
  
  return endpoints.filter(e => 
    e.path.toLowerCase().includes(query) ||
    e.title.toLowerCase().includes(query) ||
    e.desc.toLowerCase().includes(query)
  );
});

const handleScroll = () => {
  const cards = document.querySelectorAll('.endpoint-card');
  cards.forEach((card, index) => {
    const rect = card.getBoundingClientRect();
    if (rect.top <= 100 && rect.bottom >= 100) {
      activeIndex.value = index;
    }
  });
};

onMounted(() => {
  setTimeout(() => loading.value = false, 200);
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
.font-mono {
  font-family: 'Roboto Mono', 'Fira Code', 'Consolas', monospace !important;
}

.bg-black {
  background-color: #0d1117 !important;
  border: 1px solid #30363d !important;
}

.text-red { color: #ff7b72 !important; }
.text-yellow { color: #d2a8ff !important; }
.text-green { color: #7ee787 !important; }
.text-cyan { color: #a5d6ff !important; }
.text-orange { color: #ffa657 !important; }
.text-purple { color: #d2a8ff !important; }
.text-grey { color: #8b949e !important; }

.v-chip {
  font-weight: 600 !important;
}

.absolute {
  position: absolute !important;
}

.relative {
  position: relative !important;
}

.v-card.endpoint-card {
  transition: all 0.2s ease;
  border-color: #30363d;
}

.v-card.endpoint-card:hover {
  border-color: #58a6ff;
  transform: translateY(-2px);
}

pre {
  line-height: 1.4;
  tab-size: 2;
}

.gap-1 {
  gap: 4px;
}
</style>
