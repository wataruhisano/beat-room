export type ApiStatus = "ok" | "unavailable";

export async function fetchApiStatus(
  apiURL: string | undefined,
): Promise<ApiStatus> {
  if (!apiURL) {
    return "unavailable";
  }

  try {
    const res = await fetch(`${apiURL}/healthz`, {
      cache: "no-store",
      signal: AbortSignal.timeout(2000),
    });
    if (!res.ok) {
      console.error(`API health check failed: ${res.status}`);
      return "unavailable";
    }
    return "ok";
  } catch (err) {
    console.error("API health check failed", err);
    return "unavailable";
  }
}
