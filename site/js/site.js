get_list_with_active_nodes();
window.setInterval(get_list_with_active_nodes, 5000);
const mq = window.matchMedia('(prefers-color-scheme: dark)')
if (mq.matches) {
    document.documentElement.classList.add('dark')
}
function toggleTheme() {
    const root = document.documentElement
    root.classList.toggle('dark')
    if (window.ccd) {
        const isDark = root.classList.contains('dark')
        window.ccd.changeTheme(isDark ? 'dark' : 'light')
    }
}

function get_list_with_active_nodes() {
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open('GET', "v1/collect/status", true);
    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState === 4) {
            if(xmlhttp.status === 200) {
                let obj = JSON.parse(xmlhttp.responseText);
                let data = obj.data;
                if (data === null) {
                    document.getElementById("running-nodes").innerHTML = "Data is not currently being collected";
                    return
                }
                let result = "<strong>Data is currently being collected for:</strong><ul>";
                for (let fsym in data){
                    let nodes = data[fsym];
                    let litag = "<li>" + fsym + " to: "
                    for (let tsym in nodes){
                        if (nodes[tsym]["interval"] === undefined) {
                            nodes[tsym]["interval"] = "interval not set (wss)"
                        }
                        litag += tsym+ ":" + nodes[tsym]["interval"] + ", ";
                    }
                    result += litag.slice(0, -2) + "</li>";
                }
                result += "</ul>"
                document.getElementById("running-nodes").innerHTML = result
            }
        }
    };
    xmlhttp.send(null);
}