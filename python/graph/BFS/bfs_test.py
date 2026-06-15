"""
Unit tests for bfs.py — generic BFS, level-order BFS, and grid BFS
"""
import sys
import os

import pytest

sys.path.insert(0, os.path.dirname(__file__))
from bfs import bfs, bfs_levels, bfs_grid


# ── bfs ───────────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("start, graph, want", [
    # linear chain: all nodes reachable from start
    (0, {0: [1], 1: [2], 2: []}, {0, 1, 2}),
    # disconnected graph: only the reachable component is returned
    (0, {0: [1], 1: [0], 2: [3], 3: [2]}, {0, 1}),
    # single node with no neighbors
    ("a", {"a": []}, {"a"}),
    # undirected graph: all nodes reachable
    (1, {1: [2, 3], 2: [1, 4], 3: [1], 4: [2]}, {1, 2, 3, 4}),
    # cycle: must not loop forever
    (0, {0: [1], 1: [2], 2: [0]}, {0, 1, 2}),
    # star graph: all leaves reachable from center
    (0, {0: [1, 2, 3], 1: [], 2: [], 3: []}, {0, 1, 2, 3}),
])
def test_bfs(start, graph, want):
    assert bfs(start, graph) == want


# ── bfs_levels ────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("start, target, graph, want", [
    # target is the start node itself
    (0, 0, {0: [1], 1: []}, 0),
    # target is one hop away
    (0, 1, {0: [1], 1: []}, 1),
    # target is two hops away
    (0, 2, {0: [1], 1: [2], 2: []}, 2),
    # target is unreachable: disconnected component
    (0, 3, {0: [1], 1: [2], 2: [0], 3: []}, -1),
    # two paths to target; BFS finds the shorter one (2 hops)
    (0, 3, {0: [1, 2], 1: [3], 2: [3], 3: []}, 2),
    # longer chain: target 4 hops out
    (0, 4, {0: [1], 1: [2], 2: [3], 3: [4], 4: []}, 4),
])
def test_bfs_levels(start, target, graph, want):
    assert bfs_levels(start, target, graph) == want


# ── bfs_grid ──────────────────────────────────────────────────────────────────

@pytest.mark.parametrize("grid, start, want", [
    # 1×1 open cell: only start is visited
    (
        [["."]],
        (0, 0),
        {(0, 0)},
    ),
    # 1×3 open row: all three cells reachable
    (
        [[".", ".", "."]],
        (0, 0),
        {(0, 0), (0, 1), (0, 2)},
    ),
    # wall immediately to the right: BFS cannot cross
    (
        [[".", "#", "."]],
        (0, 0),
        {(0, 0)},
    ),
    # 3×3 fully open grid: all 9 cells reachable
    (
        [[".", ".", "."], [".", ".", "."], [".", ".", "."]],
        (0, 0),
        {(r, c) for r in range(3) for c in range(3)},
    ),
    # horizontal wall in the middle: bottom row still reachable via right column
    (
        [[".", ".", "."], ["#", "#", "."], [".", ".", "."]],
        (0, 0),
        {(0, 0), (0, 1), (0, 2), (1, 2), (2, 0), (2, 1), (2, 2)},
    ),
    # start is surrounded by walls: only start visited
    (
        [["#", "#", "#"], ["#", ".", "#"], ["#", "#", "#"]],
        (1, 1),
        {(1, 1)},
    ),
])
def test_bfs_grid(grid, start, want):
    assert bfs_grid(grid, start) == want
