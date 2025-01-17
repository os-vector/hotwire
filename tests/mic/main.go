package main

import (
	"encoding/binary"
	"fmt"
	"hotwire/pkg/audioproc"
	"hotwire/pkg/stt"
	"log"
	"os"

	"github.com/gordonklaus/portaudio"
)

func startMicStream() (chan []byte, error) {
	out := make(chan []byte, 32)
	err := portaudio.Initialize()
	if err != nil {
		return nil, err
	}

	// we want 160 frames each read (mono, 16-bit, 16khz -> 320 bytes)
	numFrames := 160
	in := make([]int16, numFrames)

	stream, err := portaudio.OpenDefaultStream(
		1,
		0,
		16000,
		numFrames,
		in,
	)
	if err != nil {
		return nil, err
	}

	err = stream.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(out)
		defer stream.Stop()
		defer stream.Close()
		for {
			err := stream.Read()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			buf := make([]byte, 2*numFrames)
			for i, sample := range in {
				binary.LittleEndian.PutUint16(buf[i*2:], uint16(sample))
			}
			out <- buf
		}
	}()

	return out, nil
}

func main() {
	ch, err := startMicStream()
	origFile, _ := os.Create("./test_unprocessed.pcm")
	processedFile, _ := os.Create("test_processed.pcm")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening for mic data at 16khz, 320 bytes per chunk...")
	spStr := stt.Init("test")
	aProc, _ := audioproc.NewAudioProcessor(16000, 200.0, 3)

	for chunk := range ch {
		origFile.Write(chunk)
		decoded := aProc.ProcessAudio(chunk)
		processedFile.Write(decoded)
		if spStr.DetectEndOfSpeech(decoded) {
			origFile.Close()
			processedFile.Close()
			return
		}
	}
}

func bytesToInt16(b []byte) []int16 {
	n := len(b) / 2
	out := make([]int16, n)
	for i := 0; i < n; i++ {
		out[i] = int16(b[2*i]) | int16(b[2*i+1])<<8
	}
	return out
}
