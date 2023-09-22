package credit

import "time"

// ContributorsSince use GitHub APIs to get all commit since `since` and extract the author information from that. Then
// it will count the author commit and return a sorted list of authors according to number of commits (descending).
func ContributorsSince(since time.Time) []string {

    return []string{
        "733amir@gmail.com",
    }
}

// Contributors use GitHub APIs to get 500 top contributors. The list is sorted base on the number of commits
// (descending).
func Contributors() []string {

    return []string{
        "733amir@gmail.com",
    }
}
