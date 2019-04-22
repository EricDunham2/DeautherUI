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
        <div class="l1 col-100 vhc" >
            <img src="/static/images/0_lW3k05c1RPFYo9JF.jpg" style="width:100%;">
            <div class="main-screen-text" style="position: absolute; top: 400px; font-size: 40px; -webkit-background-clip: text; -webkit-text-fill-color: transparent;">Hack the Planet</div>

        
        <!--<div class="alert-title">
                
                
                
                
                
                
                
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
    </div>
    `
});

