import sys, os
sys.path.insert(0, os.path.dirname(__file__))

from unittest.mock import MagicMock, call
from rippling_logger_system import (
    Logger,
    FilterHandler,
    TruncationHandler,
    UppercaseHandler,
    ArrayHandler,
    DatabaseHandler,
    HttpClient,
    RemoteAPIHandler,
)


def test_filter_handler(capsys):
    h = FilterHandler("world")
    h.process("hello world")
    assert capsys.readouterr().out.rstrip("\n") == "hello "


def test_truncation_handler(capsys):
    h = TruncationHandler(5)
    h.process("hello world")
    assert capsys.readouterr().out.strip() == "hello"


def test_uppercase_handler(capsys):
    h = UppercaseHandler()
    h.process("hello world")
    assert capsys.readouterr().out.strip() == "HELLO WORLD"


def test_array_handler_no_print(capsys):
    h = ArrayHandler()
    h.process("hello world")
    assert h.logs == ["hello world"]
    assert capsys.readouterr().out == ""


def test_array_handler_accumulates():
    h = ArrayHandler()
    h.process("first")
    h.process("second")
    assert h.logs == ["first", "second"]


def test_logger_fan_out(capsys):
    array_handler = ArrayHandler()
    logger = Logger([
        FilterHandler("world"),
        TruncationHandler(5),
        UppercaseHandler(),
        array_handler,
    ])

    logger.log("hello world")

    lines = capsys.readouterr().out.splitlines()
    assert lines[0] == "hello "       # filter removed "world"
    assert lines[1] == "hello"        # truncated to 5
    assert lines[2] == "HELLO WORLD"  # uppercased original
    assert array_handler.logs == ["hello world"]  # stored original, no print


def test_logger_each_handler_gets_original(capsys):
    # Confirms fan-out: FilterHandler does not affect what TruncationHandler sees
    array_handler = ArrayHandler()
    logger = Logger([FilterHandler("hello "), TruncationHandler(3), array_handler])
    logger.log("hello world")

    lines = capsys.readouterr().out.splitlines()
    assert lines[0] == "world"   # "hello " removed
    assert lines[1] == "hel"     # original truncated, not the filtered version
    assert array_handler.logs == ["hello world"]


# --- DatabaseHandler ---

def test_database_handler_inserts_message():
    conn = MagicMock()
    h = DatabaseHandler(connection_factory=lambda: conn)
    h.process("hello world")

    conn.execute.assert_called_once_with("INSERT INTO logs (message) VALUES (?)", ("hello world",))
    conn.commit.assert_called_once()
    conn.close.assert_called_once()


def test_database_handler_new_connection_per_call():
    factory = MagicMock(return_value=MagicMock())
    h = DatabaseHandler(connection_factory=factory)
    h.process("first")
    h.process("second")

    assert factory.call_count == 2


def test_database_handler_closes_on_error():
    conn = MagicMock()
    conn.execute.side_effect = RuntimeError("db error")
    h = DatabaseHandler(connection_factory=lambda: conn)

    try:
        h.process("hello")
    except RuntimeError:
        pass

    conn.close.assert_called_once()


# --- RemoteAPIHandler ---

class _FakeHttpClient(HttpClient):
    def __init__(self):
        self.calls: list[tuple[str, dict]] = []

    def post(self, url: str, payload: dict) -> None:
        self.calls.append((url, payload))


def test_remote_api_handler_posts_message():
    client = _FakeHttpClient()
    h = RemoteAPIHandler(endpoint="http://example.com/logs", client=client)
    h.process("hello world")

    assert client.calls == [("http://example.com/logs", {"message": "hello world"})]


def test_remote_api_handler_uses_configured_endpoint():
    client = _FakeHttpClient()
    h = RemoteAPIHandler(endpoint="http://other.com/ingest", client=client)
    h.process("test")

    assert client.calls[0][0] == "http://other.com/ingest"


def test_remote_api_handler_each_call_posts_independently():
    client = _FakeHttpClient()
    h = RemoteAPIHandler(endpoint="http://example.com/logs", client=client)
    h.process("first")
    h.process("second")

    assert client.calls == [
        ("http://example.com/logs", {"message": "first"}),
        ("http://example.com/logs", {"message": "second"}),
    ]
