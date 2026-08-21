export {}
declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        dateFormater: Function,
        dateSecondFormater: Function,
        date1Formater: Function,
        goldFormater: Function,
        goldSignedFormater: Function
        percentFormatter: Function,
        boolFormatter: Function,
        goldTenThousand: Function,
        imageUrl: Function,
        copyText: Function,
        openIMG: Function,
        exportAction: Function,
        percentage: Function,
    }
}
