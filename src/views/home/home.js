/*var vm = this;*/

Vue.component('home', {
    template: `
    <div class="flex-container col-100 no-touch-top vhc">
        <div class="alert l1 col-60">
            <div class="alert-title">
                <div class="col-100 hc"><i class="material-icons icon-lg" style="width: 96px; height: 96px; font-size: 96px;">warning</i></div>
                <div class="col-100 hc">Warning!</div>
            </div>
            <br>
            <br>
            <div class="alert-msg hc col-100 tc">
                <span class="hc col-80">
                    This application is meant for testing purposes only and to explore the capabilities of the ESP8266
                </span>
            </div>
        </div>
    </div>
    `
});