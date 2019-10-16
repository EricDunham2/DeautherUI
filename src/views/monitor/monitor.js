var vm = this;

Vue.component('monitor', {
    data: function () {
        return {
            packets: [],
            cht: null,
            filtersShown: {
                src: false,
                dst: false,
                rssi: false,
                channel: false,
                pkt_type: false,
                vendor:false,
            },
            filters: {
                src: null,
                dst: null,
                rssi: null,
                channel: null,
                pkt_type: null,
                vendor:null,
            }
        }
    },
    methods: {
        showFilter: function(filter) {

            this.filtersShown[filter] = true;
            this.filterActive = true;

            document.getElementById(`panel-input-${filter}`).focus();
        },
        hideFilter: function(filter) {
            if (this.filters[filter] == null || this.filters[filter] == "") {
                this.filtersShown[filter] = false;
                this.filterActive = false;

                Object.values(this.filtersShown).forEach(val => {
                    if (val) {
                        this.filterActive = true;
                    }
                });
            }
        },
        filterChanged: function() {
            custom_input();
        },
        _filterData: function(data) {
            var copy = JSON.parse(JSON.stringify(data));

            Object.keys(this.filters).forEach(fil => {
                if (this.filters[fil] != null) {
                    copy.forEach((ele, idx, arr) => {
                        if (ele[fil] == null || !ele[fil].toString().toLowerCase().includes(this.filters[fil].toString().toLowerCase())) {
                            copy[idx] = null;
                        }
                    });

                    copy = copy.filter(function (el) {
                        return el != null;
                    });
                }
            });



            return copy;
        },
        _getPackets: function() {
            axios
                .get("/getPackets")
                .then(this._setPackets)
        },
        _setPackets: function (response) {
            if (!response.data) { return; }
            var tempPackets = []
            
            response.data.forEach(pkt => {
               tempPackets.unshift(pkt);
            });

            if (this.filterActive) { 
                tempPackets = this._filterData(tempPackets);
            }

            this.packets = tempPackets;
        },
        _updateChart: function () {
            /*var self = this;

            axios
                .get("/getPacketCount")
                .then(handle)

            function handle(response) {
                self.cht.data.labels.push("");
                self.cht.data.datasets.forEach((dataset) => {
                    dataset.data.push(response.data);
                });

                self.cht.update();
            }*/
        },
        _createChart: function () {
            /*var ctx = document.getElementById('packetMonitor').getContext('2d');

            this.cht = new Chart(ctx, {
                type: "line",
                data: {
                    labels: [],
                    datasets: [{
                        label: 'Packets Per Second',
                        backgroundColor: '#353a44',
                        borderColor: '#b7bdd1',
                        data: []
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        xAxes: [{
                            gridLines: {
                                color: "rgba(0, 0, 0, 0)",
                            },
                            ticks: {
                                max: 10,
                            }
                        }],
                        yAxes: [{
                            gridLines: {
                                color: "rgba(0, 0, 0, 0)",
                            },
                            ticks: {
                                beginAtZero:0,
                                display: false,
                            }
                        }]
                    },
                    elements: {
                        point:{
                            radius: 0
                        }
                    }
                },
            })*/
        }
    },
    beforeMount() {
        this.packets = [];

        this.monitorInterval = setInterval(() => {
            this._getPackets();
        }, 300);
    },
    beforeDestroy() {},
    mounted() {
        var self = this;
        this._createChart();

        setInterval(function () {
            self._updateChart();
        }, 1000);
    },
    template: `
    <div class="flex-container col-100 no-touch-top vhc" style="max-height:100%; overflow-y:scroll;" id="monitor">
        <div class="table col-80">
            <div class="panel-content vhc" v-if="!packets" style="height:110vh; background: rgba(21,21,21,.7); position:fixed; top: 0px !important;">
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

            <!--<canvas id="packetMonitor" style="max-height: 150px;" class="col-80"></canvas>-->

            <div class="col-100 vc col-header" style="color: #efefef; border-radius: 2px; padding-top: 20px;" v-if="packets">
                <!--<div class="col-5 vhc" style="flex-grow:1;"></div>-->
                <div class="col-20 vhc" style="flex-grow:1;" v-on:click="showFilter('src')">
                    <span v-if="!filtersShown['src']" >Source</span>
                    <div class="input-group" v-if="filtersShown['src']" style="width:100%;">
                        <label for="src" id="panel-label" class="dyn-input-label">Source</label> 
                        <input style="width:100%;" type="text" id="panel-input-src" name="src" class="dyn-input" v-model="filters['src']" v-on:change="filterChanged('src')" v-on:blur="hideFilter('src')">
                    </div>
                </div>
                <div class="col-20 vhc" style="flex-grow:1;" v-on:click="showFilter('dst')">
                    <span v-if="!filtersShown['dst']" >Destination</span>
                    <div class="input-group" v-if="filtersShown['dst']" style="width:100%;">
                        <label for="dst" id="panel-label" class="dyn-input-label">Destination</label> 
                        <input style="width:100%;" type="text" id="panel-input-dst" name="dst" class="dyn-input" v-model="filters['dst']" v-on:change="filterChanged('dst')" v-on:blur="hideFilter('dst')">
                    </div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;" v-on:click="showFilter('rssi')">
                    <span v-if="!filtersShown['rssi']" >RSSI</span>
                    <div class="input-group" v-if="filtersShown['rssi']" style="width:100%;">
                        <label for="rssi" id="panel-label" class="dyn-input-label">RSSI</label> 
                        <input style="width:100%;" type="text" id="panel-input-rssi" name="rssi" class="dyn-input" v-model="filters['rssi']" v-on:change="filterChanged('rssi')" v-on:blur="hideFilter('rssi')">
                    </div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;" v-on:click="showFilter('channel')">
                    <span v-if="!filtersShown['channel']">Ch</span>
                    <div style="width:100%;" class="input-group" v-if="filtersShown['channel']">
                        <label for="channel" id="panel-label" class="dyn-input-label">Ch</label>
                        <input style="width:100%;" type="text" id="panel-input-channel" name="channel" class="dyn-input" v-model="filters['channel']" v-on:change="filterChanged('channel')" v-on:blur="hideFilter('channel')">
                    </div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;" v-on:click="showFilter('type')">
                    <span v-if="!filtersShown['type']">Type</span>
                    <div style="width:100%;" class="input-group" v-if="filtersShown['pkt_type']">
                        <label for="type" id="panel-label" class="dyn-input-label">Type</label> 
                        <input style="width:100%;" type="text" id="panel-input-pkt_type" name="type" class="dyn-input" v-model="filters['type']" v-on:change="filterChanged('pkt_type')" v-on:blur="hideFilter('pkt_type')">
                    </div>
                </div>
                <div class="col-10 vhc" style="flex-grow:1;" v-on:click="showFilter('enc')">
                    <span v-if="!filtersShown['enc']">Enc</span>
                    <div class="input-group" v-if="filtersShown['enc']" style="width:100%;">
                        <label for="enc" id="panel-label" class="dyn-input-label">Enc</label> 
                        <input style="width:100%;" type="text" id="panel-input-enc" name="enc" class="dyn-input" v-model="filters['enc']" v-on:change="filterChanged('enc')" v-on:blur="hideFilter('enc')">
                    </div>
                </div>
                <div class="col-30 vhc" style="flex-grow:1;" v-on:click="showFilter('vendor')">
                    <span v-if="!filtersShown['vendor']">Vendor</span>
                    <div class="input-group" v-if="filtersShown['vendor']" style="width:100%;">
                        <label for="vendor" id="panel-label" class="dyn-input-label">Vendor</label> 
                        <input style="width:100%;" type="text" id="panel-input-vendor" name="vendor" class="dyn-input" v-model="filters['vendor']" v-on:change="filterChanged('vendor')" v-on:blur="hideFilter('vendor')">
                    </div>
                </div>
            </div>

            <div class="col-100 card-row" style="flex-grow:1;" v-for="packet in packets">
                <!--<div class="col-5 vhc"></div>-->
                <div class="col-20 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Source:&nbsp</div>
                    <div v-text="packet.src"></div>
                </div>
                <div class="col-20 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Destination:&nbsp</div>
                    <div v-text="packet.dst"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">RSSI:&nbsp</div>
                    <div v-text="packet.rssi"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Ch:&nbsp</div>
                    <div v-text="packet.channel"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Type:&nbsp</div>
                    <div v-text="packet.pkt_type"></div>
                </div>
                <div class="col-10 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Type:&nbsp</div>
                    <div v-text="packet.Enc"></div>
                </div>
                <div class="col-30 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Vendor:&nbsp</div>
                    <div v-text="packet.vendor"></div>
                </div>
            </div>
        </div>
    </div>
    `
});