new Vue({
    el: '#index',
    data: {
        component: "home",
    },
    methods: {
        setComponent: function(value){
            this.component = value;
            document.getElementById("sidebar-toggle").checked = false;
        },
        _sniffPackets: function() {
            axios
                .get("/sniffPackets")
                .then(this._packetHandler)
        },
        _packetHandler(result) {
            if (!result || !result.data) {
                console.log(result);
            }
        }

    },
    beforeMount() {
        this._sniffPackets();

    },
    beforeDestroy() {
    }
});