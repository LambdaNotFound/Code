"""
We plan to offer a corporate credit card to companies that employees can use for business expenses.
Managers can set policies on the cards so that employees do not misuse the card or exceed allowances.
We're going to be building the rules engine that supports this product.
"""

from abc import ABC, abstractmethod
from collections import defaultdict
from dataclasses import dataclass, field


# ---------------------------------------------------------------------------
# Data models
# ---------------------------------------------------------------------------

@dataclass
class Expense:
    expense_id: str
    trip_id: str
    amount_usd: float
    expense_type: str
    vendor_type: str
    vendor_name: str
    violating_rules: list[str] = field(default_factory=list)


@dataclass
class Trip:
    trip_id: str
    expenses: list[Expense]
    violating_rules: list[str] = field(default_factory=list)


@dataclass
class EvaluationResult:
    expenses: list[Expense]
    trips: list[Trip]


# ---------------------------------------------------------------------------
# Strategy pattern — Rule hierarchy
# ---------------------------------------------------------------------------

class Rule(ABC):
    def __init__(self, rule_id: str):
        self.rule_id = rule_id

    @abstractmethod
    def violating_entity_ids(self, expenses: list[Expense]) -> set[str]:
        """Return prefixed IDs of entities that violate this rule.
        Prefix: 'expense:<id>' or 'trip:<id>'
        """


class FieldRule(Rule):
    """Per-expense rule: checks a single field against a threshold."""

    def __init__(self, rule_id: str, field: str, operator: str, threshold: float | str):
        super().__init__(rule_id)
        self.field = field
        self.operator = operator
        self.threshold = threshold

    def _check(self, expense: Expense) -> bool:
        value = getattr(expense, self.field)
        if self.operator == "gt":
            return float(value) > float(self.threshold)
        if self.operator == "eq":
            return value == self.threshold
        if self.operator == "neq":
            return value != self.threshold
        raise ValueError(f"Unknown operator: {self.operator}")

    def violating_entity_ids(self, expenses: list[Expense]) -> set[str]:
        return {f"expense:{exp.expense_id}" for exp in expenses if self._check(exp)}


class TripAggregateRule(Rule):
    """Per-trip rule: flags the trip when the sum of (filtered) expenses exceeds cap."""

    def __init__(self, rule_id: str, cap: float,
                 filter_field: str | None = None, filter_value: str | None = None):
        super().__init__(rule_id)
        self.cap = cap
        self.filter_field = filter_field
        self.filter_value = filter_value

    def violating_entity_ids(self, expenses: list[Expense]) -> set[str]:
        trips: dict[str, list[Expense]] = defaultdict(list)
        for exp in expenses:
            if self.filter_field is None or getattr(exp, self.filter_field) == self.filter_value:
                trips[exp.trip_id].append(exp)
        return {
            f"trip:{trip_id}"
            for trip_id, exps in trips.items()
            if sum(e.amount_usd for e in exps) > self.cap
        }


class CompositeRule(Rule):
    """Combines sub-rules with AND (intersection) or OR (union)."""

    def __init__(self, rule_id: str, rules: list[Rule], operator: str = "AND"):
        super().__init__(rule_id)
        self.rules = rules
        self.operator = operator

    def violating_entity_ids(self, expenses: list[Expense]) -> set[str]:
        sets = [r.violating_entity_ids(expenses) for r in self.rules]
        if self.operator == "AND":
            return set.intersection(*sets)
        return set.union(*sets)


# ---------------------------------------------------------------------------
# Factory pattern
# ---------------------------------------------------------------------------

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


# ---------------------------------------------------------------------------
# Engine
# ---------------------------------------------------------------------------

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


# ---------------------------------------------------------------------------
# Default rules from the prompt
# ---------------------------------------------------------------------------

DEFAULT_RULES: list[Rule] = [
    CompositeRule("restaurant_cap", operator="AND", rules=[
        FieldRule("_vendor_restaurant", field="vendor_type",   operator="eq",  threshold="restaurant"),
        FieldRule("_amount_75",         field="amount_usd",    operator="gt",  threshold=75),
        FieldRule("_not_meal",          field="expense_type",  operator="neq", threshold="meal"),
    ]),
    FieldRule("no_airfare",       field="expense_type", operator="eq", threshold="airfare"),
    FieldRule("no_entertainment", field="expense_type", operator="eq", threshold="entertainment"),
    FieldRule("expense_cap_250",  field="amount_usd",   operator="gt", threshold=250),
    TripAggregateRule("trip_cap_2000", cap=2000),
    TripAggregateRule("meal_cap_200",  cap=200,
                      filter_field="vendor_type", filter_value="restaurant"),
]
