package pb

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func digest(content []byte) string {
    hasher := sha1.New()
    hasher.Write(content)
    return hex.EncodeToString(hasher.Sum(nil))
}

func unhexMixin(digest string, length int) (string, error) {
    // got slice of digest, len=length if length < len(digest). 
    // if length > len(digest), then padding use 0
    if len(digest) > length {
        digest = digest[len(digest)-length:]
    }
    // + means right align, 0 means use 0 to pad, 
    // * means using the first argument(in here is length) as the number of digits
    digestPaded := fmt.Sprintf("%+0*s", length, digest) 
    digestPadedUnhexed, err := hex.DecodeString(digestPaded)
    if err != nil {
        return "", err
    }
    res := base64.URLEncoding.EncodeToString(digestPadedUnhexed)
    return res, nil
}

func long(digest string) (string, error) {
    return unhexMixin(digest, 42)
}

func short(digest string) (string, error) {
    return unhexMixin(digest, 6)
}

