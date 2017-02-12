package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

var (
	ffmpegCmd = new(ffmpegOut)
)

type ffmpegOut struct {
	outFile string
	header  string
	Video   string
	Audio0  string
	Audio1  string
	ffArgs  []string // This will hold the complete command line sent to ffmpeg
}

func (buff *ffmpegOut) genVideoConversion() error {

	temp := fmt.Sprintf(" -map 0:%d -c:v", media.masterVideoStream.Index)
	switch media.outVideo {

	case "copy":
		buff.Video = fmt.Sprintf(" copy")
	case "convert":
		buff.Video = fmt.Sprintf("libx264 -preset slow -crf 20 -profile:v high -level 4.0 ")
	default:
		err := errors.New("unknown or not set Video settings\n")
		return err
	} // end of switch
	ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, temp, buff.Video)
	return err
}

func (buff *ffmpegOut) genAudioConversion() error {

	//
	// This is the best case senerio.  We have an aac and an ac3 stream
	//
	if media.outAudio0 == "copy" && media.outAudio1 == "copy" {
		map0 := fmt.Sprintf("0:%d", media.aacAudioStream.Index)
		buff.Audio0 = fmt.Sprintf("-c:a:0")
		map1 := fmt.Sprintf("0:%d", media.masterAudioStream.Index)
		buff.Audio1 = fmt.Sprintf("-c:a:1")
		ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, " -map", map0, buff.Audio0, "copy", "-map", map1, buff.Audio1, "copy")

		return err
	} else if media.outAudio0 == "convert" && media.outAudio1 == "convert" {
		//Handle situation where all codecs need to be generated
		//TODO Handle audio1 situation
		if media.masterAudioStream.Channels > 2 {
			//TODO: this is broken
			// Output a 2 channel aac stream from the master audio which is a multichannel audio
			// the asplit filter takes the input and splits it into 2 dual streams.  one called 2ch and another called 6ch
			// The 2ch is fed into the pan filter and the output is placed into the "aac" pad
			//
			buff.Audio0 = fmt.Sprintf("-filter_complex [0:%d]asplit[2ch][6ch];[2ch]pan=stereo|FL=FC+0.6FL+0.2BL|FR=FC+0.6FR+0.2BR[aac]", media.masterAudioStream.Index)
			ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, buff.Audio0)
			// Now append the output pad mappings [aac] and [6ch]
			// if the 6ch is NOT ac3 we will have to transcode it.
			ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-map [aac]", "-map", "[6ch]", "-c:a:0", "aac", "-c:a:1", "ac3")
		} else {
			// if the master audio is NOT aac and is only 2 channel
			buff.Audio0 = fmt.Sprintf("-map 0:%d -c:a:0 aac -b:a 256k", media.masterAudioStream.Index)
			ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, buff.Audio0)
		}
		return err
	}

	switch media.outAudio0 {
	case "copy":
		// This is a situation where the only audio available is an aac 2 channel
		map0 := fmt.Sprintf("0:%d", media.masterAudioStream.Index)
		buff.Audio0 = fmt.Sprintf("-c:a:0")
		ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-map", map0, buff.Audio0, "copy")

	case "convert":
		// we need to figure out if this is multichannel
		map0 := fmt.Sprintf("0:%d", media.masterAudioStream.Index)
		ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-map", map0, "-c:a:0", "aac", "-b:a:0", "256k")

	} // end of switch

	switch media.outAudio1 {
	case "copy":
		// This situation is when we have an ac3 codec and we have to generate the aac codec
		map1 := fmt.Sprintf("0:%d", media.masterAudioStream.Index)
		ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-map", map1, "-c:a:1", "copy")
	case "convert":
		//this is currently handled in case where numb audio streams > 2
		// we need to figure out if this is multichannel
		//		fmt.Println("Not converting the second audio yet, this is not yet implemented")
		//		buff.Audio1 = fmt.Sprintf("-c:a:1 copy ")

	} // end of switch
	return err
}
func (buff *ffmpegOut) setupHeader() {

	ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-map_metadata 0:g")
	if *try == true {
		ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-t 00:00:10")
	}

}

func convertSource(in string, output string) {

	suggestConvSettings(in)
	// Media object is now setup

	// Do we need to handle additional debugging?
	if *debug == true {
		ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-v trace")
	}
	temp := fmt.Sprintf(" -i %s", in)
	ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, "-hide_banner", "-y", temp)

	if output != "" {
		ffmpegCmd.outFile = output
	} else {
		ffmpegCmd.outFile = fmt.Sprintf(" %s.mp4", in)
	}

	ffmpegCmd.setupHeader()
	ffmpegCmd.genVideoConversion()
	ffmpegCmd.genAudioConversion()
	ffmpegCmd.ffArgs = append(ffmpegCmd.ffArgs, ffmpegCmd.outFile)

	// Format String to send to ffmpeg
	fmt.Printf("%v", ffmpegCmd.ffArgs)

	err := checkFFmpegVersion()
	if err != nil {
		fmt.Println("Not sending commands to ffmpeg because: %s\n", err)
		return
	}

	color.Blue("\n\nType: %T\n%#v\n\n", ffmpegCmd.ffArgs, ffmpegCmd.ffArgs)

	for i := 0; i < len(ffmpegCmd.ffArgs); i++ {
		fmt.Printf("%s\n", ffmpegCmd.ffArgs[i])
	}

	_, err = callFFmpeg(ffmpegCmd)
	if err != nil {
		fmt.Println("Error executing ffmpeg call\n")
	}

}

//   works
// 	cmd := exec.Command("ffmpeg", "-hide_banner", "-y", "-i", "/Users/snoby/Public/public/JL.mkv",
// 		"-t", "00:00:10", "-report", "-loglevel", "verbose",
// 		"-map", "0:0", "-map", "0:1", "-map", "0:1",
// 		"-c:v", "copy",
// 		"-c:a:0", "aac", "-b:a:0", "256k",
// 		"-c:a:1", "copy",
// 		"/Users/snoby/result.mp4")
//
// 	cmd := exec.Command("ffmpeg", ffmpegCmd.ffArgs...)
// 	fmt.Printf("\n%v\n", cmd)
//
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
//
// 	cmd.Run()
