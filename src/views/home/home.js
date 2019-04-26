/*var vm = this;*/

Vue.component('home', {
    data: function() {
        return {
            particle_config: {
                "particles": {
                    "number": {
                    "value": 100,
                    "density": {
                        "enable": true,
                        "value_area": 800
                    }
                    },
                    "color": {
                    "value": "#04f886"
                    },
                    "shape": {
                    "type": "circle",
                    "stroke": {
                        "width": 0,
                        "color": "#000000"
                    },
                    "polygon": {
                        "nb_sides": 5
                    },
                    "image": {
                        "src": "img/github.svg",
                        "width": 100,
                        "height": 100
                    }
                    },
                    "opacity": {
                    "value": 0.5,
                    "random": false,
                    "anim": {
                        "enable": false,
                        "speed": 0.5,
                        "opacity_min": 0.1,
                        "sync": false
                    }
                    },
                    "size": {
                    "value": 3,
                    "random": true,
                    "anim": {
                        "enable": false,
                        "speed": 10,
                        "size_min": 0.1,
                        "sync": false
                    }
                    },
                    "line_linked": {
                    "enable": true,
                    "distance": 250,
                    "color": "#CC14AB",
                    "opacity": 0.4,
                    "width": 1
                    },
                    "move": {
                    "enable": true,
                    "speed": 3,
                    "direction": "none",
                    "random": false,
                    "straight": false,
                    "out_mode": "out",
                    "bounce": false,
                    "attract": {
                        "enable": false,
                        "rotateX": 600,
                        "rotateY": 1200
                    }
                    }
                },
                "interactivity": {
                    "detect_on": "canvas",
                    "events": {
                    "onhover": {
                        "enable": true,
                        "mode": "repulse"
                        },
                    "onclick": {
                        "enable": true,
                        "mode": "bubble"
                        },
                        "resize": true
                    },
                    "modes": {
                    "grab": {
                        "distance": 140,
                        "line_linked": {
                        "opacity": 1
                        }
                    },
                    "bubble": {
                        "distance": 400,
                        "size": 5,
                        "duration": 2,
                        "opacity": 8,
                        "speed": 3
                    },
                    "repulse": {
                        "distance": 150,
                        "duration": 0.6
                    },
                    "push": {
                        "particles_nb": 4
                    },
                    "remove": {
                        "particles_nb": 2
                    }
                    }
                },
                "retina_detect": true
            }
        }
    },
    methods: {
        particles() {
            particlesJS("particles-js", this.particle_config);
        }
    },
    mounted() {
        this.particles();
    },
    beforeMount() {
        var css = '.navbar { background: rgba(0,0,0,.8) !important; } body { background: black;} .nav-menu { background:rgba(0,0,0,.8); }',
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
        <div class="l1 col-100 vhc"  id="particles-js" style="position:absolute; top: 0;">
            <!--<img src="/static/images/0_lW3k05c1RPFYo9JF.jpg" style="width:100%;">-->
            <!--<div class="main-screen-text" style="position: absolute; font-size: 40px; -webkit-background-clip: text; -webkit-text-fill-color: transparent;">Hack the Planet</div>-->

        
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

