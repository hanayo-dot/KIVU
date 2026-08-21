document.addEventListener("DOMContentLoaded", () => {
  // 1. Handle Login Form Submission & Validation
  const loginForm = document.getElementById("login-form");
  if (loginForm) {
    loginForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const emailInput = document.getElementById("login-email").value;
      const passwordInput = document.getElementById("login-password").value;
      const errorBox = document.getElementById("login-error");

      // Simple mock authentication check
      if (
        emailInput === "farmer@lakeview.com" &&
        passwordInput === "password123"
      ) {
        localStorage.setItem("auth", "true");
        localStorage.setItem("userName", "Farmer John");
        window.location.href = "index.html";
      } else {
        if (errorBox) errorBox.style.display = "block";
      }
    });
  }

  // 2. Enforce Session Security Guard across protected pages
  if (!window.location.pathname.includes("login.html")) {
    if (localStorage.getItem("auth") !== "true") {
      window.location.href = "login.html";
    }
  }

  // 3. Handle Observation Form Submission & Toast Animation on Reports Page
  const obsForm = document.getElementById("observation-form");
  if (obsForm) {
    obsForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const toast = document.querySelector(".toast-popup");
      if (toast) {
        toast.style.display = "flex";
        setTimeout(() => {
          toast.style.display = "none";
        }, 4000);
      }
    });
  }

  // 4. Interactive Filter Tabs (e.g., Alerts page: All, High, Medium, Low)
  const filterTabs = document.querySelectorAll(".filter-tabs span");
  filterTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      filterTabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
    });
  });

  // 5. Interactive Time Tabs (e.g., Analytics/Dashboard: 24H, 7D, 30D, etc.)
  const timeTabs = document.querySelectorAll(".time-tabs span");
  timeTabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      timeTabs.forEach((t) => t.classList.remove("active"));
      tab.classList.add("active");
    });
  });

  // 6. Interactive Category Cards on Reports Page
  const catCards = document.querySelectorAll(".cat-card-option");
  catCards.forEach((card) => {
    card.addEventListener("click", () => {
      catCards.forEach((c) => c.classList.remove("active-red"));
      card.classList.add("active-red");
    });
  });

  // 7. Interactive Severity Selector on Reports Page
  const severityItems = document.querySelectorAll(".severity-seg-item");
  severityItems.forEach((item) => {
    item.addEventListener("click", () => {
      severityItems.forEach((i) => i.classList.remove("active"));
      item.classList.add("active");
    });
  });

  // 8. Handle Logout from Settings Page
  const settingsLogoutBtn = document.getElementById("settings-logout-btn");
  if (settingsLogoutBtn) {
    settingsLogoutBtn.addEventListener("click", () => {
      localStorage.removeItem("auth");
      localStorage.removeItem("userName");
      window.location.href = "login.html";
    });
  }
});
