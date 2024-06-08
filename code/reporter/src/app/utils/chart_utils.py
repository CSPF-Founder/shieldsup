"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from collections import OrderedDict
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.patches import Circle
import matplotlib as mpl

mpl.use("Agg")


def create_chart(input_data: dict, file_path: str, chart_type="pie"):
    colors = []
    severity_colors = {
        "Critical": "#ff2654",
        "High": "#ff6a38",
        "Medium": "#e6e943",
        "Low": "#36bf58",
        # "Info": "#0773b8",
    }

    new_entries = {}
    for severity, entry in input_data.items():
        if entry is None:
            new_entries.update({severity: entry})
            continue

        severity_color = severity_colors.get(severity)
        if severity_color and (entry != 0 or chart_type == "bar"):
            colors.append(severity_color)
            new_entries.update({severity: entry})

    label = new_entries.keys()
    values = new_entries.values()

    if chart_type == "pie":
        plt.suptitle("Distribution of Vulnerabilities", fontsize=14, fontweight="bold")
        plt.pie(
            values,  # type:ignore
            labels=label,  # type:ignore
            colors=colors,
            autopct="%1.1f%%",
            shadow=True,
            startangle=90,
        )
        plt.axis("equal")
    elif chart_type == "donut":
        plt.suptitle("Distribution of Vulnerabilities", fontsize=14, fontweight="bold")
        plt.pie(
            values,  # type:ignore
            labels=label,  # type:ignore
            colors=colors,
            autopct="%1.1f%%",
            shadow=True,
            startangle=90,
        )
        centre_circle = Circle((0, 0), 0.80, fc="white")
        fig = plt.gcf()
        fig.gca().add_artist(centre_circle)
        plt.axis("equal")
    elif chart_type == "bar":
        plt.suptitle("Number of Vulnerabilities", fontsize=14, fontweight="bold")
        N = len(label)
        ind = np.arange(N)
        width = 0.6
        vals = values
        plt.ylabel("Count of Vulnerabilities")
        plt.xlabel("Severity")
        plt.bar(ind, vals, width, color=colors, align="center")  # type:ignore
        plt.xticks(ind, label)  # type:ignore
    plt.savefig(file_path)
    plt.clf()


if __name__ == "__main__":
    # entries = json.OrderedDict([('High', 10), ('Medium', 4),('Low', 14)])
    entries = OrderedDict([("Critical", 1), ("High", 0), ("Medium", 0), ("Low", 14)])
    create_chart(entries, "/tmp/bar.png", chart_type="bar")
    create_chart(entries, "/tmp/pie.png")
