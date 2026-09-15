import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{js,ts,jsx,tsx}", "./components/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: { background: "#FFFFFF", text: "#000000", primary: "#2563EB" },
    },
  },
  plugins: [],
};

export default config;
