import sys, os
sys.path.insert(0, os.path.dirname(__file__))
from rippling_music_player import MusicPlayer

import pytest


def test_add_song_returns_incremental_ids():
    mp = MusicPlayer()
    assert mp.add_song("Song A") == 1
    assert mp.add_song("Song B") == 2
    assert mp.add_song("Song C") == 3


def test_play_song_increments_play_count():
    mp = MusicPlayer()
    sid = mp.add_song("Blinding Lights")
    mp.play_song(sid, user_id=1)
    mp.play_song(sid, user_id=2)
    assert mp._songs[sid].play_count == 2


def test_play_song_deduplicates_users():
    mp = MusicPlayer()
    sid = mp.add_song("Shape of You")
    mp.play_song(sid, user_id=1)
    mp.play_song(sid, user_id=1)  # same user twice
    mp.play_song(sid, user_id=2)
    song = mp._songs[sid]
    assert song.play_count == 3
    assert song.unique_listener_count == 2


def test_analytics_summary_ordering(capsys):
    mp = MusicPlayer()
    a = mp.add_song("A")
    b = mp.add_song("B")
    c = mp.add_song("C")

    # B has 3 unique listeners, A has 2, C has none
    mp.play_song(b, 1)
    mp.play_song(b, 2)
    mp.play_song(b, 3)
    mp.play_song(a, 1)
    mp.play_song(a, 2)
    # Same user replaying A doesn't change unique count
    mp.play_song(a, 1)

    mp.print_analytics_summary()
    output = capsys.readouterr().out
    table = output.split("-" * 50)[-1]

    b_pos = table.index("B")
    a_pos = table.index("A")
    c_pos = table.index("C")
    assert b_pos < a_pos < c_pos


def test_analytics_summary_tie_broken_by_id(capsys):
    mp = MusicPlayer()
    x = mp.add_song("X")
    y = mp.add_song("Y")
    # Each played by exactly 1 unique user
    mp.play_song(x, 1)
    mp.play_song(y, 2)

    mp.print_analytics_summary()
    output = capsys.readouterr().out
    table = output.split("-" * 50)[-1]

    # Both have unique_listener_count=1; lower id (X) should rank first
    x_pos = table.index("X")
    y_pos = table.index("Y")
    assert x_pos < y_pos


def test_analytics_summary_zero_plays(capsys):
    mp = MusicPlayer()
    mp.add_song("Silent Song")
    mp.print_analytics_summary()
    output = capsys.readouterr().out
    assert "Silent Song" in output
    assert "0" in output


# ---------------------------------------------------------------------------
# print_analytics_summary_heap
# ---------------------------------------------------------------------------

def test_heap_summary_ordering(capsys):
    mp = MusicPlayer()
    a = mp.add_song("A")
    b = mp.add_song("B")
    c = mp.add_song("C")

    mp.play_song(b, 1)
    mp.play_song(b, 2)
    mp.play_song(b, 3)
    mp.play_song(a, 1)
    mp.play_song(a, 2)

    mp.print_analytics_summary_heap(top_k=3)
    output = capsys.readouterr().out
    table = output.split("-" * 50)[-1]

    assert table.index("B") < table.index("A") < table.index("C")


def test_heap_summary_top_k_limits_results(capsys):
    mp = MusicPlayer()
    for i in range(5):
        sid = mp.add_song(f"Song{i}")
        for u in range(5 - i):       # Song0 gets 5 listeners, Song4 gets 1
            mp.play_song(sid, u)

    mp.print_analytics_summary_heap(top_k=3)
    output = capsys.readouterr().out

    assert "Song0" in output
    assert "Song1" in output
    assert "Song2" in output
    assert "Song3" not in output
    assert "Song4" not in output


def test_heap_summary_matches_sort_for_full_k(capsys):
    mp = MusicPlayer()
    a = mp.add_song("Alpha")
    b = mp.add_song("Beta")
    mp.play_song(a, 1)
    mp.play_song(b, 2)
    mp.play_song(b, 3)

    mp.print_analytics_summary()
    sort_out = capsys.readouterr().out.split("-" * 50)[-1]

    mp.print_analytics_summary_heap(top_k=2)
    heap_out = capsys.readouterr().out.split("-" * 50)[-1]

    # Same ranking when top_k covers all songs
    assert sort_out.index("Beta") < sort_out.index("Alpha")
    assert heap_out.index("Beta") < heap_out.index("Alpha")


# ---------------------------------------------------------------------------
# last_three_played_song_titles
# ---------------------------------------------------------------------------

def test_last_three_basic_order():
    mp = MusicPlayer()
    a = mp.add_song("Alpha")
    b = mp.add_song("Beta")
    c = mp.add_song("Gamma")

    mp.play_song(a, 1)
    mp.play_song(b, 1)
    mp.play_song(c, 1)

    # Most recent first
    assert mp.last_three_played_song_titles(1) == ["Gamma", "Beta", "Alpha"]


def test_last_three_caps_at_three():
    mp = MusicPlayer()
    ids = [mp.add_song(f"Song{i}") for i in range(5)]
    for sid in ids:
        mp.play_song(sid, 1)

    # Only the last 3 plays remain: Song4, Song3, Song2
    assert mp.last_three_played_song_titles(1) == ["Song4", "Song3", "Song2"]


def test_last_three_deduplicates_moves_to_front():
    mp = MusicPlayer()
    a = mp.add_song("Alpha")
    b = mp.add_song("Beta")
    c = mp.add_song("Gamma")

    mp.play_song(a, 1)
    mp.play_song(b, 1)
    mp.play_song(c, 1)
    # Replay Alpha — it should move to front, not duplicate
    mp.play_song(a, 1)

    assert mp.last_three_played_song_titles(1) == ["Alpha", "Gamma", "Beta"]


def test_last_three_unknown_user_returns_empty():
    mp = MusicPlayer()
    mp.add_song("Alpha")
    assert mp.last_three_played_song_titles(999) == []


def test_last_three_users_are_independent():
    mp = MusicPlayer()
    a = mp.add_song("Alpha")
    b = mp.add_song("Beta")

    mp.play_song(a, 1)
    mp.play_song(b, 2)

    assert mp.last_three_played_song_titles(1) == ["Alpha"]
    assert mp.last_three_played_song_titles(2) == ["Beta"]


def test_last_three_replay_promotes_evicted_song():
    # Play order (oldest → newest):
    #   Bohemian Rhapsody → Hello, Goodbye → Stairway to Heaven → Hello, Goodbye
    #
    # After 3 plays:    [Stairway, Hello, Bohemian]
    # Hello replayed:   remove Hello → [Stairway, Bohemian]
    #                   appendleft    → [Hello, Stairway, Bohemian]
    mp = MusicPlayer()
    hello = mp.add_song("Hello, Goodbye")
    stairway = mp.add_song("Stairway to Heaven")
    bohemian = mp.add_song("Bohemian Rhapsody")

    mp.play_song(bohemian, 1)
    mp.play_song(hello, 1)
    mp.play_song(stairway, 1)
    mp.play_song(hello, 1)   # replay — moves Hello back to front

    assert mp.last_three_played_song_titles(1) == [
        "Hello, Goodbye",       # most recent
        "Stairway to Heaven",   # second most recent
        "Bohemian Rhapsody",    # third most recent
    ]
