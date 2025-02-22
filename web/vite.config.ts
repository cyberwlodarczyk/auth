import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { Agent } from "node:https";
import { fileURLToPath, URL } from "node:url";

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      $lib: fileURLToPath(new URL("./src/lib", import.meta.url)),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "https://localhost:4000",
        agent: new Agent({
          rejectUnauthorized: false,
        }),
        rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
