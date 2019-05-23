package main

import (
	"html/template"
	"net/http"
	"os"
	"net"
	"go.bug.st/serial.v1"
	//"os/exec"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"encoding/json"
	"time"
	//"strconv"
	"io/ioutil"
	"log"
	//"reflect"
)

//const MacAddr [6]byte = [6]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} //*Replace me with self address used for the Probe and Beacon Requests
const BAUD int = 115200 //Read from a file
const PORT string = "/dev/ttyAMA0"
const MSG_START byte = 0x02 
const MSG_END byte = 0x03
const CONFIG_PATH string = "./src/static/configs/config.json"

type Settings struct {
	Ap struct {
		Ssid    string `json:"ssid"`
		Passwd  string `json:"passwd"`
		Channel int    `json:"channel"`
		Hidden  bool   `json:"hidden"`
	} `json:"accesspoint"`
	ApScanner struct {
		Interval int  `json:"interval"`
		Deep     bool `json:"deep"`
		Async    bool `json:"async"`
		Channel  int  `json:"channel"`
		Hop      bool `json:"hop"`
	} `json:"apScanner"`
	PktScanner struct {
		Interval int  `json:"interval"`
		Channel  int  `json:"channel"`
		Hop      bool `json:"hop"`
	} `json:"packetScanner"`
}

/*type AccessPoint struct {
	Ecn int //`json:"ecn"`
	Etype string //`json:"etype"`
	SSID string //`json:"ssid"`
	RSSI int //`json:"rssi"`
	Mac [6]byte //`json:"mac"`
	Stations []Station //`json:"stations"`
}*/

type AccessPoint struct {
	Ssid     string    	`json:"ssid"`
	Enc  	 int    	`json:"enc"`
	Rssi	 int 		`json:"rssi"`
	Bssid    string 	`json:"bssid"`
	Channel  int 		`json:"channel"`
	Hidden   bool 		`json:"hidden"`
	Stations []Station	`json:"stations"`
}

type Packet struct {
	Rssi     	int    `json:"rssi"`
	Channel  	int    `json:"channel"`
	PktType  	string `json:"pkt_type"`
	Src      	string `json:"src"`
	Dst      	string `json:"dst"`
	Bssid    	string `json:"bssid"`
	Ssid     	string `json:"ssid"`
}

type Station struct {
	Name 	string 	`json:"name"`
	Vendor 	string 	`json:"vendor"`
	Mac 	[]byte	`json:"mac"`
	MacStr  string 	`json:"macStr"`
	AP 		[]byte 	`json:"ap"`
	APStr	string 	`json:"apStr"`
}

var (
	AvailableAccesspoints map[string]AccessPoint
	router * mux.Router
	//APL []AccessPoint //Depcrated, Remove .
	writeQueue [][]byte
	readQueue []string

	packets []Packet
	logs []string
	sConn serial.Port
	sniffing bool
	settings Settings
)

func main() {
	router = mux.NewRouter()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	sniffing = false

	//channel = 11
	//log_message(fmt.Sprintf("Default channel set to %d...", channel))

	settingsDat := readConfig()
	AvailableAccesspoints = make(map[string]AccessPoint)
	json.Unmarshal(settingsDat, &settings)
	log_message(fmt.Sprintf("Default settings set to %s...", settings))

	initRoutes()
	initMockData()

	if initSerial() {
		go ActionHandler();
	}

	go handlerData();

	http.ListenAndServe(":8081", router)
}


func initSerial() bool{
	log_message("Opening serial connection...")

	ports, err := serial.GetPortsList()

	if err != nil {
		log_message(err.Error())
		return false
	}

	if len(ports) == 0 {
		log_message("No serial ports were found...")
		return false
	}

	mode := &serial.Mode { BaudRate: BAUD }

	sConn, err = serial.Open(PORT, mode)

	if err != nil {
		log_message(err.Error())
		return false
	}

	log_message("Serial connection opened...")
	return true
}

