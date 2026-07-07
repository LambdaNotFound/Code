package interview

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

/*
We are writing software to analyze logs for toll booths on a highway. This highway is a divided highway with limited access; the only way on to or off of the highway is through a toll booth.

There are three types of toll booths:
* ENTRY toll booths, where a car goes through a booth as it enters the highway.
* EXIT toll booths, where a car goes through a booth as it exits the highway.
* MAINROAD (M in the diagram), which have sensors that record a license plate as a car drives through at full speed.

        Exit Booth                         Entry Booth
            |                                   |
            |                                   |
             \                                 /
---<------------<---------M---------<-----------<---------<----
                                         (West-bound side)

===============================================================

                                         (East-bound side)
------>--------->---------M--------->--------->--------->------
             /                                 \
            |                                   |
            |                                   |
        Entry Booth                         Exit Booth

For our first task:
1-1) Read through and understand the code and comments below. Feel free to run the code and tests.
1-2) The tests are not passing due to a bug in the code. Make the necessary changes to HighwayLogEntry to fix the bug.
*/

/*
We are interested in how many people are using the highway, and so we would like to count how many complete journeys are taken in the log file.

A complete journey consists of:
* A driver entering the highway through an ENTRY toll booth.
* The driver passing through some number of MAINROAD toll booths (possibly 0).
* The driver exiting the highway through an EXIT toll booth.

For example, the following excerpt of log lines contains complete journeys for the cars with JOX304 and THX138:

.
.
.
90750.191 JOX304 250E ENTRY
91081.684 JOX304 260E MAINROAD
91082.101 THX138 110E ENTRY
91483.251 JOX304 270E MAINROAD
91873.920 THX138 120E MAINROAD
91874.493 JOX304 280E EXIT
.
.
91982.102 THX138 290E EXIT
92301.302 THX138 300E ENTRY
92371.302 THX138 310E EXIT
.

HashMap key: platenumber value: int
                       1) check the BoothType:
                         if ENTRY +1
                         if EXIT -1
                       2) if HashMap[placeNumber] == 0
                          this car completed 1 journey => totalJourneyCount++

→ This log contains 3 complete journeys:
  • JOX304: 1 journey
  • THX138: 2 journeys

You may assume that the log only contains complete journeys, and there are no missing entries.

2-1) Write a function in LogFile named CountJourneys() that returns how many
     complete journeys there are in the given LogFile.
*/

/*
We would like to catch people who are driving at unsafe speeds on the highway. To help us do that, we would like to identify journeys where a driver does either of the following:
* Drive 130 km/h or greater in any individual 10km segment of tollway.
* Drive 120 km/h or greater in any two 10km segments of tollway.

For example, consider the following log:
90750.191 JOX304 250E ENTRY
91081.684 JOX304 260E MAINROAD
91082.101 THX138 110E ENTRY
91483.251 JOX304 270E MAINROAD
91873.920 THX138 120E MAINROAD
91874.493 JOX304 280E EXIT
.
.
91982.102 THX138 290E EXIT
92301.302 THX138 300E ENTRY
92371.302 THX138 310E EXIT
.
.
93005.405 TST002 270W ENTRY
93280.609 TST002 260W EXIT

In this case, the driver of TST002 drove 10 km in 275 seconds. We can calculate
that this driver drove an average speed of ~130.91km/hr over this segment:

10 km * 3600 sec/hr
------------------- = 130.91 km/hr
      275 sec

Note that:
* A license plate may have multiple journeys in one file, and if they drive at unsafe speeds in both journeys, both should be counted.
* We do not mark speeding if they are not on the highway (i.e. for any driving between an EXIT and ENTRY event).
* Speeding is only marked once per journey. For example, if there are 4 segments 120km/h or greater, or multiple segments 130km/h or greater, the journey is only counted once.

3-1) Write a function CatchSpeeders in LogFile that returns a collection of license plates that drove at unsafe speeds during a journey in the LogFile.
     If the same license plate drives at unsafe speeds during two different journeys, the license plate should appear twice (once for each journey they drove at unsafe speeds).
*/

