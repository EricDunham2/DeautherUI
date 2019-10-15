var vm = this;

Vue.component('monitor', {
    data: function () {
        return {
            packets: [],
            cht: null,
            filtersShown: {
                source: false,
                dest: false,
                rssi: false,
                channel: false,
                type: false,
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
            this.filtersShown[filter] = !this.filtersShown[filter];
            this.filterActive = false;

            Object.values(this.filtersShown).forEach(val => {
                if (val) {
                    this.filterActive = true;
                }
            });
        },
        filterChanged: function() {
            custom_input();
        },
        _filterData: function(data) {
            Object.keys(this.filters).forEach(fil => {
                console.log(this.filters[fil])

                if (this.filters[fil] != null) {
                    data.forEach((ele, idx, arr) => {
                        console.log(ele)
                        if (!ele[fil].includes(this.filters[fil])) {
                            data.splice(idx,1);
                            console.log(data[idx]);
                        }
                    })
                }
            });

            return data;
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

            console.log(this.filterActive);

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
                <div class="col-20 vhc" style="flex-grow:1;">Source</div>
                <div class="col-20 vhc" style="flex-grow:1;">Destination</div>
                <div class="col-5 vhc" style="flex-grow:1;">RSSI</div>
                <div class="col-5 vhc" style="flex-grow:1;">CH</div>
                <div class="col-5 vhc" style="flex-grow:1;">Type</div>
                <div class="col-10 vhc" style="flex-grow:1;">Enc</div>
                <div class="col-30 vhc" style="flex-grow:1;">
                    <span v-if="!filtersShown['vendor']" v-on:click="showFilter('vendor')">Vendor</span>
                    <div class="input-group" v-if="filtersShown['vendor']">
                        <label for="vendor" id="panel-label" class="dyn-input-label">Vendor</label> 
                        <input type="text" id="panel-input" name="vendor" class="dyn-input" v-model="filters['vendor']" v-on:change="filterChanged('vendor')">
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