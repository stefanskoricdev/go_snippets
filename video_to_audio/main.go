package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Stream = ffmpeg.Stream

func main() {
	inputFile := "input_file_name.mp4"
	outputFile := "output_file_name.mp3"

	//Opts resource https://stackoverflow.com/questions/9913032/how-can-i-extract-audio-from-video-with-ffmpeg?rq=4
	cmdOpts := map[string]interface{}{
		"q:a": "0", "map": "a",
	}

	err := ffmpeg.Input(inputFile).Output(outputFile, ffmpeg.KwArgs(cmdOpts)).OverWriteOutput().Run()
	if err != nil {
		panic(err)
	}
}
