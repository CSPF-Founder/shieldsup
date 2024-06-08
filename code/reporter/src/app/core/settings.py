"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from dataclasses import dataclass


@dataclass
class AppSetting:
    """
    * Class to store app level settings
    """

    # Directories Reference
    config_dir: str
    app_data_dir: str
    logs_dir: str
    local_temp_dir: str
    user_dir: str
    main_config_path: str

    output_dir: str

    main_db_uri: str
    main_db_name: str
