import pluginVue from 'eslint-plugin-vue'
import tseslint from 'typescript-eslint'
import eslintConfigPrettier from 'eslint-config-prettier'
import boundaries from 'eslint-plugin-boundaries'

const FSD_LAYERS = ['app', 'pages', 'features', 'shared'] as const
const SLICED_LAYERS = new Set(['pages', 'features'])

export default tseslint.config(
  { ignores: ['dist', 'src/env.d.ts', 'src/shared/api/schema.d.ts'] },
  tseslint.configs.recommended,
  ...pluginVue.configs['flat/recommended'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
      },
    },
  },
  {
    rules: {
      'eqeqeq': 'error',
      'no-console': 'warn',
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/consistent-type-definitions': ['error', 'type'],
      '@typescript-eslint/no-restricted-types': [
        'error',
        {
          types: {
            '{}': { suggest: ['Record<string, unknown>'] },
          },
        },
      ],
      'no-restricted-syntax': [
        'error',
        {
          selector: 'TSEnumDeclaration',
          message: 'enum は使用禁止です。as const を使用してください。',
        },
      ],
      'prefer-arrow-callback': 'error',
      'func-style': ['error', 'expression'],
    },
  },
  // FSD boundaries
  {
    plugins: { boundaries },
    settings: {
      'import/resolver': {
        typescript: {
          project: './tsconfig.app.json',
        },
      },
      'boundaries/elements': FSD_LAYERS.map((layer) =>
        SLICED_LAYERS.has(layer)
          ? { type: layer, pattern: `src/${layer}/*/**`, capture: ['slice'] }
          : { type: layer, pattern: `src/${layer}/**` },
      ),
      'boundaries/ignore': ['**/*.test.*', '**/*.stories.*', '**/__tests__/**'],
    },
    rules: {
      'boundaries/dependencies': [
        'error',
        {
          default: 'disallow',
          rules: FSD_LAYERS.map((layer, index, layers) => ({
            from: { type: layer },
            allow: [
              ...layers.slice(index + 1).map((target) => ({ to: { type: target } })),
              ...(SLICED_LAYERS.has(layer)
                ? [{ to: { type: layer, captured: { slice: '{{from.captured.slice}}' } } }]
                : []),
            ],
          })),
        },
      ],
    },
  },
  eslintConfigPrettier,
)
