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
  if (loginForm) {
    loginForm.addEventListener("submit", (e) => {
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
  const obsForm = document.getElementById("observation-form");
  if (obsForm) {
    obsForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const toast = document.querySelector(".toast-popup");
      if (toast) {
        toast.classList.add("show");
        setTimeout(() => toast.classList.remove("show"), 4000);
      }
    });
  }

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
});