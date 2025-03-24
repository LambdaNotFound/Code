package palindrome

import "testing"

func Test_longestPalindrome(t *testing.T) {
    input := "babad"
    expect := "bab"

    if result := longestPalindrome(input); result != expect {
        t.Errorf("result = %v, expect = %v", result, expect)
    }
}
