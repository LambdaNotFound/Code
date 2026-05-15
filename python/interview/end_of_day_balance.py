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
        for k, v in balances.items():
            result.append(v)
        return result 
