if (!window.location.origin) { // Some browsers (mainly IE) do not have this property, so we need to build it manually...
  window.location.origin = window.location.protocol + '//' + window.location.hostname + (window.location.port ? (':' + window.location.port) : '');
}

var origin = window.location.origin;

// options usage example
var options = {
		debug: true,
		devel: true,
		protocols_whitelist: ['websocket', 'xdr-streaming', 'xhr-streaming', 'iframe-eventsource', 'iframe-htmlfile', 'xdr-polling', 'xhr-polling', 'iframe-xhr-polling', 'jsonp-polling']
};

const rtLoyal = 0
const rtMinion = 1
const rtPercival = 2
const rtMerlin = 3
const rtOberon = 4
const rtMordred = 5
const rtAssassin = 6
const rtMorgana = 7

var roles = {};
roles[rtLoyal] = "Loyal servant of Arthur"
roles[rtMinion] = "Minion of Mordred"
roles[rtPercival] = "Percival"
roles[rtMerlin] = "Merlin"
roles[rtOberon] = "Oberon"
roles[rtMordred] = "Mordred"
roles[rtAssassin] = "Assassin"
roles[rtMorgana] = "Morgana"

var sock = new SockJS(origin+'/avalon', undefined, options);

function sendName(txt) {
    var name = prompt(txt, "");
    console.log("Sent " + name)
    sock.send(name)
}

sock.onopen = function() {
    sendName("Please enter your name")
};

var clients = [];

function displayClients() {
    var ul = document.getElementById("playerlist");
    ul.innerHTML = '';


    for (var i = 0; i < clients.length; i++) {
        var li = document.createElement("li");
        li.appendChild(document.createTextNode(clients[i]));
        ul.appendChild(li);
    }
}

sock.onmessage = function(e) {
    if (e.data.startsWith("CONNECT:")) {
        clients.push(e.data.split(":")[1]);
        displayClients();
    } else if (e.data.startsWith("DISCONNECT:")) {
        var index = clients.indexOf(e.data.split(":")[1]);
        if (index > -1) {
            clients.splice(index, 1);
        }
        displayClients();
    } else if (e.data.startsWith("INFO:")) {
        document.getElementById('info').innerHTML = e.data.split(":")[1]
    } else if (e.data.startsWith("ROLE:")) {
        var role = e.data.split(":")[1]; 
        document.getElementById('role').innerHTML = roles[parseInt(role)]
    } else if (e.data.startsWith("INVALID")) {
        sendName("Invalid name")
    } else if (e.data.startsWith("GO:")) {
        for (i = 2; i <= 7; i++) {
            document.getElementById(i.toString()).checked = false;
        }
        if (e.data.indexOf(",") !== -1) {
            var specials = e.data.split(":")[1].split(",")
            for (i = 0; i < specials.length; i++) {
                if (specials[i] !== "-1") {
                    console.log(specials[i])
                    document.getElementById(specials[i]).checked = true;
                }
            }
        }
    }
};

sock.onclose = function() {
	document.getElementById("status").innerHTML = "connection closed";
};

function getRole(role) {
    return document.getElementById(role.toString()).checked ? role : -1;
}

function go() {
    var merlin = getRole(rtMerlin);
    var percival = getRole(rtPercival);
    var mordred = getRole(rtMordred);
    var morgana = getRole(rtMorgana);
    var assassin = getRole(rtAssassin);
    var oberon = getRole(rtOberon)
    sock.send("GO:" + merlin + "," + percival + "," + mordred + "," + morgana + "," + assassin + "," + oberon);
}
