package parser_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/lestrrat/go-graphql/format"
	"github.com/lestrrat/go-graphql/parser"
	"github.com/stretchr/testify/assert"
)

func parseSuccess(src string) (string, func(*testing.T)) {
	return src, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		p := parser.New()
		doc, err := p.Parse(ctx, []byte(src))
		if !assert.NoError(t, err, "parser.Parse should be successful") {
			return
		}

		var buf bytes.Buffer
		if !assert.NoError(t, format.GraphQL(&buf, doc), "format.GraphQL should be successful") {
			return
		}

		var expected string
		if strings.HasPrefix(src, "{") {
			expected = "query " + src
		} else {
			expected = src
		}

		if !assert.Equal(t, expected, buf.String(), "formatted code should be identical") {
			t.Logf("%s", buf.String())
			return
		}
	}
}

func TestParse(t *testing.T) {
	t.Run(parseSuccess(`query {
  me {
    name
  }
}`))
	t.Run(parseSuccess(`query HeroNameAndFriends($episode: Episode) {
  hero(episode: $episode) {
    name
    friends {
      name
    }
  }
}`))
	t.Run(parseSuccess(`{
  human(id: "1000") {
    name
    height
  }
}`))
	t.Run(parseSuccess(`{
  human(id: "1000") {
    name
    height(unit: FOOT)
  }
}`))
	t.Run(parseSuccess(`{
  empireHero: hero(episode: EMPIRE) {
    name
  }
  jediHero: hero(episode: JEDI) {
    name
  }
}`))
	t.Run(parseSuccess(`{
  leftComparison: hero(episode: EMPIRE) {
    ...comparisonFields
  }
  rightComparison: hero(episode: JEDI) {
    ...comparisonFields
  }
}

fragment comparisonFields on Character {
  name
  appearsIn
  friends {
    name
  }
}`))
	t.Run(parseSuccess(`query Hero($episode: Episode, $withFriends: Boolean!) {
  hero(episode: $episode) {
    name
    friends @include(if: $withFriends) {
      name
    }
  }
}`))
	t.Run(parseSuccess(`mutation CreateReviewForEpisode($ep: Episode!, $review: ReviewInput!) {
  createReview(episode: $ep, review: $review) {
    stars
    commentary
  }
}`))
t.Run(parseSuccess(`query HeroForEpisode($ep: Episode!) {
  hero(episode: $ep) {
    name
    ... on Droid {
      primaryFunction
    }
    ... on Human {
      height
    }
  }
}`))
	t.Run(parseSuccess(`{
  search(text: "an") {
    __typename
    ... on Human {
      name
    }
    ... on Droid {
      name
    }
    ... on Starship {
      name
    }
  }
}`))
}
 
