import sys, os
sys.path.insert(0, os.path.dirname(__file__))

from rippling_logger_system import (
    Logger,
    FilterHandler,
    TruncationHandler,
    UppercaseHandler,
    ArrayHandler,
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
