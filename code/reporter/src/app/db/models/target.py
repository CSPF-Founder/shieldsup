from datetime import datetime


from app.db.models.base_model import BaseDocumentModel
from app.enums.main import ScanStatus
from app.utils import ip_utils


class Target(BaseDocumentModel):
    customer_username: str
    target_address: str
    flag: int = 0
    scan_status: ScanStatus.ENUM
    created_at: datetime

    scanner_ip: str | None = None
    scanner_username: str | None = None
    scan_started_time: datetime | None = None
    scan_completed_time: datetime | None = None

    overall_cvss_score: float | None = None
    cvss_score_by_host: dict | None = None

    def get_scan_started_time(self):
        if not self.scan_started_time:
            return ""
        return self.scan_started_time.strftime("%d-%m-%Y %I:%M%p")

    def get_scan_completed_time(self):
        if not self.scan_completed_time:
            return ""
        return self.scan_completed_time.strftime("%d-%m-%Y %I:%M%p")

    def is_ip_range(self):
        if ip_utils.is_valid_cidr_by_regex(self.target_address):
            return True
