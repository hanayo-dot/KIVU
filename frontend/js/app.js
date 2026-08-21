document.addEventListener("DOMContentLoaded", () => {
  // Initialize Lucide Icons
  if (typeof lucide !== "undefined") {
    lucide.createIcons();
  }

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
        const animatedElements = document.querySelectorAll(".animate-up");
        animatedElements.forEach((el) => {
          el.style.animation = "none";
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

  // 3. Handle Top-Bar / Standard Login Form Submission
  const standardLoginForm = document.getElementById("login-form");
  if (standardLoginForm) {
    standardLoginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const phoneInput =
        document.getElementById("phone") ||
        document.getElementById("login-email");
      const passwordInput =
        document.getElementById("password") ||
        document.getElementById("login-password");
      const errorMsg =
        document.getElementById("error-msg") ||
        document.getElementById("login-error");

      if (errorMsg) errorMsg.style.display = "none";

      const phone = phoneInput ? phoneInput.value : "";
      const password = passwordInput ? passwordInput.value : "";

      // Extract a clean display name dynamically from the input email/phone if desired
      let derivedName = "Farmer John";
      if (phone && phone.includes("@")) {
        let namePart = phone.split("@")[0];
        derivedName = namePart
          .split(".")
          .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
          .join(" ");
      } else if (phone) {
        derivedName = `User ${phone.slice(-4)}`;
      }

      if (typeof API !== "undefined") {
        const res = await API.login(phone, password);
        if (res.success) {
          localStorage.setItem("farmerName", res.name || derivedName);
          localStorage.setItem("auth", "true");
          window.location.href = "index.html";
        } else {
          if (errorMsg) {
            errorMsg.innerText =
              res.error || "Authentication failed. Please check credentials.";
            errorMsg.style.display = "block";
          }
        }
      } else {
        localStorage.setItem("farmerName", derivedName);
        localStorage.setItem("auth", "true");
        window.location.href = "index.html";
      }
    });
  }

  // 4. Handle Logout
  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("auth");
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("farmerName");
      localStorage.removeItem("farmerEmail");
      updateAuthUI(false);
    });
  }

  // 5. Enforce Session Security Guard across protected pages
  if (!window.location.pathname.includes("login.html")) {
    const token = localStorage.getItem("access_token");
    const isAuth = localStorage.getItem("auth");
    if (!token && isAuth !== "true") {
      // Optional: uncomment below if you want strict redirection to login
      // window.location.href = "login.html";
    }
  }

  // 6. Update Farmer Profile in Top Bar and Greetings dynamically using IDs
  const savedFarmerName = localStorage.getItem("farmerName");
  if (savedFarmerName) {
    const topbarUserEl = document.getElementById("topbar-user-name");
    const greetingEl = document.getElementById("dashboard-greeting-title");

    if (topbarUserEl) {
      topbarUserEl.innerText = savedFarmerName;
    }
    if (greetingEl) {
      let firstName = savedFarmerName.split(" ")[0];
      greetingEl.innerText = `Good Morning, ${firstName}`;
    }
  }

  // 7. Handle Observation Form Submission & Toast Animation (Reports Page)
  const obsForm = document.getElementById("observation-form");
  if (obsForm) {
    obsForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const toast =
        document.querySelector(".toast-popup") ||
        document.getElementById("success-toast");
      if (toast) {
        toast.style.display = "flex";
        setTimeout(() => (toast.style.display = "none"), 4000);
      }
    });
  }

  // 8. Interactive Filter Tabs (Alerts page: All, High, Medium, Low)
  const filterTabs = document.querySelectorAll(
    ".filter-tabs span, .filter-tab",
  );
  filterTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      filterTabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
    });
  });

  // 9. Interactive Time Tabs (Analytics/Dashboard: 24H, 7D, 30D, etc.)
  const timeTabs = document.querySelectorAll(".time-tabs span");
  timeTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      timeTabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
    });
  });

  // 10. Interactive Category Cards on Reports Page
  const catCards = document.querySelectorAll(".cat-card-option");
  catCards.forEach((card) => {
    card.addEventListener("click", () => {
      catCards.forEach((c) => c.classList.remove("active-red"));
      card.classList.add("active-red");
    });
  });

  // 11. Interactive Severity Selector on Reports Page
  const severityItems = document.querySelectorAll(".severity-seg-item");
  severityItems.forEach((item) => {
    item.addEventListener("click", () => {
      severityItems.forEach((i) => {
        i.classList.remove("active");
        i.style.background = "transparent";
        i.style.color = "var(--text-muted)";
      });
      item.classList.add("active");
      item.style.background = "rgba(34,197,94,0.2)";
      item.style.color = "var(--accent-emerald)";
    });
  });

  // 12. Handle Logout from Settings Page
  const settingsLogoutBtn = document.getElementById("settings-logout-btn");
  if (settingsLogoutBtn) {
    settingsLogoutBtn.addEventListener("click", () => {
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("farmerName");
      localStorage.removeItem("farmerEmail");
      localStorage.removeItem("auth");
      window.location.href = "login.html";
    });
  }
});

// Category Selector Toggle Helper (Optional global fallback)
function categoryToggle(element) {
  document.querySelectorAll(".cat-card-option").forEach((btn) => {
    btn.classList.remove("active-red");
  });
  element.classList.add("active-red");
}
