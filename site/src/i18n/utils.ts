import en from "./en.json";
import id from "./id.json";
import type { AstroGlobal } from "astro";

export type Lang = "en" | "id";

const translations: Record<Lang, Record<string, string>> = { en, id };

export function t(lang: Lang, key: string): string {
  return translations[lang]?.[key] ?? translations.en[key] ?? key;
}

export function useTranslations(lang: Lang) {
  return (key: string) => t(lang, key);
}

export function path(url: URL, target: string): string {
  const base = import.meta.env.BASE_URL.replace(/\/$/, "");
  return `${base}${target}`;
}
