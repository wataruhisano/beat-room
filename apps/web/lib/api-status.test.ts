import { afterEach, describe, expect, it, vi } from "vitest";
import { fetchApiStatus } from "./api-status";

describe("fetchApiStatus", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  it("returns ok when the API responds with 200", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue(new Response(null, { status: 200 }));
    vi.stubGlobal("fetch", fetchMock);

    await expect(fetchApiStatus("http://api.test")).resolves.toBe("ok");
    expect(fetchMock).toHaveBeenCalledWith(
      "http://api.test/healthz",
      expect.anything(),
    );
  });

  it("returns unavailable when the API responds with an error", async () => {
    vi.spyOn(console, "error").mockImplementation(() => {});
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(new Response(null, { status: 503 })),
    );

    await expect(fetchApiStatus("http://api.test")).resolves.toBe(
      "unavailable",
    );
  });

  it("returns unavailable when the request fails", async () => {
    vi.spyOn(console, "error").mockImplementation(() => {});
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new TypeError("fetch failed")),
    );

    await expect(fetchApiStatus("http://api.test")).resolves.toBe(
      "unavailable",
    );
  });

  it("returns unavailable when API_URL is not set", async () => {
    await expect(fetchApiStatus(undefined)).resolves.toBe("unavailable");
  });
});
