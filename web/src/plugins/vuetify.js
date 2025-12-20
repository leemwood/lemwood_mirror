// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#2563EB', // Royal Blue
          secondary: '#10B981', // Emerald
          surface: '#FFFFFF',
          'surface-light': '#F8FAFC',
          background: '#F1F5F9', // Slate 100
          error: '#EF4444',
          info: '#3B82F6',
          success: '#22C55E',
          warning: '#F59E0B',
        },
      },
      dark: {
        dark: true,
        colors: {
            primary: '#60A5FA', // Blue 400
            secondary: '#34D399', // Emerald 400
            surface: '#1E293B', // Slate 800
            'surface-light': '#334155', // Slate 700
            background: '#0F172A', // Slate 900
            error: '#F87171',
            info: '#60A5FA',
            success: '#4ADE80',
            warning: '#FBBF24',
        }
      }
    },
  },
})
