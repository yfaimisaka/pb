package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"unsafe"

	"github.com/sirupsen/logrus"
)

var myLog *logrus.Logger

func init() {
    initLog()
}

func initLog() {
    myLog = logrus.New()
    env := os.Getenv("pbenv")
    myLog.Out = os.Stdout
    myLog.Level = logrus.DebugLevel
    if env == "production" {
        myLog.Level = logrus.WarnLevel
    } 
    myLog.SetFormatter(&logrus.TextFormatter{
        ForceColors: true,
        ForceQuote: true,    
        TimestampFormat: "2006-01-02 15:04:05",  
        FullTimestamp:true,    
	})
}


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

func long(content string) (string, error) {
    return unhexMixin(digest(stringToBytes(content)), 42)
}

func short(content string) (string, error) {
    return unhexMixin(digest(stringToBytes(content)), 6)
}

// useful utils to convert between string and byte
func bytesToString(bytes []byte) string {
    return unsafe.String(&bytes[0], len(bytes))
}

func stringToBytes(str string) []byte {
    return unsafe.Slice(unsafe.StringData(str), len(str))
}
