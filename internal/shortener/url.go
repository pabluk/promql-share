package url

import (
	"bytes"
	"github.com/go-redis/redis"
)

func encode(nb int64, buf *bytes.Buffer, base string) {
	l := int64(len(base))
	if nb/l != 0 {
		encode(nb/l, buf, base)
	}
	buf.WriteByte(base[nb%l])
}

func generateShareID(urlData string, serverAddr string) (shareID string) {
	var buf bytes.Buffer
	client := redis.NewClient(&redis.Options{
		Addr:     serverAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := client.Incr("counter").Result()
	if err != nil {
		panic(err)
	}
	encode(val, &buf, "h6sCQgcPqJNrSzlWn5TfeBK8HxiY1Z7kIap3yXODMbvVt0m2udjRU4GEFoL9Aw")
	shareID = buf.String()

	err = client.Set(shareID, urlData, 0).Err()
	if err != nil {
		panic(err)
	}

	return shareID
}
