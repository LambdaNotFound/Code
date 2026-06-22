import sys, os
sys.path.insert(0, os.path.dirname(__file__))

import pytest
from rippling_camel_card import (
    Hand,
    HandEvaluator,
    Game,
    RangeComparator,
    MaxCompletionStrategy,
    MinCompletionStrategy,
    ExtensibleGame,
    DEFAULT_RULES,
    DEFAULT_CARD_ORDER,
    HighCardRule,
    OnePairRule,
    TwoPairRule,
    ThreeOfAKindRule,
    FullHouseRule,
    FourOfAKindRule,
    FiveOfAKindRule,
)


@pytest.fixture
def game() -> Game:
    return Game(HandEvaluator(DEFAULT_RULES), DEFAULT_CARD_ORDER)


# --- Hand.count_values ---

@pytest.mark.parametrize("cards,expected", [
    (["A", "A", "A", "A", "A"], [5]),
    (["A", "A", "A", "A", "K"], [1, 4]),
    (["A", "A", "A", "K", "K"], [2, 3]),
    (["A", "A", "A", "K", "Q"], [1, 1, 3]),
    (["A", "A", "K", "K", "Q"], [1, 2, 2]),
    (["A", "A", "K", "Q", "J"], [1, 1, 1, 2]),
    (["A", "K", "Q", "J", "T"], [1, 1, 1, 1, 1]),
])
def test_hand_count_values(cards, expected):
    assert Hand(cards).count_values() == expected


# --- HandRule.matches ---

@pytest.mark.parametrize("rule,cards,expected", [
    (FiveOfAKindRule(),  ["A", "A", "A", "A", "A"], True),
    (FiveOfAKindRule(),  ["A", "A", "A", "A", "K"], False),
    (FourOfAKindRule(),  ["A", "A", "A", "A", "K"], True),
    (FourOfAKindRule(),  ["A", "A", "A", "K", "K"], False),
    (FullHouseRule(),    ["A", "A", "A", "K", "K"], True),
    (FullHouseRule(),    ["A", "A", "A", "K", "Q"], False),
    (ThreeOfAKindRule(), ["A", "A", "A", "K", "Q"], True),
    (ThreeOfAKindRule(), ["A", "A", "K", "K", "Q"], False),
    (TwoPairRule(),      ["A", "A", "K", "K", "Q"], True),
    (TwoPairRule(),      ["A", "A", "K", "Q", "J"], False),
    (OnePairRule(),      ["A", "A", "K", "Q", "J"], True),
    (OnePairRule(),      ["A", "K", "Q", "J", "T"], False),
    (HighCardRule(),     ["A", "K", "Q", "J", "T"], True),
    (HighCardRule(),     ["A", "A", "K", "Q", "J"], False),
])
def test_rule_matches(rule, cards, expected):
    assert rule.matches(Hand(cards)) == expected


# --- HandEvaluator.evaluate ---

@pytest.mark.parametrize("cards,expected_rank", [
    (["A", "A", "A", "A", "A"], 7),
    (["A", "A", "A", "A", "K"], 6),
    (["A", "A", "A", "K", "K"], 5),
    (["A", "A", "A", "K", "Q"], 4),
    (["A", "A", "K", "K", "Q"], 3),
    (["A", "A", "K", "Q", "J"], 2),
    (["A", "K", "Q", "J", "T"], 1),
])
def test_evaluator_rank(cards, expected_rank):
    evaluator = HandEvaluator(DEFAULT_RULES)
    assert evaluator.evaluate(Hand(cards)) == expected_rank


# --- Game.compare: higher type wins ---

@pytest.mark.parametrize("h1_cards,h2_cards,expected", [
    # five-of-a-kind beats four-of-a-kind
    (["A","A","A","A","A"], ["K","K","K","K","Q"], "player1"),
    # four-of-a-kind beats full house
    (["2","2","2","2","3"], ["A","A","A","K","K"], "player1"),
    # full house beats three-of-a-kind
    (["Q","Q","Q","J","J"], ["A","A","A","K","Q"], "player1"),
    # three-of-a-kind beats two pair
    (["5","5","5","2","3"], ["A","A","K","K","Q"], "player1"),
    # two pair beats one pair
    (["A","A","K","K","Q"], ["Q","Q","J","T","9"], "player1"),
    # one pair beats high card
    (["A","K","Q","J","J"], ["A","K","Q","J","T"], "player1"),
])
def test_compare_by_type(game, h1_cards, h2_cards, expected):
    assert game.compare(Hand(h1_cards), Hand(h2_cards)) == expected


