package helper

import (
	"log"
	"os/exec"
)

func Shell(cmd string) ([]byte, error) {
	log.Println(cmd)
	proc := exec.Command("sh", "-c", cmd)
	out, err := proc.Output()
	return out, err
}

// func GetFreeSpace(disk string) {
// }

// func GetFolderSize(disk string, folder string) {
// 	out, err := process(fmt.Sprintf("du -s %s", disk+folder+"/*"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
