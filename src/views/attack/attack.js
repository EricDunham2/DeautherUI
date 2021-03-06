/*var vm = this;*/

Vue.component('attack', {
    //el: '#content',
    data: function() {
        return {
            accesspoints: null,
            selectedStations: null,
            selectedAccesspoint: null,
            isAttacking: false,
            isScanning: true,
            selectionTitle: "Accesspoints",
            hoveredItem: null,
            scanId: null,
            scanInterval: 5000,
            sentPackets: 0,
        }
    },
    methods: {
        _getPackets: function() {
            axios
            .get("/getPackets")
            .then(this._setPackets)
        },
        _setPackets: function(response) {
            if (!response.data) { return; }
        },
        _getAccesspoints: function () {
            axios
                .get("/accesspoints")
                .then(this._setAccesspoints)
        },
        _setAccesspoints: function (response) {
            if (!response.data) {
                return;
            }

            Object.values(response.data).forEach(ap => {
                //ap.fmtMac = this._createMacString(ap.bssid);

                if (!ap.stations || !ap.stations.length) {
                    return;
                }
                
                ap.LastSeen = Date.now();

                ap.stations.forEach(sta => {
                    if (this.selectedStations != null && this.selectedStations.find(st => st.mac === sta.mac)) { 
                        sta.selected = true;
                        return;
                    }

                    sta.selected = false; //Changes so that when rescans it doesnt overwrite old values
                   //sta.fmtMac = this._createMacString(sta.mac);
                    //sta.ap = ap.bssid;
                });
            });

            this.accesspoints = Object.values(response.data);

            this.accesspoints.forEach((ap, index, arr) => {
                if (Date.now() - ap.LastSeen > 300000) {
                    this.accesspoints.splice(index, 1);
                }
            })
        },
        _createMacString: function (arr) {
            return `${this._toHex(arr[0])}:${this._toHex(arr[1])}:${this._toHex(arr[2])}:${this._toHex(arr[3])}:${this._toHex(arr[4])}:${this._toHex(arr[5])}`;
        },
        _toHex: function (byte) {
            return Number(byte).toString(16);
        },
        setSelectedAccesspoint: function (ap) {
            this.selectedAccesspoint = ap;

            if (!this.selectedAccesspoint) {
                this.selectionTitle = "Accesspoints";
            } else {
                this.selectionTitle = "Stations";
            }
        },
        selectStation: function (sta) {
            sta.selected = !sta.selected;

            if (sta.selected) {
                sta.ssid = this.selectedAccesspoint.ssid;
                sta.rssi = this.selectedAccesspoint.rssi;
                sta.enc = this.selectedAccesspoint.enc;

                if (!this.selectedStations) {
                    this.selectedStations = [];
                }

                this.selectedStations.push(sta);
                console.log(this.selectedStations);
            } else {
                var index = this.selectedStations.findIndex(st => st.mac === sta.mac);
                this.selectedStations.splice(index, 1);
            }
        },
        onTileHover: function (item) {
            this.hoveredItem = item;
            if (!this.hoveredItem.rssi) {
                this.hoveredItem.ssid = this.selectedAccesspoint.ssid;
                this.hoveredItem.rssi = this.selectedAccesspoint.rssi;
                this.hoveredItem.enc = this.selectedAccesspoint.enc;
            }
        },
        offTileHover: function () {
            this.hoveredItem = null;
        },
        toggleScan: function () {
            isScanning = !isScanning;

            if (!isScanning) {
                clearInterval(this.scanId);
                return;
            }

            this.scanId = setInterval(function () {
                this._getAccesspoints();
            }, scanInterval);
        },
        toggleAttack: function () {
            this.isAttacking = !this.isAttacking;

            if (this.isAttacking) {
                /*var stationMacs = [];

                this.selectedStations.forEach(sta => {
                    stationMacs.push(sta.Mac);
                });*/
                this.attackInterval = setInterval(() => {
                    var data = JSON.stringify(this.selectedStations);
                    axios.post('/deauth', data);
                }, 500);
            } else {
                clearInterval(this.attackInterval);
            }
        }
    },
    beforeMount() {
        this.packets = [];

        setInterval(() => {
            this._getAccesspoints();
        }, 2000);

       setInterval(() => {
            this._getPackets();
        }, 2000);
    },
    template: `
    <div class="col-100">
        <div class="flex-container col-100" id="content">
            <div class="flex-container col-100 no-touch-top" style="overflow:none;">
                <div class="panel-content vhc" v-if="isAttacking || !accesspoints" style="height: 110vh;background: rgba(21, 21, 21, 0.7);position: fixed; top: 0px !important;">
                    <div>
                        <div class="row clearfix">
                            <div class="square one"></div> 
                            <div class="square two"></div>
                            <div class="square three"></div>
                        </div>

                        <div class="row clearfix">
                            <div class="square eight"></div> 
                            <div class="square nine"></div>
                            <div class="square four"></div>
                        </div>

                        <div class="row clearfix">
                            <div class="square seven"></div> 
                            <div class="square six"></div>
                            <div class="square five"></div>
                        </div>
                    </div>
                </div>
                <div class="panel col-100">
                    <!--<div class="panel-header" style="position:sticky; top:0; color:crimson; height:45px;">
                        <div style="height:45px;">
                            <div v-text="hoveredItem.fmtMac" v-if="hoveredItem"></div>
                            <div v-text="hoveredItem.ssid" v-if="hoveredItem"></div>
                        </div>
                    </div>-->
                    <div class="panel-content flex-container" style="height:calc(100% - 155px)">
                        <div class="panel col-70 no-touch-top main-content-col">
                            <div class="panel-header vhc" id="targetsTitle" v-text="selectionTitle"></div>
                            <div class="panel-content col-90 hc" id="targetList" style="overflow:auto; height:calc(100% - 75px)">
                                <template v-if="!selectedAccesspoint">
                                    <div
                                        class="ap-information-panel__medium l5 clickable no-touch-table panel target"
                                        v-on:click="setSelectedAccesspoint(ap)"
                                        v-on:mouseover="onTileHover(ap)"
                                        v-on:mouseout="offTileHover()"
                                        v-for="ap in accesspoints"
                                    >
                                        <div class="panel-header vhc">
                                            <div v-text="ap.ssid"></div>
                                        </div>
                                        <div class="panel-content vhc">
                                            <img v-if="ap.rssi == 0" class="icon-md" src="/static/images/wifi_0_grey.png">
                                            <img v-else-if="ap.rssi > 0 && ap.rssi < 25" class="icon-md" src="/static/images/wifi_1_grey.png">	
                                            <img v-else-if="ap.rssi >= 25 && ap.rssi < 50" class="icon-md" src="/static/images/wifi_2_grey.png">	
                                            <img v-else-if="ap.rssi >= 50 && ap.rssi < 75" class="icon-md" src="/static/images/wifi_3_grey.png">	
                                            <img v-else-if="ap.rssi >= 75 && ap.rssi <= 100" class="icon-md" src="/static/images/wifi_4_grey.png">	
                                        </div>
                                        <div class="panel-footer tc">
                                            <div v-text="ap.bssid"></div>
                                        </div>
                                    </div>
                                </template>
                                <template v-else>
                                    <div class="ap-information-panel__medium l3 clickable no-touch-table action" style="background: #23096f; color: #05fc7c;" v-on:click="setSelectedAccesspoint()">
                                        <div class="panel col-100">
                                            <div class="panel-header col-80 vhc">Back</div>
                                            <div class="panel-content vhc">
                                                <i class="material-icons icon-md vhc">arrow_back</i>
                                            </div>
                                        </div>
                                    </div>
                                    <div
                                        class="ap-information-panel__medium target l5 clickable no-touch-table panel"
                                        v-on:click="selectStation(sta)"
                                        v-on:mouseover="onTileHover(sta)"
                                        v-on:mouseout="offTileHover()"
                                        v-for="sta in selectedAccesspoint.stations"
                                        v-if="!sta.selected"
                                    >
                                        <div class="panel-header vhc">
                                            <div v-text="sta.fmtMac"></div>
                                        </div>
                                        <div class="panel-content vhc" style="color:slateblue">
                                            <i class="material-icons icon-md vhc">devices</i>
                                        </div>
                                        <div class="panel-footer tc">
                                            <div v-text="sta.vendor"></div>
                                        </div>
                                    </div>
                                </template>
                            </div>
                        </div>
                        <div class="panel col-30 no-touch-top main-content-col">
                            <div class="panel-header vhc">Targets</div>
                            <div class="panel-content col-90 hc" id="selectedTargets" style="overflow:auto; height:calc(100% - 75px)">
                                <div class="ap-information-panel__medium l5 clickable no-touch-table target" v-on:click="selectStation(sta)" v-for="sta in selectedStations"
                                v-on:mouseover="onTileHover(sta)" v-on:mouseout="offTileHover()">
                                    <div class="panel col-100">
                                        <div class="panel-header vhc">
                                            <div v-text="sta.fmtMac"></div>
                                        </div>
                                        <div class="panel-content vhc">
                                            <i class="material-icons icon-md vhc devices-icon" style="color:slateblue;">devices</i>
                                        </div>
                                        <div class="panel-footer tc">
                                            <div v-text="sta.vendor"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="footer vhc" style="position:fixed; bottom:0; height:max-content; background:#112;">
                <div class="col-100 vhc" v-if="(selectedAccesspoint && hoveredItem) || (hoveredItem != null && hoveredItem.selected)" style="color:greenyellow;">
                    <div style="width: 40%; padding: 5px 5px 0 5px; font-size:13px;" class="vhc">VENDOR: <span v-text="hoveredItem.vendor"> </span></div>
                    <div style="width: 40%; padding: 5px 5px 0 5px; font-size:13px;" class="vhc">MAC: <span v-text="hoveredItem.fmtMac"></span></div>
                    <div style="width: 40%; padding: 5px 5px 0 5px; font-size:13px;" class="vhc">SSID: <span v-text="hoveredItem.ssid"></span></div>
                    <div style="width: 40%; padding: 5px 5px 0 5px; font-size:13px" class="vhc">RSSI: <span v-text="hoveredItem.rssi"></span></div>
                    <div v-if="selectedAccesspoint" style="width: 40%; padding: 5px 5px 0 5px; font-size:13px;" class="vhc">ENCRYPTION: <span v-text="selectedAccesspoint.enc"></span></div>
                </div>
                <div class="vhc col-100">

                    <div class="input-group vhc">
                        <div style="font-size:13px;" class="checkbox-label">
                            <span v-if="!isAttacking">Attack</span>
                            <span v-if="isAttacking">Attacking</span>
                        </div>
                        <label class="switch" for="deauth-checkbox">
                            <input type="checkbox" id="deauth-checkbox" @change="toggleAttack()"/>
                            <div class="slider round"></div>
                        </label>
                    </div>
                </div>
            </div>
        </div>
    </div>
    `
});