# --- Game.compare: same type, tie-break by dealt order ---

@pytest.mark.parametrize("h1_cards,h2_cards,expected", [
    # both high card; first card A vs K → player1
    (["A","Q","J","T","9"], ["K","Q","J","T","9"], "player1"),
    # both high card; first cards match, second A vs Q → player1
    (["K","A","J","T","9"], ["K","Q","J","T","9"], "player1"),
    # both one pair; pair of Ks vs pair of Qs → player1
    (["K","K","A","J","T"], ["Q","Q","A","J","T"], "player1"),
    # both one pair; same pair, kicker decides at position 2
    (["K","K","A","J","T"], ["K","K","Q","J","T"], "player1"),
    # both full house; trip rank decides (position 0 of dealt order)
    (["A","A","A","K","K"], ["K","K","K","A","A"], "player1"),
    # tie-break favors player2 when their card at first difference is higher
    (["3","K","Q","J","T"], ["3","A","Q","J","T"], "player2"),
])
def test_compare_tiebreak(game, h1_cards, h2_cards, expected):
    assert game.compare(Hand(h1_cards), Hand(h2_cards)) == expected


# --- Game.compare: exact tie ---

def test_compare_tie(game):
    cards = ["A", "K", "Q", "J", "T"]
    assert game.compare(Hand(cards), Hand(cards)) == "tie"


# --- Game._card_value respects card order ---

@pytest.mark.parametrize("low,high", [
    ("2", "3"), ("9", "T"), ("T", "J"), ("J", "Q"), ("Q", "K"), ("K", "A"),
])
def test_card_order(game, low, high):
    assert game._card_value(low) < game._card_value(high)


# ── Part 2: Incomplete hands ────────────────────────────────────────────────

@pytest.fixture
def rc(game) -> RangeComparator:
    return RangeComparator(game)


# --- MaxCompletionStrategy ---

@pytest.mark.parametrize("cards,expected_cards", [
    # already complete — unchanged
    (["A","K","Q","J","T"], ["A","K","Q","J","T"]),
    # single card — five of a kind
    (["9"], ["9","9","9","9","9"]),
    # two cards, tied count → prefer higher face value (K over 9)
    (["9","K"], ["9","K","K","K","K"]),
    # three cards, clear mode
    (["9","9","K"], ["9","9","K","9","9"]),
    # empty hand — fill with "A" (highest in card_order)
    ([], ["A","A","A","A","A"]),
])
def test_max_completion(cards, expected_cards):
    result = MaxCompletionStrategy().complete(Hand(cards), DEFAULT_CARD_ORDER)
    assert result.cards == expected_cards


# --- MinCompletionStrategy ---

@pytest.mark.parametrize("cards,expected_cards", [
    # already complete — unchanged
    (["A","K","Q","J","T"], ["A","K","Q","J","T"]),
    # single 9 → fill with four lowest unseen: 2,3,4,5
    (["9"], ["9","2","3","4","5"]),
    # existing pair stays; fill with two lowest unseen
    (["9","9","K"], ["9","9","K","2","3"]),
    # three distinct cards → fill with two lowest unseen
    (["A","K","Q"], ["A","K","Q","2","3"]),
    # empty hand — fill with five lowest cards
    ([], ["2","3","4","5","6"]),
])
def test_min_completion(cards, expected_cards):
    result = MinCompletionStrategy().complete(Hand(cards), DEFAULT_CARD_ORDER)
    assert result.cards == expected_cards


# --- RangeComparator.compare ---

@pytest.mark.parametrize("h1_cards,h2_cards,expected", [
    # classic UNKNOWN: [9,9,9,9] vs [9] — B can also complete to five-of-a-kind
    (["9","9","9","9"], ["9"], "UNKNOWN"),
    # A definitively wins: complete five-of-a-kind Aces vs incomplete single 2
    # even worst A (four-of-a-kind Aces) beats best B (five 2s) on tiebreaker — wait:
    # actually A is already [A,A,A,A,A], min=max same; B max=[2,2,2,2,2]; five 2s vs five As
    # → both rank 7, card[0]: 2 vs A → player2... so UNKNOWN; use complete hand instead:
    (["A","A","A","A","A"], ["2"], "player1"),
    # B definitively wins: single 2 vs complete five Aces
    (["2"], ["A","A","A","A","A"], "player2"),
    # A definitively wins: four Aces vs two 2-3; best B = [2,3,3,3,3] rank 6; min A = [A,A,A,A,2] rank 6
    # tiebreaker: A vs 2 → player1
    (["A","A","A","A"], ["2","3"], "player1"),
    # both complete — delegates correctly to Game.compare
    (["A","A","A","A","A"], ["K","K","K","K","K"], "player1"),
    (["2","2","2","2","2"], ["A","A","A","A","A"], "player2"),
    # both incomplete, ranges overlap → UNKNOWN
    (["A","K"], ["Q","J"], "UNKNOWN"),
])
def test_range_comparator(rc, h1_cards, h2_cards, expected):
    assert rc.compare(Hand(h1_cards), Hand(h2_cards)) == expected


