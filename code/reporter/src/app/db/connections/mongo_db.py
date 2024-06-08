"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from pymongo.database import Database as PyMongoDatabase
from pymongo import MongoClient

from app import core_app


class MongoDatabaseWrapper(object):
    def __init__(self, *, db_uri, db_name):
        self.db_name = db_name
        self.db_uri = db_uri

        self.client: MongoClient | None = None
        self.db: PyMongoDatabase | None = None

    @property
    def port(self):
        return self._port

    @port.setter
    def port(self, value):
        self._port = int(value)

    def connect(self):
        try:
            if not self.client:
                self.client = MongoClient(self.db_uri)
                self.db = self.client.get_database(self.db_name)
        except Exception as e:
            core_app.logger.exception(e)

    def __enter__(self):
        self.connect()
        return self.db

    def __exit__(self, exc_type, exc_val, exc_tb):
        # if exc_type is not None:
        #     print(exc_type, exc_val, exc_tb)
        # self.close()

        # don't close the connection for mongodb
        # - no need to close as pymongo does pooling
        pass

    def close(self):
        if isinstance(self.client, MongoClient):
            self.client.close()
