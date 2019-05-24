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

            console.log(response.data)

            Object.values(response.data).forEach(ap => {
                //ap.fmtMac = this._createMacString(ap.bssid);

                if (!ap.stations || !ap.stations.length) {
                    return;
                }
                ap.stations.forEach(sta => {
                    sta.selected = false; //Changes so that when rescans it doesnt overwrite old values
                   //sta.fmtMac = this._createMacString(sta.mac);
                    sta.ap = ap.bssid;
                });
            });

            this.accesspoints = Object.values(response.data);
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
                if (!this.selectedStations) {
                    this.selectedStations = [];
                }
                this.selectedStations.push(sta);
            } else {
                var index = this.selectedStations.findIndex(st => st.mac === sta.mac);
                this.selectedStations.splice(index, 1);
            }
        },
        onTileHover: function (item) {
            this.hoveredItem = item;
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
                <div class="panel-content vhc" v-if="isAttacking || !accesspoints" style="height:110vh; background: rgba(21,21,21,.7); position:absolute;">
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
                                        class="ap-information-panel__large l5 clickable no-touch-table panel"
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
                                        <div class="panel-footer vhc">
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
                                        class="ap-information-panel__medium l5 clickable no-touch-table panel"
                                        v-on:click="selectStation(sta)"
                                        v-on:mouseover="onTileHover(sta)"
                                        v-on:mouseout="offTileHover()"
                                        v-for="sta in selectedAccesspoint.stations"
                                        v-if="!sta.selected"
                                    >
                                        <div class="panel-header vhc">
                                            <div v-text="sta.fmtMac"></div>
                                        </div>
                                        <div class="panel-content vhc">
                                            <i class="material-icons icon-md vhc">devices</i>
                                        </div>
                                        <div class="panel-footer vhc">
                                            <div v-text="sta.vendor"></div>
                                        </div>
                                    </div>
                                </template>
                            </div>
                        </div>
                        <div class="panel col-30 no-touch-top main-content-col">
                            <div class="panel-header vhc">Targets</div>
                            <div class="panel-content col-90 hc" id="selectedTargets" style="overflow:auto; height:calc(100% - 75px)">
                                <div class="ap-information-panel__medium l5 clickable no-touch-table shadow" v-on:click="selectStation(sta)" v-for="sta in selectedStations"
                                v-on:mouseover="onTileHover(sta)" v-on:mouseout="offTileHover()">
                                    <div class="panel col-100">
                                        <div class="panel-header vhc">
                                            <div v-text="sta.fmtMac"></div>
                                        </div>
                                        <div class="panel-content vhc">
                                            <i class="material-icons icon-md vhc">devices</i>
                                        </div>
                                        <div class="panel-footer vhc">
                                            <div v-text="sta.vendor"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="footer l2 vhc" style="position:fixed; bottom:0; height:max-content;">
                <div class="vhc col-100">
                    <div class="vhc" style="width:100%;">
                        <input id="deauth" type="checkbox" role="button" class="toggle-btn" v-model="isAttacking" @click="toggleAttack()"/>
                        <label for="deauth" class="toggle-lbl vh-center" style="width:100%; margin:0;" ><span class="v-center" style="text-transform: uppercase;">Start</span></label>
                    </div>
                </div>
            </div>
            <!--<div v-if="isAttacking" class="vhc" id="statusbar" style="height:30px; position: absolute; width:100%; background:crimson; color:black; font-size: 15px;">
                [<span ></span>>]
            </div>-->
        </div>
    </div>
    `
});