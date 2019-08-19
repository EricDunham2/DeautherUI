Vue.component('logs', {
    data: function() {
        return {
            logs:[],
        }
    },
    methods: {
        _getLogs: function() {
            axios
            .get("/getLogs")
            .then(this._setLogs)
        },
        _setLogs: function(response) {
            if (!response.data) { return; }
            this.logs = [];
            
            response.data.forEach(log => {
                var parts = log.split("]")
                this.logs.unshift({time:`${parts[0]}]`, data:parts[1]});
            });
        }
    },
    beforeMount() {
        this.logInterval = setInterval(() => {
            this._getLogs();
        }, 1000);
    },
    beforeDestroy() {
        clearInterval(this.logInterval);
    },
    template: `
        <div class="flex-container col-100 no-touch-top vhc" style="max-height:100%; overflow-y:scroll;" id="logger">
            <div class="col-80" id="logs" style="padding-left:10px;" v-for="log in logs">
                <span v-text="log.time" style="color:#1b6b52; font-weight: bold;"></span>
                <span v-text="log.data" style="color:#b7bdd1;"></span>
            </div>
        </div>
    `
});