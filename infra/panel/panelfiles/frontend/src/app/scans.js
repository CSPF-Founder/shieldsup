import {
  redirectToLogin,
  showError,
  showSuccess,
  resetInputForm,
  loadingBox,
  hideLoadingBox,
  requestWithCSRFToken,
  ready,
} from "./main.js";

import "bootstrap";
import "datatables.net";
import "datatables.net-bs4";
import "datatables.net-responsive";
import "datatables.net-responsive-bs4";

// Only Jquery Dependency
$(document).ready(function () {
  $(".table").DataTable({
    responsive: {
      details: {
        responsive: true,
        type: "none",
        target: "",
      },
    },
    order: [[0, "desc"]],
    language: {
      search: "",
    },
  });

  $(".dataTables_filter input").attr("placeholder", "Search...");

  load_target_table(false);

  scansTable.on("draw.dt", function (e, settings) {
    generateTargetsToUpdate();
    if (any_scan_running) {
      status_check_interval = STATUS_CHECK_INTERVAL_FOR_RUNNING_SCANS;
    } else {
      status_check_interval = STATUS_CHECK_INTERVAL_FOR_NO_SCANS;
    }

    setTimeout(checkStatus, status_check_interval);
  });

  function destroyDataTableRow(tableReference, row) {
    var table = $(tableReference, row.child());

    table.detach();
    table.DataTable().destroy();

    // And then hide the row
    row.child.hide();
    row.child.remove();
    row.remove().draw(false); // set draw to false to prevent the table from losing pagination
  }

  // $('.dataTables_filter input').attr("placeholder", "Search...");

  $("body").on("click", ".delete-target", function (e) {
    // $(".delete-target").on("click", function(e) {
    e.preventDefault();
    if (!confirm("Are you sure want to delete the scan")) {
      return;
    }
    // let entryID = $(this).data("id");
    //get datatable row id
    let entryID = $(this).closest("tr").attr("id");

    $.ajax({
      type: "DELETE",
      url: "/targets/" + entryID,
      headers: {
        "X-CSRF-Token": CSRF_TOKEN,
        "X-Requested-With": "XMLHttpRequest",
      },
      disableLoading: true,
      success: function (res) {
        if (res.redirect) {
          redirectToLogin(res.redirect);
        }
        if (res.error) {
          showError(res.error);
        } else if (res.success) {
          showSuccess(res.success);
          let row = scansTable.row("#" + entryID);
          destroyDataTableRow(scansTable, row);
        }
      },
      error: function (jqXHR, textStatus, errorThrown) {
        if (jqXHR.status === 422) {
          var response = jqXHR.responseJSON;
          if (response && response.error) {
            showError(response.error);
          } else {
            showError("Unprocessable Entity: Invalid request parameters.");
          }
        } else {
          showError("An unexpected error occurred: " + textStatus);
        }
      },
    });
  });
});

// add Scan
ready(function () {
  const addForm = document.getElementById("add-scan-form");
  if (!addForm) {
    return;
  }

  const addButton = document.getElementById("add-scan-btn");

  // Add Scan Form Submit
  addForm.addEventListener("submit", function (event) {
    addButton.disabled = true;

    event.preventDefault();
    loadingBox();

    const formData = new FormData(addForm);

    requestWithCSRFToken("/targets/add", {
      method: "POST",
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

        if (data.success) {
          resetInputForm("#add-scan-form");
          showSuccess(data.success);
        } else if (data.redirect) {
          redirectToLogin(data.redirect);
        }

        addButton.disabled = false;
      })
      .catch((error) => {
        console.log(error);
        hideLoadingBox();
        showError(error.message);
        addButton.disabled = false;
      });
  });
});

//View Scans Page function
const STATUS_CHECK_INTERVAL_FOR_NO_SCANS = 10 * 1000; // 10 seconds
const STATUS_CHECK_INTERVAL_FOR_RUNNING_SCANS = 60 * 1000; // 60 seconds
let status_check_interval = STATUS_CHECK_INTERVAL_FOR_RUNNING_SCANS;
let any_scan_running = false;

let scansTable = undefined;

let targetsToUpdate = [];
let statusCheckIntervalId = undefined;

