from abc import ABC, abstractmethod


class LogHandler(ABC):
    @abstractmethod
    def process(self, s: str) -> None: ...


class FilterHandler(LogHandler):
    def __init__(self, filter_str: str):
        self._filter = filter_str

    def process(self, s: str) -> None:  # T: O(n), S: O(n) — str.replace builds a new string of length n
        print(s.replace(self._filter, ""))


class TruncationHandler(LogHandler):
    def __init__(self, n: int):
        self._n = n

    def process(self, s: str) -> None:  # T: O(n), S: O(n) — slice creates a new string of up to n chars
        print(s[: self._n])


class UppercaseHandler(LogHandler):
    def process(self, s: str) -> None:  # T: O(n), S: O(n) — upper() allocates a new string of length n
        print(s.upper())


class ArrayHandler(LogHandler):
    def __init__(self):
        self.logs: list[str] = []

    def process(self, s: str) -> None:  # T: O(1) amortized, S: O(n) total across all appended strings
        self.logs.append(s)


class Logger:
    def __init__(self, handlers: list[LogHandler]):
        self._handlers = handlers

    def log(self, s: str) -> None:  # T: O(k·n), S: O(1) — calls k handlers each doing O(n) work on string of length n
        for handler in self._handlers:
            handler.process(s)
