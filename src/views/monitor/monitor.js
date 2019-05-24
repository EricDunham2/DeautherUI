Vue.component('monitor', {
    data: function() {
        return {
            packets:[],
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
            this.packets = [];

            console.log(response.data);

            response.data.forEach(pkt => {
                this.packets.push(pkt);          
            });
        }
    },
    beforeMount() {
        this.packets = [];

        this.monitorInterval = setInterval(() => {
            this._getPackets();
        }, 300);
    },
    beforeDestroy() {
    },
    template : `
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
            <div class="col-100 vc col-header" style="color: #9e16c3;; border-radius: 2px; padding-top: 20px;" v-if="packets">
                <!--<div class="col-5 vhc" style="flex-grow:1;"></div>-->
                <div class="col-30 vhc" style="flex-grow:1;">Source</div>
                <div class="col-30 vhc" style="flex-grow:1;">Destination</div>
                <div class="col-5 vhc" style="flex-grow:1;">RSSI</div>
                <div class="col-5 vhc" style="flex-grow:1;">CH</div>
                <div class="col-5 vhc" style="flex-grow:1;">Type</div>
                <!--<div class="col-10 vhc" style="flex-grow:1;">Enc</div>-->
                <div class="col-5 vhc" style="flex-grow:1;">Vendor</div>
            </div>

            <div class="col-100 card-row" style="flex-grow:1;" v-for="packet in packets">
                <!--<div class="col-5 vhc"></div>-->
                <div class="col-30 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Source:</div>
                    <div v-text="packet.src"></div>
                </div>
                <div class="col-30 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Destination:</div>
                    <div v-text="packet.dst"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">RSSI:</div>
                    <div v-text="packet.rssi"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Ch:</div>
                    <div v-text="packet.channel"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Type:</div>
                    <div v-text="packet.pkt_type"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1; overflow:hidden;">
                    <div class="card-row-header">Vendor:</div>
                    <div v-text="packet.vendor"></div>
                </div>
            </div>
        </div>
    </div>
    `
});