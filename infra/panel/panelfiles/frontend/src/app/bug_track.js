import {
  showError,
  showSuccess,
  redirectToLogin,
  requestWithCSRFToken,
  hideLoadingBox,
  resetInputForm,
  loadingBox,
  virtualFormSubmit,
  ready,
} from "./main.js";

import "bootstrap";
import "datatables.net";
import "datatables.net-bs4";
import "datatables.net-responsive";
import "datatables.net-responsive-bs4";

import datepicker from "js-datepicker";

ready(function () {
  const bugtrackTable = $("#bugtrack-table").DataTable({
    aaSorting: [],
    responsive: {
      details: {
        responsive: true,
        display: $.fn.dataTable.Responsive.display.childRowImmediate,
        type: "none",
        target: "",
      },
    },
    language: {
      search: "",
    },
  });

  function destroyDataTableRow(tableReference, row) {
    var table = $(tableReference, row.child());
    console.log(row.child());
    table.detach();
    table.DataTable().destroy();

    // And then hide the row
    row.child.hide();
    row.child.remove();
    row.remove().draw();
  }

  $(".dataTables_filter input").attr("placeholder", "Search...");

  $("#bugtrack-table").on("click", ".delete-entry", function (e) {
    e.preventDefault();
    if (!confirm("Are you sure want to delete the entry")) {
      return;
    }
    let entryID = $(this).data("id");

    loadingBox();

    requestWithCSRFToken("/bug-track/" + entryID, {
      method: "DELETE",
    })
      .then((response) =>
        response.json().then((data) => ({ ok: response.ok, data }))
      )
      .then(({ ok, data }) => {
        hideLoadingBox();

        if (!ok) {
          throw new Error(data.error || "Error occurred");
        }

        if (data.redirect) {
          redirectToLogin(data.redirect);
        }
        if (data.error) {
          showError(data.error);
        } else if (data.success) {
          showSuccess(data.success);

          destroyDataTableRow(
            bugtrackTable,
            bugtrackTable.row("#bugtrack_row_" + entryID)
          );
        } else {
          showError("Unexpected Error !");
        }
      })
      .catch(function (error) {
        hideLoadingBox();
        showError(error.message);
      });
  });
});

//Update Bugtrack Form
ready(function () {
  const updateForm = document.getElementById("update-bugtrack-form");
  if (!updateForm) {
    return;
  }

  const updateBugButton = document.getElementById("update-bugtrack");
  updateBugButton.addEventListener("click", function (e) {
    const formData = new FormData(updateForm);

    let bugID = document.getElementById("bug-id").value;

    e.preventDefault();
    requestWithCSRFToken("/bug-track/" + bugID, {
      method: "PATCH",
      body: formData,
    })
      .then((response) =>
        response.json().then((data) => ({ ok: response.ok, data }))
      )
      .then(({ ok, data }) => {
        hideLoadingBox();

        if (!ok) {
          throw new Error(data.error || "Error occurred");
        }
        if (data.redirect) {
          redirectToLogin(data.redirect);
        }
        if (data.error) {
          showError(data.error);
        } else if (data.success) {
          showSuccess(data.success);
        } else {
          showError("Unexpected Error !");
        }
      })
      .catch(function (error) {
        hideLoadingBox();
        showError(error.message);
      });
  });
});

//Export Bugtrack
ready(function () {
  const exportButton = document.getElementById("export-button");
  if (!exportButton) {
    return;
  }
  exportButton.addEventListener("click", function (e) {
    virtualFormSubmit("/bug-track/export-bugs", {
      [`${CSRF_NAME}`]: CSRF_TOKEN,
    });
  });
});

//Add Bug Track
ready(function () {
  function formatDate(date) {
    let d = date.toLocaleDateString("en-us", { day: "2-digit" });
    let month = date.toLocaleDateString("en-us", { month: "2-digit" });
    let year = date.toLocaleDateString("en-us", { year: "numeric" });
    return `${year}-${month}-${d}`;
  }

  var today_date = new Date(document.getElementById("today_date").value);

  datepicker("#found_date", {
    maxDate: today_date,
    formatter: (input, date, instance) => {
      input.value = formatDate(date);
    },
  });
  datepicker("#revalidated_date", {
    maxDate: today_date,
    formatter: (input, date, instance) => {
      input.value = formatDate(date);
    },
  });

  const addBugTrack = document.getElementById("add-bugtrack-form");
  if (!addBugTrack) {
    return;
  }

  addBugTrack.addEventListener("submit", function (e) {
    const formData = new FormData(addBugTrack);

    e.preventDefault();
    loadingBox();

    requestWithCSRFToken("/bug-track/add-from-scanresult", {
      method: "POST",
      body: formData,
    })
      .then((response) =>
        response.json().then((data) => ({ ok: response.ok, data }))
      )
      .then(({ ok, data }) => {
        if (!ok) {
          throw new Error(data.error || "Error occurred");
        }

        if (data.redirect) {
          redirectToLogin(data.redirect);
        }

        if (data.error) {
          showError(data.error);
        } else if (data.success) {
          showSuccess(data.success);
          resetInputForm("#add-bugtrack-form");
        }
      })
      .catch(function (error) {
        hideLoadingBox();
        showError(error.message);
      });
  });
});
