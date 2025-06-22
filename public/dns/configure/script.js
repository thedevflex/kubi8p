document.addEventListener("DOMContentLoaded", function () {
  const statusContainer = document.getElementById("status-container");
  const statusMessage = document.getElementById("status-message");
  const loadingSpinner = document.getElementById("loading-spinner");
  const errorMessage = document.getElementById("error-message");
  const urlInputContainer = document.getElementById("url-input-container");
  const cnameInstruction = document.getElementById("cname-instruction");
  const configureButton = document.getElementById("configure-button");
  const verifyButton = document.getElementById("verify-button");

  const subdomainInput = document.getElementById("subdomain");
  const domainInput = document.getElementById("domain");
  const cnameSubdomain = document.getElementById("cname-subdomain");
  const cnameValue = document.getElementById("cname-value");

  let dnsExternalIP = "";
  let generatedPrefix = "";

  function showLoading(message = "Loading DNS information...") {
    statusContainer.style.display = "block";
    loadingSpinner.style.display = "block";
    statusMessage.style.display = "none";
    errorMessage.style.display = "none";
    urlInputContainer.style.display = "none";
    cnameInstruction.style.display = "none";
    configureButton.style.display = "none";
    verifyButton.style.display = "none";

    const loadingText = loadingSpinner.querySelector("p");
    if (loadingText) {
      loadingText.textContent = message;
    }
  }

  function showError(message) {
    statusContainer.style.display = "block";
    loadingSpinner.style.display = "none";
    statusMessage.style.display = "block";
    statusMessage.textContent = message;
    statusMessage.style.color = "#ff6b6b";
    errorMessage.style.display = "none";
    urlInputContainer.style.display = "none";
    cnameInstruction.style.display = "none";
    configureButton.style.display = "none";
    verifyButton.style.display = "none";
  }

  function showSuccess(message) {
    statusContainer.style.display = "block";
    loadingSpinner.style.display = "none";
    statusMessage.style.display = "block";
    statusMessage.textContent = message;
    statusMessage.style.color = "#0f0";
    errorMessage.style.display = "none";
  }

  function showURLInput() {
    statusContainer.style.display = "none";
    errorMessage.style.display = "none";
    urlInputContainer.style.display = "block";
    configureButton.style.display = "block";
  }

  function showCNAMEInstructions() {
    cnameInstruction.style.display = "block";
    verifyButton.style.display = "block";
    configureButton.style.display = "none";
  }

  function updateCNAMEValues() {
    const subdomain = subdomainInput.value.trim();
    const domain = domainInput.value.trim();

    if (subdomain && domain) {
      cnameSubdomain.textContent = subdomain;
      cnameValue.textContent = dnsExternalIP;
    }
  }

  async function loadDNSStatus() {
    showLoading("Loading DNS information...");

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
          dnsExternalIP = data.external_ip;
          // Generate a random prefix (you might want to get this from the backend)
          generatedPrefix = generateRandomPrefix();

          showURLInput();
          return true;
        } else {
          showError(
            "DNS service is not ready. Please go back to the setup page."
          );
          return false;
        }
      } else {
        showError(
          "Failed to load DNS status. Please check the logs in your cluster."
        );
        return false;
      }
    } catch (error) {
      showError(
        "Network error while loading DNS status. Please check the logs in your cluster."
      );
      return false;
    }
  }

  function generateRandomPrefix() {
    // Generate a random 8-character alphanumeric string
    const chars = "abcdefghijklmnopqrstuvwxyz0123456789";
    let result = "";
    for (let i = 0; i < 8; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  }

  async function configureDNS() {
    const subdomain = subdomainInput.value.trim();
    const domain = domainInput.value.trim();

    if (!subdomain || !domain) {
      errorMessage.textContent =
        "Please fill in both subdomain and domain fields.";
      errorMessage.style.display = "block";
      return;
    }

    // Validate domain format
    const domainRegex =
      /^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]?\.[a-zA-Z]{2,}$/;
    if (!domainRegex.test(domain)) {
      errorMessage.textContent =
        "Please enter a valid domain (e.g., yourdomain.com)";
      errorMessage.style.display = "block";
      return;
    }

    configureButton.disabled = true;
    configureButton.textContent = "Configuring...";
    errorMessage.style.display = "none";

    try {
      const response = await fetch("/api/configure-dns", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          prefix: generatedPrefix,
          domain: domain,
        }),
      });

      if (response.ok) {
        showCNAMEInstructions();
        updateCNAMEValues();
      } else {
        throw new Error("Configuration failed");
      }
    } catch (error) {
      errorMessage.textContent = "Failed to configure DNS. Please try again.";
      errorMessage.style.display = "block";
    } finally {
      configureButton.disabled = false;
      configureButton.textContent = "Configure DNS";
    }
  }

  async function verifyIntegration() {
    verifyButton.disabled = true;
    verifyButton.textContent = "Verifying...";
    errorMessage.style.display = "none";

    try {
      const response = await fetch("/api/verify-dns", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.status === 200) {
        showSuccess("DNS integration verified successfully! Redirecting...");
        setTimeout(() => {
          window.location.href = "/complete";
        }, 2000);
      } else {
        throw new Error("Verification failed");
      }
    } catch (error) {
      errorMessage.textContent =
        "DNS verification failed. Please check your DNS configuration and try again.";
      errorMessage.style.display = "block";
    } finally {
      verifyButton.disabled = false;
      verifyButton.textContent = "Verify Integration";
    }
  }

  // Event listeners
  if (configureButton) {
    configureButton.addEventListener("click", configureDNS);
  }

  if (verifyButton) {
    verifyButton.addEventListener("click", verifyIntegration);
  }

  // Update CNAME values when inputs change
  if (subdomainInput) {
    subdomainInput.addEventListener("input", updateCNAMEValues);
  }

  if (domainInput) {
    domainInput.addEventListener("input", updateCNAMEValues);
  }

  // Load DNS status on page load
  loadDNSStatus();
});
