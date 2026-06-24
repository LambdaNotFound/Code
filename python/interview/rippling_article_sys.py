import heapq
from collections import defaultdict

# ── Design ────────────────────────────────────────────────────────────────────
#
# Data model (in-memory, no thread safety):
#   _articles     dict[int, str]                article_id → name; IDs auto-increment from 1
#   _scores       dict[int, int]                article_id → net score (upvotes − downvotes);
#                                               initialized to 0 on add_article so unvoted
#                                               articles appear in get_top_k results
#   _votes        dict[(user_id, article_id), str]  current vote per (user, article): "up"|"down"
#   _user_flips   dict[int, list[int]]          user_id → last ≤3 unique flipped article IDs,
#                                               ordered oldest→newest
#
# Flip semantics: a flip is a direction change on a previously voted article (up→down or
#   down→up). A first vote is never a flip. Re-voting the same direction is a no-op.
#
# _user_flips invariant: list is bounded at 3 unique article IDs.
#   On each flip of article X: remove X if already present (O(1), list ≤3), append X,
#   trim front if len > 3. print_last_three_flips reverses → O(1) retrieval.
#
# get_top_k: heapq.nlargest → O(n log k) vs O(n log n) full sort.
#   Tie-break by article_id ascending (lower ID = higher rank) for determinism.
#
# ─────────────────────────────────────────────────────────────────────────────


class ArticleSystem:
    def __init__(self):
        self._articles: dict[int, str] = {}
        self._next_id: int = 1
        self._scores: dict[int, int] = {}
        self._votes: dict[tuple[int, int], str] = {}
        self._user_flips: dict[int, list[int]] = defaultdict(list)

    def add_article(self, article_name: str) -> int:  # T: O(1), S: O(1)
        article_id = self._next_id
        self._articles[article_id] = article_name
        self._scores[article_id] = 0
        self._next_id += 1
        return article_id

    def _cast_vote(self, article_id: int, user_id: int, direction: str) -> None:
        key = (user_id, article_id)
        current = self._votes.get(key)

        if current == direction:
            return  # re-vote same direction: no-op

        if current is None:
            # First vote: update score, no flip
            self._scores[article_id] += 1 if direction == "up" else -1
        else:
            # Direction change: up→down = −2, down→up = +2
            self._scores[article_id] += 2 if direction == "up" else -2

            # Record flip: remove article if already in list, append, cap at 3
            flips = self._user_flips[user_id]
            if article_id in flips:
                flips.remove(article_id)   # O(1) — list bounded at 3
            flips.append(article_id)
            if len(flips) > 3:
                flips.pop(0)               # O(1) — list bounded at 4 before trim

        self._votes[key] = direction

    def upvote_article(self, article_id: int, user_id: int) -> None:   # T: O(1)
        self._cast_vote(article_id, user_id, "up")

    def downvote_article(self, article_id: int, user_id: int) -> None:  # T: O(1)
        self._cast_vote(article_id, user_id, "down")

    def print_last_three_flips(self, user_id: int) -> list[str]:  # T: O(1), S: O(1)
        flips = self._user_flips.get(user_id, [])
        return [self._articles[aid] for aid in reversed(flips)]

    def get_top_k(self, k: int) -> list[tuple[str, int]]:  # T: O(n log k), S: O(k)
        # nlargest key: primary = score descending, tie-break = article_id ascending
        top = heapq.nlargest(k, self._scores.items(), key=lambda x: (x[1], -x[0]))
        return [(self._articles[aid], score) for aid, score in top]
