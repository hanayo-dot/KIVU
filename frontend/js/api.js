// KIVU API Integration Client

const API_BASE = ""; // Relative URL since frontend is served by Gin on port 8080

const API = {
  // Helper for authorized fetch
  async request(endpoint, options = {}) {
    const token = localStorage.getItem("access_token");
    const headers = {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers,
    });

    if (response.status === 401 && !endpoint.includes("/auth/")) {
      localStorage.removeItem("access_token");
      localStorage.removeItem("farmer");
      localStorage.removeItem("auth");
      window.location.href = "login.html";
      return null;
    }

    return response.json();
  },

  // Auth API
  async login(phoneNumber, password) {
    const res = await this.request("/auth/login", {
      method: "POST",
      body: JSON.stringify({ phone_number: phoneNumber, password }),
    });

    if (res && res.data && res.data.access_token) {
      localStorage.setItem("access_token", res.data.access_token);
      localStorage.setItem("refresh_token", res.data.refresh_token);
      localStorage.setItem("farmer", JSON.stringify(res.data.farmer));
      localStorage.setItem("auth", "true");
      return { success: true, data: res.data };
    }
    return { success: false, error: res ? res.error : "Login failed" };
  },

  // Farm Dashboard API
  async getFarmDashboard(farmId) {
    return this.request(`/farms/${farmId}/dashboard`);
  },

  // Alerts API
  async getAlerts(scope = "") {
    const query = scope ? `?scope=${scope}` : "";
    return this.request(`/alerts${query}`);
  },

  // Lake Zones Spatial API
  async getLakeZones() {
    return this.request("/lake/zones");
  },

  // Submit Observation / Telemetry
  async submitCageReading(cageId, reading) {
    return this.request(`/cages/${cageId}/readings`, {
      method: "POST",
      body: JSON.stringify(reading),
    });
  },
};
