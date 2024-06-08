"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""


class CrudBase:
    """
    Base crud class
    """

    def __init__(self, db_session, collection_name) -> None:
        self.db_session = db_session
        self.collection = db_session[collection_name]
