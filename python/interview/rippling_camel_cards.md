# Camel Cards — Design

## Problem Summary

Two players each hold a 5-card hand (strings). Determine which player wins, or if it's a tie.
Winning logic:
1. The hand with the higher **type rank** wins.
2. If types are equal, compare cards **left-to-right** by face value; the first difference decides.

---

## Class Design

### `Hand`

Owns the card data and exposes frequency counts for rule evaluation.

```
Hand
  cards: list[str]           # ordered as dealt, used for tie-breaking
  _counts: Counter[str, int] # hashmap: card label -> frequency
  count_values() -> list[int] # sorted frequencies, e.g. [2,3] for full house
```

The `Counter` (hashmap) is the single source of truth for pattern matching.
`cards` is kept in original order for left-to-right tie-breaking.

---

### `HandRule` (abstract strategy)

Each concrete rule encapsulates one hand type.

```
HandRule (ABC)
  rank: int      # numerical strength; higher = stronger
  name: str      # human-readable label
  matches(hand: Hand) -> bool
```

Rules are **open for extension** — add a new subclass to support new hand types without touching existing code.

#### Concrete rules (weakest → strongest)

| Class              | rank | Pattern (`sorted(count_values())`) |
|--------------------|------|--------------------------------------|
| `HighCardRule`     | 1    | `[1,1,1,1,1]`                       |
| `OnePairRule`      | 2    | `[1,1,1,2]`                         |
| `TwoPairRule`      | 3    | `[1,2,2]`                           |
| `ThreeOfAKindRule` | 4    | `[1,1,3]`                           |
| `FullHouseRule`    | 5    | `[2,3]`                             |
| `FourOfAKindRule`  | 6    | `[1,4]`                             |
| `FiveOfAKindRule`  | 7    | `[5]`                               |

Pattern matching is pure: `sorted(hand.count_values()) == expected_pattern`.

---

### `HandEvaluator`

Holds the ordered rule set and scores a hand.

```
HandEvaluator
  _rules: list[HandRule]   # sorted descending by rank at construction

  evaluate(hand: Hand) -> int
    # iterate rules highest-first; return rank of first match
    # falls back to HighCardRule rank (1) if nothing matches
```

Injecting `_rules` at construction makes the evaluator configurable — swap in a different rule set for a variant game.

---

### `Game`

Top-level coordinator. Owns card face-value ordering for tie-breaking.

```
Game
  _evaluator: HandEvaluator
  _card_order: str   # e.g. "23456789TJQKA" (left = weakest)

  compare(h1: Hand, h2: Hand) -> str   # "player1" | "player2" | "tie"

  _card_value(card: str) -> int
    # index into _card_order; higher index = stronger card
```

`compare` algorithm:
1. Score both hands with `_evaluator.evaluate()`.
2. If scores differ, return the player with the higher score.
3. If equal, iterate paired cards left-to-right; return the player whose card value is first higher.
4. If all cards match, return `"tie"`.

---

## Sequence: comparing two hands

```
Game.compare(h1, h2)
  │
  ├─ HandEvaluator.evaluate(h1)
  │     └─ FiveOfAKindRule.matches(h1)?  rank=7  (check sorted counts == [5])
  │     └─ FourOfAKindRule.matches(h1)?  rank=6  ...
  │     └─ ...first match wins
  │
  ├─ HandEvaluator.evaluate(h2)  (same)
  │
  ├─ rank(h1) vs rank(h2)  →  done if different
  │
  └─ card-by-card: Game._card_value(h1.cards[i]) vs _card_value(h2.cards[i])
```

---

## Extension points

- **New hand type** (e.g. "Straight"): subclass `HandRule`, set a rank between existing ranks or renumber, register in the rule list passed to `HandEvaluator`.
- **Wild cards / jokers**: override `Hand.count_values()` to fold joker counts into the most-frequent card before returning.
- **Different card ordering** (e.g. Ace-low): pass a different `_card_order` string to `Game`.
- **N-card hands**: `HandRule.matches` operates on `count_values()` only; resize `Hand.cards` without touching rules.

---

## Part 2: Incomplete Hands

Hands may arrive with fewer than 5 cards due to packet loss. The missing cards are unknown, so we cannot evaluate a single strength — instead we evaluate a **range** `[min_strength, max_strength]` for each hand.

Decision rule:
- If `min(A) > max(B)` → A wins definitively (A beats B even in the worst case for A vs best case for B)
- If `min(B) > max(A)` → B wins definitively
- Otherwise → `UNKNOWN`

