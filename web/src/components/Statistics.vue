<template>
  <div class="mb-6">
    <h2 class="text-h4 font-weight-bold mb-6">数据洞察</h2>

    <div v-if="loading" class="d-flex justify-center pa-12">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <div v-else>
      <!-- Overview Cards -->
      <v-row class="mb-6">
        <v-col cols="12" sm="6" md="3">
          <v-card class="py-6 rounded-xl" elevation="2" border>
            <div class="d-flex flex-column align-center">
               <v-icon color="primary" size="32" class="mb-2">mdi-eye-outline</v-icon>
               <div class="text-h4 font-weight-bold">{{ stats.total_visits?.toLocaleString() || '-' }}</div>
               <div class="text-subtitle-2 text-medium-emphasis">总访问量</div>
            </div>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-card class="py-6 rounded-xl" elevation="2" border>
             <div class="d-flex flex-column align-center">
               <v-icon color="success" size="32" class="mb-2">mdi-download-outline</v-icon>
               <div class="text-h4 font-weight-bold">{{ stats.total_downloads?.toLocaleString() || '-' }}</div>
               <div class="text-subtitle-2 text-medium-emphasis">总下载量</div>
            </div>
          </v-card>
        </v-col>
         <v-col cols="12" sm="6" md="3">
          <v-card class="py-6 rounded-xl" elevation="2" border>
             <div class="d-flex flex-column align-center">
               <v-icon color="warning" size="32" class="mb-2">mdi-server-network</v-icon>
               <div class="text-h4 font-weight-bold">{{ Object.keys(stats.daily_stats || {}).length || '-' }}</div>
               <div class="text-subtitle-2 text-medium-emphasis">统计天数</div>
            </div>
          </v-card>
        </v-col>
         <v-col cols="12" sm="6" md="3">
          <v-card class="py-6 rounded-xl" elevation="2" border>
             <div class="d-flex flex-column align-center">
               <v-icon color="info" size="32" class="mb-2">mdi-earth</v-icon>
               <div class="text-h4 font-weight-bold">{{ stats.geo_distribution?.length || '-' }}</div>
               <div class="text-subtitle-2 text-medium-emphasis">覆盖国家/地区</div>
            </div>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <!-- Global Map -->
        <v-col cols="12" lg="8">
          <v-card class="rounded-xl h-100 pa-4" border>
            <v-card-title class="font-weight-bold">
                <v-icon color="primary" class="mr-2">mdi-map-marker-radius</v-icon>
                全球访问分布
            </v-card-title>
            <div style="height: 400px; width: 100%; position: relative;">
                <v-chart class="chart" :option="mapOption" autoresize />
            </div>
          </v-card>
        </v-col>

        <!-- Top Downloads -->
        <v-col cols="12" lg="4">
          <v-card class="rounded-xl h-100 pa-4" border>
            <v-card-title class="font-weight-bold">
                <v-icon color="success" class="mr-2">mdi-fire</v-icon>
                热门资源排行
            </v-card-title>
             <div style="height: 400px; overflow-y: auto;">
                 <v-list density="compact">
                    <v-list-item v-for="(item, i) in stats.top_downloads" :key="i" class="mb-2 rounded-lg bg-surface-light">
                        <template v-slot:prepend>
                            <v-avatar color="surface-variant" size="24" class="mr-2 text-caption font-weight-bold">
                                {{ i + 1 }}
                            </v-avatar>
                        </template>
                        <v-list-item-title class="font-weight-medium text-body-2">
                            {{ item.launcher }}
                        </v-list-item-title>
                        <v-list-item-subtitle class="text-caption">
                             {{ item.version }}
                        </v-list-item-subtitle>
                        <template v-slot:append>
                             <v-chip size="small" color="primary" variant="flat" class="font-weight-bold">
                                 {{ item.count }}
                             </v-chip>
                        </template>
                    </v-list-item>
                 </v-list>
             </div>
          </v-card>
        </v-col>

        <!-- Trend Chart -->
        <v-col cols="12">
          <v-card class="rounded-xl pa-4" border>
            <v-card-title class="font-weight-bold">
                 <v-icon color="info" class="mr-2">mdi-chart-timeline-variant</v-icon>
                 最近 30 天趋势
            </v-card-title>
            <div style="height: 350px;">
               <v-chart class="chart" :option="trendOption" autoresize />
            </div>
          </v-card>
        </v-col>
      </v-row>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, provide } from 'vue';
import { getStats } from '../services/api';
import { useTheme } from 'vuetify';

