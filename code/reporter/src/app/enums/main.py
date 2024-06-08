from app.enums import CustomEnumClass
from enum import IntEnum


class ScanStatus:
    class ENUM(IntEnum):
        YET_TO_START = 0
        INITIATING_SCAN = 1
        SCAN_STARTED = 2
        RETRIEVED = 3
        REPORT_GENERATED = 4
        SCAN_FAILED = 999


class ApiScanStatus:
    class ENUM(IntEnum):
        NOT_COMPLETED = 1
        COMPLETED = 3

        DOES_NOT_EXIST = 501


class SeverityIndex(CustomEnumClass):
    class ENUM(IntEnum):
        CRITICAL = 1
        HIGH = 2
        MEDIUM = 3
        LOW = 4
        INFO = 5

    _meta_dictionary = {
        ENUM.CRITICAL: "Critical",
        ENUM.HIGH: "High",
        ENUM.MEDIUM: "Medium",
        ENUM.LOW: "Low",
        ENUM.INFO: "Info",
    }


def str_to_severity_index(severity: str) -> SeverityIndex.ENUM:
    severity = severity.lower()
    if severity == "critical":
        return SeverityIndex.ENUM.CRITICAL
    elif severity == "high":
        return SeverityIndex.ENUM.HIGH
    elif severity == "medium":
        return SeverityIndex.ENUM.MEDIUM
    elif severity == "low":
        return SeverityIndex.ENUM.LOW
    return SeverityIndex.ENUM.INFO
