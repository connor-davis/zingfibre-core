exports.apps = [
  {
    name: "zingfibrecore-api",
    script: "/home/connor/kalimbu/cmd/api/main.go",
    interpreter: "go",
    interpreter_args: "run",
  },
  {
    name: "zingfibrecore-app",
    script: "serve",
    env: {
      PM2_SERVE_PATH: "/home/connor/zingfibre-upgrade/frontend/dist",
      PM2_SERVE_PORT: 3000,
      PM2_SERVE_SPA: "true",
      PM2_SERVE_HOMEPAGE: "/index.html",
    },
  },
];
