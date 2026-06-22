import sys, os
sys.path.insert(0, os.path.dirname(__file__))
from rippling_expense_rule_engine import (
    Expense, Trip, EvaluationResult,
    FieldRule, TripAggregateRule, CompositeRule,
    RuleFactory, ExpenseRuleEngine, DEFAULT_RULES,
)

import pytest


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def make_expense(expense_id, trip_id, amount_usd, expense_type="meal",
                 vendor_type="cafe", vendor_name="Cafe"):
    return Expense(expense_id=expense_id, trip_id=trip_id, amount_usd=amount_usd,
                   expense_type=expense_type, vendor_type=vendor_type, vendor_name=vendor_name)


def trip_by_id(result: EvaluationResult, trip_id: str) -> Trip:
    return next(t for t in result.trips if t.trip_id == trip_id)


def expense_by_id(result: EvaluationResult, expense_id: str) -> Expense:
    return next(e for e in result.expenses if e.expense_id == expense_id)


# ---------------------------------------------------------------------------
# FieldRule
# ---------------------------------------------------------------------------

class TestFieldRule:
    def test_eq_flags_matching_expense(self):
        rule = FieldRule("no_airfare", field="expense_type", operator="eq", threshold="airfare")
        expenses = [
            make_expense("e1", "t1", 100, expense_type="airfare"),
            make_expense("e2", "t1", 100, expense_type="meal"),
        ]
        assert rule.violating_entity_ids(expenses) == {"expense:e1"}

    def test_gt_flags_over_threshold(self):
        rule = FieldRule("cap_250", field="amount_usd", operator="gt", threshold=250)
        expenses = [
            make_expense("e1", "t1", 300),
            make_expense("e2", "t1", 250),  # equal, not over
            make_expense("e3", "t1", 100),
        ]
        assert rule.violating_entity_ids(expenses) == {"expense:e1"}

    def test_no_violations_returns_empty(self):
        rule = FieldRule("no_entertainment", field="expense_type", operator="eq", threshold="entertainment")
        expenses = [make_expense("e1", "t1", 50, expense_type="meal")]
        assert rule.violating_entity_ids(expenses) == set()


# ---------------------------------------------------------------------------
# TripAggregateRule
# ---------------------------------------------------------------------------

class TestTripAggregateRule:
    def test_trip_total_exceeds_cap(self):
        rule = TripAggregateRule("trip_cap_2000", cap=2000)
        expenses = [
            make_expense("e1", "t1", 1000),
            make_expense("e2", "t1", 1001),  # t1 total: 2001 → violated
            make_expense("e3", "t2", 500),   # t2 total: 500 → ok
        ]
        assert rule.violating_entity_ids(expenses) == {"trip:t1"}

    def test_trip_exactly_at_cap_not_flagged(self):
        rule = TripAggregateRule("trip_cap_2000", cap=2000)
        expenses = [make_expense("e1", "t1", 1000), make_expense("e2", "t1", 1000)]
        assert rule.violating_entity_ids(expenses) == set()

    def test_filtered_aggregate_meal_cap(self):
        rule = TripAggregateRule("meal_cap_200", cap=200,
                                 filter_field="vendor_type", filter_value="restaurant")
        expenses = [
            make_expense("e1", "t1", 120, vendor_type="restaurant"),
            make_expense("e2", "t1", 100, vendor_type="restaurant"),  # total meals: 220 → violated
            make_expense("e3", "t1", 500, vendor_type="cafe"),        # not counted toward meal cap
        ]
        assert rule.violating_entity_ids(expenses) == {"trip:t1"}

    def test_trip_violations_do_not_flag_individual_expenses(self):
        rule = TripAggregateRule("trip_cap_2000", cap=2000)
        expenses = [make_expense("e1", "t1", 1500), make_expense("e2", "t1", 600)]
        engine = ExpenseRuleEngine([rule])
        result = engine.evaluate(expenses)
        assert expense_by_id(result, "e1").violating_rules == []
        assert expense_by_id(result, "e2").violating_rules == []
        assert "trip_cap_2000" in trip_by_id(result, "t1").violating_rules

    def test_multiple_trips_only_violating_flagged(self):
        rule = TripAggregateRule("trip_cap_2000", cap=2000)
        expenses = [
            make_expense("e1", "t1", 2500),  # t1 violated
            make_expense("e2", "t2", 500),   # t2 ok
        ]
        assert rule.violating_entity_ids(expenses) == {"trip:t1"}


# ---------------------------------------------------------------------------
# CompositeRule
# ---------------------------------------------------------------------------

