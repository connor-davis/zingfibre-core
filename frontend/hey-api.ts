import { type UserConfig } from '@hey-api/openapi-ts';

export default {
  input: {
    path: 'http://localhost:4000/api/api-spec',
  },
  output: {
    lint: 'eslint',
    path: 'src/api-client',
  },
  plugins: [
    '@tanstack/react-query',
    '@hey-api/client-fetch',
    '@hey-api/typescript',
    '@hey-api/schemas',
  ],
} as UserConfig;
