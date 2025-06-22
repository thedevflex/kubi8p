document.addEventListener("DOMContentLoaded", function () {
  const initiateButton = document.getElementById("initiate-button");
  const checkAgainButton = document.getElementById("check-again-button");
  const statusContainer = document.getElementById("status-container");
  const statusMessage = document.getElementById("status-message");
  const loadingSpinner = document.getElementById("loading-spinner");
  const errorMessage = document.getElementById("error-message");

  let pollingInterval;
  let hasCheckedAgain = false;
  const POLLING_DURATION = 20000; // 20 seconds
  const POLLING_INTERVAL = 2000; // 2 seconds

  function showLoading() {
    statusContainer.style.display = "block";
    loadingSpinner.style.display = "block";
    statusMessage.style.display = "none";
    errorMessage.style.display = "none";
    initiateButton.style.display = "none";
    checkAgainButton.style.display = "none";
  }

  function showError(message) {
    statusContainer.style.display = "block";
    loadingSpinner.style.display = "none";
    statusMessage.style.display = "block";
    statusMessage.textContent = message;
    statusMessage.style.color = "#ff6b6b";
    errorMessage.style.display = "none";
    initiateButton.style.display = "none";
    checkAgainButton.style.display = "block";
  }

  function showSuccess(message) {
    statusContainer.style.display = "block";
    loadingSpinner.style.display = "none";
    statusMessage.style.display = "block";
    statusMessage.textContent = message;
    statusMessage.style.color = "#0f0";
    errorMessage.style.display = "none";
    initiateButton.style.display = "none";
    checkAgainButton.style.display = "none";
  }

  function resetUI() {
    statusContainer.style.display = "none";
    errorMessage.style.display = "none";
    initiateButton.style.display = "block";
    checkAgainButton.style.display = "none";
    hasCheckedAgain = false;
  }

  async function initiateDNS() {
    showLoading();

    try {
      const response = await fetch("/api/initiate-dns", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to initiate DNS configuration");
      }

      // Start polling for status
      startPolling();
    } catch (error) {
      // Even if initiation fails, continue polling for 20 seconds
      startPolling();
    }
  }

  async function checkDNSStatus() {
    try {
      const response = await fetch("/api/get-dns-status", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        const data = await response.json();

        if (data.status === "ready" && data.external_ip) {
          clearInterval(pollingInterval);
          showSuccess(`DNS service is ready! External IP: ${data.external_ip}`);

          // Redirect to configure page after a short delay
          setTimeout(() => {
            window.location.href = "/dns/configure";
          }, 2000);
          return true;
        } else if (data.status === "pending") {
          // Continue polling
          return false;
        } else {
          // Service is in error state
          clearInterval(pollingInterval);
          showError(
            "DNS service configuration failed. Please check the logs in your cluster."
          );
          return true;
        }
      } else {
        // API error
        clearInterval(pollingInterval);
        showError(
          "Failed to check DNS status. Please check the logs in your cluster."
        );
        return true;
      }
    } catch (error) {
      // Network error
      clearInterval(pollingInterval);
      showError(
        "Network error while checking DNS status. Please check the logs in your cluster."
      );
      return true;
    }
  }

  function startPolling() {
    let elapsedTime = 0;

    pollingInterval = setInterval(async () => {
      elapsedTime += POLLING_INTERVAL;

      const isComplete = await checkDNSStatus();

      if (isComplete || elapsedTime >= POLLING_DURATION) {
        clearInterval(pollingInterval);

        if (!isComplete) {
          // Timeout reached
          showError(
            "DNS configuration is taking longer than expected. Please check the logs in your cluster."
          );
        }
      }
    }, POLLING_INTERVAL);
  }

  // Handle initiate button click
  if (initiateButton) {
    initiateButton.addEventListener("click", initiateDNS);
  }

  // Handle check again button click
  if (checkAgainButton) {
    checkAgainButton.addEventListener("click", function () {
      if (!hasCheckedAgain) {
        hasCheckedAgain = true;
        checkAgainButton.disabled = true;
        checkAgainButton.textContent = "Checking...";

        checkDNSStatus().then((isComplete) => {
          if (!isComplete) {
            // If still not complete, show error
            showError(
              "DNS service is still not ready. Please check the logs in your cluster."
            );
          }
          checkAgainButton.disabled = false;
          checkAgainButton.textContent = "Check Again";
        });
      }
    });
  }
});