### `CompletionStrategy` (abstract)

Encapsulates how to fill the missing cards. Added cards are appended to `hand.cards` in the remaining dealt positions (positions `k..4`).

```
CompletionStrategy (ABC)
  complete(hand: Hand, card_order: str) -> Hand
    # returns a new 5-card Hand with missing slots filled
```

`card_order` is passed in so strategies know what "highest" and "lowest" face values are.

---

### `MaxCompletionStrategy`

Goal: produce the strongest possible complete hand.

**Type rank step** — fill remaining slots with copies of the highest-frequency existing card. This concentrates counts and maximises the type rank.
- If the hand is empty, default to the highest card in `card_order` (e.g. `"A"`).
- If multiple cards are tied for highest frequency, prefer the one with the higher card value (better tiebreaker for the appended positions).

**Tiebreaker step** — since added cards land at the end of the dealt order (positions `k..4`), use the highest-value card for those slots. This only differs from the type-rank step when the optimal fill card for type rank is not the highest-value card.

Example:
```
[9, K]  +3 slots → concentrate on K (higher value, same count as 9)
         → [9, K, K, K, K] = four-of-a-kind, rank 6
```

---

### `MinCompletionStrategy`

Goal: produce the weakest possible complete hand.

**Type rank step** — fill remaining slots with distinct cards not yet present in the hand, picked lowest-value first. This spreads the count histogram and minimises the achievable type rank.
- Because a 5-card hand uses at most 5 of the 13 available labels, there are always enough unseen labels to spread into.
- Already-existing pairs/trips cannot be un-made; the min is floored by what's already dealt.

**Tiebreaker step** — use the lowest-value cards for the appended positions to produce the weakest tiebreaker.

Example:
```
[9, 9, K]  +2 slots → add two unseen cards, lowest first: 2, 3
            → [9, 9, K, 2, 3] = one pair, rank 2  (existing pair cannot be avoided)
[A, K, Q]  +2 slots → add two unseen cards: 2, 3
            → [A, K, Q, 2, 3] = high card, rank 1
```

---

### `RangeComparator`

Top-level coordinator for incomplete hands. Reuses `Game.compare` to compare the four boundary hands.

```
RangeComparator
  _game: Game
  _max_strategy: MaxCompletionStrategy
  _min_strategy: MinCompletionStrategy

  compare(h1: Hand, h2: Hand) -> str   # "player1" | "player2" | "UNKNOWN"
```

Algorithm:
```
r1 = 5 - len(h1.cards)
r2 = 5 - len(h2.cards)

min1 = _min_strategy.complete(h1, card_order)
max1 = _max_strategy.complete(h1, card_order)
min2 = _min_strategy.complete(h2, card_order)
max2 = _max_strategy.complete(h2, card_order)

# A wins definitively: A's worst still beats B's best
if _game.compare(min1, max2) == "player1":  return "player1"

# B wins definitively: B's worst still beats A's best
if _game.compare(max1, min2) == "player2":  return "player2"

return "UNKNOWN"
```

---

### Worked example

`[9,9,9,9]` vs `[9]`:

| Hand  | Min completion         | Type   | Max completion         | Type            |
|-------|------------------------|--------|------------------------|-----------------|
| A     | `[9,9,9,9,2]`         | rank 6 | `[9,9,9,9,9]`         | rank 7          |
| B     | `[9,2,3,4,5]`         | rank 1 | `[9,9,9,9,9]`         | rank 7          |

- `compare(min_A, max_B)` = `compare([9,9,9,9,2], [9,9,9,9,9])` → rank 6 vs 7 → "player2" ≠ "player1" → no
- `compare(max_A, min_B)` = `compare([9,9,9,9,9], [9,2,3,4,5])` → rank 7 vs 1 → "player1" ≠ "player2" → no
- Result: **UNKNOWN** ✓ (B could complete to five-of-a-kind and win)

---

### Integration with Part 1

`RangeComparator` composes over `Game` — it does not replace it. Complete (5-card) hands can still be compared with `Game.compare` directly. `RangeComparator` delegates all actual comparisons to `Game.compare`; it only adds the range logic on top.

---

## Part 3: Extensibility API

Hand-type rankings are no longer hardcoded. The caller registers types at runtime and controls the evaluation order per call.

### API surface

```
add_type(type_id: str, matches_fn: Callable[[Hand], bool]) -> None
evaluate(hand1: Hand, hand2: Hand, evaluation_order: list[str]) -> str
```

