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
    {
      dates: true,
      name: '@hey-api/transformers',
      bigInt: true,
    },
    {
      enums: 'javascript',
      name: '@hey-api/typescript',
    },
    {
      name: '@hey-api/sdk',
      transformer: true,
    },
  ],
} as UserConfig;
