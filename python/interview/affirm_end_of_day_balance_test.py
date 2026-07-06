import pytest

from interview.affirm_end_of_day_balance import (
    Solution,  # type: ignore[import-not-found]
)


@pytest.mark.parametrize(
    "transactions,initial_balance,expected",
    [
        (
            [
                ["Alice", "Bob", "50"],
                ["Bob", "Charlie", "30"],
                ["Charlie", "Alice", "20"],
                ["Alice", "David", "70"],
            ],
            [["Alice", "100"], ["Bob", "50"], ["Charlie", "75"], ["David", "25"]],
            [0, 70, 85, 95],
        ),
        (
            [
                ["Alice", "Bob", "300"],
                ["Charlie", "David", "400"],
                ["Eve", "Frank", "200"],
                ["George", "Hank", "100"],
                ["Alice", "Charlie", "150"],
                ["David", "Eve", "300"],
                ["Frank", "George", "250"],
                ["Bob", "Hank", "100"],
                ["Charlie", "Alice", "200"],
                ["Eve", "David", "150"],
            ],
            [
                ["Alice", "1000"],
                ["Bob", "600"],
                ["Charlie", "500"],
                ["David", "700"],
                ["Eve", "600"],
                ["Frank", "400"],
                ["George", "200"],
                ["Hank", "100"],
                ["Ivy", "300"],
                ["Jack", "400"],
            ],
            [750, 800, 50, 950, 550, 350, 350, 300, 300, 400],
        ),
        (
            [
                ["Alice", "Bob", "150"],
                ["Charlie", "Alice", "200"],
                ["David", "Bob", "50"],
                ["Eve", "Charlie", "100"],
                ["Bob", "Eve", "300"],
                ["David", "Alice", "50"],
                ["Alice", "Eve", "100"],
                ["Charlie", "Bob", "150"],
            ],
            [
                ["Alice", "500"],
                ["Bob", "300"],
                ["Charlie", "200"],
                ["David", "150"],
                ["Eve", "250"],
            ],
            [500, 350, -50, 50, 550],
        ),
    ],
)
def test_get_end_of_day_balance(transactions, initial_balance, expected):
    assert Solution().getEndOfDayBalance(transactions, initial_balance) == expected


# ── 465. Optimal Account Balancing ────────────────────────────────────────────


@pytest.mark.parametrize(
    "transactions, want",
    [
        # no transactions → nothing to settle
        ([], 0),
        # one net creditor, two debtors — cannot merge
        ([[0, 1, 10], [2, 0, 5]], 2),
        # perfect cancellation: net balances are [+4, -4] → one transfer
        ([[0, 1, 10], [1, 0, 1], [1, 2, 5], [2, 0, 5]], 1),
        # circular: every net balance is 0
        ([[0, 1, 1], [1, 2, 1], [2, 0, 1]], 0),
        # single-direction debt: 0 owes 1
        ([[0, 1, 5]], 1),
        # one creditor, two debtors with different amounts
        ([[1, 3, 15], [2, 3, 10], [3, 1, 2]], 2),
        # two independent pairs settle independently
        ([[0, 1, 5], [2, 3, 7]], 2),
    ],
)
def test_min_transfers(transactions, want):
    assert Solution().minTransfers(transactions) == want
