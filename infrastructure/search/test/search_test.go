package search

import (
	"testing"

	"github.com/expectedsh/go-sonic/sonic"
)

const pswd = "SecretPassword"

func TestSearch(t *testing.T) {
	ingester, err := sonic.NewIngester("localhost", 1491, pswd)
	if err != nil {
		panic(err)
	}

	// I will ignore all errors for demonstration purposes

	_ = ingester.BulkPush("movies", "general", 3, []sonic.IngestBulkRecord{
		{Object: "id:6ab56b4kk3", Text: "世界大戰"},
		{Object: "id:5hg67f8dg5", Text: "蜘蛛 俠"},
		{Object: "id:1m2n3b4vf6", Text: "蝙蝠俠"},
		{Object: "id:68d96h5h9d0", Text: "這是另一部電影"},
	}, sonic.LangAutoDetect)

	search, err := sonic.NewSearch("localhost", 1491, pswd)
	if err != nil {
		panic(err)
	}

	results, _ := search.Query("movies", "general", "俠", 10, 0, sonic.LangAutoDetect)

	t.Log(results)

	// Search with LANG set to "none" and "eng"

	_ = ingester.FlushCollection("movies")
	_ = ingester.BulkPush("movies", "general", 3, []sonic.IngestBulkRecord{
		{Object: "id:6ab56b4kk3", Text: "世界大戰"},
		{Object: "id:5hg67f8dg5", Text: "蜘蛛 俠"},
		{Object: "id:1m2n3b4vf6", Text: "蝙蝠俠"},
		{Object: "id:68d96h5h9d0", Text: "這是另一部電影"},
	}, sonic.LangNone)

	results, _ = search.Query("movies", "general", "這是", 10, 0, sonic.LangNone)
	t.Log(results)
	// [id:68d96h5h9d0]

	// English stop words should be encountered by Sonic now
	results, _ = search.Query("movies", "general", "這是", 10, 0, sonic.LangEng)
	t.Log(results)


	_ = ingester.FlushCollection("w")
	_ = ingester.BulkPush("w", "words", 3, []sonic.IngestBulkRecord{
		{Object: "id:6ab56b4kk3", Text: "put-si"},
		{Object: "id:5hg67f8dg5", Text: "kau-si"},
		{Object: "id:1m2n3b4vf6", Text: "ping-si在kau-si"},
		{Object: "id:68d96h5h9d0", Text: "kau si put si"},
	}, sonic.LangNone)

	results, _ = search.Query("w", "words", "put-si", 10, 0, sonic.LangNone)
	t.Log(results)

	results, _ = search.Query("w", "words", "put si", 10, 0, sonic.LangEng)
	t.Log(results)
}
