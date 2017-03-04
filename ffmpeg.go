package mkv2Appletv

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/fatih/color"
)

func checkVersion() {
	out, err := exec.Command("ffmpeg", "-version").CombinedOutput()
	if err != nil {
		// Not found
		fmt.Println(err)
		fmt.Println("Unable to find ffmpeg in your path")
	} else {
		// Find out what version of ffmpeg that is installed
		version := regexp.MustCompile(`ffmpeg.version.\d\S*`)
		re := regexp.MustCompile("ffmpeg.version 3")

		bFound := re.MatchString(string(out[:30]))
		if bFound == false {
			//found but wrong version
			fmt.Printf("Requires ffmpeg >= 3.x.x I found version %s", version)
			fmt.Println("Go here to download binaries for your machine: https://ffmpeg.org/download.html")
			fmt.Println("I would recommend the static compiled version")
			err = errors.New("Wrong version of ffmpeg found ")
		} else {
			fmt.Printf("Found the correct version of ffmpeg\n")
		}
		// found and right version
		fmt.Printf("%s\n", out)
	}

	// Now check the ffprobe version
	probeout, err := exec.Command("ffprobe", "-version").CombinedOutput()
	if err != nil {
		// Not found
		fmt.Println(err)
		fmt.Println("Unable to find ffmpeg in your path")
	} else {
		// Find out what version of ffmpeg that is installed
		version := regexp.MustCompile(`ffprobe.version.\d\S*`)
		re := regexp.MustCompile("ffprobe.version 3")

		bFound := re.MatchString(string(probeout[:30]))
		if bFound == false {
			//found but wrong version
			fmt.Printf("Requires ffprobe >= 3.x.x I found version %s", version)
			fmt.Println("Go here to download binaries for your machine: https://ffmpeg.org/download.html")
			fmt.Println("I would recommend the static compiled version")
			err = errors.New("Wrong version of ffmpeg found ")
		} else {
			fmt.Println("Found the correct version of ffprobe")
		}
		fmt.Printf("%s", probeout)
	}
}
func checkFFmpegVersion() error {

	out, err := exec.Command("ffmpeg", "-version").CombinedOutput()
	if err != nil {
		// Not found
		fmt.Println(err)
	} else {
		// Find out what version of ffmpeg that is installed
		version := regexp.MustCompile(`ffmpeg.version.\d\S*`)
		re := regexp.MustCompile("ffmpeg.version 3")

		bFound := re.MatchString(string(out[:30]))
		if bFound == false {
			//found but wrong version
			fmt.Printf("Requires ffmpeg >= 3.x.x I found version %s", version)
			fmt.Println("Go here to download binaries for your machine: https://ffmpeg.org/download.html")
			fmt.Println("I would recommend the static compiled version")
			err = errors.New("Wrong version of ffmpeg found ")
		}
		// found and right version
	}
	return err
}

// add stats: -stats
func callFFmpeg(ffmpegCmd *ffmpegOut) (string, error) {

	color.Yellow("Starting conversion")
	cmd := exec.Command("ffmpeg", ffmpegCmd.ffArgs...)
	//Debug
	//fmt.Printf("\n%#v\n", cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return "Success", err
}
