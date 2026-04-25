import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import checker from "vite-plugin-checker";
import path from "path";

export default defineConfig({
  plugins: [
    vue(),
    checker({
      typescript: {
        tsconfigPath: "./tsconfig.app.json",
      },
      vueTsc: {
        tsconfigPath: "./tsconfig.app.json",
      },
    }),
  ],
  server: {
    port: 6502,
    open: true,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8086",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
  build: {
    outDir: path.resolve(__dirname, "../golang-backend/public"),
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            return "vendor";
          }
        },
      },
    },
  },
});
