### Overview
Design a in-memory music player class for tracking songs and play times and ranking the popular songs, using OO design principles. 

### Data Model

Song
-title string
-id int // incremental number
-play_count // int for counting times played
-played_by // a set for tracking unique users

User
-id int
-name string

### API
add_song (song_title [string]) [integer]- A song is given an incremental integer ID when it's added, starting with 1.

play_song (song_ID [integer], user_ID [integer])- Assume any user ID is valid, and that the given song ID will have been added.

print_analytics_summary () This is used for a report, created once per day for our company's Analytics department.
- The summary should be sorted (descending) by the number of unique users who have played each song.
- The summary should include the song titles, and the number of unique users, but the formatting does not matter.

last_three_played_song_titles(user_ID [integer]) - Returns the titles of the last three unique played songs for the given user (ordered, most recent first).

### High Level Design
MusicPlayer class stores songs and users. it has 2 hashmaps (key is the song's id and user's id respectively).

add_song() creates a song with incremental id (lenght of the hashmap of songs +1)

play_song() add the user to the song's played_by set (using set for dedup), and increase the play_count

last_three_played_song_titles() use a queue (last_three_played_songs) to track the last played 3 songs by id. this queue (array) is on user entity. when play song by a user, if the song is in the last_three_played_songs array, move it to the beginning of the array. otherwise pop the last id of the array and insert the id in the beginning. Keep the FIFO ordering of the songs id.