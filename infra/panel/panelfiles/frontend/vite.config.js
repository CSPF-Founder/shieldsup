import { defineConfig } from "vite";
import path from "path";
import inject from "@rollup/plugin-inject";

export default defineConfig({
  base: "",
  root: path.resolve(__dirname, "src"),

  resolve: {
    alias: {
      "~coreui": path.resolve(__dirname, "node_modules/@coreui/coreui-pro"),
    },
  },
  build: {
    minify: true,
    manifest: true,
    rollupOptions: {
      input: {
        scans: "./src/app/scans.js",
        update: "./src/app/updater.js",
        bug_track: "./src/app/bug_track.js",
        scan_result: "./src/app/scan_result.js",
        main: "./src/app/main.js",
        app: "./src/scss/app.scss",
      },
    },
    outDir: "../static",
  },
  plugins: [
    inject({
      include: "**/*.js", // Only include JavaScript files
      exclude: "**/*.scss", // Exclude SCSS files
      $: "jquery",
      jQuery: "jquery",
    }),
  ],
});