function generateTargetsToUpdate() {
  // one time should be called after the datatable is loaded
  let any_scan_started = false;
  // console.log("generateTargetsToUpdate");
  targetsToUpdate = [];
  scansTable.rows().every(function () {
    let row = this.data();
    const scanStatus = row.scan_status;
    if (scanStatus === TARGET_STATUS.YET_TO_START) {
      targetsToUpdate.push({
        id: row.id,
        scan_status: scanStatus,
      });
    } else if (scanStatus === TARGET_STATUS.SCAN_STARTED) {
      targetsToUpdate.push({
        id: row.id,
        scan_status: scanStatus,
      });
      any_scan_started = true;
    }
  });

  if (any_scan_started) {
    status_check_interval = STATUS_CHECK_INTERVAL_FOR_RUNNING_SCANS; // 60 seconds
    any_scan_running = true;
  } else {
    any_scan_running = false;
    status_check_interval = STATUS_CHECK_INTERVAL_FOR_NO_SCANS; // 10 seconds
  }
}

function load_target_table(redraw) {
  var ajax_url = "/targets/list";
  var alreadyExecuted = false;

  let ajaxData = {
    ...baseAjaxData,
  };
  scansTable = $("#scan-list").DataTable({
    // "order": [
    //     [0, "desc"]
    // ],
    aaSorting: [],
    searching: false,
    destroy: true,
    serverSide: true,
    paging: true,
    initComplete: function () {
      $(".dt-button").removeClass("dt-button");
      $(".dataTables_length").addClass("control-label pt-3");
    },
    processing: false,
    // "language": {
    //     processing: "<span class='fa-stack fa-lg'  ><i class='fa fa-spinner fa-spin fa-stack-2x fa-fw'></i></span>&emsp;Processing ..."
    // },
    responsive: {
      details: {
        responsive: true,
        display: $.fn.dataTable.Responsive.display.childRowImmediate,
        type: "none",
        target: "",
      },
    },
    serverSide: true,
    stateSave: false,
    lengthMenu: [
      [10, 50, 100],
      [10, 50, 100],
    ],
    ajax: {
      url: ajax_url,
      dataType: "json",
      type: "POST",
      data: ajaxData,
      dataSrc: "records",
      error: function (xhr, error, code) {
        if (xhr.status == 303) {
          redirectToLogin();
        } else {
          showError(
            "Unable to fetch data. Please try again later. If the problem persists, please contact support."
          );
        }
      },
    },
    columns: [
      {
        data: "target_address",
      },
      {
        data: "scan_status_text",
      },
      {
        data: "scan_started_time",
      },
      {
        data: "scan_completed_time",
      },
      {
        data: "action",
      },
    ],
    rowId: "id",
  });
}

function sendStatusCheckRequest() {
  let targetIds = targetsToUpdate.map(function (t, index, arr) {
    return t.id;
  });
  $.ajax({
    type: "POST",
    url: "check-multi-status",
    disableLoading: true,
    data: {
      target_ids: targetIds,
    },
    success: function (res) {
      if (res.data) {
        for (let index in res.data) {
          let entry = res.data[index];
          let targetId = entry.id;
          var row = scansTable.row("#" + targetId);
          var data = row.data();

          if (data === undefined) {
            continue;
          }

          let actions = "";

          if (entry.scan_status === undefined) {
            continue;
          }

          // update scan_status in the target list
          targetsToUpdate = targetsToUpdate.map(function (t, index, arr) {
            if (t.id == targetId) {
              t.scan_status = entry.scan_status;
            }
            return t;
          });

          if (entry.scan_status === TARGET_STATUS.SCAN_STARTED) {
            // update the datatable row find by target id
            data.scan_started_time = entry.scan_started_time;
            data.scan_status_text =
              '<span class="spinner-border spinner-border-sm text-primary" aria-hidden="true"></span> <span role="status">Scanning...</span>';

            any_scan_running = true;
          } else {
            // update the datatable row find by target id
            data.scan_status_text = entry.scan_status_text;
          }

          if (entry.scan_status == TARGET_STATUS.REPORT_GENERATED) {
            // update the datatable row find by target id
            data.scan_completed_time = entry.scan_completed_time;
            //enable actions priamry-info
            actions =
              '<a href="/targets/' +
              targetId +
              '/report" class="btn btn-sm btn-primary m-1 report-button">Report</a>';
            actions +=
              '<a href="/targets/' +
              targetId +
              '/scan-results" class="btn btn-sm btn-primary m-1 alerts-button">Alerts</a>';
          } else {
            // disable report and alerts button
            actions =
              '<a href="/targets/' +
              targetId +
              '/report" class="btn btn-sm btn-dark m-1 report-button disabled" disabled>Report</a>';
            actions +=
              '<a href="/targets/' +
              targetId +
              '/scan-results" class="btn btn-sm btn-dark m-1 alerts-button disabled" disabled>Alerts</a>';

            actions =
              '<a href="/targets/' +
              targetId +
              '/report" class="btn btn-sm btn-dark m-1 report-button disabled" disabled>Report</a>';
            actions +=
              '<a href="/targets/' +
              targetId +
              '/scan-results" class="btn btn-sm btn-dark m-1 alerts-button disabled" disabled>Alerts</a>';
          }

          if (
            entry.scan_status === TARGET_STATUS.SCAN_FAILED ||
            entry.scan_status === TARGET_STATUS.REPORT_GENERATED
          ) {
            // Remove the target from the list if the scan is completed or failed
            targetsToUpdate = targetsToUpdate.filter(function (t, index, arr) {
              return t.id != targetId;
            });
          }

          if (entry.scan_status === TARGET_STATUS.SCAN_STARTED) {
            // if it is scanning then disable the delete button
            actions +=
              '<a href="#" class="btn btn-sm btn-dark text-white m-1 delete-target disabled" disabled>Delete</a>';
          } else {
            actions +=
              '<a href="#" class="btn btn-sm btn-danger text-white m-1 delete-target">Delete</a>';
          }

          if (actions !== "") {
            data.action = actions;
          }

          row.scan_status = entry.scan_status;

          row.data(data);
        }
      }
    },
  }).always(function () {
    recheckStatus();
  });
}

