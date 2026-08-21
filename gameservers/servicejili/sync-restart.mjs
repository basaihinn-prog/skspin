#!/usr/bin/env zx

import 'zx/globals'
import {getServices, readServicesInput, formatData, readHostsInput, runRemote, dumpGameSpinData} from './run-comm.mjs'

$.verbose = true

const cwd = process.cwd()
console.log("cwd", cwd)
let debugmode = false

// cd('../')

// $`ping www.baidu.com`.pipe(process.stdout)

 
const host = debugmode? 'doudou-test': await readHostsInput()
await $`rsync -vc *.mjs services.yaml root@${host}:${cwd}`

while (true) {
    const choices = debugmode? ['jili_2_csh'] : await readServicesInput()
    console.log(choices)
    // break

    for (const s of choices) {
        // await syncOne(s)
        await dumpGameSpinData(s)
        await runRemote(host, s)
    }

    console.log(choices.join(' '))
}
