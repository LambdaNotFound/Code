import sys, os
sys.path.insert(0, os.path.dirname(__file__))

import pytest
from rippling_document_processor import (
    DocumentProcessor, LengthHandler, WordCountHandler, CommonWordHandler,
)


class TestLengthHandler:
    @pytest.mark.parametrize("doc, expected", [
        ("hello world", 10),
        ("Hello, World!", 10),
        ("", 0),
        ("   ", 0),
        ("abc 123!!", 3),  # digits and punctuation excluded
    ])
    def test_length(self, doc, expected):
        assert LengthHandler().process(doc) == expected


class TestWordCountHandler:
    @pytest.mark.parametrize("doc, expected", [
        ("hello world", 2),
        ("hello, world!", 2),
        ("", 0),
        ("one two three", 3),
        ("  spaces   between  ", 2),
    ])
    def test_word_count(self, doc, expected):
        assert WordCountHandler().process(doc) == expected


class TestCommonWordHandler:
    @pytest.mark.parametrize("doc, top_k, expected", [
        ("the cat sat on the mat", 1, ["the"]),
        # "the"→2; cat/mat/on/sat→1 each; top 3 after "the": cat < mat < on alphabetically
        ("the cat sat on the mat", 3, ["the", "cat", "mat"]),
        # fewer than top_k unique words → return what's available
        ("one", 3, ["one"]),
        # case-insensitive
        ("Hello hello HELLO", 1, ["hello"]),
        # punctuation stripped: "it's" → "its"
        ("it's it's nice", 1, ["its"]),
    ])
    def test_common_word(self, doc, top_k, expected):
        assert CommonWordHandler(top_k).process(doc) == expected


class TestDocumentProcessor:
    def test_multiple_handlers_run_independently(self):
        dp = DocumentProcessor([LengthHandler(), WordCountHandler()])
        assert dp.process("hello, world!") == [10, 2]

    def test_single_handler(self):
        dp = DocumentProcessor([CommonWordHandler(1)])
        assert dp.process("foo foo bar") == [["foo"]]