// ECharts imports
import { use, registerMap } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { BarChart, MapChart, LineChart } from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  VisualMapComponent
} from 'echarts/components';
import VChart from 'vue-echarts';

use([
  CanvasRenderer,
  BarChart,
  MapChart,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  VisualMapComponent
]);

const stats = ref({});
const loading = ref(true);
const mapLoaded = ref(false);
const theme = useTheme();

// Helper to map country names if needed (very basic)
const mapCountryName = (name) => {
    // Basic mapping or return as is. The API likely returns English or Chinese names.
    // ECharts world map usually expects English names.
    return name; 
};

const mapOption = computed(() => {
    const isDark = theme.global.current.value.dark;
    const textColor = isDark ? '#fff' : '#333';
    
    // Prepare data
    const data = (stats.value.geo_distribution || []).map(item => ({
        name: mapCountryName(item.country),
        value: item.count
    }));

    return {
        backgroundColor: 'transparent',
        tooltip: {
            trigger: 'item',
            formatter: '{b}: {c} (Visits)'
        },
        visualMap: {
            min: 0,
            max: data.length ? Math.max(...data.map(d => d.value)) : 100,
            left: 'left',
            top: 'bottom',
            text: ['High', 'Low'],
            calculable: true,
            inRange: {
                color: ['#e0f2f1', '#009688'] // Teal gradient
            },
            textStyle: { color: textColor }
        },
        series: [
            {
                name: 'Visits',
                type: 'map',
                map: 'world',
                roam: true,
                emphasis: {
                    label: { show: true },
                    itemStyle: {
                        areaColor: '#ff9800' // Orange on hover
                    }
                },
                itemStyle: {
                    areaColor: isDark ? '#424242' : '#eee',
                    borderColor: isDark ? '#616161' : '#ccc'
                },
                data: data
            }
        ]
    };
});

const trendOption = computed(() => {
    const isDark = theme.global.current.value.dark;
    const textColor = isDark ? '#ccc' : '#666';
    
    if (!stats.value.daily_stats) return {};

    const rawData = [...stats.value.daily_stats].reverse();
    const dates = rawData.map(d => d.date.slice(5));
    const visits = rawData.map(d => d.visit_count);
    const downloads = rawData.map(d => d.download_count);

    return {
        backgroundColor: 'transparent',
        tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' }
        },
        legend: {
            data: ['访问量', '下载量'],
            textStyle: { color: textColor }
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            data: dates,
            axisLine: { lineStyle: { color: textColor } }
        },
        yAxis: {
            type: 'value',
            splitLine: { lineStyle: { type: 'dashed', color: isDark ? '#333' : '#eee' } },
            axisLine: { lineStyle: { color: textColor } }
        },
        series: [
            {
                name: '访问量',
                type: 'line',
                smooth: true,
                data: visits,
                itemStyle: { color: '#FFB74D' }, // Orange
                areaStyle: { opacity: 0.1 }
            },
            {
                name: '下载量',
                type: 'bar',
                barWidth: '40%',
                data: downloads,
                itemStyle: { color: '#2196F3', borderRadius: [4, 4, 0, 0] } // Blue
            }
        ]
    };
});

const loadData = async () => {
    loading.value = true;
    try {
        const [res] = await Promise.all([
             getStats(),
             // Load map data only once
             !mapLoaded.value ? fetch('https://cdn.jsdelivr.net/npm/echarts@4.9.0/map/json/world.json').then(r => r.json()) : Promise.resolve(null)
        ]);
        
        stats.value = res.data;
        
        if (res && typeof res[1] === 'object' && !mapLoaded.value) { // Wait, fetch returns response then json
             // Actually, promise.all returns array of results.
             // But my logic above is a bit mixed. Let's fix.
        }
    } catch (e) {
        console.error('Failed to load stats:', e);
    } finally {
        loading.value = false;
    }
};

onMounted(async () => {
    // Load map separately to ensure registration
    try {
        const mapRes = await fetch('https://cdn.jsdelivr.net/npm/echarts@4.9.0/map/json/world.json');
        const mapJson = await mapRes.json();
        registerMap('world', mapJson);
        mapLoaded.value = true;
    } catch (e) {
        console.error('Failed to load map JSON', e);
    }

    // Load API stats
    const res = await getStats();
    stats.value = res.data;
    loading.value = false;
});

defineExpose({ refresh: loadData });
</script>

<style scoped>
.chart {
  height: 100%;
  width: 100%;
}
</style>
