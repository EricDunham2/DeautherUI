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
            <div class="col-100 vc col-header" style="color: #cc14ab; border-bottom: 1px solid #673AB7; border-radius: 2px; padding-top: 20px;">
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
                <div class="col-30 vhc" style="flex-grow:1;">
                    <div class="card-row-header">Source:</div>
                    <div v-text="packet.src"></div>
                </div>
                <div class="col-30 vhc" style="flex-grow:1;">
                    <div class="card-row-header">Destination:</div>
                    <div v-text="packet.dst"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;">
                    <div class="card-row-header">RSSI:</div>
                    <div v-text="packet.rssi"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;">
                    <div class="card-row-header">Ch:</div>
                    <div v-text="packet.channel"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;">
                    <div class="card-row-header">Type:</div>
                    <div v-text="packet.pkt_type"></div>
                </div>
                <div class="col-5 vhc" style="flex-grow:1;">
                    <div class="card-row-header">Vendor:</div>
                    <div v-text="packet.vendor"></div>
                </div>
            </div>
        </div>
    </div>
    `
});






















/*Card format
//Hide col-header
//Each packet will be formatted as such.
.card {
    width:100%;
    height:250px;
}


<div class="card">
    //n number of rows
    <div class="row"><div class="pull-left">HEADER</div><div class="pull-right"></div></div>
</div>

Standard Format
    //This will be hard coded in html
    <div class="row vc col-header">
        <div class="col-5 vhc"></div>
        <div class="col-15 vhc">Source</div>
        <div class="col-15 vhc">Destination</div>
        <div class="col-5 vhc">RSSI</div>
        <div class="col-5 vhc">CH</div>
        <div class="col-10 vhc">Type</div>
        <div class="col-10 vhc">Enc</div>
        <div class="col-30 vhc">SSID</div>
    </div>

    //This will be dynamic
    <div class="row vc">
        <div class="col-5 vhc">
            <!--<label class="container vhc">
                <input type="checkbox" checked="checked">
                <span class="checkmark"></span>
            </label>-->
        </div>
        <div class="col-15 vhc">AF:2B:4C:98:82</div>
        <div class="col-15 vhc">AF:2B:4C:98:82</div>
        <div class="col-5 vhc">100</div>
        <div class="col-5 vhc">4</div>
        <div class="col-10 vhc">BEACON</div>
        <div class="col-10 vhc">WPA2</div>
        <div class="col-30 vhc">Some SSID Name That can be 32 222</div>
    </div>
*/