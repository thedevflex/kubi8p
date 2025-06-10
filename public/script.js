document.addEventListener("DOMContentLoaded", function () {
  const dbTypeRadios = document.querySelectorAll('input[name="db-type"]');
  const externalDbSection = document.getElementById("external-db-section");
  const connectionTypeSelect = document.getElementById("connection-type");
  const secretSection = document.getElementById("secret-section");
  const stringSection = document.getElementById("string-section");
  const nextButton = document.getElementById("next-button");
  const errorMessage = document.getElementById("error-message");

  // Show/hide external DB section based on radio selection
  dbTypeRadios.forEach((radio) => {
    radio.addEventListener("change", function () {
      externalDbSection.style.display =
        this.value === "external" ? "block" : "none";
      validateForm();
    });
  });

  // Toggle between secret and string sections
  connectionTypeSelect.addEventListener("change", function () {
    if (this.value === "secret") {
      secretSection.classList.add("active");
      stringSection.classList.remove("active");
    } else {
      secretSection.classList.remove("active");
      stringSection.classList.add("active");
    }
    validateForm();
  });

  // Validate secret name and path format
  const secretNameInput = document.getElementById("secret-name");
  const secretPathInput = document.getElementById("secret-path");

  secretNameInput.addEventListener("input", validateForm);
  secretPathInput.addEventListener("input", validateForm);

  // Validate connection string
  const connectionStringInput = document.getElementById("connection-string");
  connectionStringInput.addEventListener("input", validateForm);

  function validateForm() {
    const selectedDbType = document.querySelector(
      'input[name="db-type"]:checked'
    ).value;

    if (selectedDbType === "new") {
      nextButton.disabled = false;
      return;
    }

    const connectionType = connectionTypeSelect.value;
    if (connectionType === "secret") {
      const secretName = secretNameInput.value.trim();
      const secretPath = secretPathInput.value.trim();
      // Basic validation for secret name (k8s secret name format)
      const isValidName =
        /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/.test(
          secretName
        );
      // Basic validation for secret path format
      const isValidPath = /^\.data\.[a-zA-Z0-9_-]+$/.test(secretPath);
      nextButton.disabled = !(isValidName && isValidPath);
    } else {
      const connectionString = connectionStringInput.value.trim();
      // Basic validation for connection string
      const isValid = connectionString.length > 0;
      nextButton.disabled = !isValid;
    }
  }

  // Handle next button click
  nextButton.addEventListener("click", async function () {
    const selectedDbType = document.querySelector(
      'input[name="db-type"]:checked'
    ).value;
    let payload = {
      type: selectedDbType,
    };

    if (selectedDbType === "external") {
      const connectionType = connectionTypeSelect.value;
      if (connectionType === "secret") {
        payload.connectionType = "secret";
        payload.secretName = secretNameInput.value.trim();
        payload.secretPath = secretPathInput.value.trim();
      } else {
        payload.connectionType = "string";
        payload.connectionString = connectionStringInput.value.trim();
      }
    }

    // Disable button and show loading state
    nextButton.disabled = true;
    nextButton.classList.add("loading");
    errorMessage.style.display = "none";

    // // For new DB:
    // {
    //     "type": "new"
    // }

    // // For external DB with secret:
    // {
    //     "type": "external",
    //     "connectionType": "secret",
    //     "secretName": "my-db-secret",
    //     "secretPath": ".data.your_secret"
    // }

    // // For external DB with connection string:
    // {
    //     "type": "external",
    //     "connectionType": "string",
    //     "connectionString": "your_connection_string"
    // }

    try {
      const response = await fetch("/api/initiate-db", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        window.location.href = "/webhook/setup/index.html";
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

  // Initial validation
  validateForm();
});
