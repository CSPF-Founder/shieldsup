"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from configparser import ConfigParser
import os
import time

from app.core.config import AppConfig
from app.core.settings import AppSetting
from app.core.logger import AppLogger


class App:
    """Core app object

    Args:
        object (_type_): _description_
    """

    def __init__(self):
        self.config = AppConfig()
        self.settings = self.__load_settings()

        self.logger = AppLogger(
            os.path.join(self.settings.logs_dir, "app_info.log"),
            os.path.join(self.settings.logs_dir, "app_error.log"),
            log_name="core_app.logger",
        )

    def get_main_config_dir(self):
        """
        Calculate/format main config directory

        Returns:
            _type_: _description_
        """
        return os.path.join(self.config.app_path, "config")

    def get_main_config_path(self):
        """
        Calculate/format main config directory

        Returns:
            _type_: _description_
        """
        return os.path.join(self.get_main_config_dir(), "app.conf")

    def __load_settings(self) -> AppSetting:
        """
        * Load config_parser_object from app.conf and put in dictionary
        """
        # interpolation is disabled
        config_obj = ConfigParser(interpolation=None)
        config_obj.read(self.get_main_config_path())

        # Directorie path reference
        app_data_dir = os.path.join(self.config.app_path, "app_data")
        local_temp_dir = os.path.join(self.config.app_path, "local_temp")
        if not os.path.exists(local_temp_dir):
            os.makedirs(local_temp_dir)

        logs_dir = os.path.join(self.config.app_path, "logs")
        user_dir = os.path.join(self.config.app_path, "user_data")

        output_dir = config_obj.get("MAIN", "output_dir")

        main_db_uri = config_obj.get("MAIN_DATABASE", "db_uri")
        main_db_name = config_obj.get("MAIN_DATABASE", "db_name")

        # Initialize Settings
        return AppSetting(
            config_dir=self.get_main_config_dir(),
            app_data_dir=app_data_dir,
            logs_dir=logs_dir,
            local_temp_dir=local_temp_dir,
            user_dir=user_dir,
            main_config_path=self.get_main_config_path(),
            output_dir=output_dir,
            main_db_uri=main_db_uri,
            main_db_name=main_db_name,
        )


def create_app():
    """
    * Function to create core app object
    """
    print("Initiating App")

    # Set Timezone
    os.environ["TZ"] = "Asia/Calcutta"
    time.tzset()

    app_obj = App()
    return app_obj


core_app = create_app()
