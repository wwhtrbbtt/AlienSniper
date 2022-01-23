const t = getToken();
const add_logs = (acc) => {
    let content = ""

    // loop through the logs
    for (const x of acc.sends) {
        // for each log (for each past snipe)

        // assume a background of red
        let bg = "error"

        // create HTML for the requests
        let logHTML = ""
        x.content.forEach(l => {
            console.log(l)
            if (l.statuscode == 200) { bg = "success" }
            logHTML += `
            <p>
                <span class="text-${(l.statuscode == 200) ? "green" : "red" }-500">[${l.statuscode}]</span>
                <span>${l.timestamp}</span>
            </p>
            `
        })

        // create HTML for the log thing
        content += `<div class="bg-${bg} p-2 rounded-md shadow mt-4"><details>
            <summary>
                <h1 class="text-md font-mono">${x.email}</h1>
                <h2 class="text-sm font-mono">${x.ip}</h2>
            </summary>
            <div class="font-mono text-sm mt-2 p-3 bg-neutral ">
                <p>
                    ${logHTML}
                </p>
            </div>
        </details></div>`
    }


    statusC = (acc.success) ? "Yes" : "No"
    bgC = (acc.success) ? "green" : "red"

    return `<div id="${acc.name}" class="modal modal-closed">

    <div class="modal-box">
        <h1 class="text-2xl">Logs for
            <span class="kbd">${acc.name}</span>
        </h1>

        <p class="text-md">${acc.requests} requests</p>
        <p class="text-md">Delay: ${acc.delay}</p>
        <p class="text-md">Success: ${acc.success}</p>

        <div class="m-2 p-5 ">

        ${content}  
    
        </div>
        <div class="modal-action">
            <label onclick="modalClose('${acc.name}', 'modal-open')" class="btn">Done</label>
        </div>
    </div>
    </div>
    
    <tr class="hover" onclick="modalOpen('${acc.name}', 'modal-open')">
    <td class="row-data">${acc.name}</td>
    <td class="row-data">${acc.requests}</td>
    <td class="row-data">${acc.delay}</td>
    <td class="row-data">
        <span class="text-${bgC}-500">
        ${statusC}
    </span>
    </td>
</tr>`
}


// make new connection
let socket = new WebSocket(`ws://${t.ip}:${t.port}/ws`)

// send auth packet on open
socket.onopen = event => {
    console.log('Connected to server', event);
    socket.send(new Packet('auth', { auth: t.token }).toJson());
    socket.send(new Packet('get_state', {}).toJson());
}

// handle incoming packets
socket.onmessage = (event) => {
    let packet = JSON.parse(event.data);

    switch (packet.type) {
        case 'error':
            console.log(packet.content.error);
            break;
        case 'auth':
            console.log(packet.content.auth);
            break;
        case 'state_response':
            accs = packet.content.state.logs
            for (const x of accs) {
                document.getElementById("table1").innerHTML += add_logs(x);
            }
            break;
        case 'config':
            console.log(packet.content.config);
            break;
        case 'add_task_response':
            add_task_html(packet.content.task);
            // console.log(packet.content.response);
            break;
        default:
            console.log(packet);
    }
}