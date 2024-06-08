from enum import Enum


class CustomEnumClass:
    _meta_dictionary = {}

    @classmethod
    def get_string(cls, value):
        if isinstance(value, Enum):
            value = value.value
        return cls._meta_dictionary.get(value)

    @classmethod
    def get_list(cls):
        return cls._meta_dictionary.keys()

    @classmethod
    def get_dictionary(cls):
        return cls._meta_dictionary
