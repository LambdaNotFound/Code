from ast import List


class Solution:
    """
    140. Word Break II

    Given a string s and a dictionary of strings wordDict, add spaces in s to construct a sentence where each word is a valid dictionary word. Return all such possible sentences in any order.

    Note that the same word in the dictionary may be reused multiple times in the segmentation.

    Example 1:

    Input: s = "catsanddog", wordDict = ["cat","cats","and","sand","dog"]
    Output: ["cats and dog","cat sand dog"]
    Example 2:

    Input: s = "pineapplepenapple", wordDict = ["apple","pen","applepen","pine","pineapple"]
    Output: ["pine apple pen apple","pineapple pen apple","pine applepen apple"]
    Explanation: Note that you are allowed to reuse a dictionary word.
    Example 3:

    Input: s = "catsandog", wordDict = ["cats","dog","sand","and","cat"]
    Output: []
    """

    def wordBreak(self, s: str, wordDict: List[str]) -> List[str]:
        word_set = set(wordDict)

        n = len(s)
        dp = [False] * (n + 1)
        table = [[] for _ in range(n + 1)]

        for i in range(1, n + 1):
            substr = s[0:i]
            if substr in word_set:
                dp[i] = True
                table[i].append(substr)

            for j in range(i):
                substr = s[j:i]
                if substr in word_set and dp[j]:
                    dp[i] = True
                    for sentence in table[j]:
                        table[i].append(sentence + " " + substr)

        return table[n]
