package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/gookit/color"
)

func Red(format string, args ...interface{}) {
	printer(color.Red, format, args...)
}

func Blue(format string, args ...interface{}) {
	printer(color.Blue, format, args...)
}

func Green(format string, args ...interface{}) {
	printer(color.Green, format, args...)
}

func Yellow(format string, args ...interface{}) {
	printer(color.Yellow, format, args...)
}

func LightGreen(format string, args ...interface{}) {
	printer2(color.S256(106), format, args...)
}

func LightRed(format string, args ...interface{}) {
	printer2(color.S256(9), format, args...)
}

func NotBetter(format string, args ...interface{}) {
	printer2(color.S256(132), format, args...)
}

func Olive(format string, args ...interface{}) {
	printer2(color.S256(11), format, args...)
}

func LightBlue(format string, args ...interface{}) {
	printer2(color.S256(14), format, args...)
}

func printer(fn color.Color, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	fn.Printf("%s %s\n", time.Now().Format("15:04"), line)
	log.Println(line)
}

func printer2(fn *color.Style256, format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	fn.Printf("%s %s\n", time.Now().Format("15:04"), line)
	log.Println(line)
}
