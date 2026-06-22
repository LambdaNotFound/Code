# Expense Rule Engine — Design Plan

## Data models

```python
@dataclass
class Expense:
    expense_id: str
    trip_id: str
    amount_usd: float
    expense_type: str       # e.g. "client_hosting", "airfare", "entertainment"
    vendor_type: str        # e.g. "restaurant"
    vendor_name: str
    violating_rules: list[str] = field(default_factory=list)

@dataclass
class Trip:
    trip_id: str
    expenses: list[Expense]
    violating_rules: list[str] = field(default_factory=list)
```

`Trip` is derived — not passed in, built from the expense list by grouping on `trip_id`.

---

## Return type

`evaluate_rules` returns expenses and trips with `violating_rules` populated.

```python
@dataclass
class EvaluationResult:
    expenses: list[Expense]
    trips: list[Trip]
```

Per-expense rules → populate `expense.violating_rules`  
Trip aggregate rules → populate `trip.violating_rules` only (individual expenses are not flagged)

---

## Strategy pattern — Rule hierarchy

Each `Rule` subclass is a **strategy**: it encapsulates one evaluation algorithm behind a uniform interface. `ExpenseRuleEngine` calls the interface without knowing the concrete type.

```
Rule (abstract strategy)
├── FieldRule          — checks a single field on one expense → returns expense:* IDs
├── TripAggregateRule  — checks a sum over a trip → returns trip:* IDs
└── CompositeRule      — combines sub-rules with AND / OR
```

```python
class Rule(ABC):
    rule_id: str

    @abstractmethod
    def violating_entity_ids(self, expenses: list[Expense]) -> set[str]:
        """Return prefixed IDs of entities that violate this rule.
        Prefix convention: 'expense:<id>' or 'trip:<id>'
        """
```

The prefix makes each returned ID self-describing — the engine routes violations to the right entity without inspecting the rule type.

### `FieldRule` — per-expense strategy

Returns `expense:<id>` for each expense that fails the field check.

```python
@dataclass
class FieldRule(Rule):
    field: str               # "amount_usd" | "expense_type" | "vendor_type"
    operator: str            # "gt" | "eq"
    threshold: float | str

    def violating_entity_ids(self, expenses):
        return {f"expense:{exp.expense_id}" for exp in expenses if self._check(exp)}
```

Example: `FieldRule("no_airfare", field="expense_type", operator="eq", threshold="airfare")`

### `TripAggregateRule` — per-trip strategy

Groups expenses by `trip_id`, sums `amount_usd`, returns `trip:<id>` for each trip that exceeds the cap. Individual expenses are **not** flagged.

```python
@dataclass
class TripAggregateRule(Rule):
    cap: float
    filter_field: str | None = None   # e.g. "vendor_type" to cap meals only
    filter_value: str | None = None   # e.g. "restaurant"

    def violating_entity_ids(self, expenses):
        trips = defaultdict(list)
        for exp in expenses:
            if self.filter_field is None or getattr(exp, self.filter_field) == self.filter_value:
                trips[exp.trip_id].append(exp)
        return {
            f"trip:{trip_id}"
            for trip_id, exps in trips.items()
            if sum(e.amount_usd for e in exps) > self.cap
        }
```

### `CompositeRule` — structural strategy

Combines sub-rules with boolean logic. Enables expressing "restaurant AND amount > 75" without adding filter hacks to `FieldRule`. Sub-rules must return the same entity prefix for intersection/union to be meaningful.

```python
@dataclass
class CompositeRule(Rule):
    rules: list[Rule]
    operator: str = "AND"   # "AND" → intersection | "OR" → union

    def violating_entity_ids(self, expenses):
        sets = [r.violating_entity_ids(expenses) for r in self.rules]
        if self.operator == "AND":
            return set.intersection(*sets)
        return set.union(*sets)
```

---

## Factory pattern — `RuleFactory`

Decouples rule creation from rule types. New rule types register themselves; API-driven creation passes a type string + kwargs.

```python
class RuleFactory:
    _registry: dict[str, type[Rule]] = {}

    @classmethod
    def register(cls, rule_type: str, rule_cls: type[Rule]):
        cls._registry[rule_type] = rule_cls

    @classmethod
    def create(cls, rule_type: str, **kwargs) -> Rule:
        rule_cls = cls._registry.get(rule_type)
        if not rule_cls:
            raise ValueError(f"Unknown rule type: {rule_type}")
        return rule_cls(**kwargs)

RuleFactory.register("field", FieldRule)
RuleFactory.register("trip_aggregate", TripAggregateRule)
RuleFactory.register("composite", CompositeRule)
```

---

## `ExpenseRuleEngine`

```python
class ExpenseRuleEngine:
    def __init__(self, rules: list[Rule]):
        self.rules = rules

    def evaluate(self, expenses: list[Expense]) -> EvaluationResult:
        trips: dict[str, Trip] = {}
        for exp in expenses:
            if exp.trip_id not in trips:
                trips[exp.trip_id] = Trip(trip_id=exp.trip_id, expenses=[])
            trips[exp.trip_id].expenses.append(exp)

        expense_index = {exp.expense_id: exp for exp in expenses}

        for rule in self.rules:
            for entity_id in rule.violating_entity_ids(expenses):
                prefix, id_ = entity_id.split(":", 1)
                if prefix == "expense" and id_ in expense_index:
                    expense_index[id_].violating_rules.append(rule.rule_id)
                elif prefix == "trip" and id_ in trips:
                    trips[id_].violating_rules.append(rule.rule_id)

        return EvaluationResult(expenses=expenses, trips=list(trips.values()))
```

T: O(R × N) where R = number of rules, N = number of expenses  
S: O(N) for trip grouping and expense index

---

## Concrete rules from the prompt

```python
rules = [
    # restaurant cap: AND composite of two FieldRules
    CompositeRule("restaurant_cap", operator="AND", rules=[
        FieldRule("_vendor_restaurant", field="vendor_type", operator="eq", threshold="restaurant"),
        FieldRule("_amount_75",         field="amount_usd",  operator="gt", threshold=75),
    ]),
    FieldRule("no_airfare",       field="expense_type", operator="eq", threshold="airfare"),
    FieldRule("no_entertainment", field="expense_type", operator="eq", threshold="entertainment"),
    FieldRule("expense_cap_250",  field="amount_usd",  operator="gt", threshold=250),
    TripAggregateRule("trip_cap_2000", cap=2000),
    TripAggregateRule("meal_cap_200",  cap=200,
                      filter_field="vendor_type", filter_value="restaurant"),
]
```

---

## Extensibility

| Need | How |
|---|---|
| New operator (e.g. `gte`, `contains`) | Extend `FieldRule.operator` |
| New aggregate (e.g. per-employee monthly cap) | New `Rule` subclass + `RuleFactory.register` |
| API-driven rule creation | `RuleFactory.create(rule_type, **json_payload)` |
| Complex logic (NOT, nested AND/OR) | Nest `CompositeRule` instances |
