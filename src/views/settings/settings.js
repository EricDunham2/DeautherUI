Vue.component('settings', {
    data: function () {
        return {
            accesspoint: {
                ssid: null,
                passwd: null,
                channel: null,
                hidden: null,
            },
            scanner: {
                interval: null,
                deep: null,
                async: null,
                channel: null,
                hop: null
            }
        }
    },
    methods: {
        _getConfig: function () {
            axios
                .get("/getConfig")
                .then(this._setConfig)
        },
        _setConfig: function (result) {
            if (!result || !result.data) {
                return;
            }

            this.accesspoint = result.data.accesspoint;
            this.scanner = result.data.scanner;

            $(function () {
                custom_input();
            });
        },
        saveConfig: function () {
            var data = {
                accesspoint: this.accesspoint,
                scanner: this.scanner
            };

            axios
                .post('/setConfig', JSON.stringify(data))
                .then(this.configConfirmation);
        },
        configConfirmation: function(result) {
            if (!result || result.status !== 200) {
                toastr("Something happened while saving, please try again...", "error", 1000);
                return;
            }

            toastr("Save successful", "success", 1000);
        }
    },
    mounted() {
        this._getConfig();
    },
    template: `
        <div class="flex-container col-100 hc no-touch-top">
            <div class="panel col-50 no-touch-top">
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
                        <input type="text" placeholder="1-12" id="panel-input" name="channel" class="dyn-input" v-model="accesspoint.channel">
                    </div>
                    <div class="input-group">
                        <label for="hidden" id="panel-label" class="dyn-input-label">Hidden</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="hidden" class="dyn-input" v-model="accesspoint.hidden">
                    </div>
                </div>
            </div>
            <div class="panel col-50 no-touch-top">
                <div class="panel-header">Scanner</div>
                <div class="panel-content">
                    <div class="input-group">
                        <label for="interval" id="panel-label" class="dyn-input-label">Interval</label> 
                        <input type="text" placeholder="1000" id="panel-input" name="interval" class="dyn-input" v-model="scanner.interval">
                    </div>
                    <div class="input-group">
                        <label for="deep" id="panel-label" class="dyn-input-label">Deep</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="deep" class="dyn-input" v-model="scanner.deep">
                    </div>
                    <div class="input-group">
                        <label for="async" id="panel-label" class="dyn-input-label">Async</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="async" class="dyn-input" v-model="scanner.async">
                    </div>
                    <div class="input-group">
                        <label for="channel" id="panel-label" class="dyn-input-label">Channel</label> 
                        <input type="text" placeholder="1-12" id="panel-input" name="channel" class="dyn-input" v-model="scanner.channel">
                    </div>
                    <div class="input-group">
                        <label for="hop" id="panel-label" class="dyn-input-label">Hop</label> 
                        <input type="text" placeholder="false/true" id="panel-input" name="hop" class="dyn-input" v-model="scanner.hop">
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