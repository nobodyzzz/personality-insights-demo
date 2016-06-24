package stackongo

import (
	"testing"
)

func TestAllTags(t *testing.T) {
	dummy_server := returnDummyResponseForPath("/2.0/tags", dummyTagsResponse, t)
	defer closeDummyServer(dummy_server)

	session := NewSession("stackoverflow")
	tags, err := session.AllTags(map[string]string{"sort": "votes", "order": "desc", "page": "1"})

	if err != nil {
		t.Error(err.Error())
	}

	if len(tags.Items) != 3 {
		t.Error("Number of items wrong.")
	}

	if tags.Items[0].Name != "c#" {
		t.Error("Name invalid.")
	}

	if tags.Items[0].Count != 261768 {
		t.Error("Tag count invalid.")
	}

	if tags.Items[0].Has_synonyms != true {
		t.Error("boolean invalid.")
	}

}

func TestTagsForUsers(t *testing.T) {
	dummy_server := returnDummyResponseForPath("/2.0/users/1;2;3/tags", dummyTagsResponse, t)
	defer closeDummyServer(dummy_server)

	session := NewSession("stackoverflow")
	_, err := session.TagsForUsers([]int{1, 2, 3}, map[string]string{"sort": "votes", "order": "desc", "page": "1"})

	if err != nil {
		t.Error(err.Error())
	}

}

func TestRelatedTags(t *testing.T) {
	dummy_server := returnDummyResponseForPath("/2.0/tags/tag1;tag2;tag3/tags", dummyTagsResponse, t)
	defer closeDummyServer(dummy_server)

	session := NewSession("stackoverflow")
	_, err := session.RelatedTags([]string{"tag1", "tag2", "tag3"}, map[string]string{"sort": "votes", "order": "desc", "page": "1"})

	if err != nil {
		t.Error(err.Error())
	}

}

//Test Data

var dummyTagsResponse string = `
{
  "items": [
    {
      "name": "c#",
      "count": 261768,
      "is_required": false,
      "is_moderator_only": false,
      "has_synonyms": true
    },
    {
      "name": "java",
      "count": 202323,
      "is_required": false,
      "is_moderator_only": false,
      "has_synonyms": true
    },
    {
      "name": "php",
      "count": 187267,
      "is_required": false,
      "is_moderator_only": false,
      "has_synonyms": true
    }
  ]
}
`
