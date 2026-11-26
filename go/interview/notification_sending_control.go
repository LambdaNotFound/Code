package interview

/*
You're implementing the logic for a service that performs rate limiting on the notifications we send out.

We don't want to send out notifications from a service more than K days in a row.

You should implement the following functions that will be called repeatedly: shouldSendNotification(sendingService: str, day: int) -> bool

Conditions & Parameters:

Days are non-decreasing (i.e, they should handle case where a service receives multiple notifications on the same day).
Service arguments are english-language character only strings like "livestreams", "giveaways", "refunds".
*/
type RateLimiter struct {
	k       int
	lastDay map[string]int // service → last day seen
	streak  map[string]int // service → current consecutive day count

	cooldown       map[string]int // per-service cooldown duration in days
	lastBlockedDay map[string]int // day the service was last blocked
}

func NewRateLimiter(k int) *RateLimiter {
	return &RateLimiter{
		k:       k,
		lastDay: make(map[string]int),
		streak:  make(map[string]int),
	}
}

func (r *RateLimiter) ShouldSendNotification(service string, day int) bool {
	last, seen := r.lastDay[service]

	if !seen {
		r.lastDay[service] = day
		r.streak[service] = 1
		return true
	}

	if day == last {
		// same day: never update state, just check
		return r.streak[service] <= r.k
	}

	if day == last+1 {
		// consecutive day: check FIRST, update only if allowed
		if r.streak[service]+1 <= r.k {
			r.lastDay[service] = day
			r.streak[service]++
			return true
		}
		return false // blocked — state unchanged
	}

	// gap: always reset and allow
	r.lastDay[service] = day
	r.streak[service] = 1
	return true
}

func (r *RateLimiter) ShouldSendNotification2(service string, day int) bool {
	// check cooldown first — regardless of streak state
	if blockedDay, wasBlocked := r.lastBlockedDay[service]; wasBlocked {
		cooldownDays := r.cooldown[service]
		if day <= blockedDay+cooldownDays {
			return false // still within cooldown window
		}
		// cooldown expired — clear blocked state
		delete(r.lastBlockedDay, service)
	}

	last, seen := r.lastDay[service]

	if !seen {
		r.lastDay[service] = day
		r.streak[service] = 1
		return true
	}

	if day == last {
		// same day: state unchanged, check streak
		return r.streak[service] <= r.k
	}

	if day == last+1 {
		// consecutive day: check before updating
		if r.streak[service]+1 <= r.k {
			r.lastDay[service] = day
			r.streak[service]++
			return true
		}
		// blocked — record blocked day, leave streak/lastDay unchanged
		r.lastBlockedDay[service] = day
		return false
	}

	// gap in days: reset streak regardless of cooldown expiry
	r.lastDay[service] = day
	r.streak[service] = 1
	return true
}

// queue based approach
/*
func (r *RateLimiter) ShouldSendNotification3(service string, day int) bool {
    // check cooldown window first
    if blockedDay, wasBlocked := r.blocked[service]; wasBlocked {
        if day <= blockedDay+r.getCooldown(service) {
            return false                      // still cooling down
        }
        delete(r.blocked, service)            // cooldown expired
    }

    queue := r.queues[service]

    // same day as last send: don't enqueue, just check
    if len(queue) > 0 && queue[len(queue)-1] == day {
        return true                           // already counted this day
    }

    // gap detected: reset queue entirely
    if len(queue) > 0 && day > queue[len(queue)-1]+1 {
        r.queues[service] = []int{day}
        return true
    }

    // consecutive day: check if window is full
    if len(queue) >= r.k {
        // queue full and new day is consecutive → blocked
        r.blocked[service] = day
        return false
    }

    // enqueue and allow
    r.queues[service] = append(r.queues[service], day)
    return true
}

func (r *RateLimiter) getCooldown(service string) int {
    if days, found := r.cooldown[service]; found {
        return days
    }
    return 0
}
*/