**`add_type`** — registers a custom hand type under `type_id`. If `type_id` already exists, the new function overwrites it (last write wins). This allows callers to replace built-in types or add entirely new ones without touching any existing class.

**`evaluate`** — compares two hands using a caller-supplied ranking hierarchy:
- `evaluation_order` is **weakest-to-strongest**: `evaluation_order[0]` = rank 1, `evaluation_order[-1]` = highest rank.
- A hand's score = `index + 1` of the first type in the list (scanning strongest-first) whose `matches_fn` returns `True`. Unmatched hand scores 0.
- Higher score wins; equal scores fall back to card-by-card tiebreaking (same left-to-right logic as `Game.compare`).
- Unknown `type_id` in `evaluation_order` raises `ValueError` — fail fast rather than silently skip.

---

### `ExtensibleGame`

Wraps `Game` for card ordering and tiebreaking. Does not inherit from or replace `Game`.

```
ExtensibleGame
  _game: Game                                      # card ordering + tiebreak logic
  _type_registry: dict[str, Callable[[Hand], bool]]  # type_id → matches_fn

  add_type(type_id: str, matches_fn: Callable[[Hand], bool]) -> None
    registry[type_id] = matches_fn

  evaluate(hand1: Hand, hand2: Hand, evaluation_order: list[str]) -> str
    score1 = _score(hand1, evaluation_order)
    score2 = _score(hand2, evaluation_order)
    if score1 != score2: return "player1" if score1 > score2 else "player2"
    card-by-card tiebreak via _game._card_value  →  "player1" | "player2" | "tie"

  _score(hand: Hand, evaluation_order: list[str]) -> int
    iterate evaluation_order in reverse (strongest → weakest)
      raise ValueError if type_id not in registry
      if registry[type_id](hand): return index + 1
    return 0   # no type matched
```

---

### Pre-registering built-in types

At construction, `ExtensibleGame` auto-registers every rule from `DEFAULT_RULES` using each rule's `name` as the `type_id`:

```
for rule in DEFAULT_RULES:
    add_type(rule.name, rule.matches)
```

Built-in `type_id` strings: `"High Card"`, `"One Pair"`, `"Two Pair"`, `"Three of a Kind"`, `"Full House"`, `"Four of a Kind"`, `"Five of a Kind"`.

A caller can overwrite any of these (e.g. redefine `"Five of a Kind"` to include joker rules) or skip them entirely by not listing them in `evaluation_order`.

---

### Worked example

Register a custom `"Flush"` type and place it between Full House and Four of a Kind:

```python
game = ExtensibleGame(Game(HandEvaluator(DEFAULT_RULES)))

def is_flush(hand):
    suits = [c[-1] for c in hand.cards]   # assumes "9H", "KD" etc.
    return len(set(suits)) == 1

game.add_type("Flush", is_flush)

order = [
    "High Card", "One Pair", "Two Pair",
    "Three of a Kind", "Full House",
    "Flush",               # inserted between Full House (5) and Four of a Kind (6)
    "Four of a Kind", "Five of a Kind",
]

result = game.evaluate(hand1, hand2, order)
```

Rank table under this custom order:

| type_id          | index | score |
|------------------|-------|-------|
| High Card        | 0     | 1     |
| One Pair         | 1     | 2     |
| Two Pair         | 2     | 3     |
| Three of a Kind  | 3     | 4     |
| Full House       | 4     | 5     |
| Flush            | 5     | 6     |
| Four of a Kind   | 6     | 7     |
| Five of a Kind   | 7     | 8     |

No existing class was modified — `Flush` was added via `add_type` and ranked by its position in `evaluation_order`.

---

### Design decisions

| Decision | Choice | Reason |
|----------|--------|--------|
| `evaluation_order` direction | weakest-to-strongest (index 0 = lowest rank) | Mirrors the natural "rank ladder" mental model; rank = index + 1 is unambiguous |
| Duplicate `type_id` in `add_type` | overwrite | Simplest contract; caller controls the registry |
| Unknown `type_id` in `evaluate` | raise `ValueError` | Silent skip would hide typos and produce wrong results |
| Tiebreaking | delegated to `_game._card_value` | Reuses the card-ordering already owned by `Game`; no duplication |
| `ExtensibleGame` vs modifying `Game` | new class | Keeps `Game` stable; `ExtensibleGame` is an opt-in layer |
