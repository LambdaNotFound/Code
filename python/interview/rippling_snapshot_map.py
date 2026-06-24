
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
