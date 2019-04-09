Vue.component('monitor', {
    data: {
        function() {
            return {
                packets:[],
            }
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
            
            response.data.forEach(pkt => {
                try {
                    this.packets.push(JSON.parse(pkt));
                } catch (err) {
                    console.error(err);
                    console.error(`Failed to parse packet ${pkt}`);
                }
            });
        }
    },
    beforeMount() {
        this.monitorInterval = setInterval(() => {
            this._getPackets();
        }, 300);
    },
    beforeDestroy() {
        clearInterval(this.monitorInterval);
    },
    template : `
    <div class="flex-container col-100 no-touch-top vhc" style="max-height:100%; overflow-y:scroll;" id="monitor">
        <div class="table col-80">
            <div class="row vc col-header" style="color: #cc14ab; border-bottom: 1px solid #673AB7; border-radius: 2px; margin: 10px; padding: 5px;">
                <div class="col-5 vhc"></div>
                <div class="col-20 vhc">Source</div>
                <div class="col-20 vhc">Destination</div>
                <div class="col-5 vhc">RSSI</div>
                <div class="col-5 vhc">CH</div>
                <div class="col-10 vhc">Type</div>
                <!--<div class="col-10 vhc">Enc</div>-->
                <div class="col-30 vhc">Vendor</div>
            </div>

            <div class="row vc" v-for="packet in packets">
                <div class="col-5 vhc">
                    <!--<label class="container vhc">
                <input type="checkbox" checked="checked">
                <span class="checkmark"></span>
            </label>-->
                </div>
                <div class="col-20 vhc" v-text="packet.addr1"></div>
                <div class="col-20 vhc" v-text="packet.addr2"></div>
                <div class="col-5 vhc" v-text="packet.rssi"></div>
                <div class="col-5 vhc" v-text="packet.channel"></div>
                <div class="col-10 vhc" v-text="packet.pkt_type"></div>
                <div class="col-30 vhc" v-text="packet.pkt_type"></div>
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