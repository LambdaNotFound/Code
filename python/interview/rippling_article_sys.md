We are building the vote management service for an online news platform, where users can upvote or downvote published articles.
We are interested in what makes users change their minds.

Let's start with the following functionality:

add_article(article_name [string]): Each article is given an incremental integer ID when it's added, starting with 1.
upvote_article(article_id [integer], user_id [integer]): Assume any user ID is valid, and that the given article ID will have been added.
downvote_article(article_id [integer], user_id [integer]): Assume any user ID is valid, and that the given article ID will have been added.
print_last_three_flips(user_id [integer]): The titles of the last three unique articles for which the given user changed their vote, either from upvote to downvote or downvote to upvote.
For our MVP, consider performance as we will eventually support millions of articles and users.
However, let's not worry about thread safety or persistence for now - store data in memory.

Let's prioritize solving the problem for the last three articles for now.
We can tackle extensibility at a later stage.

Part 2
Implement get_top_k(k):

Each upvote on an article contributes +1 point; each downvote contributes −1 point.
Based on the current scores, return the top k articles with the highest total scores.