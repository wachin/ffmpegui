package ffprobe

import (
	"reflect"
	"testing"
	"time"
)

func init() {
	SetLogger()
}

func Test_Mac_getDevices(t *testing.T) {
	want := []string{"0  Built-in Microphone"}
	got := GetFfmpegDevices(&macProber).Audios
	if !reflect.DeepEqual(got, want) {
		t.Errorf("proberKeys.getAudios() = %v, want %v", got, want)
	}
}

func Test_Mac_getCmd(t *testing.T) {
	macprober := GetPlatformProber()
	loadCommonConfig(cfgname)
	want := []string{"ffmpeg", "-benchmark", "-y", "-loglevel", "verbose", "-thread_queue_size", "512", "-f", "avfoundation", "-i", "1:none", "-map", "0:v", "-c:v", "vp9", "0.webm"}
	opts := Options{VidIdx: 1, AudIdx: -1, Container: 1}
	// Preset: "webm - vp9 default with no audio"}
	SetOptions(opts)
	if got := getCommand(macprober); !reflect.DeepEqual(got, want) {
		t.Errorf("getCommand = %#v, want %v", got, want)
	}
}

func Test_getPlatformProber(t *testing.T) {
	tests := []struct {
		name string
		want Prober
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPlatformProber(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPlatformProber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFfmpegDevices(t *testing.T) {
	tests := []struct {
		name  string
		dtype string
		want  string
	}{
		{"mac", "audio", "0  Built-in Microphone"},
		{"macvideo", "video", "0  FaceTime HD Camera"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFfmpegDeviceType(&macProber, tt.dtype)
			if !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("parseFfmpegDevices() = %#v, want %v", got, tt.want)
			}
		})
	}
}

func Test_StartProcessFail(t *testing.T) {
	prober := GetPlatformProber()
	config.Ffconf.Ffcmdprefix = "pytho"
	scanner, _ := StartEncode(prober)
	if scanner != nil {
		t.Errorf("expected process fail")
	}
}

func Test_ProcessInterrupt(t *testing.T) {
	prober := GetPlatformProber()
	config.Ffconf.Ffcmdprefix = "sleep 10"
	tbeg := time.Now().UnixNano()
	StartEncode(prober)
	if !StopEncode() || (time.Now().UnixNano()-tbeg > 1e9) {
		t.Error("process interrupt failed or too slow")
	}
}

func Test_StartProcessOutput(t *testing.T) {
	prober := GetPlatformProber()
	config.Ffconf.Ffcmdprefix = "ls asdf1234"
	scanner, _ := StartEncode(prober)
	var ffout string
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			txt := scanner.Text()
			ffout += txt
		}
		done <- true
	}()
	<-done
	if ffout != "ls: asdf1234: No such file or directory" {
		t.Error("wrong process output" + ffout)
	}
}
