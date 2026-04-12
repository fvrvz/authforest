import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import { escapeSvelte, mdsvex } from "mdsvex";
import { createHighlighter } from "shiki";

const highlighter = await createHighlighter({
  themes: ["github-light", "github-dark"],
  langs: [
    "javascript",
    "typescript",
    "tsx",
    "bash",
    "json",
    "yaml",
    "python",
    "go",
    "ini",
    "html",
    "angular-ts",
    "bash",
  ],
});

/** @type {import('@sveltejs/kit').Config} */
const config = {
  // Consult https://svelte.dev/docs/kit/integrations
  // for more information about preprocessors
  preprocess: [
    vitePreprocess(),
    mdsvex({
      extensions: [".svx", ".md"],
      highlight: {
        highlighter: async (code, lang) => {
          const html = escapeSvelte(
            highlighter.codeToHtml(code, {
              lang: lang || "text",
              themes: { light: "github-light", dark: "github-dark" },
            }),
          );
          return `{@html \`${html}\`}`;
        },
      },
    }),
  ],
  kit: {
    // adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
    // If your environment is not supported, or you settled on a specific environment, switch out the adapter.
    // See https://svelte.dev/docs/kit/adapters for more information about adapters.
    adapter: adapter({ fallback: "index.html" }),

    paths: {
      base: process.env.PUBLIC_BASE_PATH || "",
    },
  },
  extensions: [".svelte", ".svx", ".md"],
};

export default config;
