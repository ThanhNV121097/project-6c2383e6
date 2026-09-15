export type Greeting = {
  text: string;
};

const STORAGE_KEY = "hello-world-acceptance-6:greeting";
const DEFAULT_GREETING = "Hello, World!";

export const mockGreeting: Greeting = {
  text: DEFAULT_GREETING,
};

export function readStoredGreeting(): Greeting {
  if (typeof window === "undefined") return mockGreeting;
  return { text: window.localStorage.getItem(STORAGE_KEY) ?? DEFAULT_GREETING };
}

export function saveStoredGreeting(text: string): Greeting {
  const greeting = { text };
  window.localStorage.setItem(STORAGE_KEY, greeting.text);
  return greeting;
}