func initMockData() {
	/*sta1 := Station{Vendor:"Apple", Mac:[6]byte{52,0,0,0,0,1}}
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
	APL = append(APL, ap15);*/

	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 1", "enc": 2, "rssi": 10, "bssid":"11:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 2", "enc": 2, "rssi": 20, "bssid":"12:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 3", "enc": 2, "rssi": 40, "bssid":"13:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 4", "enc": 2, "rssi": 80, "bssid":"14:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 5", "enc": 2, "rssi": 100, "bssid":"15:11:11:11:11:11", "channel": 11, "hidden": false}`)

	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "11:11:11:11:11:11", "bssid": "11:11:11:11:11:11", "ssid": "Some Cool SSID 1"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "12:11:11:11:11:11", "bssid": "12:11:11:11:11:11", "ssid": "Some Cool SSID 2"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "13:11:11:11:11:11", "bssid": "13:11:11:11:11:11", "ssid": "Some Cool SSID 3"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "14:11:11:11:11:11", "bssid": "14:11:11:11:11:11", "ssid": "Some Cool SSID 4"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "15:11:11:11:11:11", "bssid": "15:11:11:11:11:11", "ssid": "Some Cool SSID 5"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "14:11:11:11:11:11", "bssid": "14:11:11:11:11:11", "ssid": "Some Cool SSID 4"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "13:11:11:11:11:11", "bssid": "13:11:11:11:11:11", "ssid": "Some Cool SSID 3"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "12:11:11:11:11:11", "bssid": "12:11:11:11:11:11", "ssid": "Some Cool SSID 2"}`)
	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "11:11:11:11:11:11", "bssid": "11:11:11:11:11:11", "ssid": "Some Cool SSID 1"}`)
}

func initRoutes() {
	resources := packr.NewBox("./src/static")
	templates := packr.NewBox("./src/views")

	log_message("Initiating routes...")

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(resources)))
	router.PathPrefix("/templates").Handler(http.StripPrefix("/templates/", http.FileServer(templates)))

	router.HandleFunc("/accesspoints", getAccessPoints).Methods("GET")
	router.HandleFunc("/getPackets", getPackets).Methods("GET")
	router.HandleFunc("/sniffPackets", sniffPackets).Methods("GET")
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
	var cmd string = fmt.Sprintf("scan -async=%t -hidden=%t -channel=%d -hop=%t", settings.ApScanner.Async, settings.ApScanner.Deep, settings.ApScanner.Channel, settings.ApScanner.Hop)
	fmt.Println(cmd)
	//TODO Uncomment me for prod
	//writeQueue = append(writeQueue, cmd)

	re, _ := json.Marshal(AvailableAccesspoints)

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func handlerData() {
	for {
		if (len(readQueue) == 0) { 
			time.Sleep(200)
			continue; 
		}

		var buffer []byte = []byte(readQueue[0])
		readQueue = readQueue[1:]
		
		bufferMap := make(map[string]interface{})
		err := json.Unmarshal(buffer, &bufferMap)

		if err != nil { 
			log_message(err.Error())
			continue
		}

		var datType string = fmt.Sprintf("%v", bufferMap["data_type"])

		switch datType {
			case "packet" :
				packet := Packet{}

				if err =  json.Unmarshal(buffer, &packet); err != nil {
					log_message(err.Error())
					break
				}

				packets = append(packets, packet)

				if (packet.PktType != "MGMT") { 
					break
				}
				
				station := Station{}
				//station.Vendor = vendorLookup(packet.Ssid)

				var mac net.HardwareAddr

				if packet.Src != packet.Bssid {
					mac, err = net.ParseMAC(packet.Src)
				} else if packet.Dst != packet.Bssid {
					mac, err = net.ParseMAC(packet.Src)
				}

				if mac != nil {
					station.Mac = []byte(mac)
				}

				if err != nil {
					log_message(err.Error())
					break
				}

				if _, ok := AvailableAccesspoints[packet.Bssid]; ok {
					apMac, err := net.ParseMAC(packet.Bssid)
					station.AP = []byte(apMac)

					if err != nil {
						log_message(err.Error())
						break
					}

					Ap := AvailableAccesspoints[packet.Bssid]

					Ap.Stations = append(Ap.Stations, station)
					AvailableAccesspoints[packet.Bssid] = Ap
				}

				break
			case "accesspoint" :
				ap := AccessPoint{}

				if err = json.Unmarshal(buffer, &ap); err != nil {
					log_message(err.Error())
					break
				}

				if val, ok := AvailableAccesspoints[ap.Bssid]; ok {
					val.Rssi = ap.Rssi
					val.Ssid = ap.Ssid
					val.Channel = ap.Channel
					val.Hidden = ap.Hidden
				} else {
					ap.Stations = []Station{}
					AvailableAccesspoints[ap.Bssid] = ap
				}

				break
			default :
				break
		}
	}
}

