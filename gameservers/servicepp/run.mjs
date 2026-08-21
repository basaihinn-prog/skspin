#!/usr/bin/env zx

//import 'zx/globals'
import {status, getServices, readServicesInput, procName} from './run-comm.mjs'


while (true) {
    const services = await getServices()
    console.log(services)

    const choices = await readServicesInput(services)

    for (const s of choices) {
        await run(s)
    }
}

async function run(service) {
    console.log("run service", chalk.green(service))
    // const pth = path.join("../service", service)
    const pth = service 
    if (await fs.pathExists(pth)) {
        // console.log(`build ${service}`)
        await $`go build  -gcflags="all=-N -l"  -o ../bin/ ./${pth}`
    }

    cd('../bin')
    await $`pkill -x ${procName(service)}`.nothrow()
    await sleep(1000)
    await $`launch ./${service}`

    cd('../servicepp')
}
