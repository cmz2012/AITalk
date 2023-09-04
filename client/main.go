package main

import (
	"encoding/binary"
	"fmt"
	"github.com/go-audio/aiff"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/google/uuid"
	"github.com/gordonklaus/portaudio"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8888/chat", nil)
	if err != nil {
		logrus.Errorf("%v", err)
		return
	}
	for {
		// 生成wav
		name := _aiff()
		wavName := aiff2wav(name)
		bs, err := ioutil.ReadFile(wavName)
		if err != nil {
			logrus.Errorf("%v", err)
			break
		}
		// send to server
		err = c.WriteMessage(websocket.BinaryMessage, bs)
		if err != nil {
			logrus.Errorf("%v", err)
			break
		}

		// receive reply
		_, bs, err = c.ReadMessage()
		if err != nil {
			logrus.Errorf("%v", err)
			break
		}
		// write reply file
		replyName := "./tmp/" + time.Now().Format(time.Kitchen) + "-reply.wav"
		err = ioutil.WriteFile(replyName, bs, fs.ModePerm)
		if err != nil {
			logrus.Errorf("%v", err)
			break
		}
		// play reply file
		cmd := exec.Command("afplay", replyName)
		err = cmd.Run()
		if err != nil {
			logrus.Errorf("%v", err)
			break
		}
	}
}

func aiff2wav(fileName string) (outPath string) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Invalid path", fileName, err)
		os.Exit(1)
	}
	defer f.Close()

	d := aiff.NewDecoder(f)
	if !d.IsValidFile() {
		fmt.Println("invalid AIFF file")
		os.Exit(1)
	}

	outPath = strings.ReplaceAll(fileName, "aiff", "wav")
	of, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Failed to create", outPath)
		panic(err)
	}
	defer of.Close()

	e := wav.NewEncoder(of, int(d.SampleRate), int(d.BitDepth), int(d.NumChans), 1)

	format := &audio.Format{
		NumChannels: int(d.NumChans),
		SampleRate:  int(d.SampleRate),
	}

	bufferSize := 1000000
	buf := &audio.IntBuffer{Data: make([]int, bufferSize), Format: format}
	var n int
	for err == nil {
		n, err = d.PCMBuffer(buf)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		if n != len(buf.Data) {
			buf.Data = buf.Data[:n]
		}
		if err := e.Write(buf); err != nil {
			panic(err)
		}
	}

	if err := e.Close(); err != nil {
		panic(err)
	}
	fmt.Printf("Aiff file converted to %s\n", outPath)
	return
}

func _aiff() (fileName string) {
	fmt.Println("Recording.  Press Ctrl-C to start.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	select {
	case <-sig:
		fmt.Println("Recording.  Press Ctrl-C to stop.")
	}

	uu, err := uuid.NewUUID()
	chk(err)
	fileName = uu.String()
	if !strings.HasSuffix(fileName, ".aiff") {
		fileName += ".aiff"
	}
	fileName = "./client/tmp/" + fileName
	f, err := os.Create(fileName)
	chk(err)

	// form chunk
	_, err = f.WriteString("FORM")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(0))) //total bytes
	_, err = f.WriteString("AIFF")
	chk(err)

	// common chunk
	_, err = f.WriteString("COMM")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(18)))                  //size
	chk(binary.Write(f, binary.BigEndian, int16(1)))                   //channels
	chk(binary.Write(f, binary.BigEndian, int32(0)))                   //number of samples
	chk(binary.Write(f, binary.BigEndian, int16(32)))                  //bits per sample
	_, err = f.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) //80-bit sample rate 44100
	chk(err)

	// sound chunk
	_, err = f.WriteString("SSND")
	chk(err)
	chk(binary.Write(f, binary.BigEndian, int32(0))) //size
	chk(binary.Write(f, binary.BigEndian, int32(0))) //offset
	chk(binary.Write(f, binary.BigEndian, int32(0))) //block
	nSamples := 0
	defer func() {
		// fill in missing sizes
		totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples
		_, err = f.Seek(4, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(totalBytes)))
		_, err = f.Seek(22, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(nSamples)))
		_, err = f.Seek(42, 0)
		chk(err)
		chk(binary.Write(f, binary.BigEndian, int32(4*nSamples+8)))
		chk(f.Close())
	}()

	portaudio.Initialize()
	defer portaudio.Terminate()
	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	for {
		chk(stream.Read())
		chk(binary.Write(f, binary.BigEndian, in))
		nSamples += len(in)
		select {
		case <-sig:
			return
		default:
		}
	}
	chk(stream.Stop())
	return
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
