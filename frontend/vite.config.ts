import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// The backend: http://backend:8080 inside Docker (set in docker-compose.yml), localhost otherwise.
const apiProxyTarget = process.env.API_PROXY_TARGET || "http://localhost:8080";

export default defineConfig({
  plugins: [vue()],
  server: {
    host: "0.0.0.0",
    port: 5173,
    // Opened directly on port 5173 (without nginx): pass the API and WebSockets on to the backend,
    // so the website and the game can always use their own address.
    proxy: {
      "/api": apiProxyTarget,
      "/health": apiProxyTarget,
      "/transcriptions": apiProxyTarget,
      "/realtime": apiProxyTarget,
      "/ws": { target: apiProxyTarget, ws: true },
    },
  },
});
