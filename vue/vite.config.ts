import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue({}), tailwindcss()],
  server: {
    proxy: {
      "/api": "http://localhost:6969/",
    },
  },
  build: {
    outDir: "../local_server/build/eve_traiders/web/pages",
    assetsDir: "../static", // статика будет в web/static
    emptyOutDir: true,
    copyPublicDir: false,
  },
});
