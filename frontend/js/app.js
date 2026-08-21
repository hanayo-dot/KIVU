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
      // Show User Profile, Hide Login
      if (loginForm) loginForm.style.display = "none";
      if (profileSection) profileSection.style.display = "flex";
      if (farmSelector) farmSelector.style.display = "block";
      
      // Unlock Dashboard & restart animations
      if (dashboardLayout) {
        dashboardLayout.classList.remove("dashboard-content-locked");
        
        // Retrigger animations by resetting the DOM nodes briefly
        const animatedElements = document.querySelectorAll('.animate-up');
        animatedElements.forEach(el => {
          el.style.animation = 'none';
          el.offsetHeight; // Trigger reflow
          el.style.animation = null; 
        });
      }
    } else {
      // Show Login, Hide User Profile
      if (loginForm) loginForm.style.display = "flex";
      if (profileSection) profileSection.style.display = "none";
      if (farmSelector) farmSelector.style.display = "none";
      
      // Lock Dashboard (Blurs the background)
      if (dashboardLayout) {
        dashboardLayout.classList.add("dashboard-content-locked");
      }
    }
  };

  // 3. Check initial state on page load
  const isAuthenticated = localStorage.getItem("auth") === "true";
  updateAuthUI(isAuthenticated);

  // 4. Handle Top-Bar Login
  if (loginForm) {
    loginForm.addEventListener("submit", (e) => {
      e.preventDefault();
      localStorage.setItem("auth", "true");
      updateAuthUI(true);
    });
  }

  // 5. Handle Logout
  if (logoutBtn) {
    logoutBtn.addEventListener("click", () => {
      localStorage.removeItem("auth");
      updateAuthUI(false);
    });
  }

  // 6. Handle Observation Form Submission & Toast Animation (Reports Page)
  const obsForm = document.getElementById("observation-form");
  if (obsForm) {
    obsForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const toast = document.querySelector(".toast-popup");
      if (toast) {
        toast.classList.add("show");
        setTimeout(() => {
          toast.classList.remove("show");
        }, 4000);
      }
    });
  }
});

// 7. Category Selector Toggle (Reports Page)
function categoryToggle(element) {
  document.querySelectorAll(".cat-card-option").forEach((btn) => {
    btn.classList.remove("active-red");
  });
  element.classList.add("active-red");
}