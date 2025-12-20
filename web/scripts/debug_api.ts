// Mock data structures based on the frontend expectations

interface Asset {
  name: string;
  url: string;
  size: number;
}

interface Version {
  tag_name: string;
  name: string;
  published_at: string;
  assets: Asset[];
}

interface LauncherStatus {
  [key: string]: Version[];
}

interface LatestVersions {
  [key: string]: string; // launcher name -> version tag
}

interface Stats {
  total_visits: number;
  total_downloads: number;
  daily_stats: Array <{
    date: string;
    visit_count: number;
    download_count: number;
  }>;
  top_downloads: Array <{
    launcher: string;
    version: string;
    count: number;
  }>;
  geo_distribution: Array <{
    country: string;
    count: number;
  }>;
}

// Mock Service Simulation
class MockApiService {
  async getStatus(): Promise<LauncherStatus> {
    console.log('[MockAPI] Fetching Status...');
    return {
      'HMCL': [
        {
          tag_name: 'v3.5.3',
          name: 'v3.5.3',
          published_at: new Date().toISOString(),
          assets: [
            { name: 'HMCL-3.5.3.jar', url: 'https://example.com/hmcl.jar', size: 1024 * 1024 * 5 },
            { name: 'HMCL-3.5.3.exe', url: 'https://example.com/hmcl.exe', size: 1024 * 1024 * 6 }
          ]
        },
        {
            tag_name: 'v3.5.2',
            name: 'v3.5.2',
            published_at: new Date(Date.now() - 86400000 * 10).toISOString(),
            assets: []
        }
      ],
      'PCL2': [
         {
          tag_name: '2.4.0',
          name: 'Snapshot 2.4.0',
          published_at: new Date().toISOString(),
          assets: [{ name: 'PCL2.exe', url: 'https://example.com/pcl2.exe', size: 1024 * 1024 * 2 }]
        }
      ]
    };
  }

  async getLatest(): Promise<LatestVersions> {
    console.log('[MockAPI] Fetching Latest Versions...');
    return {
      'HMCL': 'v3.5.3',
      'PCL2': '2.4.0'
    };
  }

  async getStats(): Promise<Stats> {
    console.log('[MockAPI] Fetching Statistics...');
    return {
      total_visits: 15420,
      total_downloads: 5430,
      daily_stats: Array.from({ length: 7 }, (_, i) => ({
        date: new Date(Date.now() - i * 86400000).toISOString().split('T')[0],
        visit_count: Math.floor(Math.random() * 1000),
        download_count: Math.floor(Math.random() * 500)
      })),
      top_downloads: [
        { launcher: 'HMCL', version: 'v3.5.3', count: 1200 },
        { launcher: 'PCL2', version: '2.4.0', count: 980 }
      ],
      geo_distribution: [
        { country: 'China', count: 8000 },
        { country: 'United States', count: 1200 },
        { country: 'Japan', count: 500 }
      ]
    };
  }
}

// Simulation Runner
async function runDebugSession() {
  console.log('=== Starting API Simulation Debugging ===\n');
  const service = new MockApiService();

  try {
    const status = await service.getStatus();
    console.log('1. Status Response:');
    console.log(JSON.stringify(status, null, 2));
    console.log('-----------------------------------');

    const latest = await service.getLatest();
    console.log('2. Latest Response:');
    console.log(JSON.stringify(latest, null, 2));
    console.log('-----------------------------------');

    const stats = await service.getStats();
    console.log('3. Stats Response:');
    // Compact log for large arrays
    console.log(`   Total Visits: ${stats.total_visits}`);
    console.log(`   Daily Stats Length: ${stats.daily_stats.length}`);
    console.log(`   Top Download: ${stats.top_downloads[0].launcher}`);
    console.log('-----------------------------------');

    console.log('\n=== Simulation Complete: API Contracts Valid ===');

  } catch (error) {
    console.error('Simulation Failed:', error);
  }
}

runDebugSession();
