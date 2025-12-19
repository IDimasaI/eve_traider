import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue({}), tailwindcss()],
  server: {
    proxy: {
      "/api": "http://localhost:6970/",
    },
  },
  build: {
    outDir: "../build/web/pages",
    assetsDir: "../static", // статика будет в web/static
    emptyOutDir: true,
    copyPublicDir: false,
  },
});
