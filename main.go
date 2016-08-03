package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"time"

	"github.com/ninjasphere/go-castv2"
	"github.com/ninjasphere/go-castv2/controllers"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var filesVar = flag.String("files", "", "Comma separated list of files.")
var deviceVar = flag.String("devicename", "", "Name of chromecast device")
var portVar = flag.String("port", "50041", "Port the server will listen on.")

func main() {
	flag.Parse()

	files := strings.Split(*filesVar, ",")
	deviceName := *deviceVar
	port := fmt.Sprintf(":%v", *portVar)

	log.Printf("deviceName: %v", deviceName)
	log.Printf("files: %v", files)
	log.Printf("port: %v", port)

	ctx := context.Background()
	devices := listChromecastsWithTimeout(ctx, 5*time.Second)

	log.Printf("Found: %+v", devices)

	device, err := findDeviceName(devices, deviceName)
	if err != nil {
		panic(err)
	}

	client, err := castv2.NewClient(device.Addr, device.Port)
	if err != nil {
		log.Fatalf("Failed to connect to chromecast %v: %v", device.Addr, err)
	}

	//_ = controllers.NewHeartbeatController(client, "Tr@n$p0rt-0", "Tr@n$p0rt-0")

	heartbeat := controllers.NewHeartbeatController(client, "sender-0", "receiver-0")
	heartbeat.Start()

	connection := controllers.NewConnectionController(client, "sender-0", "receiver-0")
	connection.Connect()
	defer connection.Close()

	//media, err := controllers.NewMediaController(client, "sender-0", "receiver-0")

	receiver := controllers.NewReceiverController(client, "sender-0", "receiver-0")

	response, err := receiver.GetStatus(time.Second * 5)
	log.Printf("Status response=%+v err=%v", response, err)

	res, err := receiver.LaunchApplication(5*time.Second, controllers.DefaultMediaPlayerAppID)
	log.Printf("Status response=%+v err=%v", res, err)

	http.ListenAndServe(port, serveFile(files[0]))
}

func findDeviceName(devices []ChromecastDevice, name string) (ChromecastDevice, error) {
	for _, device := range devices {
		if device.Name == name {
			return device, nil
		}
	}
	return ChromecastDevice{}, errors.New(fmt.Sprintf("No device with name '%v' found", name))
}

func getMimeFromFilename(filename string) string {
	return mime.TypeByExtension(filepath.Ext(filename))
}

func serveFile(path string) http.HandlerFunc {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	totalSize := stat.Size()
	mimeType := getMimeFromFilename(stat.Name())

	log.Printf("Serve file path=%v size=%v mime=%v", path, totalSize, mimeType)

	return func(rw http.ResponseWriter, req *http.Request) {
		http.ServeFile(rw, req, path)
		//
		//	rw.Header().Set("Content-Type", mimeType)
		//	rw.Header().Set("Access-Control-Allow-Origin", "*");
		//
		//	hRange := req.Header.Get("Range")
		//	if hRange == "" {
		//		rw.Header().Set("Content-Length", fmt.Sprintf("%v", totalSize));
		//		rw.Write()
		//		//return fs.createReadStream(filePath).pipe(res);
		//
		//	}
		//
	}
}
