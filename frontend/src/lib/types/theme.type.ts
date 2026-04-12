export const THEME = {
  LIGHT: "light",
  DARK: "dark",
} as const;

export type ThemeKey = keyof typeof THEME;
export type ThemeValue = (typeof THEME)[ThemeKey];
