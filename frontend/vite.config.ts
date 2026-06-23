import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// Minimal Vite config for the Call Booking SPA.
// The API base URL is configured via VITE_API_BASE_URL (see .env.example),
// not here, so the same build can target the Prism mock or a real backend.
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
  },
});
