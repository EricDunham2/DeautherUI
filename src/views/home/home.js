/*var vm = this;*/

Vue.component('home', {
    beforeMount() {
        var css = '.navbar { background: transparent !important; } body { background: black;}',
        head = document.head || document.getElementsByTagName('head')[0],
        style = document.createElement('style');
        style.id = "removeMe";

        head.appendChild(style);
        
        style.type = 'text/css';
        if (style.styleSheet){
        // This is required for IE8 and below.
        style.styleSheet.cssText = css;
        } else {
        style.appendChild(document.createTextNode(css));
        }
    },
    beforeDestroy() {
        var head = document.head;
        Array.from(head.children).forEach(child => { if (child.id === "removeMe") child.remove()} )
    },
    template: `
    <div class="flex-container col-100 vhc">
        <div class="l1 col-100" style="background-image: url('/static/images/0_lW3k05c1RPFYo9JF.jpg'); height: 100vw; min-height: 750px; background-attachment: fixed; background-position: center; background-repeat: no-repeat;">                <!--<div class="alert-title">
                    <div class="col-100 hc"><i class="material-icons icon-lg" style="width: 96px; height: 96px; font-size: 96px;">warning</i></div>
                    <div class="col-100 hc">Warning!</div>
                </div>
                <br>
                <br>
                <div class="alert-msg hc col-100 tc">
                    <span class="hc col-80">
                        This application is meant for testing purposes only and to explore the capabilities of the ESP8266
                    </span>
                </div>-->
        </div>
        <div style="font-size: 40px; background: transparent; background: linear-gradient(to left, #ff5770, #e4428d, #c42da8, #9e16c3, #6501de, #9e16c3, #c42da8, #e4428d, #ff5770); -webkit-background-clip: text; -webkit-text-fill-color: transparent;">Hack the Planet</div>
    </div>
    `
});