document.addEventListener("DOMContentLoaded", function () {
  const tabs = document.querySelectorAll(".tab");
  const tabContents = document.querySelectorAll(".tab-content");
  const nextButton = document.getElementById("next-button");
  const errorMessage = document.getElementById("error-message");

  // Tab switching functionality
  tabs.forEach((tab) => {
    tab.addEventListener("click", () => {
      // Remove active class from all tabs and contents
      tabs.forEach((t) => t.classList.remove("active"));
      tabContents.forEach((content) => content.classList.remove("active"));

      // Add active class to clicked tab
      tab.classList.add("active");

      // Show corresponding content
      const tabId = tab.getAttribute("data-tab");
      document.getElementById(`${tabId}-content`).classList.add("active");
    });
  });

  // Handle next button click for webhook setup
  if (nextButton) {
    nextButton.addEventListener("click", async function () {
      const webhookSecret = document
        .getElementById("webhook-secret")
        .value.trim();

      const payload = {
        webhook_secret: webhookSecret || null,
      };

      // Disable button and show loading state
      nextButton.disabled = true;
      nextButton.classList.add("loading");
      errorMessage.style.display = "none";

      try {
        const response = await fetch("/api/initiate-webhook", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
        });

        if (response.ok) {
          window.location.href = "/dns/setup";
        } else {
          throw new Error("Server error");
        }
      } catch (error) {
        errorMessage.textContent =
          "We encountered some error, look into the pod logs";
        errorMessage.style.display = "block";
        nextButton.disabled = false;
        nextButton.classList.remove("loading");
      }
    });
  }
});
