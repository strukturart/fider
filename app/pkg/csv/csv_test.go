package csv_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/entities"
	"github.com/getfider/fider/app/models/enum"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/csv"
)

func TestExportPostsToCSV_Empty(t *testing.T) {
	RegisterT(t)

	posts := []*entities.Post{}
	expected, err := ioutil.ReadFile("./testdata/empty.csv")
	Expect(err).IsNil()
	actual, err := csv.FromPosts(posts)
	Expect(err).IsNil()
	Expect(actual).Equals(expected)
}

func TestExportPostsToCSV_OnePost(t *testing.T) {
	RegisterT(t)

	posts := []*entities.Post{
		declinedPost,
	}

	expected, err := ioutil.ReadFile("./testdata/one-post.csv")
	Expect(err).IsNil()
	actual, err := csv.FromPosts(posts)
	Expect(err).IsNil()
	Expect(actual).Equals(expected)
}

func TestExportPostsToCSV_MorePosts(t *testing.T) {
	RegisterT(t)

	posts := []*entities.Post{
		declinedPost,
		openPost,
		duplicatePost,
	}

	expected, err := ioutil.ReadFile("./testdata/more-posts.csv")
	Expect(err).IsNil()
	actual, err := csv.FromPosts(posts)
	Expect(err).IsNil()
	Expect(actual).Equals(expected)
}

var declinedPost = &entities.Post{
	Number:      10,
	Title:       "Go is fast",
	Description: "Very tiny description",
	CreatedAt:   time.Date(2018, 3, 23, 19, 33, 22, 0, time.UTC),
	User: &entities.User{
		Name: "Faceless",
	},
	VotesCount:    4,
	CommentsCount: 2,
	Status:        enum.PostDeclined,
	Response: &entities.PostResponse{
		Text:        "Nothing we need to do",
		RespondedAt: time.Date(2018, 4, 4, 19, 48, 10, 0, time.UTC),
		User: &entities.User{
			Name: "John Snow",
		},
	},
	Tags: []string{"easy", "ignored"},
}

var openPost = &entities.Post{
	Number:      15,
	Title:       "Go is great",
	Description: "",
	CreatedAt:   time.Date(2018, 2, 21, 15, 51, 35, 0, time.UTC),
	User: &entities.User{
		Name: "Someone else",
	},
	VotesCount:    4,
	CommentsCount: 2,
	Status:        enum.PostOpen,
}

var duplicatePost = &entities.Post{
	Number:      20,
	Title:       "Go is easy",
	Description: "",
	CreatedAt:   time.Date(2018, 1, 12, 1, 46, 59, 0, time.UTC),
	User: &entities.User{
		Name: "Faceless",
	},
	VotesCount:    4,
	CommentsCount: 2,
	Status:        enum.PostDuplicate,
	Response: &entities.PostResponse{
		Text:        "This has already been suggested",
		RespondedAt: time.Date(2018, 3, 17, 10, 15, 42, 0, time.UTC),
		User: &entities.User{
			Name: "Arya Stark",
		},
		Original: &entities.OriginalPost{
			Number: 99,
			Title:  "Go is very easy",
		},
	},
	Tags: []string{"this-tag-has,comma"},
}
