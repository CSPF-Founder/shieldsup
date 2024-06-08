from app.enums.main import SeverityIndex


def calculate_cvss_by_severity(severity: SeverityIndex.ENUM) -> float:
    if severity == SeverityIndex.ENUM.CRITICAL:
        return 9.0
    elif severity == SeverityIndex.ENUM.HIGH:
        return 7.0
    elif severity == SeverityIndex.ENUM.MEDIUM:
        return 4.0
    elif severity == SeverityIndex.ENUM.LOW:
        return 1.0
    else:
        return 0.0
