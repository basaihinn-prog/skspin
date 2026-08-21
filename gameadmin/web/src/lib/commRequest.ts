import {AddSystemLogParam, Log, LOG_PRIMARY_MODULE} from "@/api/adminpb/log";
import {Client} from "@/lib/client";

export const AddLog = async (OperateMod:number, OperateType:number, oldData:any, newData:any, appId: string, contentOptions?: any) => {




    let queryData:AddSystemLogParam = <AddSystemLogParam>{
        OperateMod,
        OperateType,
        OperateContent: JSON.stringify({
            oldData,
            newData,
            appId,
            ...contentOptions
        }),
    }
    await Client.Do(Log.AddSystemLog, queryData)

}


// const getIpAdress = () => {
//
//     var myPeerConnection = window.RTCPeerConnection || window?.mozRTCPeerConnection || window?.webkitRTCPeerConnection;
//     var pc = new myPeerConnection({iceServers:[]}), noop = function(){};
//     debugger
//     var localIPs = {};
//     pc.createDataChannel("");
//     pc.createOffer().then(function (sdp) {
//         sdp.sdp.split('\n').forEach(function (line) {
//             if (line.indexOf('candidate') < 0) return;
//             line.match(ipRegex).forEach(function (ip) {
//                 localIPs[ip] = true;
//             });
//         });
//         pc.setLocalDescription(sdp, noop, noop);
//     }).catch(function (reason) {
//         console.log(reason);
//     });
//     var ipRegex = /([0-9]{1,3}(\.[0-9]{1,3}){3})/;
//     var ips = Object.keys(localIPs);
//
//
//     return ips[0]
// }
