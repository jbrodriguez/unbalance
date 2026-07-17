package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/gookit/color"
)

const (
	tagInfo  = "[info]"
	tagWarn  = "[warn]"
	tagError = "[error]"
)

func Red(format string, args ...interface{}) {
	printer(color.Red, tagError, format, args...)
}

func Blue(format string, args ...interface{}) {
	printer(color.Blue, tagInfo, format, args...)
}

func Green(format string, args ...interface{}) {
	printer(color.Green, tagInfo, format, args...)
}

func Yellow(format string, args ...interface{}) {
	printer(color.Yellow, tagWarn, format, args...)
}

func LightGreen(format string, args ...interface{}) {
	printer2(color.S256(106), tagInfo, format, args...)
}

func LightRed(format string, args ...interface{}) {
	printer2(color.S256(9), tagError, format, args...)
}

func NotBetter(format string, args ...interface{}) {
	printer2(color.S256(132), tagInfo, format, args...)
}

func Olive(format string, args ...interface{}) {
	printer2(color.S256(11), tagInfo, format, args...)
}

func LightBlue(format string, args ...interface{}) {
	printer2(color.S256(14), tagInfo, format, args...)
}

func printer(fn color.Color, tag string, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	fn.Printf("%s %s\n", time.Now().Format("15:04"), line)
	log.Println(tag, line)
}

func printer2(fn *color.Style256, tag string, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	fn.Printf("%s %s\n", time.Now().Format("15:04"), line)
	log.Println(tag, line)
}
