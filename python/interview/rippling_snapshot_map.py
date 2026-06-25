from bisect import bisect_right

# Design:
#   _map       dict[str, int]        current key-value state
#   _snapshots dict[int, dict]       snap_id → frozen copy of _map at snapshot time
#   _next_id   int                   monotonic counter; snap_id starts at 0
#
# TakeSnapshot copies the entire current map in O(n) and stores it under _next_id,
# then increments the counter. Snapshots are immutable — subsequent put/delete only
# affect _map, never _snapshots.
#
# Get(k, snap_id) raises KeyError for both unknown snap_id and missing key within
# the snapshot, matching the spec.


class SnapMap:
    def __init__(self):
        self._map: dict[str, int] = {}
        self._snapshots: dict[int, dict[str, int]] = {}
        self._next_id: int = 0

    def get(self, k: str, snap_id: int | None = None) -> int:  # T: O(1)
        if snap_id is None:
            if k not in self._map:
                raise KeyError(k)
            return self._map[k]

        if snap_id not in self._snapshots:
            raise KeyError(snap_id)
        snapshot = self._snapshots[snap_id]
        if k not in snapshot:
            raise KeyError(k)
        return snapshot[k]

    def put(self, k: str, v: int) -> None:  # T: O(1)
        self._map[k] = v

    def delete(self, k: str) -> None:  # T: O(1)
        if k not in self._map:
            raise KeyError(k)
        del self._map[k]

    def take_snapshot(self) -> int:  # T: O(n), S: O(n) — full copy of current map
        snap_id = self._next_id
        self._snapshots[snap_id] = dict(self._map)
        self._next_id += 1
        return snap_id


_DELETED = object()

# Design:
#   _current  dict[str, int]                    live state for O(1) get/put/delete
#   _history  dict[str, list[tuple[int, any]]]  key → sorted [(snap_id, val)] pairs; at most S entries per key
#   _snap_id  int                               current snapshot context; increments on take_snapshot
#
# Key decision: avoid O(N) full copy on snapshot by storing per-key deltas tagged with snap_id.
# get(k, snap_id) binary-searches the short history list (≤ S entries).
# With <1% keys updated per snapshot and ~100 snapshots: space is O(N) not O(N*S).


class SnapMap2:
    def __init__(self):
        self._current: dict[str, int] = {}
        self._history: dict[str, list[tuple[int, any]]] = {}
        self._snap_id: int = 0

    def get(self, k: str, snap_id: int | None = None) -> int:  # T: O(1) current, O(log S) snapshot; S: O(1)
        if snap_id is None:
            if k not in self._current:
                raise KeyError(k)
            return self._current[k]

        hist = self._history.get(k)
        if not hist:
            raise KeyError(k)
        idx = bisect_right(hist, (snap_id, float("inf"))) - 1
        if idx < 0:
            raise KeyError(k)
        _, val = hist[idx]
        if val is _DELETED:
            raise KeyError(k)
        return val

    def put(self, k: str, v: int) -> None:  # T: O(1); S: O(1) amortized — one entry per (key, snapshot)
        self._current[k] = v
        hist = self._history.setdefault(k, [])
        if hist and hist[-1][0] == self._snap_id:
            hist[-1] = (self._snap_id, v)
        else:
            hist.append((self._snap_id, v))

    def delete(self, k: str) -> None:  # T: O(1); S: O(1) amortized
        if k not in self._current:
            raise KeyError(k)
        del self._current[k]
        hist = self._history[k]
        if hist[-1][0] == self._snap_id:
            hist[-1] = (self._snap_id, _DELETED)
        else:
            hist.append((self._snap_id, _DELETED))

    def take_snapshot(self) -> int:  # T: O(1), S: O(1) — no copy, just increment counter
        snap = self._snap_id
        self._snap_id += 1
        return snap
