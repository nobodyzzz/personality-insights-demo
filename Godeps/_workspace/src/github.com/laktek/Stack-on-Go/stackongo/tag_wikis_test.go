package stackongo

import (
	"testing"
)

func TestWikisForTags(t *testing.T) {
	dummy_server := returnDummyResponseForPath("/2.0/tags/tag1;tag2;tag3/wikis", dummyTagWikisResponse, t)
	defer closeDummyServer(dummy_server)

	session := NewSession("stackoverflow")
	tag_wikis, err := session.WikisForTags([]string{"tag1", "tag2", "tag3"}, map[string]string{"sort": "votes", "order": "desc", "page": "1"})

	if err != nil {
		t.Error(err.Error())
	}

	if len(tag_wikis.Items) != 1 {
		t.Error("Number of items wrong.")
	}

	if tag_wikis.Items[0].Tag_name != "go" {
		t.Error("Tag name invalid.")
	}

	if tag_wikis.Items[0].Body_last_edit_date != 1322081597 {
		t.Error("last edit date invalid.")
	}

}

//Test Data

var dummyTagWikisResponse string = `
{
  "items": [
    {
      "tag_name": "go",
      "excerpt": "Go is a general-purpose programming language designed by Google.",
      "body_last_edit_date": 1322081597,
      "excerpt_last_edit_date": 1322081452
    }
  ]
}
`