class TestCompositeRule:
    def test_and_requires_both_conditions(self):
        rule = CompositeRule("restaurant_cap", operator="AND", rules=[
            FieldRule("_vendor", field="vendor_type", operator="eq", threshold="restaurant"),
            FieldRule("_amount", field="amount_usd",  operator="gt", threshold=75),
        ])
        expenses = [
            make_expense("e1", "t1", 80,  vendor_type="restaurant"),  # both → flagged
            make_expense("e2", "t1", 50,  vendor_type="restaurant"),  # amount ok → not flagged
            make_expense("e3", "t1", 100, vendor_type="cafe"),         # not restaurant → not flagged
        ]
        assert rule.violating_entity_ids(expenses) == {"expense:e1"}

    def test_or_flags_either_condition(self):
        rule = CompositeRule("blocked_types", operator="OR", rules=[
            FieldRule("_airfare",       field="expense_type", operator="eq", threshold="airfare"),
            FieldRule("_entertainment", field="expense_type", operator="eq", threshold="entertainment"),
        ])
        expenses = [
            make_expense("e1", "t1", 100, expense_type="airfare"),
            make_expense("e2", "t1", 100, expense_type="entertainment"),
            make_expense("e3", "t1", 100, expense_type="meal"),
        ]
        assert rule.violating_entity_ids(expenses) == {"expense:e1", "expense:e2"}


# ---------------------------------------------------------------------------
# RuleFactory
# ---------------------------------------------------------------------------

class TestRuleFactory:
    def test_creates_field_rule(self):
        rule = RuleFactory.create("field", rule_id="r1", field="expense_type",
                                  operator="eq", threshold="airfare")
        assert isinstance(rule, FieldRule)
        assert rule.rule_id == "r1"

    def test_creates_trip_aggregate_rule(self):
        rule = RuleFactory.create("trip_aggregate", rule_id="r2", cap=2000)
        assert isinstance(rule, TripAggregateRule)
        assert rule.cap == 2000

    def test_creates_composite_rule(self):
        sub = FieldRule("sub", field="expense_type", operator="eq", threshold="airfare")
        rule = RuleFactory.create("composite", rule_id="r3", rules=[sub], operator="OR")
        assert isinstance(rule, CompositeRule)

    def test_unknown_type_raises(self):
        with pytest.raises(ValueError, match="Unknown rule type"):
            RuleFactory.create("nonexistent", rule_id="r4")


# ---------------------------------------------------------------------------
# ExpenseRuleEngine — integration
# ---------------------------------------------------------------------------

class TestExpenseRuleEngine:
    def test_per_expense_violation_attached_to_expense(self):
        rule = FieldRule("no_airfare", field="expense_type", operator="eq", threshold="airfare")
        expenses = [make_expense("e1", "t1", 100, expense_type="airfare")]
        result = ExpenseRuleEngine([rule]).evaluate(expenses)
        assert expense_by_id(result, "e1").violating_rules == ["no_airfare"]

    def test_expense_with_multiple_violations(self):
        rules = [
            FieldRule("no_airfare",      field="expense_type", operator="eq", threshold="airfare"),
            FieldRule("expense_cap_250", field="amount_usd",   operator="gt", threshold=250),
        ]
        expenses = [make_expense("e1", "t1", 300, expense_type="airfare")]
        result = ExpenseRuleEngine(rules).evaluate(expenses)
        assert set(expense_by_id(result, "e1").violating_rules) == {"no_airfare", "expense_cap_250"}

    def test_clean_expense_has_no_violations(self):
        expenses = [make_expense("e1", "t1", 50, expense_type="meal", vendor_type="cafe")]
        result = ExpenseRuleEngine(DEFAULT_RULES).evaluate(expenses)
        assert expense_by_id(result, "e1").violating_rules == []
        assert trip_by_id(result, "t1").violating_rules == []

    def test_default_rules_integration(self):
        expenses = [
            make_expense("e1", "t1", 80,   expense_type="client_hosting", vendor_type="restaurant"),
            make_expense("e2", "t1", 400,   expense_type="meal",           vendor_type="cafe"),
            make_expense("e3", "t2", 100,   expense_type="airfare"),
            make_expense("e4", "t2", 50,    expense_type="entertainment"),
            make_expense("e5", "t3", 120,   expense_type="meal",           vendor_type="restaurant"),
            make_expense("e6", "t3", 100,   expense_type="meal",           vendor_type="restaurant"),
        ]
        result = ExpenseRuleEngine(DEFAULT_RULES).evaluate(expenses)

        # e1: restaurant AND amount > 75 → restaurant_cap
        assert "restaurant_cap" in expense_by_id(result, "e1").violating_rules
        # e2: amount > 250 → expense_cap_250
        assert "expense_cap_250" in expense_by_id(result, "e2").violating_rules
        # e3: airfare → no_airfare
        assert "no_airfare" in expense_by_id(result, "e3").violating_rules
        # e4: entertainment → no_entertainment
        assert "no_entertainment" in expense_by_id(result, "e4").violating_rules
        # t1: 80 + 400 = 480, under 2000 → no trip violation
        assert trip_by_id(result, "t1").violating_rules == []
        # t3: meals 120 + 100 = 220 > 200 → meal_cap_200 on trip, not on expenses
        assert "meal_cap_200" in trip_by_id(result, "t3").violating_rules
        assert expense_by_id(result, "e5").violating_rules == []
        assert expense_by_id(result, "e6").violating_rules == []
