export interface LanguageConfigInterface {
    hotGame?: string
    link?: string
    contactInformation?: string
    language?: string
    festival?: string,
    animal?: string,
    newGame?: string,
    searchPrefix?: string,
    searchSuffix?: string,
    hitTitle?: string,
    hitDescription?: string,
    searchPlaceHolder?: string,
    sortList?: Array<Record<string, any>>,
    theme?: string,
    manufacturer?: string,
    knowAll?: string,
    gameList: Object
    RTPSetting: RTPSettingInterface
}


export interface RTPSettingInterface {
    title: string
    label: string
    sure: string
    cancel: string
    comm: string
}

