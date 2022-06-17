package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
)

func NowDateTimeStr() string {
	return time.Now().Format(DateTimeFormat)
}

func NowDateStr() string {
	return time.Now().Format(DateFormat)
}

func AddDateStr(num int) string {
	return time.Now().Add(time.Duration(num) * time.Hour * 24).Format(DateFormat)
}

func ParseDateTime(value string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, value, time.Local)
}

func ParseDate(value string) (time.Time, error) {
	return time.ParseInLocation(DateFormat, value, time.Local)
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Size(filename string) (int64, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func MkdirIfNotExist(dirPath string) error {
	if Exist(dirPath) {
		return nil
	}
	return os.MkdirAll(dirPath, 0755)
}

func EncodeMD5(values string) string {
	m := md5.New()
	m.Write([]byte(values))
	return hex.EncodeToString(m.Sum(nil))
}

/*read cmd stdout, stderr result*/
func DoCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var out, errOut bytes.Buffer
	cmd.Stderr = &errOut
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "do cmd: %s error: %s", cmd.String(), errOut.String())
	}
	return nil
}

func CtxCmd(ctx context.Context, command string) error {
	if command == "" {
		return errors.New("do ctx command empty")
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	var out, errOut bytes.Buffer
	cmd.Stderr = &errOut
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "do cmd: %s error: %s", cmd.String(), errOut.String())
	}
	return nil
}

func Copy(src, dest string) error {
	cmd := fmt.Sprintf("cp %s %s", src, dest)
	return DoCmd("/bin/sh", "-c", cmd)
}

func RemoveReplica(src []string) (result []string) {
	length := len(src)
	if length == 0 {
		return
	}

	tempMap := make(map[string]bool, length)
	for _, e := range src {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return
}

func SplitNameExt(filePath string) (string, string) {
	if filePath == "" {
		return "", ""
	}
	baseName := path.Base(filePath)
	ext := path.Ext(filePath)
	name := baseName[0 : len(baseName)-len(ext)]
	return name, ext
}

const (
	ShortToDay     = "ToDay"
	ShortYesterDay = "YesterDay"
	ShortWeek      = "Week"
	ShortMonth     = "Month"
)

// ToDay 今天 YesterDay 昨天 Week 最近一周 最近月 Month
func ShortTimeToStartEnd(short string) (start, end string) {
	now := time.Now()
	switch short {
	case ShortToDay:
		start = now.Format(DateFormat)
		end = now.AddDate(0, 0, 1).Format(DateFormat)
	case ShortYesterDay:
		start = now.AddDate(0, 0, -1).Format(DateFormat)
		end = now.Format(DateFormat)
	case ShortWeek:
		start = now.AddDate(0, 0, -6).Format(DateFormat)
		end = now.Format(DateTimeFormat)
	case ShortMonth:
		start = now.AddDate(0, 0, -29).Format(DateFormat)
		end = now.Format(DateTimeFormat)
	}
	return
}

const CharNum = "Aa1Bb2Cc3Dd4Ee5Ff6Gg7Hh8Ii0JjKk9Ll8Mm7Nn6Oo5Pp4Qq3Rr2Ss1Tt2Uu3Vv4Ww0XxaYy1Zzb"
const CharNumLength = 77

func RandomStr(n int) string {
	var rest = make([]byte, n, n)
	for i := range rest {
		rest[i] = CharNum[rand.Int63()%CharNumLength]
	}
	return string(rest)
}

// unicodeo 转 utf8
func UnicodeToUtf8(source string) (string, error) {
	rs := []rune(source)
	var temp string

	for _, r := range rs {
		if int(r) < 128 {
			temp += string(r)
		} else {
			temp += "\\u" + strconv.FormatInt(int64(r), 16) // json
		}
	}

	rest, err := strconv.Unquote(`"` + temp + `"`)
	return rest, err
}
