import sys, os
sys.path.insert(0, os.path.dirname(__file__))

import pytest
from rippling_snapshot_map import SnapMap


@pytest.fixture
def m():
    return SnapMap()


class TestBasicMapOps:
    def test_put_and_get(self, m):
        m.put("a", 1)
        assert m.get("a") == 1

    def test_put_overwrites(self, m):
        m.put("a", 1)
        m.put("a", 2)
        assert m.get("a") == 2

    def test_get_missing_key_raises(self, m):
        with pytest.raises(KeyError):
            m.get("x")

    def test_delete_removes_key(self, m):
        m.put("a", 1)
        m.delete("a")
        with pytest.raises(KeyError):
            m.get("a")

    def test_delete_missing_key_raises(self, m):
        with pytest.raises(KeyError):
            m.delete("x")


class TestTakeSnapshot:
    def test_snap_ids_start_at_zero(self, m):
        assert m.take_snapshot() == 0

    def test_snap_ids_increment(self, m):
        assert m.take_snapshot() == 0
        assert m.take_snapshot() == 1
        assert m.take_snapshot() == 2

    def test_snapshot_captures_current_state(self, m):
        m.put("a", 10)
        snap = m.take_snapshot()
        assert m.get("a", snap) == 10

    def test_snapshot_is_immutable_after_put(self, m):
        m.put("a", 10)
        snap = m.take_snapshot()
        m.put("a", 99)
        assert m.get("a", snap) == 10   # snapshot unchanged
        assert m.get("a") == 99         # current map updated

    def test_snapshot_is_immutable_after_delete(self, m):
        m.put("a", 10)
        snap = m.take_snapshot()
        m.delete("a")
        assert m.get("a", snap) == 10   # still visible in snapshot
        with pytest.raises(KeyError):
            m.get("a")                  # gone from current map

    def test_snapshot_of_empty_map(self, m):
        snap = m.take_snapshot()
        with pytest.raises(KeyError):
            m.get("a", snap)

    def test_multiple_snapshots_are_independent(self, m):
        m.put("a", 1)
        snap0 = m.take_snapshot()
        m.put("a", 2)
        snap1 = m.take_snapshot()
        assert m.get("a", snap0) == 1
        assert m.get("a", snap1) == 2

    def test_get_invalid_snap_id_raises(self, m):
        with pytest.raises(KeyError):
            m.get("a", 99)

    def test_get_missing_key_in_snapshot_raises(self, m):
        snap = m.take_snapshot()
        with pytest.raises(KeyError):
            m.get("missing", snap)
