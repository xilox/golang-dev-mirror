package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"net/url"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/hr3lxphr6j/bililive-go/src/instance"
)

func GetFFmpegPath() (string, error) {
	path, err := exec.LookPath("ffmpeg")
	if errors.Is(err, exec.ErrDot) {
		// put ffmpeg.exe and binary like bililive-windows-amd64.exe to the same folder is allowed
		path, err = exec.LookPath("./ffmpeg")
	}
	return path, err
}

func IsFFmpegExist() bool {
	_, err := GetFFmpegPath()
	return err == nil
}

func GetMd5String(b []byte) string {
	md5Obj := md5.New()
	md5Obj.Write(b)
	return hex.EncodeToString(md5Obj.Sum(nil))
}

var (
	lowercaseRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	uppercaseRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lettersRunes   = append(lowercaseRunes, uppercaseRunes...)
	digitsRunes    = []rune("0123456789")
	allRunes       = append(lettersRunes, digitsRunes...)
)

func GenRandomName(n int) string {
	b := make([]rune, n)
	b[0] = lowercaseRunes[rand.Intn(len(lowercaseRunes))]
	for i := 1; i < n; i++ {
		b[i] = allRunes[rand.Intn(len(allRunes))]
	}
	return string(b)
}

func GenRandomString(length int, validChars string) string {
	b := make([]string, length)
	chars := strings.Split(validChars, "")
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return strings.Join(b, "")
}

func Match1(re, str string) string {
	reg, err := regexp.Compile(re)
	if err != nil {
		return ""
	}
	match := reg.FindStringSubmatch(str)
	if match == nil || len(match) < 2 {
		return ""
	}
	return match[1]
}

func GenUrls(strs ...string) ([]*url.URL, error) {
	urls := make([]*url.URL, 0, len(strs))
	for _, str := range strs {
		u, err := url.Parse(str)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

func PrintStack(ctx context.Context) {
	inst := instance.GetInstance(ctx)
	logger := inst.Logger
	logger.Debugf(string(debug.Stack()))
}