type HighwayLogEntry struct {
	/**
	 * Represents an entry from a single log line. Log lines look like this in the file:
	 *
	 * 34400.409 SXY288 210E ENTRY
	 *
	 * Where:
	 * * 34400.409 is the timestamp in seconds since the software was started.
	 * * SXY288 is the license plate of the vehicle passing through the toll booth.
	 * * 210E is the location and traffic direction of the toll booth. Here, the toll
	 *     booth is at 210 kilometers from the start of the tollway, and the E indicates
	 *     that the toll booth was on the east-bound traffic side. Tollbooths are placed
	 *     every ten kilometers.
	 * * ENTRY indicates which type of toll booth the vehicle went through. This is one of
	 *     "ENTRY", "EXIT", or "MAINROAD".
	 **/
	Timestamp    float64
	LicensePlate string
	BoothType    string
	Location     int
	Direction    string
}

func NewLogEntry(logLine string) HighwayLogEntry {
	// fmt.Println("logLine: ", logLine)
	tokens := strings.Fields(logLine)
	location := tokens[2][:len(tokens[2])-1]
	// fmt.Println("location: ", location)
	directionLetter := tokens[2][len(tokens[2])-1]
	// fmt.Println("directionLetter: ", directionLetter)

	direction := ""
	if directionLetter == 'E' {
		direction = "EAST"
	} else if directionLetter == 'W' {
		direction = "WEST"
	} else {
		panic("Invalid direction letter")
	}

	timestamp, _ := strconv.ParseFloat(tokens[0], 64)

	return HighwayLogEntry{
		Timestamp:    timestamp, // tokens[0]
		LicensePlate: tokens[1],
		BoothType:    tokens[3],
		Location:     parseInt(location),
		Direction:    direction,
	}
}

func (entry HighwayLogEntry) String() string {
	return fmt.Sprintf("<HighwayLogEntry timestamp: %f  license: %s  location: %d  direction: %s  booth type: %s>",
		entry.Timestamp, entry.LicensePlate, entry.Location, entry.Direction, entry.BoothType)
}

type LogFile struct {
	/*
	 * Represents a file containing a number of log lines, converted to HighwayLogEntry
	 * objects.
	 */
	LogEntries []HighwayLogEntry
}

func (l *LogFile) CountJourneys() int {
	journeyByLicensePlate := map[string]int{}
	journeyCount := 0

	for _, log := range l.LogEntries {
		// if _, found := journeyByLicensePlate[log.LicensePlate]; !found {
		if log.BoothType == "ENTRY" {
			journeyByLicensePlate[log.LicensePlate]++
		}
		if log.BoothType == "EXIT" {
			journeyByLicensePlate[log.LicensePlate]--
		}

		if journeyByLicensePlate[log.LicensePlate] == 0 {
			journeyCount++
		}
	}

	return journeyCount
}

func (l *LogFile) CatchSpeeders() []string {
	flagged := map[string]bool{} // already added to list this journey
	speedersList := []string{}
	lastTS := map[string]float64{} // previous booth timestamp
	highCount := map[string]int{}  // segments >= 120 km/h this journey

	for _, entry := range l.LogEntries {
		if _, found := lastTS[entry.LicensePlate]; !found {
			lastTS[entry.LicensePlate] = entry.Timestamp
			continue
		}

		elapsed := entry.Timestamp - lastTS[entry.LicensePlate]
		speed := 10.0 * 3600 / elapsed
		lastTS[entry.LicensePlate] = entry.Timestamp // Bug 2 fix

		if speed >= 120 {
			highCount[entry.LicensePlate]++
		}
		isSpeeder := speed >= 130 || highCount[entry.LicensePlate] >= 2 // Bug 4 fix
		if isSpeeder && !flagged[entry.LicensePlate] {
			speedersList = append(speedersList, entry.LicensePlate)
			flagged[entry.LicensePlate] = true
		}

		if entry.BoothType == "EXIT" {
			delete(lastTS, entry.LicensePlate) // Bug 3 fix
			delete(flagged, entry.LicensePlate)
			delete(highCount, entry.LicensePlate)
		}
	}
	return speedersList
}

func NewLogFile(fileName string) *LogFile {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
	}
	defer file.Close()

	logFile := &LogFile{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		logEntry := NewLogEntry(line)
		logFile.LogEntries = append(logFile.LogEntries, logEntry)
	}
	if scanner.Err() != nil {
		fmt.Printf("Error reading log file: %v\n", scanner.Err())
	}

	return logFile
}

func (file *LogFile) Size() int {
	return len(file.LogEntries)
}

