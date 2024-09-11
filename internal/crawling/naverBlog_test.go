package crawling_test

import (
	"log"
	"testing"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/stretchr/testify/assert"
)

func TestParseNaverMapUrl(t *testing.T) {

	blogId := "mardukas"

	// this is for naver map v1
	logNo1 := "220436478787"
	result1 := crawling.ParseNaverMapUrl(blogId, logNo1)
	assert.Equal(t, []string{"https://map.naver.com/p/entry/place/11860445"}, result1)
	log.Println("Result1: ", result1)

	// this is for naver map v2
	logNo2 := "222271782231"
	result2 := crawling.ParseNaverMapUrl(blogId, logNo2)
	assert.Equal(t, []string{"https://map.naver.com/p/entry/place/19877488"}, result2)

	// this is for naver map between v1 and v2
	logNo3 := "221197009442"
	result3 := crawling.ParseNaverMapUrl(blogId, logNo3)
	log.Println("Result3: ", result3)
	assert.Equal(t, []string{"https://map.naver.com/p/entry/place/11888725"}, result3)

	logNo4 := "223358501558"
	result4 := crawling.ParseNaverMapUrl(blogId, logNo4)
	log.Println("Result4: ", result4)
	assert.Equal(t, []string(nil), result4)

	logNo5 := "222992469510"
	result5 := crawling.ParseNaverMapUrl(blogId, logNo5)
	log.Println("Result5: ", result5)
	assert.Equal(t, []string(nil), result5)
}
