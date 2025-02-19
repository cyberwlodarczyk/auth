import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import { Agent } from "node:https";

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
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
