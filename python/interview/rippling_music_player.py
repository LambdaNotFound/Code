"""
Design an in-memory music player for tracking songs and play times,
and ranking popular songs using OO design principles.
"""

import heapq
from collections import deque
from dataclasses import dataclass, field


# ---------------------------------------------------------------------------
# Data models
# ---------------------------------------------------------------------------

@dataclass
class Song:
    id: int
    title: str
    play_count: int = 0
    played_by: set = field(default_factory=set)

    @property
    def unique_listener_count(self) -> int:
        return len(self.played_by)


@dataclass
class User:
    id: int
    name: str = ""
    last_three_played_songs: deque = field(default_factory=lambda: deque(maxlen=3))  # index 0 = most recent


# ---------------------------------------------------------------------------
# MusicPlayer
# ---------------------------------------------------------------------------

class MusicPlayer:
    def __init__(self):
        self._songs: dict[int, Song] = {}   # song_id -> Song
        self._users: dict[int, User] = {}   # user_id -> User

    def add_song(self, song_title: str) -> int:
        # Time: O(1)  Space: O(1)
        song_id = len(self._songs) + 1
        self._songs[song_id] = Song(id=song_id, title=song_title)
        return song_id

    def play_song(self, song_id: int, user_id: int) -> None:
        # Time: O(1) — all list operations are bounded by the fixed cap of 3
        # Space: O(U) per song for played_by, O(1) per user for the 3-element queue
        song = self._songs[song_id]
        song.play_count += 1
        song.played_by.add(user_id)

        if user_id not in self._users:
            self._users[user_id] = User(id=user_id)

        q = self._users[user_id].last_three_played_songs
        if song_id in q:
            q.remove(song_id)           # deduplicate: pull out existing entry
        q.appendleft(song_id)           # most recent goes to the front; maxlen=3 auto-evicts the tail

    def last_three_played_song_titles(self, user_id: int) -> list[str]:
        # Time: O(1) — at most 3 lookups
        # Space: O(1)
        user = self._users.get(user_id)
        if not user:
            return []
        return [self._songs[sid].title for sid in user.last_three_played_songs]

    def print_analytics_summary(self) -> None:
        # Time: O(n log n) — full sort of all n songs
        # Space: O(n) — sorted() allocates a new list
        #
        # Trade-off vs heap approach:
        #   Use this when the report must rank ALL songs, or when n is small
        #   enough that the full sort cost is acceptable. Since this is called
        #   once per day while play_song runs millions of times, O(n log n)
        #   here is almost always fine. Prefer this for correctness and clarity.
        songs_by_listeners = sorted(
            self._songs.values(),
            key=lambda s: (-s.unique_listener_count, s.id),
        )

        print("=== Daily Analytics Summary ===")
        print(f"Total songs: {len(self._songs)}")
        print()
        print(f"{'Rank':<6}{'Title':<30}{'Unique Listeners'}")
        print("-" * 50)
        for rank, song in enumerate(songs_by_listeners, start=1):
            print(f"{rank:<6}{song.title:<30}{song.unique_listener_count}")

    def print_analytics_summary_heap(self, top_k: int = 50) -> None:
        # Time: O(n log k) — heapq.nlargest scans all n songs but only keeps
        #   a min-heap of size k, so each push/pop is O(log k) not O(log n).
        #   When k << n this is a meaningful win (e.g. top-50 out of 10M songs).
        #   When k == n it degrades to O(n log n), same as full sort.
        # Space: O(k) — only the top-k songs are kept in memory at once
        #
        # Trade-off vs sort approach:
        #   Use this when the report only needs a leaderboard (top-K songs) and
        #   n is large. If k is small relative to n the savings are significant.
        #   Downside: songs outside the top-k are silently omitted, so this is
        #   wrong when a full ranking is required.
        # nlargest picks the highest key values, so use (count, -id):
        # higher count wins; for equal counts, smaller id wins (larger -id).
        top = heapq.nlargest(
            top_k,
            self._songs.values(),
            key=lambda s: (s.unique_listener_count, -s.id),
        )

        print(f"=== Daily Analytics Summary (Top {top_k}) ===")
        print(f"Total songs: {len(self._songs)}")
        print()
        print(f"{'Rank':<6}{'Title':<30}{'Unique Listeners'}")
        print("-" * 50)
        for rank, song in enumerate(top, start=1):
            print(f"{rank:<6}{song.title:<30}{song.unique_listener_count}")
