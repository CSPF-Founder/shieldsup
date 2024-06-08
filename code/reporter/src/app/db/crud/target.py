from datetime import datetime
from bson.objectid import ObjectId

from app.db.models.target import Target
from app.db.crud.base import CrudBase
from app.enums.main import ScanStatus


class CrudTarget(CrudBase):
    COLLECTION_NAME = "targets"

    def __init__(self, db_session) -> None:
        super().__init__(db_session, self.COLLECTION_NAME)

    def find_by_id(self, target_id: ObjectId) -> Target | None:
        """
        Find Target by id

        Args:
            target_id (ObjectId): _description_

        Raises:
            Exception: _description_

        Returns:
            Target | None: _description_
        """
        if not target_id or not isinstance(target_id, ObjectId):
            raise Exception("Invalid Mongodb object id given")

        row = self.collection.find_one({"_id": target_id})
        if not row:
            return None

        # Change _id to id
        row.update({"id": row["_id"]})
        return Target.model_validate(row)

    def update_scan_status(self, target: Target):
        """
        Update Scan Status by target

        Args:
            target (Target): _description_

        Raises:
            Exception: _description_
        """
        if not target or not isinstance(target, Target):
            raise Exception("Invalid Target object given")

        now = datetime.now()
        to_update = {}
        if target.scan_status == ScanStatus.ENUM.SCAN_STARTED:
            target.scan_started_time = now
            to_update["scan_started_time"] = now
        elif target.scan_status == ScanStatus.ENUM.REPORT_GENERATED:
            target.scan_completed_time = now
            to_update["scan_completed_time"] = now

        to_update["scan_status"] = target.scan_status
        self.collection.update_one({"_id": target.id}, {"$set": to_update})

    def update_scan_status_by_id(self, target_id: ObjectId, scan_status):
        if scan_status is None:
            raise Exception("Invalid scan_status")

        if not target_id or not isinstance(target_id, ObjectId):
            raise Exception("Invalid Mongodb object id given")

        now = datetime.now()

        to_update = {}
        if scan_status == ScanStatus.ENUM.SCAN_STARTED:
            to_update["scan_started_time"] = now
        elif scan_status == ScanStatus.ENUM.REPORT_GENERATED:
            to_update["scan_completed_time"] = now

        to_update["scan_status"] = scan_status

        self.collection.update_one({"_id": target_id}, {"$set": to_update})

    def mark_as_completed(self, target: Target):
        """
        Update Scan Status by target

        Args:
            target (Target): _description_

        Raises:
            Exception: _description_
        """
        if not target or not isinstance(target, Target):
            raise Exception("Invalid Target object given")

        now = datetime.now()
        target.scan_status = ScanStatus.ENUM.REPORT_GENERATED
        target.scan_completed_time = now

        to_update = {
            "scan_completed_time": now,
            "scan_status": target.scan_status,
            "overall_cvss_score": target.overall_cvss_score,
        }

        if target.is_ip_range():
            to_update["cvss_score_by_host"] = target.cvss_score_by_host

        to_update["scan_status"] = target.scan_status
        self.collection.update_one({"_id": target.id}, {"$set": to_update})

    def get_all(self) -> list[Target]:
        """
        Get all targets

        Returns:
            list[Target]: _description_
        """
        rows = self.collection.find()
        targets = []
        for row in rows:
            row.update({"id": row["_id"]})
            targets.append(Target.model_validate(row))
        return targets
