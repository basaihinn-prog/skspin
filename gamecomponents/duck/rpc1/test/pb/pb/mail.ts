export enum Status{
    Error=0,
}
export interface SubMessage{
    Pid:number// PID 
}
export interface SendMailReq{
    Pid:number// PID
    Title:string
    Content:string// content
    Status:Status
    Sub:SubMessage
}
export interface SendMailResp{
}

export class MailServer {
    // fasong
    // asd
    public static SendMail(client : any, req :SendMailReq) : Promise<[SendMailResp,any]> {// hello
        return client.sendAsync("Mail", "SendMail", req)
    }
}


