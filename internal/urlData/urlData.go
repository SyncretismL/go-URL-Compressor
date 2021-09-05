package urlData

import (
	"fmt"
	"math"
	"strings"
)

type URLData struct {
	ID            int64
	URL           string
	URLCompressed string
}

type URLDatas interface {
	SetURL(url *URLData) error
	SetURLCompressed(url *URLData) error
	GetFullURL(url *URLData) error
}

func (url *URLData) URLCompressing() {
	alph := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	digits := make([]int, 0, 0)

	num := float64(url.ID)
	for z := 0; z < 10; z++ {
		remainder := math.Mod(num, 63.0)
		digits = append(digits, int(remainder))
		num = num / 63.0

	}

	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	for _, v := range digits {
		url.URLCompressed += fmt.Sprintf("%c", alph[v])
	}

	url.ReplaceQuery()
}

func (url *URLData) ReplaceQuery() {
	separated := strings.Split(url.URL, "/")
	separated[len(separated)-1] = url.URLCompressed
	url.URLCompressed = strings.Join(separated, "/")
}
