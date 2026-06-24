import sys, os
sys.path.insert(0, os.path.dirname(__file__))

import pytest
from rippling_article_sys import ArticleSystem


@pytest.fixture
def sys2():
    """Two-article system used across most tests."""
    s = ArticleSystem()
    s.add_article("Alpha")   # id=1
    s.add_article("Beta")    # id=2
    return s


class TestAddArticle:
    def test_ids_increment_from_one(self):
        s = ArticleSystem()
        assert s.add_article("A") == 1
        assert s.add_article("B") == 2
        assert s.add_article("C") == 3


class TestScores:
    def test_upvote_increments_score(self, sys2):
        sys2.upvote_article(1, 1)
        assert sys2.get_top_k(1) == [("Alpha", 1)]

    def test_downvote_decrements_score(self, sys2):
        sys2.downvote_article(1, 1)
        assert sys2.get_top_k(1) == [("Beta", 0)]  # Alpha=-1, Beta=0 → Beta wins

    def test_same_vote_is_noop(self, sys2):
        sys2.upvote_article(1, 1)
        sys2.upvote_article(1, 1)
        assert sys2._scores[1] == 1

    def test_flip_up_to_down_adjusts_score_by_two(self, sys2):
        sys2.upvote_article(1, 1)    # score=+1
        sys2.downvote_article(1, 1)  # score=-1
        assert sys2._scores[1] == -1

    def test_flip_down_to_up_adjusts_score_by_two(self, sys2):
        sys2.downvote_article(1, 1)  # score=-1
        sys2.upvote_article(1, 1)    # score=+1
        assert sys2._scores[1] == 1

    def test_multiple_users_accumulate(self, sys2):
        sys2.upvote_article(1, 1)
        sys2.upvote_article(1, 2)
        sys2.downvote_article(1, 3)
        assert sys2._scores[1] == 1  # +1 +1 -1


class TestPrintLastThreeFlips:
    def test_no_flips_returns_empty(self, sys2):
        assert sys2.print_last_three_flips(1) == []

    def test_first_vote_is_not_a_flip(self, sys2):
        sys2.upvote_article(1, 1)
        assert sys2.print_last_three_flips(1) == []

    def test_single_flip_recorded(self, sys2):
        sys2.upvote_article(1, 1)
        sys2.downvote_article(1, 1)
        assert sys2.print_last_three_flips(1) == ["Alpha"]

    def test_returns_most_recent_first(self, sys2):
        sys2.add_article("Gamma")    # id=3
        # user 1: flip on Alpha, then flip on Beta
        sys2.upvote_article(1, 1);   sys2.downvote_article(1, 1)   # flip Alpha
        sys2.upvote_article(2, 1);   sys2.downvote_article(2, 1)   # flip Beta
        assert sys2.print_last_three_flips(1) == ["Beta", "Alpha"]

    def test_capped_at_three(self, sys2):
        sys2.add_article("Gamma")    # id=3
        sys2.add_article("Delta")    # id=4
        for aid in [1, 2, 3, 4]:
            sys2.upvote_article(aid, 1)
            sys2.downvote_article(aid, 1)
        # flips in order: 1,2,3,4 → last 3 unique = [4,3,2]
        assert sys2.print_last_three_flips(1) == ["Delta", "Gamma", "Beta"]

    def test_deduplicates_repeated_flips_on_same_article(self, sys2):
        sys2.add_article("Gamma")    # id=3
        # flip Alpha, flip Beta, flip Alpha again
        sys2.upvote_article(1, 1);   sys2.downvote_article(1, 1)   # flip Alpha
        sys2.upvote_article(2, 1);   sys2.downvote_article(2, 1)   # flip Beta
        sys2.upvote_article(1, 1)                                   # flip Alpha again
        # _user_flips should be [Beta, Alpha] (Alpha moved to most recent)
        assert sys2.print_last_three_flips(1) == ["Alpha", "Beta"]

    def test_different_users_independent(self, sys2):
        sys2.upvote_article(1, 1);   sys2.downvote_article(1, 1)   # user 1 flips Alpha
        sys2.upvote_article(2, 2);   sys2.downvote_article(2, 2)   # user 2 flips Beta
        assert sys2.print_last_three_flips(1) == ["Alpha"]
        assert sys2.print_last_three_flips(2) == ["Beta"]


class TestGetTopK:
    def test_returns_k_articles(self, sys2):
        sys2.upvote_article(1, 1)
        sys2.upvote_article(1, 2)
        sys2.upvote_article(2, 1)
        result = sys2.get_top_k(2)
        assert len(result) == 2
        assert result[0] == ("Alpha", 2)
        assert result[1] == ("Beta", 1)

    def test_unvoted_articles_appear_with_zero_score(self, sys2):
        result = sys2.get_top_k(2)
        assert ("Alpha", 0) in result
        assert ("Beta", 0) in result

    def test_k_larger_than_articles_returns_all(self, sys2):
        assert len(sys2.get_top_k(10)) == 2

    def test_tie_broken_by_lower_article_id(self, sys2):
        # both articles score 0; Alpha (id=1) should rank above Beta (id=2)
        result = sys2.get_top_k(2)
        assert result[0][0] == "Alpha"

    def test_negative_scores_ranked_last(self, sys2):
        sys2.downvote_article(1, 1)
        sys2.downvote_article(1, 2)
        result = sys2.get_top_k(2)
        assert result[0] == ("Beta", 0)
        assert result[1] == ("Alpha", -2)
