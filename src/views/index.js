new Vue({
	el: '#index',
	data: {
		component: "home",
	},
	methods: {
		setComponent: function(value){
			this.component = value;
			document.getElementById("sidebar-toggle").checked = false;
			
			var navbar = document.getElementsByClassName("navbar")[0];
			navbar.style.background = "#1b2127"
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
		},
		scrollHandler() {
			var body = document.getElementsByTagName("body")[0];
			var navbar = document.getElementsByClassName("navbar")[0];
		
			if (body.scrollTop >= 50) {
				navbar.style.background = "#1b2127";//"rgba(0,0,0,.9)";
			} else {
				navbar.style.background = "#1b2127";
			}
		}
	},
	beforeMount() {
		this._sniffPackets();
	},
	beforeDestroy() {
	},
	created() {
		window.addEventListener('scroll', this.scrollHandler);
	}
});

function toastr(message, type, timeout) {
	var body = document.getElementsByTagName("body")[0];

	var t = document.createElement("div");
	t.setAttribute("class", `toast ${type} vhc tc`);
	t.innerHTML = message;

	body.appendChild(t);

	setTimeout(function () {
		removeToast(t);
	}, timeout)
}

function removeToast(el) {
	return new Promise((resolve, reject) => {
		var body = document.getElementsByTagName("body")[0];
		var inter = setInterval(function () {
			var opa = parseFloat(getComputedStyle(el).opacity)
			el.style.opacity = opa - 0.05;

			if (opa <= 0) {
				clearInterval(inter);
				body.removeChild(el);
				resolve(true);
			}
		}, 100);
	});
}

function fadeOut(el) {
	return new Promise((resolve, reject) => {
		var inter = setInterval(function () {
			var opa = parseFloat(getComputedStyle(el).opacity)
			el.style.opacity = opa - 0.05;

			if (opa <= 0) {
				clearInterval(inter);
				el.style.display = "none";
				resolve(true);
				//body.removeChild(el);
			}
		}, 100);
	});
}

function fadeIn(el) {
	//var body = document.getElementsByTagName("body")[0];
	return new Promise((resolve, reject) => {
		var inter = setInterval(function () {
			var opa = parseFloat(getComputedStyle(el).opacity)
			el.style.opacity = opa + 0.1;

			if (opa >= 1) {
				clearInterval(inter);
				el.style.display = "initial";
				toastr("The endpoint specified could not be found.")
				resolve(true);
				//body.removeChild(el);
			}
		}, 100);
	});
}

function custom_input() {
	$(".input-group input").focus(function () {
		$(this).parent(".input-group").each(function () {
			$("label", this).css({
				"font-size": "13px",
				"color": "#b7bdd1"
			})
		})
	}).blur(function () {
		if ($(this).val() === "") {

			$(this).css({
				"background": "#333333",
			})

			$(this).parent(".input-group").each(function () {
				$("label", this).css({
					"font-size": "15px",
				})
			});
		} else {
			$(this).css({
				"box-shadow": "none",
				"background": "#353a44",
			})

			$(this).parent(".input-group").each(function () {
				$("label", this).css({
					"color": "#CC14AB"
				})
			})
		}
	});

	Array.from($(".input-group input")).forEach( i => {
		i.focus();
		i.blur();
	});
}