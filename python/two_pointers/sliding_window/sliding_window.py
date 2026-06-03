from collections import Counter


class Solution:
    """
    76. Minimum Window Substring
    """

    def minWindow(self, s: str, t: str) -> str:
        if not s or not t:
            return ""

        freq = Counter(t)  # how many of each char the window still needs
        needed = len(t)  # total chars still needed, duplicates included
        left = 0
        min_len, result = float("inf"), ""

        for right in range(len(s)):
            char = s[right]
            freq[char] -= 1
            # freq[char] >= 0 means this was a genuinely-needed copy (not surplus),
            # so it closes part of the deficit.
            if freq[char] >= 0:
                needed -= 1

            # Window now covers all of t — shrink from the left to minimize it.
            while needed == 0:
                # Skip leading surplus (freq < 0 means we hold more than required).
                while freq[s[left]] < 0:
                    freq[s[left]] += 1
                    left += 1

                if right - left + 1 < min_len:
                    min_len = right - left + 1
                    result = s[left : right + 1]

                # Give back the leftmost required char, reopening the deficit.
                freq[s[left]] += 1
                needed += 1
                left += 1

        return result
