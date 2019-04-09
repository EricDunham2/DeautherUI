new Vue({
    el: '#index',
    data: {
        ssid: null,
        password: null,
        channel: null,
        hidden: null,
        interval: null,
        hiddenScan: null,
        async: null,
        channelScan: null,
        component: "home",
    },
    methods: {
        setComponent: function(value){
            this.component = value;
            document.getElementById("sidebar-toggle").checked = false;
        }
    },
    beforeMount() {
    },
    beforeDestroy() {
    }

});