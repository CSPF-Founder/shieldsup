from bson.objectid import ObjectId

from app.db.crud.base import CrudBase
from app.db.models.target import Target


class CrudScanResult(CrudBase):
    COLLECTION_NAME = "scan_results"

    def __init__(self, db_session) -> None:
        super().__init__(db_session, self.COLLECTION_NAME)

    def add(self, records: list) -> int:
        if not records:
            return 0

        inserted_result = self.collection.insert_many(records)
        if not inserted_result:
            return 0

        return len(inserted_result.inserted_ids)

    def exists(self, filters: dict) -> bool:
        if not filters:
            return False

        result = self.collection.find_one(filters)
        if not result:
            return False

        return True

    def get_list_by_target(self, target: Target):
        result_cursor = self.collection.find({"target_id": target.id}).sort("severity")

        entries = []
        if not result_cursor:
            return
        for doc in result_cursor:
            entries.append(doc)

        return entries

    def get_list_by_target_id(self, target_id):
        if not target_id:
            raise Exception("Empty Mongodb object id given")

        if not ObjectId.is_valid(target_id):
            raise Exception("Invalid Mongodb object id given")

        if type(target_id) != ObjectId:
            target_id = ObjectId(target_id)

        result_cursor = self.collection.find({"target_id": target_id}).sort("severity")

        entries = []
        if not result_cursor:
            return
        for doc in result_cursor:
            entries.append(doc)

        return entries
