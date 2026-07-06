from typing import List, Optional

"""
End Of Day Balance

Input:
transactions = [["Alice","Bob","50"],["Bob","Charlie","30"],["Charlie","Alice","20"],["Alice","David","70"]]

initialBalance = [["Alice","100"],["Bob","50"],["Charlie","75"],["David","25"]]

Output: [0, 70, 85, 95]

"""

class Solution:
    def getEndOfDayBalance(self, transactions: List[List[str]], initialBalance: List[List[str]]) -> List[int]:
        balances: dict[str, int] = {}
        for _, list in enumerate(initialBalance):
            name, balance = list[0], int(list[1])
            balances[name] = balances.get(name, 0) + balance        
        for _, txn in enumerate(transactions):
            src, dst, amount = txn[0], txn[1], int(txn[2])
            balances[src] -= amount
            balances[dst] += amount

        result: List[int] = []
        for _, v in balances.items():
            result.append(v)
        return result

    """
    465. Optimal Account Balancing
    """
    # minimum transfers = (number of non-zero balances) - (maximum perfect cancellations)
    def minTransfers(self, transactions: List[List[int]]) -> int:
        balances: dict[int, int] = {}
        for src, dst, amount in transactions:
            balances[src] = balances.get(src, 0) - amount
            balances[dst] = balances.get(dst, 0) + amount

        debts: list[int] = [v for v in balances.values() if v != 0]

        def dfs(start: int) -> int:
            # skip already-settled balances
            while start < len(debts) and debts[start] == 0:
                start += 1

            # base case: all balances settled
            if start == len(debts):
                return 0

            min_tx = float('inf')
            for j in range(start + 1, len(debts)):
                # pruning: skip same sign
                if debts[start] * debts[j] > 0:
                    continue
                # apply
                debts[j] += debts[start]
                min_tx = min(min_tx, 1 + dfs(start + 1))
                # undo
                debts[j] -= debts[start]

            return min_tx

        return dfs(0)
        
