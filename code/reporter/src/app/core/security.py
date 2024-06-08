import re
from xml.sax.saxutils import escape


def sanitize_xss(input_string):
    # escape() takes care of &, < and >.
    if not input_string:
        return ""
    characters_to_replace = {
        ")": "&#41;",
        "(": '&#40;',
        '"': "&quot;",
        "+": "&#43;",
        "{": '&#123;',
        "}": '&#125;',
        "[": '&#91;',
        "]": '&#93;',
        "'": "&apos;",
        "<": "&lt;",
        ">": "&gt;",
    }
    return escape(input_string, characters_to_replace)


def is_valid_table_name(table_name):
    if table_name and re.match(r"^[A-Za-z][a-zA-Z0-9._-]{1,40}$", table_name):
        return True
    return False