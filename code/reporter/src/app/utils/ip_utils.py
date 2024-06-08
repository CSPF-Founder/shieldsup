import re

import ipaddress


def cidr_to_ip_list(input_data: str) -> list | None:
    if not input_data or "/" not in input_data:
        return None

    try:
        range_ip_parse_result = ipaddress.ip_network(input_data, strict=False)

        if range_ip_parse_result:
            ip_list = [str(ip) for ip in range_ip_parse_result]
            return ip_list
    except Exception:
        return None
    return None


def get_ip_count_if_range(input_data: str) -> int:
    """
    * Get IP count if range

    Args:
        input_data (str): _description_

    Returns:
        int: _description_
    """
    if not input_data or "/" not in input_data:
        return 0

    try:
        ip_list = cidr_to_ip_list(input_data)
        if not ip_list:
            return 0

        return len(ip_list)
    except Exception:
        return 0
    return 0


def is_valid_cidr_by_regex(cidr_str: str) -> bool:
    # Define the regular expression pattern for CIDR notation
    pattern = (
        r"^(?:(?:25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|"
        r"[0-1]?[0-9][0-9]?)\/(?:[1-9]|[1-2][0-9]|3[0-2])$"
    )

    # Use the re.match function to check if the input matches the pattern
    return bool(re.match(pattern, cidr_str))
