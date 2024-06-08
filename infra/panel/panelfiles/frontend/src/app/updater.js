import {
  showError,
  hideLoadingBox,
  redirectToLogin,
  requestWithCSRFToken,
  ready,
} from "./main.js";

let updateStatusInterval = null;
const STATUS_CHECK_INTERVAL = 15 * 1000; // 15 seconds

ready(function () {
  // call checkstatus every 5 seconds
  updateStatusInterval = setInterval(checkStatus, STATUS_CHECK_INTERVAL);
});

function checkStatus() {
  const updateButton = document.getElementById("update-button");
  const updateStatus = document.getElementById("update-status");
  const lastUpdate = document.getElementById("last-updated");

  requestWithCSRFToken("/update/status", {
    method: "GET",
  })
    .then((response) => {
      if (response.redirected) {
        redirectToLogin(response.redirect);
      }
      return response.json().then((data) => ({ ok: response.ok, data }));
    })
    .then(({ ok, data }) => {
      hideLoadingBox();

      if (!ok) {
        throw new Error(data.error || "Error occurred");
      }

      if (data.yet_to_update !== undefined && data.yet_to_update) {
        updateButton.disabled = false;
      } else if (data.status !== undefined) {
        if (data.status == 0) {
          // updating
          updateButton.disabled = true;
          updateStatus.innerHTML =
            "Update Status: <span class='ms-2 px-2 badge bg-warning'>Updating...</span>";
        } else if (data.status == 1) {
          // updated
          updateButton.disabled = false;
          updateStatus.innerHTML =
            "Update Status: <span class='ms-2 px-2 badge bg-success'>Updated</span>";
          lastUpdate.innerHTML =
            "Last Updated: <span class='ms-2'>" +
            data.last_updated_time +
            "</span>";
          clearInterval(updateStatusInterval);
        } else if (data.status == 2) {
          // update failed
          updateButton.disabled = false;
          updateStatus.innerHTML =
            "Update Status: <span class='ms-2 px-2 badge bg-danger'>Update Failed</span>";
          clearInterval(updateStatusInterval);
        }
      } else {
        updateButton.disabled = false;
        showError("Unable to check the update status");
        clearInterval(updateStatusInterval);
      }
    })
    .catch((error) => {
      updateButton.disabled = false;
      hideLoadingBox();
      showError(error.message);
    });
}

ready(function () {
  // update button clicked
  const updateButton = document.getElementById("update-button");
  const updateStatus = document.getElementById("update-status");
  updateButton.addEventListener("click", function (e) {
    e.preventDefault();
    updateButton.disabled = true;
    // $("#update-button").html("Updating...");
    requestWithCSRFToken("/update/start", {
      method: "POST",
    })
      .then((response) =>
        response.json().then((data) => ({ ok: response.ok, data }))
      )
      .then(({ ok, data }) => {
        hideLoadingBox();

        if (!ok) {
          throw new Error(data.error || "Error occurred");
        }
        if (data.success) {
          updateStatus.innerHTML =
            "Update Status: <span class='ms-2 px-2 badge bg-warning'>Updating...</span>";
          // $("#updated-date").html("Last Updated: <span class='ms-2'>--</span>");
          // call checkstatus every 5 seconds
          updateStatusInterval = setInterval(
            checkStatus,
            STATUS_CHECK_INTERVAL
          );
        } else {
          updateButton.disabled = false;
          updateButton.innerHTML = "Update";
          showError(data.error);
        }
      })
      .catch((error) => {
        hideLoadingBox();
        showError(error.message);
        updateButton.disabled = false;
      });
  });
});
