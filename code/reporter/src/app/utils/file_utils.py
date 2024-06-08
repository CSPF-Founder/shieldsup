"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""
import os
from shutil import move
from tempfile import mkstemp


def replace_str_in_file(filename, str_to_search, replace_with):
    """
    Method that replaces the specified string with another string in a file

    :param filename:    File Name
    :param str_to_search:    Keyword that needs to be replaced
    :param replace_with: New Keyword

    """
    # Create temp file
    temp_file_hanle, abs_path = mkstemp()
    with open(abs_path, "w", encoding="utf-8") as new_file:
        with open(filename, "r", encoding="utf-8") as old_file:
            for line in old_file:
                new_file.write(line.replace(str_to_search, replace_with))
    # close temp file
    os.close(temp_file_hanle)
    # Remove original file
    os.remove(filename)
    # Move new file
    move(abs_path, filename)


def replace_str_list_in_file(
    file_path: str,
    meta_data: list,
):
    temp_file_hanle, abs_path = mkstemp()
    with open(abs_path, "w", encoding="utf-8") as new_file:
        with open(file_path, "r", encoding="utf-8") as old_file:
            for meta in meta_data:
                str_to_search = meta["str_to_search"]
                replace_with = meta["replace_with"]
                for line in old_file:
                    new_file.write(line.replace(str_to_search, replace_with))
    # close temp file
    os.close(temp_file_hanle)
    # Remove original file
    os.remove(file_path)
    # Move new file
    move(abs_path, file_path)
