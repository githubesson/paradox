package decode

import (
	"crypto/sha1"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func DeriveMasterKey(systemPassword []byte) []byte {
	chromeSalt := []byte("saltysalt")
	secretKey := pbkdf2.Key(systemPassword, chromeSalt, 1003, 16, sha1.New)

	return secretKey
}

func IntToBool(a int) bool {
	switch a {
	case 0, -1:
		return false
	}
	return true
}

func TimeEpochFormat(epoch int64) time.Time {
	maxTime := int64(99633311740000000)
	if epoch > maxTime {
		return time.Date(2049, 1, 1, 1, 1, 1, 1, time.Local)
	}
	t := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	d := time.Duration(epoch)
	for i := 0; i < 1000; i++ {
		t = t.Add(d)
	}
	return t
}
