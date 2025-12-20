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
          primary: '#1867C0', // Classic Vuetify Blue
          secondary: '#5CBBF6',
          surface: '#FFFFFF',
          'surface-light': '#F5F5F5',
          background: '#FFFFFF',
          error: '#FF5252',
          info: '#2196F3',
          success: '#4CAF50',
          warning: '#FB8C00',
        },
      },
      dark: {
        dark: true,
        colors: {
            primary: '#2196F3', // Blue 500
            secondary: '#424242', // Grey 800
            surface: '#121212',
            'surface-light': '#212121',
            background: '#000000',
            error: '#FF5252',
            info: '#2196F3',
            success: '#4CAF50',
            warning: '#FB8C00',
        }
      }
    },
  },
})
