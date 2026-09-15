export type Greeting = {
  text: string;
};

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "/api";

export async function readGreeting(): Promise<Greeting> {
  const response = await fetch(`${apiBase}/v1/greeting`, {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Failed to read greeting.");
  }

  return response.json();
}

export async function saveGreeting(text: string): Promise<Greeting> {
  const response = await fetch(`${apiBase}/v1/greeting`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ text }),
  });

  if (!response.ok) {
    throw new Error("Failed to save greeting.");
  }

  return response.json();
}

export async function readGreetingFromServer(): Promise<Greeting> {
  const apiOrigin = process.env.API_ORIGIN ?? "http://backend:8080";
  const response = await fetch(`${apiOrigin}/v1/greeting`, {
    cache: "no-store",
  });

  if (!response.ok) {
    throw new Error("Failed to read greeting.");
  }

  return response.json();
}
