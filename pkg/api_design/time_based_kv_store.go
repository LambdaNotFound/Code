package apidesign

/**
 * 981. Time Based Key-Value Store
 *
 * Implement the TimeMap class:
 *
 * TimeMap() Initializes the object of the data structure.
 * void set(String key, String value, int timestamp) Stores the key
 * key with the value value at the given time timestamp.
 *
 * String get(String key, int timestamp) Returns a value such that
 * set was called previously, with timestamp_prev <= timestamp.
 * If there are multiple such values, it returns the value associated with
 * the largest timestamp_prev. If there are no values, it returns "".
 *
 */
type TimeMap struct {
    KeyMap map[string][]Pair
}

type Pair struct {
    val       string
    timeStamp int
}

func ConstructorTimeMap() TimeMap {
    return TimeMap{
        KeyMap: make(map[string][]Pair),
    }
}

func (t *TimeMap) Set(key string, value string, timestamp int) {
    pair := Pair{
        val:       value,
        timeStamp: timestamp,
    }

    if _, ok := t.KeyMap[key]; !ok {
        t.KeyMap[key] = make([]Pair, 0)
    }
    t.KeyMap[key] = append(t.KeyMap[key], pair)
}

func (t *TimeMap) Get(key string, timestamp int) string {
    if _, ok := t.KeyMap[key]; !ok {
        return ""
    }
    arr := t.KeyMap[key]
    if arr[0].timeStamp > timestamp {
        return ""
    }

    left, right := 0, len(arr)
    for left < right {
        mid := left + (right-left)/2
        if arr[mid].timeStamp < timestamp { // lower_bound(), left is the first element <= target
            left = mid + 1
        } else {
            right = mid
        }
    }
    if left < len(arr) && arr[left].timeStamp == timestamp {
        return arr[left].val
    }
    return arr[left-1].val
}
