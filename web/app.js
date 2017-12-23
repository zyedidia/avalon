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

var roles = {};
roles[0] = "Loyal servant"
roles[1] = "Minion of Mordred"
roles[2] = "Percival"
roles[3] = "Merlin"
roles[4] = "Oberon"
roles[5] = "Mordred"
roles[6] = "Assassin"
roles[7] = "Morgana"

var sock = new SockJS(origin+'/avalon', undefined, options);

sock.onopen = function() {
    var name = prompt("Please enter your name", "");
    sock.send(name)
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
    }
};

sock.onclose = function() {
	document.getElementById("status").innerHTML = "connection closed";
};

function go() {
    sock.send("GO")
}
