package secretkey

import (
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"niuNiuSDKBackend/common/log"
	"time"

	"net/http"
)

type SecretKey struct {
	Id          int32     `json:"id"  xorm:"'id' pk autoincr BIGINT(20)" `
	SK          string    `json:"sk" xorm:"'sk' VARCHAR(70)"`
	CreatedTime time.Time `json:"created_time" xorm:"'created_time' created"`
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated"`
	DeletedTime time.Time `json:"deleted_time" xorm:"'deleted_time' deleted"`
}

const SKTABLE = "sk"
const StdLen = 64

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewSk(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)
	sk := SecretKey{
		SK: NewLenChars(StdLen, StdChars),
	}

	if _, err := db.Table(SKTABLE).Insert(&sk); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "get sk failed", "code": 501})
		log.Logger.Error("get sk failed", log.Any("get sk error", "get sk error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": sk,
		"code":    200,
	})
	return
}

// NewLenChars returns a new random string of the provided length, consisting of the provided byte slice of allowed characters(maximum 256).
func NewLenChars(length int, chars []byte) string {
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue // Skip this number to avoid modulo bias.
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}
