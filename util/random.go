package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var rng *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(source)
}

// 返回一个介于 min max 之间的随机的 int64 数字
func RandomInt(min, max int64) int64 {
	return min + rng.Int63n(max-min+1)
}

// 生成 n 个字符的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rng.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// 随机生成owner
func RandomOwner() string {
	return RandomString(6)
}

// 随机生成钱的数量
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// 随机产生一种货币
func RandomCurrency() string {
	currencies := []string{"RMB", "USD", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
