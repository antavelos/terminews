package main

import (
	"github.com/antavelos/terminews/db"
	"testing"
)

func TestGetContentURL(t *testing.T) {
	type test struct {
		site     db.Site
		url      string
		expected string
	}

	tests := []test{
		{
			site:     db.Site{Url: "https://blog.example.org/index.xml"},
			url:      "/2021/02/a-blog-post",
			expected: "https://blog.example.org/2021/02/a-blog-post",
		},
		{
			site:     db.Site{Url: "http://blog.example.org/index.xml"},
			url:      "2021/02/a-blog-post",
			expected: "http://blog.example.org/2021/02/a-blog-post",
		},
		{
			site:     db.Site{Url: "https://blog.example.org/feed.rss"},
			url:      "https://blog.example.org/super-blog-post",
			expected: "https://blog.example.org/super-blog-post",
		},
		{
			site:     db.Site{Url: "http://blog.example.org/feed.rss"},
			url:      "http://blog.example.org/super-blog-post",
			expected: "http://blog.example.org/super-blog-post",
		},
	}

	for _, test := range tests {
		value := getContentURL(test.site, test.url)
		if value != test.expected {
			t.Errorf("got: %s want: %s", value, test.expected)
		}
	}
}