# ── Part 3: ExtensibleGame ───────────────────────────────────────────────────

BUILTIN_ORDER = [
    "High Card", "One Pair", "Two Pair", "Three of a Kind",
    "Full House", "Four of a Kind", "Five of a Kind",
]


@pytest.fixture
def eg(game) -> ExtensibleGame:
    return ExtensibleGame(game)


# --- add_type: built-in types pre-registered ---

@pytest.mark.parametrize("h1_cards,h2_cards,expected", [
    (["A","A","A","A","A"], ["K","K","K","K","K"], "player1"),
    (["A","A","A","A","K"], ["A","A","A","K","K"], "player1"),
    (["A","K","Q","J","T"], ["A","K","Q","J","9"], "player1"),
    (["A","K","Q","J","T"], ["A","K","Q","J","T"], "tie"),
])
def test_evaluate_builtin_types(eg, h1_cards, h2_cards, expected):
    assert eg.evaluate(Hand(h1_cards), Hand(h2_cards), BUILTIN_ORDER) == expected


# --- evaluation_order controls ranking ---

def test_evaluate_custom_order_reverses_ranks(eg):
    # reverse the order so High Card is strongest and Five of a Kind is weakest
    reversed_order = list(reversed(BUILTIN_ORDER))
    h1 = Hand(["A", "K", "Q", "J", "T"])  # high card
    h2 = Hand(["A", "A", "A", "A", "A"])  # five of a kind
    # in reversed order, High Card has highest rank → player1 wins
    assert eg.evaluate(h1, h2, reversed_order) == "player1"


def test_evaluate_partial_order_only_listed_types_matter(eg):
    # only "One Pair" and "Two Pair" in the order; Five of a Kind scores 0
    partial_order = ["One Pair", "Two Pair"]
    h1 = Hand(["A", "A", "K", "K", "Q"])  # two pair  → score 2
    h2 = Hand(["A", "A", "A", "A", "A"])  # five of a kind → score 0 (not listed)
    assert eg.evaluate(h1, h2, partial_order) == "player1"


# --- add_type: register a custom type at runtime ---

def test_add_type_custom_rule(eg):
    # register "All Red" (fake suit-based rule for test purposes — checks all cards are "R")
    eg.add_type("All Same Suit", lambda hand: len(set(hand.cards)) == 1)
    custom_order = BUILTIN_ORDER + ["All Same Suit"]
    h1 = Hand(["A", "A", "A", "A", "A"])  # matches "All Same Suit" AND "Five of a Kind"
    h2 = Hand(["K", "K", "K", "K", "Q"])  # matches "Four of a Kind" only
    # "All Same Suit" has rank 8 > "Five of a Kind" rank 7
    assert eg.evaluate(h1, h2, custom_order) == "player1"


def test_add_type_overwrites_existing(eg):
    # redefine "Five of a Kind" to never match
    eg.add_type("Five of a Kind", lambda hand: False)
    h1 = Hand(["A", "A", "A", "A", "A"])  # would be five of a kind, now scores 0
    h2 = Hand(["K", "K", "K", "K", "Q"])  # four of a kind → score 6
    assert eg.evaluate(h1, h2, BUILTIN_ORDER) == "player2"


# --- unknown type_id raises ValueError ---

def test_evaluate_unknown_type_raises(eg):
    with pytest.raises(ValueError, match="Unknown type"):
        eg.evaluate(Hand(["A","K","Q","J","T"]), Hand(["A","K","Q","J","9"]), ["no_such_type"])


# --- tiebreaker still applies under custom order ---

def test_evaluate_tiebreak_within_same_type(eg):
    h1 = Hand(["A", "A", "K", "Q", "J"])  # one pair, higher pair card
    h2 = Hand(["K", "K", "A", "Q", "J"])  # one pair, lower pair card
    # both score rank 2 (one pair); tiebreak: position 0: A vs K → player1
    assert eg.evaluate(h1, h2, BUILTIN_ORDER) == "player1"
