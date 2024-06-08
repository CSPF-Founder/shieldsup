"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""


from app import core_app
from app.db.connections.mongo_db import MongoDatabaseWrapper


class MainDatabase(MongoDatabaseWrapper):
    def __init__(self) -> None:
        super().__init__(
            db_uri=core_app.settings.main_db_uri, db_name=core_app.settings.main_db_name
        )
