// background.js


chrome.runtime.onInstalled.addListener(() => {
  console.log('IVPN plugin loaded');
 
  chrome.storage.local.set({ConnectedResp: null});

  createWebSocketConnection();

});

var websocket = undefined;
function createWebSocketConnection() {
    console.log("createWebSocketConnection")

    if('WebSocket' in window){
        chrome.storage.local.get("instance", function(data) {
            connect();
        });
    }
}

//Close the websocket connection
function closeWebSocketConnection(username) {
    if (websocket != null || websocket != undefined) {
        websocket.close();
        websocket = undefined;
    }
}

function connect() {
    console.log("connect")
    if (websocket === undefined) {
        websocket = new WebSocket("ws://127.0.0.1:7294");
        
    }  

    websocket.onopen = function(e) {
        console.log("Connected to daemon!");
        
        let cmd = {Command: "Hello", GetStatus: true, Idx:1}
        websocket.send(JSON.stringify(cmd))


        let cmd2 = {Command: "AccountStatus", Idx:2}
    };

    websocket.onmessage = function(event) {
      console.log(`[message] Data received from server: ${event.data}`);

      let resp = JSON.parse(event.data)

      //chrome.browserAction.setBadgeText({ "text": "test" });

      switch (resp.Command) {
        case "ConnectedResp":
          chrome.browserAction.setIcon({path:  "/images/connected.png"});
          chrome.browserAction.setTitle({ "title": "IVPN: Connected" });  
          chrome.storage.local.set({ ConnectedResp: resp });
        break;

        case "DisconnectedResp":
          chrome.browserAction.setIcon({path:  "/images/disconnected.png"});
          chrome.browserAction.setTitle({ "title": "IVPN: Disonnected" });  
          chrome.storage.local.set({ ConnectedResp: null });
        break;

        case "VpnStateResp":
        break;

        case "AccountStatusResp":
          chrome.storage.local.set({ AccountStatusResp: resp });
        break;

        default:
        break;
      }
      
   
    };

    websocket.onclose = function(event) {
      if (event.wasClean) {
        console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
      } else {
        // e.g. server process killed or network down
        // event.code is usually 1006 in this case
        console.log('[close] Connection died');
      }
    };

    websocket.onerror = function(error) {
      console.log(`[error] ${error.message}`);
    }; 
}