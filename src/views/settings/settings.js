Vue.component('settings', {
    data: function () {
        return {
            saving: false,
            loading: false,
            accesspoint: {
                ssid: null,
                passwd: null,
                channel: null,
                hidden: null,
            },
            apScanner: {
                interval: null,
                deep: null,
                async: null,
                channel: null,
                hop: null
            },
            packetScanner: {
                interval: null,
                channel: null,
                hop: null
            },
            deauther: {
                interval: null,
                channel: null,
            }
        }
    },
    methods: {
        _getConfig: function () {
            this.loading = true;
            axios
                .get("/getConfig")
                .then(this._setConfig)
        },
        _setConfig: function (result) {
            try {
                if (!result || !result.data) {
                    return;
                }

                this.accesspoint = result.data.accesspoint;
                this.apScanner = result.data.apScanner;
                this.packetScanner = result.data.packetScanner;
                this.deauther = result.data.deauther;

                $(function () {
                    custom_input();
                });
            } catch (err) {
                console.log(err);
            } finally {
                this.loading = false;
            }

        },
        saveConfig: function () {
            //Find a better way to force int
            this.saving = true;

            this.accesspoint.channel = parseInt(this.accesspoint.channel);

            this.apScanner.channel = parseInt(this.apScanner.channel);
            this.apScanner.interval = parseInt(this.apScanner.interval);

            this.packetScanner.channel = parseInt(this.packetScanner.channel);
            this.packetScanner.interval = parseInt(this.packetScanner.interval);

            this.deauther.channel = parseInt(this.deauther.channel);
            this.deauther.interval = parseInt(this.deauther.interval);

            var data = {
                accesspoint: this.accesspoint,
                apScanner: this.apScanner,
                packetScanner: this.packetScanner,
                deauther: this.deauther
            };

            axios
                .post('/setConfig', JSON.stringify(data))
                .then(this.configConfirmation);
        },
        configConfirmation: function(result) {
            try {
                if (!result || result.status !== 200) {
                    toastr("Something happened while saving, please try again...", "error", 1000);
                    return;
                }

                toastr("Save successful", "success", 1000);
            } catch(err) {
                console.log(err)
            } finally {
                this.saving = false;
            }
        }
    },
    mounted() {
        this._getConfig();
    },
    template: `
        <div class="flex-container col-100 hc no-touch-top">
            <div class="panel-content vhc" v-if="saving || loading" style="height:110vh; background: rgba(21,21,21,.7); position:fixed; top: 0px !important;">
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
            <div class="panel col-100 no-touch-top">
                <div class="panel-header">Flash</div>
                <div class="panel-content">
                    <div class="input-group">
                        <label for="save" class="toggle-lbl vh-center" style="width:100%; margin:0;" v-on:click="saveConfig()">
                            <i class="material-icons icon-sm vhc">publish</i>
                        </label>
                    </div>
                </div>
            </div>
            <div class="panel col-50 no-touch-top" style="vertical-align: top;">
                <div class="panel-header">Accesspoint</div>
                <div class="panel-content">
                    <div class="input-group">
                        <label for="ssid" id="panel-label" class="dyn-input-label">SSID</label> 
                        <input placeholder="SSID Name" type="text" id="panel-input" name="ssid" class="dyn-input" v-model="accesspoint.ssid">
                    </div>
                    <div class="input-group">
                        <label for="password" id="panel-label" class="dyn-input-label">Password</label> 
                        <input type="text" placeholder="Some Password" id="panel-input" name="password" class="dyn-input" v-model="accesspoint.passwd">
                        <!--<i class="material-icons icon-xs" style="position: relative; color:#212121">visibility</i>-->
                    </div>
                    <div class="input-group">
                        <label for="channel" id="panel-label" class="dyn-input-label">Channel</label> 
                        <input type="number" placeholder="1-12" id="panel-input" name="channel" class="dyn-input" v-model="accesspoint.channel">
                    </div>
                    <div class="input-group">
                        <label for="hidden" id="panel-label" class="dyn-input-label">Hidden</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="hidden" class="dyn-input" v-model="accesspoint.hidden">
                    </div>
                </div>
            </div>
            <div class="panel col-50 no-touch-top " style="vertical-align: top;">
                <div class="panel-header">Accesspoint Scanner</div>
                <div class="panel-content">
                    <div class="input-group">
                        <label for="interval" id="panel-label" class="dyn-input-label">Interval</label> 
                        <input type="text" placeholder="1000" id="panel-input" name="interval" class="dyn-input" v-model="apScanner.interval">
                    </div>
                    <div class="input-group">
                        <label for="deep" id="panel-label" class="dyn-input-label">Deep</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="deep" class="dyn-input" v-model="apScanner.deep">
                    </div>
                    <div class="input-group">
                        <label for="async" id="panel-label" class="dyn-input-label">Async</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="async" class="dyn-input" v-model="apScanner.async">
                    </div>
                    <div class="input-group">
                        <label for="channel" id="panel-label" class="dyn-input-label">Channel</label> 
                        <input type="number" placeholder="1-12" id="panel-input" name="channel" class="dyn-input" v-model="apScanner.channel">
                    </div>
                    <div class="input-group">
                        <label for="hop" id="panel-label" class="dyn-input-label">Hop</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="hop" class="dyn-input" v-model="apScanner.hop">
                    </div>
                </div>
            </div>
            <div class="panel col-50 no-touch-top" style="vertical-align: top;">
                <div class="panel-header">Packet Scanner</div>
                <div class="panel-content">
                    <div class="input-group">
                        <label for="interval" id="panel-label" class="dyn-input-label">Interval</label> 
                        <input placeholder="1000" type="text" id="panel-input" name="interval" class="dyn-input" v-model="packetScanner.interval">
                    </div>
                    <div class="input-group">
                        <label for="channel" id="panel-label" class="dyn-input-label">Channel</label> 
                        <input type="number" placeholder="1-12" min="1" max="12" id="panel-input" name="channel" class="dyn-input" v-model="packetScanner.channel">
                    </div>
                    <div class="input-group">
                        <label for="hop" id="panel-label" class="dyn-input-label">Hop</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="hop" class="dyn-input" v-model="packetScanner.hop">
                    </div>
                </div>
            </div>
            <div class="panel col-50 no-touch-top" style="vertical-align: top;">
                <div class="panel-header">Deauther</div>
                <div class="panel-content">
                    <div class="input-group">
                        <label for="interval" id="panel-label" class="dyn-input-label">Interval</label> 
                        <input placeholder="1000" type="text" id="panel-input" name="interval" class="dyn-input" v-model="deauther.interval">
                    </div>
                    <div class="input-group">
                        <label for="channel" id="panel-label" class="dyn-input-label">Channel</label> 
                        <input type="number" placeholder="1-12" min="1" max="12" id="panel-input" name="channel" class="dyn-input" v-model="deauther.channel">
                    </div>
                </div>
            </div>
            <div class="panel col-100 no-touch-top">
                <div class="input-group">
                    <label for="save" class="toggle-lbl vh-center" style="width:100%; margin:0;" v-on:click="saveConfig()"><span class="v-center" style="text-transform: uppercase;">Save</span></label>
                </div>
            </div>
        </div>
    `
});