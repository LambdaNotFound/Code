package hashmap

import "sort"

/**
 * Analyze User Website Visit Pattern
 *
 */

func mostVisitedPattern(username []string, timestamp []int, website []string) []string {
	type log struct {
		username  string
		website   string
		timestamp int
	}
	logs := make([]log, len(username))
	for i := 0; i < len(username); i++ {
		logs[i] = log{
			username[i], website[i], timestamp[i],
		}
	}

	// group by user + sort by time
	logsByUser := make(map[string][]log)
	for _, log := range logs {
		logsByUser[log.username] = append(logsByUser[log.username], log)
	}
	for _, logs := range logsByUser {
		sort.Slice(logs, func(i, j int) bool {
			return logs[i].timestamp < logs[j].timestamp
		})
	}

	// extract website
	sequenceByUser := make(map[string][]string)
	for _, logs := range logsByUser {
		for _, log := range logs {
			sequenceByUser[log.username] = append(sequenceByUser[log.username], log.website)
		}
	}

	// counting patterns
	patternScore := make(map[[3]string]int)
	for _, sequence := range sequenceByUser {
		if len(sequence) < 3 {
			continue
		}

		for i := 0; i < len(sequence); i++ {
			for j := i + 1; j < len(sequence); j++ {
				for k := j + 1; k < len(sequence); k++ {
					patternScore[[3]string{sequence[i], sequence[j], sequence[k]}]++
				}
			}
		}
	}

	// picking up highest score sequences
	highest := 0
	res := [][]string{}
	for k, v := range patternScore {
		temp := []string{k[0], k[1], k[2]}
		if v > highest {
			highest = v
			res = append([][]string{}, temp)
		} else if v == highest {
			res = append(res, temp)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i][0] != res[j][0] {
			return res[i][0] < res[j][0]
		}
		if res[i][1] != res[j][1] {
			return res[i][1] < res[j][1]
		}
		return res[i][2] < res[j][2]
	})

	return res[0]
}

func mostVisitedPatternClaude(username []string, timestamp []int, website []string) []string {
	// group visits by user, sorted by timestamp
	type visit struct {
		time int
		site string
	}
	userVisits := map[string][]visit{}
	for i := range username {
		userVisits[username[i]] = append(userVisits[username[i]], visit{timestamp[i], website[i]})
	}
	for user := range userVisits {
		sort.Slice(userVisits[user], func(i, j int) bool {
			return userVisits[user][i].time < userVisits[user][j].time
		})
	}

	// count patterns — each pattern counted once per user
	patternCount := map[[3]string]int{}
	for _, visits := range userVisits {
		sites := make([]string, len(visits))
		for i, v := range visits {
			sites[i] = v.site
		}

		// deduplicate patterns for this user
		seen := map[[3]string]bool{}
		for i := 0; i < len(sites); i++ {
			for j := i + 1; j < len(sites); j++ {
				for k := j + 1; k < len(sites); k++ {
					pattern := [3]string{sites[i], sites[j], sites[k]}
					if !seen[pattern] {
						seen[pattern] = true
						patternCount[pattern]++
					}
				}
			}
		}
	}

	// find pattern with highest count, break ties lexicographically
	var best [3]string
	highestCount := 0
	for pattern, count := range patternCount {
		if count > highestCount ||
			(count == highestCount && pattern < best) { // [3]string supports < directly!
			highestCount = count
			best = pattern
		}
	}

	return []string{best[0], best[1], best[2]}
}
