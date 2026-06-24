from abc import ABC, abstractmethod
from typing import Callable
import json
import urllib.request


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


class DatabaseHandler(LogHandler):
    # connection_factory is called once per log event to open a fresh connection (per spec)
    def __init__(self, connection_factory: Callable[[], object]):
        self._factory = connection_factory

    def process(self, s: str) -> None:  # T: O(n), S: O(1) — IO-bound; n = message length
        conn = self._factory()
        try:
            conn.execute("INSERT INTO logs (message) VALUES (?)", (s,))
            conn.commit()
        finally:
            conn.close()


class HttpClient(ABC):
    @abstractmethod
    def post(self, url: str, payload: dict) -> None: ...


class DefaultHttpClient(HttpClient):
    def post(self, url: str, payload: dict) -> None:
        data = json.dumps(payload).encode()
        req = urllib.request.Request(url, data=data, headers={"Content-Type": "application/json"})
        with urllib.request.urlopen(req):
            pass


class RemoteAPIHandler(LogHandler):
    def __init__(self, endpoint: str, client: HttpClient | None = None):
        self._endpoint = endpoint
        self._client = client or DefaultHttpClient()

    def process(self, s: str) -> None:  # T: O(n), S: O(n) — IO-bound; n = message length for JSON encoding
        self._client.post(self._endpoint, {"message": s})


class Logger:
    def __init__(self, handlers: list[LogHandler]):
        self._handlers = handlers

    def log(self, s: str) -> None:  # T: O(k·n), S: O(1) — calls k handlers each doing O(n) work on string of length n
        for handler in self._handlers:
            handler.process(s)
