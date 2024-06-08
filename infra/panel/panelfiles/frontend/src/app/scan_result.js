import { ready } from "./main.js";

import "bootstrap";
import "datatables.net";
import "datatables.net-bs4";
import "datatables.net-responsive";
import "datatables.net-responsive-bs4";
import Chart from "chart.js/auto";

ready(function () {
  $(".table").DataTable({
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

  $(".dataTables_filter input").attr("placeholder", "Search...");

  $(".add-to-bugtrack").on("submit", function (e) {
    e.preventDefault();
    $("input[name='csrf_token']").remove();
    // Submit the form
    const form = $(this);
    form.unbind("submit").submit();
  });
});

let no_vulnerabilities =
  parseInt(document.getElementById("no_vulnerabilities").value) ?? 0;

const CHART_COLORS = {
  RED_COLOR: "rgb(255, 99, 132)", // Red
  BLUE_COLOR: "rgb(54, 162, 235)", // Blue
  YELLOW_COLOR: "rgb(255, 205, 86)", // Yellow
  GREEN_COLOR: "rgb(75, 192, 192)", // Green
  LIGHT_PINK_COLOR: "#D36086", // Light Pink
  LIGHT_GREEN_COLOR: "#54B399", // Light Green
  LIGHT_PURPLE_COLOR: "#9170B8", // Light Purple
  LIGHT_BLUE_COLOR: "#6092C0", // Light Blue
  LIGHT_GREY_COLOR: "#efefef", // Light Grey
  LIGHT_ORANGE_COLOR: "#FFB347", // Light Orange
  ORANGE_COLOR: "#FF9A5B", // Orange (Adjusted)
};

let cvssScore = document.getElementById("overall_cvss_score").value ?? 0.0;

ready(function () {
  /** !-- alerts distro stats chart starts --! */

  let needleValue = (cvssScore / 10) * 100;
  new Chart("cvss-score-gauge", {
    type: "doughnut",
    plugins: [
      {
        afterDraw: (chart) => {
          var dataTotal = chart.config.data.datasets[0].data.reduce(
            (a, b) => a + b,
            0
          );
          var angle = Math.PI + (1 / dataTotal) * needleValue * Math.PI;
          var ctx = chart.ctx;
          var cw = chart.canvas.offsetWidth;
          var ch = chart.canvas.offsetHeight;
          var cx = cw / 2;
          var cy = ch - 6;

          ctx.translate(cx, cy);
          ctx.rotate(angle);
          ctx.beginPath();
          ctx.moveTo(0, -3);
          ctx.lineTo(ch - 20, 0);
          ctx.lineTo(0, 3);
          ctx.fillStyle = "rgb(0, 0, 0)";
          ctx.fill();
          ctx.rotate(-angle);
          ctx.translate(-cx, -cy);
          ctx.beginPath();
          ctx.arc(cx, cy, 5, 0, Math.PI * 2);
          ctx.fill();
        },
      },
    ],
    data: {
      labels: [],
      datasets: [
        {
          data: [39, 30, 20, 11],
          needleValue: 27,
          backgroundColor: [
            CHART_COLORS.GREEN_COLOR,
            CHART_COLORS.YELLOW_COLOR,
            CHART_COLORS.ORANGE_COLOR,
            CHART_COLORS.RED_COLOR,
          ],
        },
      ],
    },
    options: {
      responsive: true,
      aspectRatio: 2,
      layout: {
        padding: {
          bottom: 3,
        },
      },
      rotation: -90,
      cutout: "50%",
      circumference: 180,
      legend: {
        display: false,
      },
      animation: {
        animateRotate: false,
        animateScale: true,
      },
      plugins: {
        tooltip: {
          enabled: false,
        },
      },
    },
  });

  let alertsDistroChartId = "alerts-distro-chart";
  var alertsDistroCtx = document
    .getElementById(alertsDistroChartId)
    .getContext("2d");

  // options
  var alertsDistroOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      labels: {
        render: "percentage",
        fontColor: "#fff",
      },
      layout: {
        padding: {
          left: 0,
          right: 0,
          top: 0,
          bottom: 0,
        },
      },
      legend: {
        fullWidth: false,
        position: "right",
        display: true,
        positionDoughnutLabelExtended: "bottomLeft",
        labels: {
          fontSize: 10,
          usePointStyle: true,
        },
      },
      positionPDExtended: "absolute",
    },
    cutout: "80%",
  };

  const spaceChar = " ";

  let alert_labels = ["Critical", "High", "Medium", "Low", "Info"];

  if (no_vulnerabilities) {
    alert_labels.push("No alerts");
  }

  if (no_vulnerabilities) {
    alert_chart_data.push(1);
  }

  let alertChartBgColor = [
    CHART_COLORS.RED_COLOR,
    CHART_COLORS.ORANGE_COLOR,
    CHART_COLORS.YELLOW_COLOR,
    CHART_COLORS.LIGHT_GREEN_COLOR,
    CHART_COLORS.LIGHT_BLUE_COLOR,
  ];

  if (no_vulnerabilities) {
    alertChartBgColor.push(CHART_COLORS.LIGHT_GREY_COLOR);
  }

  let alertsDistrodata = {
    labels: alert_labels,
    datasets: [
      {
        label: "Vulnerabilities Count",
        data: alert_chart_data,
        backgroundColor: alertChartBgColor,
        hoverOffset: 4,
        borderWidth: [1.8, 1.8, 1.8],
      },
    ],
  };
  var alertsDistroChart = new Chart(alertsDistroCtx, {
    type: "doughnut",
    data: alertsDistrodata,
    options: alertsDistroOptions,
  });
});
