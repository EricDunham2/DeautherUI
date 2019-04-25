new Vue({
    el: '#index',
    data: {
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