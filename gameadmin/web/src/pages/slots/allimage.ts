
// dyImport("NiuBi/1.png")
// export function dyImport(name) {
//     let url = new URL(`/src/assets/slots/${name}`, import.meta.url).href
//     console.log(name, url)
//     return url
// }

let imagesOrigin = import.meta.glob('/src/assets/slots/*/*.*',  { import: 'default', eager: true, })
let imagesOrigins = import.meta.glob('/src/assets/game/*/*.*',  { import: 'default', eager: true, })

// images["NiuBi"]["1"] = "NiuBi/1.png"
let images = {}

for (let k in imagesOrigin) {
    let ma = k.match(/(\w+)\/(\w+)\.(\w+)/)

    let gameName = ma[1]

    if (!images[gameName]) {
        images[gameName] = {}
    }

    images[gameName][ma[2]] = imagesOrigin[k]
}
console.log(images)

export function getImages(gameName : string) : JsMap<string, string> {
    return images[gameName]
}


// dyImport("NiuBi/1.png")
export function dyImport(name:string) :string{
    console.log(name,'dyImport------------------name');
    return imagesOrigin['/src/assets/slots/' + name]as string
}
export function gameImgImport(name:string) :string{
    let url = imagesOrigin['/src/assets/game/' + name]as string
    return url
}
