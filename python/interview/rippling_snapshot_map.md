### Overview
Deisng a key-value loopup map with snapshot functionaly

### API
1. Get(k): This method should return the value v associated with the key k. If the key k does not exist, it should raise a KeyError exception. The key k will be a string and the value v will be an integer.
2. Put(k, v): This method should add a key‑value pair to the map. The key k will be a string and the value v will be an integer.
3. Delete(k): This method should remove the key‑value pair associated with the key k from the map. If the key k does not exist, it should raise a KeyError exception.

4. TakeSnapshot(): This method should capture the current state of the map and return a unique snapshot ID (snap_id), which will be an integer; snap_id should start from 0 and increases every time when a new snapshot is taken.
5. Get(k, snap_id): This method should return the value v associated with the key k from the snapshot identified by snap_id. If the key k or snapshotId does not exist in the snapshot, it should raise a KeyError exception.

### Design spec
1. Apply OOP principle in this design. 
2. the basic map operations can be support by a map in Python.
3. For this SnapMap class, uese a hashmap to track the current map status, and use a hashmap of maps & an int counter to support the snap functionality. e.g. when take a snapshot, insert a copy of the current map to the hashmap the with the count as the key, inc count after the insertion. 
4. the snapshot is immutable, it cannot be modified with following update after taking the snapshot.

### Test Plan
1. test map operation for get, put and delete
2. test snap support by taking snapshot and revert to a certain version of the map