function checkStatus() {
  console.log("initiating Checking status");

  if (targetsToUpdate.length === 0) {
    console.log("Finished Checking status at " + new Date().toLocaleString());
    return;
  }

  console.log("Checking status at " + new Date().toLocaleString());
  sendStatusCheckRequest();
}

function recheckStatus() {
  if (targetsToUpdate.length === 0) {
    console.log("Finished Checking status at " + new Date().toLocaleString());
    return;
  }

  if (
    status_check_interval === STATUS_CHECK_INTERVAL_FOR_NO_SCANS &&
    any_scan_running
  ) {
    console.log("Resetting the status check interval to 60 seconds");
    // if any scan is running then change the status check interval to 60 seconds
    status_check_interval = STATUS_CHECK_INTERVAL_FOR_RUNNING_SCANS;
  } else {
    any_scan_running = targetsToUpdate.some(function (t, index, arr) {
      return t.scan_status == TARGET_STATUS.SCAN_STARTED;
    });
    if (
      !any_scan_running &&
      status_check_interval === STATUS_CHECK_INTERVAL_FOR_RUNNING_SCANS
    ) {
      // if no scan is running then change the status check interval to 10 seconds
      console.log("Resetting the status check interval to 10 seconds");

      status_check_interval = STATUS_CHECK_INTERVAL_FOR_NO_SCANS;
    }
  }

  setTimeout(checkStatus, status_check_interval);
}

// // Delete Target
// ready(function () {
//   const scanList = document.getElementById("scan-list");

//   if (!scanList) {
//     return;
//   }

//   scanList.addEventListener("click", function (e) {
//     let targetElement = e.target;

//     // Traverse up the DOM tree if the clicked element is not the button itself
//     while (
//       targetElement != null &&
//       !targetElement.classList.contains("delete-target")
//     ) {
//       targetElement = targetElement.parentElement;
//     }

//     // If no element with the 'delete-scan' class was found in the ancestry, exit the function
//     if (targetElement == null) return;

//     e.preventDefault();
//     if (!confirm("Are you sure want to delete the scan ?")) {
//       return;
//     }

//     let row = e.target.closest("tr");
//     let targetId = row.getAttribute("data-id"); // Adjust if data-id is stored differently
//     let url = "/targets/" + targetId;

//     loadingBox();
//     requestWithCSRFToken(url, {
//       method: "DELETE",
//     })
//       .then((response) =>
//         response.json().then((data) => ({ ok: response.ok, data }))
//       )
//       .then(({ ok, data }) => {
//         hideLoadingBox();

//         if (!ok) {
//           throw new Error(data.error || "Error occurred");
//         }

//         if (data.redirect) {
//           redirectToLogin(data.redirect);
//         } else if (data.success) {
//           showSuccess(data.success);
//           row.remove();
//         }
//       })
//       .catch((error) => {
//         hideLoadingBox();
//         showError(error.message);
//       });
//   });
// });
