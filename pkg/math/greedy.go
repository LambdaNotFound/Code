package math

/**
 * Math/Greedy Tricks
 */

/**
 * 55. Jump Game
 */
func canJump(nums []int) bool {
    reachablePosition := 0
    for i := 0; i < len(nums); i++ {
        if i > reachablePosition {
            return false
        }
        if i + nums[i] > reachablePosition {
            reachablePosition = i + nums[i]
        }
    }
    return true
}

/**
 * 45. Jump Game II
 */
func jump(nums []int) int {
    curJump, farthestJump, jumps := 0, 0, 0
    for i := 0; i < len(nums)-1; i++ {
        // push index of furthest jump during current iteration
        if i+nums[i] > farthestJump {
            farthestJump = i + nums[i]
        }

        // if current iteration is ended - setup the next one
        if i == curJump {
            jumps, curJump = jumps+1, farthestJump

            if curJump >= len(nums)-1 {
                return jumps
            }
        }
    }

    return 0
}

/**
 * Scheduling & Intervals
 */

/**
 * 134. Gas Station
 * Given two integer arrays gas and cost, return the starting gas station's index
 * if you can travel around the circuit once in the clockwise direction,
 * otherwise return -1
 */
func canCompleteCircuit(gas []int, cost []int) int {
    n := len(gas)
    fuelLeft, globalFuelLeft, start := 0, 0, 0
    for i := 0; i < n; i++ {
        globalFuelLeft += gas[i] - cost[i]
        fuelLeft += gas[i] - cost[i]
        if fuelLeft < 0 {
            start = i + 1 // always restart from the next station after failure
            fuelLeft = 0
        }
    }

    if globalFuelLeft < 0 {
        return -1
    }
    return start
}

/**
 * Tasks / Allocation
 */

/**
 * 621. Task Scheduler
 *
 * but there's a constraint: there has to be a gap of at least n intervals between two tasks with the same label.
 */
func leastInterval(tasks []byte, n int) int {
    maxFreq := 0
    dict := make(map[byte]int)

    for i := 0; i < len(tasks); i++ {
        dict[tasks[i]]++
        maxFreq = max(maxFreq, dict[tasks[i]])
    }

    maxFreqCount := 0
    // if there are tasks with equal freq, then time increases
    for _, value := range dict {
        if value == maxFreq {
            maxFreqCount++
        }
    }

    partCount := maxFreq - 1
    partLength := n - (maxFreqCount - 1) // idol slots in each part
    emptySlots := partCount * partLength
    availableTasks := len(tasks) - maxFreq*maxFreqCount
    idles := max(0, emptySlots-availableTasks)

    return len(tasks) + idles
}