func getPackets(w http.ResponseWriter, r *http.Request) {
	re, _ := json.Marshal(readQueue)

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func sniffPackets(w http.ResponseWriter, r *http.Request) {
	var cmd string = fmt.Sprintf("sniff -interval=%d -channel=%d -hop=%t", settings.PktScanner.Interval, settings.PktScanner.Channel, settings.PktScanner.Hop)

	fmt.Println(cmd)
	//TODO Uncomment me for prod
	//writeQueue = append(writeQueue, cmd)

	log_message("Sniffing packets")

	w.Header().Set("Content-Type", "application/json")
}

func startAccessPoint() {
	var cmd string = fmt.Sprintf("setup  -ssid=%q -hidden=%t -channel=%d -password=%q", settings.Ap.Ssid, settings.Ap.Hidden, settings.Ap.Channel, settings.Ap.Passwd)
	
	fmt.Println(cmd)
	//TODO Uncomment me for prod
	//writeQueue = append(writeQueue, cmd)

	log_message("Starting accesspoint packets")
}

func getConfig(w http.ResponseWriter, r *http.Request) {

	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err)  {
		log_message("Config file does not exist...")
		createConfig(`{"accesspoint":{"ssid":"Lambs to the Cosmic Slaughter","passwd":"Rick and Mortison","channel":12,"hidden":false},"apScanner":{"interval":1000,"deep":true,"async":false,"channel":12,"hop":false},"packetScanner":{"interval":1000, "channel": 12, "hop": false}}`)
	}

	var dat []byte = readConfig()

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

	json.Unmarshal(body, &settings)
	startAccessPoint()
}

func readConfig() []byte {
	file, err := os.Open(CONFIG_PATH)

	if err != nil {
		log_message(err.Error())
		return []byte{};
	}

	dat, err := ioutil.ReadAll(file)

	if err != nil {
		log_message(err.Error())
		return []byte{};
	}

	json.Unmarshal(dat, &settings)

	return dat
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

/*``
*	Pass a station and a deauth packet will be created to be sent to the station aswell as the 
*	accesspoint. The created packet will be pushed to the writeQueue to be sent
*/

func deauthStation(station Station) {
	//cmd := "send -interval=%d -channel=%d -buffer=%s"

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

	//var toStation []byte
	//var toAccesspoint []byte

	//toStation = append(toStation, []byte(cmd)...)
	//toAccesspoint = append(toAccesspoint, []byte(cmd)...)

	//Create the deauth packet that will be sent to the station aswell as the accesspoint
	//toStation = append(toStation, createPacket(ps)...)
	//toAccesspoint = append(toStation, createPacket(pa)...)

	//log_message(fmt.Sprintf("%#X", createPacket(ps)))
	//log_message(fmt.Sprintf("%#X", createPacket(pa)))

	//log_message(fmt.Sprintf("Deauthing %#X...", station.Mac))

	//Add the packets to the queue to be sent
	cmd := fmt.Sprintf("send -interval=%d -channel=%d -buffer=%s", 1000, 12, createPacket(pa))
	fmt.Println(cmd)
	//writeQueue = append(writeQueue, []byte(cmd))

	cmd = fmt.Sprintf("send -interval=%d -channel=%d -buffer=%s", 1000, 12, createPacket(ps))
	fmt.Println(cmd)

	//writeQueue = append(writeQueue, []byte(cmd))
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
	/*n, err := sConn.Write(data)

	if err != nil {
		log_message(err.Error())
	}*/
}

func ReadHandler() []byte {
	var building bool = false
	var buf []byte
	var ch []byte = make([]byte, 1)

	_, err := sConn.Read(ch)

	if err != nil {
		log_message(err.Error())
		return nil
	}

	for ch[0] != MSG_END {
		if !building && len(writeQueue) > 0 { return []byte{} }

		_, err = sConn.Read(ch)
		if err != nil {
			log_message(err.Error())
			return nil
		}

		if ch[0] == MSG_START {
			building = true;
		} else if ch[0] == MSG_END {
			break
		}else if building {
			buf = append(buf, ch[0])
		}
	}
	
	return buf
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
			byteRead := ReadHandler();
			
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

	log.Print(message)
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

	
