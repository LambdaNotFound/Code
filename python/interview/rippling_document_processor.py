import heapq
import re
from abc import ABC, abstractmethod
from collections import defaultdict


class Handler(ABC):
    @abstractmethod
    def process(self, document: str): ...


class LengthHandler(Handler):
    def process(self, document: str) -> int:  # T: O(n), S: O(1)
        return sum(1 for c in document if c.isalpha())


class WordCountHandler(Handler):
    def process(self, document: str) -> int:  # T: O(n), S: O(n)
        return len(re.sub(r"[^a-zA-Z\s]", "", document).split())


class _RevStr:
    """Inverts string comparison so a min-heap pops the lexically largest word on count ties."""

    __slots__ = ("s",)

    def __init__(self, s: str):
        self.s = s

    def __lt__(self, other: "_RevStr") -> bool:
        return self.s > other.s

    def __eq__(self, other: object) -> bool:
        return isinstance(other, _RevStr) and self.s == other.s


class CommonWordHandler(Handler):
    def __init__(self, top_k: int = 1):
        self._top_k = top_k

    def process(self, document: str) -> list[str]:  # T: O(n log k), S: O(n)
        words = re.sub(r"[^a-z\s]", "", document.lower()).split()

        freq: dict[str, int] = defaultdict(int)
        for word in words:
            freq[word] += 1

        # Fixed-size min-heap: (count, _RevStr(word))
        # Pops worst candidate = smallest count; on tie = lexically largest word
        heap: list[tuple[int, _RevStr]] = []
        for word, count in freq.items():
            heapq.heappush(heap, (count, _RevStr(word)))
            if len(heap) > self._top_k:
                heapq.heappop(heap)

        return [e.s for _, e in sorted(heap, key=lambda x: (-x[0], x[1].s))]


class DocumentProcessor:
    def __init__(self, handlers: list[Handler]):
        self._handlers = handlers

    def process(self, document: str) -> list:  # T: O(k · cost_per_handler), S: O(k)
        return [h.process(document) for h in self._handlers]
