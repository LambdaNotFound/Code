"""
Build a card game called Camel Cards, a simplified version of Poker. Two players each have a hand of 5 cards. Determine which player has the stronger hand.
Input: two string arrays representing hands. Output: which player wins (or tie/unknown).
Hand rankings (simplified poker rules, given to you in the prompt):

Evaluate hands against a set of type rules (e.g., "five of a kind", "four of a kind", "full house", "three of a kind", "two pair", "one pair", "high card" — exact set given in the problem)
Higher-ranked type wins
If same type, compare by card values in the order the player received the cards
"""

from abc import ABC, abstractmethod
from collections import Counter
from typing import Callable


class Hand:
    def __init__(self, cards: list[str]):
        # Time: O(N)  Space: O(N)
        self.cards = cards
        self._counts: Counter = Counter(cards)

    def count_values(self) -> list[int]:
        # Time: O(N log N) — sort up to N distinct counts  Space: O(N)
        return sorted(self._counts.values())


class HandRule(ABC):
    rank: int
    name: str

    @abstractmethod
    def matches(self, hand: Hand) -> bool: ...


class HighCardRule(HandRule):
    rank = 1
    name = "High Card"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [1, 1, 1, 1, 1]


class OnePairRule(HandRule):
    rank = 2
    name = "One Pair"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [1, 1, 1, 2]


class TwoPairRule(HandRule):
    rank = 3
    name = "Two Pair"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [1, 2, 2]


class ThreeOfAKindRule(HandRule):
    rank = 4
    name = "Three of a Kind"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [1, 1, 3]


class FullHouseRule(HandRule):
    rank = 5
    name = "Full House"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [2, 3]


class FourOfAKindRule(HandRule):
    rank = 6
    name = "Four of a Kind"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [1, 4]


class FiveOfAKindRule(HandRule):
    rank = 7
    name = "Five of a Kind"

    def matches(self, hand: Hand) -> bool:
        # Time: O(N log N) — delegates to count_values()  Space: O(N)
        return hand.count_values() == [5]


DEFAULT_RULES: list[HandRule] = [
    FiveOfAKindRule(),
    FourOfAKindRule(),
    FullHouseRule(),
    ThreeOfAKindRule(),
    TwoPairRule(),
    OnePairRule(),
    HighCardRule(),
]

DEFAULT_CARD_ORDER = "23456789TJQKA"


class HandEvaluator:
    def __init__(self, rules: list[HandRule]):
        # Time: O(R log R)  Space: O(R)
        self._rules = sorted(rules, key=lambda r: r.rank, reverse=True)

    def evaluate(self, hand: Hand) -> int:
        # Time: O(R · N log N) — up to R rules, each match is O(N log N)  Space: O(N)
        for rule in self._rules:
            if rule.matches(hand):
                return rule.rank
        return 0


class Game:
    def __init__(self, evaluator: HandEvaluator, card_order: str = DEFAULT_CARD_ORDER):
        self._evaluator = evaluator
        self._card_order = card_order

    def _card_value(self, card: str) -> int:
        # Time: O(C) — linear scan of card_order string  Space: O(1)
        return self._card_order.index(card)

    def compare(self, h1: Hand, h2: Hand) -> str:
        # Time: O(R · N log N + N · C) — two evaluations + O(N · C) tiebreak  Space: O(N)
        rank1 = self._evaluator.evaluate(h1)
        rank2 = self._evaluator.evaluate(h2)

        if rank1 != rank2:
            return "player1" if rank1 > rank2 else "player2"

        for c1, c2 in zip(h1.cards, h2.cards):
            v1, v2 = self._card_value(c1), self._card_value(c2)
            if v1 != v2:
                return "player1" if v1 > v2 else "player2"

        return "tie"


"""
Part 2: Incomplete hands
Hands arrive over a network and may be incomplete due to packet loss.
One or both hands might have fewer than 5 cards.
Example: [9,9,9,9] vs [9] → output UNKNOWN (can't determine winner with incomplete information).

Strategy: compute [min_strength, max_strength] for each hand, then:
  min(A) > max(B)  →  player1 wins definitively
  min(B) > max(A)  →  player2 wins definitively
  otherwise        →  UNKNOWN
"""


class CompletionStrategy(ABC):
    @abstractmethod
    def complete(self, hand: Hand, card_order: str) -> Hand: ...


