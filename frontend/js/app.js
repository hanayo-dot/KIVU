document.addEventListener("DOMContentLoaded", () => {
  // 1. Top-bar Auth Elements
  const loginForm = document.getElementById("topbar-login-form");
  const profileSection = document.getElementById("user-profile-section");
  const dashboardLayout = document.getElementById("dashboard-main-content");
  const farmSelector = document.getElementById("farm-selector");
  const logoutBtn = document.getElementById("logout-btn");

  // 2. Auth UI State Manager (Handles blur/unblur and animations)
  const updateAuthUI = (isLoggedIn) => {
    if (isLoggedIn) {
      if (loginForm) loginForm.style.display = "none";
      if (profileSection) profileSection.style.display = "flex";
      if (farmSelector) farmSelector.style.display = "block";
      
      if (dashboardLayout) {
        dashboardLayout.classList.remove("dashboard-content-locked");
        // Retrigger entrance animations
        const animatedElements = document.querySelectorAll('.animate-up');
        animatedElements.forEach(el => {
          el.style.animation = 'none';
          el.offsetHeight; 
          el.style.animation = null; 
        });
      }
    } else {
      if (loginForm) loginForm.style.display = "flex";
      if (profileSection) profileSection.style.display = "none";
      if (farmSelector) farmSelector.style.display = "none";
      
      if (dashboardLayout) {
        dashboardLayout.classList.add("dashboard-content-locked");
      }
    }
  };

  const isAuthenticated = localStorage.getItem("auth") === "true";
  updateAuthUI(isAuthenticated);

  // 3. Handle Top-Bar Login & Logout
  // 1. Handle Login Form Submission via API (or mock email fallback)
  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      localStorage.setItem("auth", "true");
      updateAuthUI(true);
    });
  }

  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("auth");
      updateAuthUI(false);
    });
  }

  // 4. Handle Observation Form Submission (Reports Page)
      const phoneInput = document.getElementById("phone") || document.getElementById("login-email");
      const passwordInput = document.getElementById("password") || document.getElementById("login-password");
      const errorMsg = document.getElementById("error-msg") || document.getElementById("login-error");

      if (errorMsg) errorMsg.style.display = "none";

      const phone = phoneInput ? phoneInput.value : "";
      const password = passwordInput ? passwordInput.value : "";

      if (typeof API !== "undefined") {
        const res = await API.login(phone, password);
        if (res.success) {
          window.location.href = "index.html";
        } else {
          if (errorMsg) {
            errorMsg.innerText = res.error || "Authentication failed. Please check credentials.";
            errorMsg.style.display = "block";
          }
        }
      } else {
        localStorage.setItem("auth", "true");
        window.location.href = "index.html";
      }
    });
  }

  // 2. Enforce Session Security Guard across protected pages
  if (!window.location.pathname.includes("login.html")) {
    const token = localStorage.getItem("access_token");
    const isAuth = localStorage.getItem("auth");
    if (!token && isAuth !== "true") {
      window.location.href = "login.html";
    }
  }

  // 3. Update Farmer Profile in Top Bar if available
  const farmerData = localStorage.getItem("farmer");
  if (farmerData) {
    try {
      const farmer = JSON.parse(farmerData);
      const profileName = document.querySelector(".user-profile span:nth-child(2)");
      const greeting = document.querySelector(".dashboard-greeting");
      if (profileName && farmer.name) profileName.innerText = farmer.name;
      if (greeting && farmer.name) greeting.innerText = `Good Day, ${farmer.name}`;
    } catch (err) {
      console.warn("Could not parse farmer data", err);
    }
  }

  // 4. Handle Observation Form Submission & Toast Animation
  const obsForm = document.getElementById("observation-form");
  if (obsForm) {
    obsForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const toast = document.querySelector(".toast-popup") || document.getElementById("success-toast");
      if (toast) {
        toast.classList.add("show");
        setTimeout(() => toast.classList.remove("show"), 4000);
      }
    });
  }

  // 5. Interactive Time Tabs (e.g., Analytics/Dashboard: 24H, 7D, 30D, etc.)
  // 5. Interactive Filter Tabs (Alerts page: All, High, Medium, Low)
  const filterTabs = document.querySelectorAll(".filter-tabs span");
  filterTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      filterTabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
    });
  });

  // 6. Interactive Time Tabs (Analytics/Dashboard: 24H, 7D, 30D, etc.)
  const timeTabs = document.querySelectorAll(".time-tabs span");
  timeTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      timeTabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
    });
  });

  // 7. Interactive Category Cards on Reports Page
  const catCards = document.querySelectorAll(".cat-card-option");
  catCards.forEach((card) => {
    card.addEventListener("click", () => {
      catCards.forEach((c) => c.classList.remove("active-red"));
      card.classList.add("active-red");
    });
  });

  // 8. Interactive Severity Selector on Reports Page
  const severityItems = document.querySelectorAll(".severity-seg-item");
  severityItems.forEach((item) => {
    item.addEventListener("click", () => {
      severityItems.forEach((i) => i.classList.remove("active"));
      item.classList.add("active");
    });
  });
});

  // 9. Handle Logout from Settings Page
  const settingsLogoutBtn = document.getElementById("settings-logout-btn");
  if (settingsLogoutBtn) {
    settingsLogoutBtn.addEventListener("click", () => {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("farmer");
      localStorage.removeItem("auth");
      window.location.href = "login.html";
    });
  }
});

// Category Selector Toggle for Reports Page
function categoryToggle(element) {
  document.querySelectorAll(".cat-btn").forEach((btn) => {
    btn.classList.remove("active");
    btn.style.background = "var(--bg-color)";
    btn.style.borderColor = "var(--card-border)";
    btn.style.color = "var(--text-muted)";
  });
  element.classList.add("active");
  element.style.background = "rgba(239, 68, 68, 0.1)";
  element.style.borderColor = "var(--accent-red)";
  element.style.color = "var(--text-main)";
}
