  chrome.storage.local.get('ConnectedResp', function(result) {
    console.log('ConnectedResp Value currently is ' + result);


    if (!result || !result.ConnectedResp) 
    {
      document.getElementById("connected").style.visibility = "hidden";
      document.getElementById("disconnected").style.visibility = "visible";
    }
    else 
    {
      
      document.getElementById("disconnected").style.visibility = "hidden";
      document.getElementById("connected").style.visibility = "visible";

      connected = result.ConnectedResp;

      if (connected.VpnType===0) 
        document.getElementById("VpnType").innerHTML = "OpenVPN";
      else
        document.getElementById("VpnType").innerHTML = "WireGuard";

      document.getElementById("ClientIP").innerHTML = connected.ClientIP;
      document.getElementById("ServerIP").innerHTML = connected.ServerIP;
      }
  });

  chrome.storage.local.get('AccountStatusResp', function(result) {
    console.log('AccountStatusResp Value currently is ' + result);
  });