class MaxCompletionStrategy(CompletionStrategy):
    """Fill remaining slots with the highest-frequency existing card.
    On count tie, prefer the higher face-value card for a better tiebreaker.
    If the hand is empty, use the highest card in card_order."""

    def complete(self, hand: Hand, card_order: str) -> Hand:
        # Time: O(N · C) — max() scans N distinct keys, each with O(C) index lookup  Space: O(N)
        r = 5 - len(hand.cards)
        if r == 0:
            return hand
        if not hand.cards:
            fill = card_order[-1]
        else:
            fill = max(
                hand._counts, key=lambda c: (hand._counts[c], card_order.index(c))
            )
        return Hand(hand.cards + [fill] * r)


class MinCompletionStrategy(CompletionStrategy):
    """Fill remaining slots with distinct cards not yet in the hand, lowest-value first.
    This spreads the count histogram to minimise the achievable type rank."""

    def complete(self, hand: Hand, card_order: str) -> Hand:
        # Time: O(C) — one pass over card_order with O(1) set lookups  Space: O(C)
        r = 5 - len(hand.cards)
        if r == 0:
            return hand
        seen = set(hand._counts)
        unseen = [c for c in card_order if c not in seen]
        # unseen always has >= r cards: at most 4 existing labels, 13 total, r <= 4
        added = (
            unseen[:r]
            if len(unseen) >= r
            else unseen + [card_order[0]] * (r - len(unseen))
        )
        return Hand(hand.cards + added)


class RangeComparator:
    def __init__(
        self,
        game: Game,
        max_strategy: MaxCompletionStrategy | None = None,
        min_strategy: MinCompletionStrategy | None = None,
    ):
        self._game = game
        self._max = max_strategy or MaxCompletionStrategy()
        self._min = min_strategy or MinCompletionStrategy()

    def compare(self, h1: Hand, h2: Hand) -> str:
        # Time: O(R · N log N + N · C) — four completions + two Game.compare calls  Space: O(N)
        order = self._game._card_order
        min1 = self._min.complete(h1, order)
        max1 = self._max.complete(h1, order)
        min2 = self._min.complete(h2, order)
        max2 = self._max.complete(h2, order)

        if self._game.compare(min1, max2) == "player1":
            return "player1"
        if self._game.compare(max1, min2) == "player2":
            return "player2"
        return "UNKNOWN"


"""
Part 3. the extensibility API:
add_type(type_id: string, matches_function: function)
evaluate(hand1, hand2, evaluation_orders: array of string)

add_type registers a new hand type with a custom matching function at runtime
evaluate takes two hands plus an ordered array of type IDs that defines the ranking hierarchy for this particular game
This means hand-type rankings are not hardcoded — the caller controls the evaluation order
"""
class ExtensibleGame:
    """Wraps Game for card ordering/tiebreaking; adds a runtime type registry
    and a caller-controlled evaluation order.

    evaluation_order convention: weakest-to-strongest (index 0 = rank 1).
    """

    def __init__(self, game: Game):
        # Time: O(R)  Space: O(R)
        self._game = game
        self._type_registry: dict[str, Callable[[Hand], bool]] = {}
        for rule in DEFAULT_RULES:
            self.add_type(rule.name, rule.matches)

    def add_type(self, type_id: str, matches_fn: Callable[[Hand], bool]) -> None:
        # Time: O(1)  Space: O(1)
        self._type_registry[type_id] = matches_fn

    def _score(self, hand: Hand, evaluation_order: list[str]) -> int:
        # Time: O(T · N log N) — up to T types, each matches_fn is O(N log N)  Space: O(N)
        for i in range(len(evaluation_order) - 1, -1, -1):
            type_id = evaluation_order[i]
            if type_id not in self._type_registry:
                raise ValueError(f"Unknown type: {type_id!r}")
            if self._type_registry[type_id](hand):
                return i + 1
        return 0

    def evaluate(self, hand1: Hand, hand2: Hand, evaluation_order: list[str]) -> str:
        # Time: O(T · N log N + N · C) — two _score calls + O(N · C) tiebreak  Space: O(N)
        score1 = self._score(hand1, evaluation_order)
        score2 = self._score(hand2, evaluation_order)

        if score1 != score2:
            return "player1" if score1 > score2 else "player2"

        for c1, c2 in zip(hand1.cards, hand2.cards):
            v1, v2 = self._game._card_value(c1), self._game._card_value(c2)
            if v1 != v2:
                return "player1" if v1 > v2 else "player2"

        return "tie"
