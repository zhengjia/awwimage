package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var page1 = `
{
  "meta": {
    "status": 200,
    "msg": "OK"
  },
  "response": [
    {
      "blog_name": "marangio",
      "id": 111415925944,
      "post_url": "http://marangio.tumblr.com/post/111415925944/lbungeejumping",
      "slug": "lbungeejumping",
      "type": "photo",
      "date": "2015-02-18 23:44:08 GMT",
      "timestamp": 1424303048,
      "state": "published",
      "format": "html",
      "reblog_key": "J558MRGe",
      "tags": [
        "bungeejumping",
        "jump",
        "mouse",
        "giraffe",
        "illustrator",
        "illustration",
        "ilustração",
        "ilustración",
        "abbildung",
        "graphic",
        "graphicart",
        "graphicillustration",
        "graphicline",
        "graphicdraw",
        "graphicdesign",
        "art",
        "draw",
        "line"
      ],
      "short_url": "http://tmblr.co/ZVKBKm1dmw1ou",
      "highlighted": [
      ],
      "note_count": 0,
      "caption": "<p>LBUNGEEJUMPING</p>",
      "image_permalink": "http://marangio.tumblr.com/image/111415925944",
      "photos": [
        {
          "caption": "",
          "alt_sizes": [
            {
              "width": 1280,
              "height": 1280,
              "url": "http://41.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_1280.jpg"
            },
            {
              "width": 500,
              "height": 500,
              "url": "http://36.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_500.jpg"
            },
            {
              "width": 400,
              "height": 400,
              "url": "http://36.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_400.jpg"
            },
            {
              "width": 250,
              "height": 250,
              "url": "http://40.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_250.jpg"
            },
            {
              "width": 100,
              "height": 100,
              "url": "http://40.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_100.jpg"
            },
            {
              "width": 75,
              "height": 75,
              "url": "http://40.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_75sq.jpg"
            }
          ],
          "original_size": {
            "width": 1280,
            "height": 1280,
            "url": "http://41.media.tumblr.com/c8dcfcb11f07801f41db2d2ff46b13ba/tumblr_njzr9kqbUI1u0ep4qo1_1280.jpg"
          }
        }
      ]
    },
    {
      "blog_name": "piecomic",
      "id": 111409833242,
      "post_url": "http://piecomic.tumblr.com/post/111409833242",
      "slug": "",
      "type": "photo",
      "date": "2015-02-18 22:30:13 GMT",
      "timestamp": 1424298613,
      "state": "published",
      "format": "html",
      "reblog_key": "bcMPdZhi",
      "tags": [
        "cartoon",
        "lol",
        "animals",
        "giraffe"
      ],
      "short_url": "http://tmblr.co/ZU7Znx1dmYoKQ",
      "highlighted": [
      ],
      "note_count": 165,
      "source_url": "http://www.piecomic.com",
      "source_title": "piecomic.com",
      "caption": "",
      "link_url": "http://www.piecomic.com",
      "image_permalink": "http://piecomic.tumblr.com/image/111409833242",
      "photos": [
        {
          "caption": "",
          "alt_sizes": [
            {
              "width": 500,
              "height": 552,
              "url": "http://40.media.tumblr.com/dfbaff69a8390c1bc3c048a2c41b7ab0/tumblr_njznudWb0D1qhnegdo1_500.jpg"
            },
            {
              "width": 400,
              "height": 442,
              "url": "http://40.media.tumblr.com/dfbaff69a8390c1bc3c048a2c41b7ab0/tumblr_njznudWb0D1qhnegdo1_400.jpg"
            },
            {
              "width": 250,
              "height": 276,
              "url": "http://40.media.tumblr.com/dfbaff69a8390c1bc3c048a2c41b7ab0/tumblr_njznudWb0D1qhnegdo1_250.jpg"
            },
            {
              "width": 100,
              "height": 110,
              "url": "http://40.media.tumblr.com/dfbaff69a8390c1bc3c048a2c41b7ab0/tumblr_njznudWb0D1qhnegdo1_100.jpg"
            },
            {
              "width": 75,
              "height": 75,
              "url": "http://40.media.tumblr.com/dfbaff69a8390c1bc3c048a2c41b7ab0/tumblr_njznudWb0D1qhnegdo1_75sq.jpg"
            }
          ],
          "original_size": {
            "width": 500,
            "height": 552,
            "url": "http://40.media.tumblr.com/dfbaff69a8390c1bc3c048a2c41b7ab0/tumblr_njznudWb0D1qhnegdo1_500.jpg"
          }
        }
      ]
    }
  ]
}
`

type testFetcher struct{}

func (fetcher testFetcher) Fetch(url string) ([]byte, error) {
	return []byte(page1), nil
}

func reset_image_mapping() {
	image_mapping = make(map[string][]string)
}

type ResponseJson struct {
	Url string
}

func TestPopulateImageMapping(t *testing.T) {
	var kind = "pug"
	var image_limit_was = image_limit
	image_limit = 20
	fetcher := testFetcher{}
	populate(kind, &fetcher)
	if _, ok := image_mapping[kind]; !ok {
		t.Error("populate failed")
	}
	if len(image_mapping[kind]) != 20 {
		t.Error("populate count failed")
	}
	log.Println(image_limit_was)
	image_limit = image_limit_was
	reset_image_mapping()
}

func TestRandom(t *testing.T) {
	var url = "http://example.com/1.jpg"
	image_mapping["pug"] = []string{url}
	ts := httptest.NewServer(getHttpHandler())
	defer ts.Close()
	res, err := http.Get(ts.URL + "/random/pug")
	if err != nil {
		t.Error("random failed", err)
	}
	if res.StatusCode != 200 {
		t.Error("random failed", res.StatusCode)
	}
	body_bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	response_json := ResponseJson{}
	json.Unmarshal(body_bytes, &response_json)
	if response_json.Url != url {
		t.Error("random failed", response_json.Url)
	}
	reset_image_mapping()
}

func TestRandomPreview(t *testing.T) {
	var url = "http://example.com/1.jpg"
	image_mapping["pug"] = []string{url}
	ts := httptest.NewServer(getHttpHandler())
	defer ts.Close()
	res, _ := http.Get(ts.URL + "/random/pug/preview")
	body_bytes, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(body_bytes) != "<html><img src='http://example.com/1.jpg' /></html>" {
		t.Error("random preview failed", string(body_bytes))
	}
	reset_image_mapping()
}
