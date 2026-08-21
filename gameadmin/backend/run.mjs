#!/usr/bin/env zx

// import 'zx/globals'

// const date = await $`date`
// await $`echo Current date is ${date}.`


const services = [
"proxy",
"admin",
"gamecenter",
"gateway",
"slots",
"stats",
"fakeapp",
"fakeapptrans",
"BaoZang",
"Hilo",
"JinNiu",
"NiuBi",
"Roma",
"RomaX",
"TuZi",
"XingYunXiang",
"YingCaiShen",
"ZhaoCaiMao",
"MaJiang",
"MaJiang2",
"caipiao",
]

let i = 0
for (const s of services) {
    i++
    console.log(chalk.red(i.toString())+ ") " +  chalk.green(s))
}


// fs.exists()
