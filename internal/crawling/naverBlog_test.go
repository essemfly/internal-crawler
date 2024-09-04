package crawling_test

import (
	"testing"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/stretchr/testify/assert"
)

func TestParseNaverMapUrl(t *testing.T) {

	blogId := "mardukas"

	// this is for naver map v1
	// logNo1 := "220436478787"
	// result1 := crawling.ParseNaverMapUrl(blogId, logNo1)

	// this is for naver map v2
	logNo2 := "222271782231"
	result2 := crawling.ParseNaverMapUrl(blogId, logNo2)

	// assert.Equal(t, []string{"https://map.naver.com/p/entry/place/11860445"}, result1)
	assert.Equal(t, []string{"https://map.naver.com/p/entry/place/19877488"}, result2)
}
