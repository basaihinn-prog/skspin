export interface AboutInterface {
    logo: string,
    link: LinkInterface[],
    connect: ConnectInterface[],
    bottomText: string,
    runStatus: RunStatusInterface
}


export interface LinkInterface extends ConnectInterface{

    href: string

}

export interface ConnectInterface {
    image: string
    label: string
    href: string
}

export interface RunStatusInterface {
    disabled: string
    enabled: string
}
