/* eslint-env node */
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

const apiProxyTarget =
  process.env.VITE_API_PROXY_TARGET || 'http://localhost:8080';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
    },
  },
});
