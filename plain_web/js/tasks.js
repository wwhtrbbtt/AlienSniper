const IP = "localhost"
const PORT = "8080"
const TOKEN = "GMWJGSAPGATLMODYLUMG"
let accounts = []


const get_task_html = (task) => {
    // const formatted = moment(task.time).format('MMMM Do YYYY, h:mm:ss a');
    return `
    <tr>
    <td>${task.name}</td> 
    <td>
        <span class="font-mono">
            ${task.timestamp || task.unix || "-"}
        </span>
    </td> 
    <td>${task.searches || '-'}</td>
    <td>${task.group || 'all'}</td>
  </tr>`
}

const add_task_html = (task) => {
    let html = get_task_html(task);
    document.getElementById("task_list").innerHTML += html;
}

const add_task_handler = () => {
    let name = document.getElementById('task_name').value.trim();
    const group = document.getElementById('task_group').value.trim();
    console.log(name)
    name = name.replace(/\t/g, '');
    console.log(name)
    const task = {
        type: "snipe",
        name: name,
        group: group || null
    }
    console.log(task);
    tasks.push(task);
    socket.send(new Packet('add_task', { task: task }).toJson());
    // add_task_html(task);
}

// make new connection
let socket = new WebSocket(`ws://${IP}:${PORT}/ws`)

// send auth packet on open
socket.onopen = event => {
    console.log('Connected to server', event);
    socket.send(new Packet('auth', { auth: TOKEN }).toJson());
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
            // console.log(packet.content.state);
            tasks = packet.content.state.tasks || [];
            tasks.forEach(task => {
                add_task_html(task);
            });
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
