

async function getPids() {
    const msg = await $`pgrep -l pp_`.quiet().nothrow()
    // console.log(msg.stdout)

    let arr = msg.stdout.split('\n')
    // console.log(arr)

    const pids = new Map()

    for(let line of arr) {
        const p = line.split(' ')
        if (p.length != 2) {
            continue
        }

        pids.set(p[1], p[0])
    }

    //console.log(pids)
    return pids
}

export async function status(services) {
    let i = 0
    console.log('-----------')
    // const services = await getServices()

    // const msg = await $`pgrep -l pg_`.quiet()
    // console.log(msg.stdout)

    const pids = await getPids()
    for (const s of services) {
        i++

        const procname = procName(s)

        // const run = await running(s)
        const run = pids.has(procname)

        const istr = chalk.blue( i.toString().padStart(3))

        let msg = ""
        if (run) {
            msg = `${istr}) ${s.padEnd(16)}  ${chalk.green('Running')}`
        } else {
            msg = `${istr}) ${s.padEnd(16)}  ${chalk.red('Stopped')}`
        }

        console.log(msg)
    }
}

export function procName(svr) {
    if (svr.length < 16) {
        return svr
    }

    return svr.substring(0, 15)
}

export async function running(svr) {
    let ok = false
    try {
        await $`pgrep -x ${procName(svr)}`.quiet()
        ok = true
    } catch {
        ok = false
    }
    return ok
}

export async function getServices() {
    // console.log("hello")
    const contents = await fs.readFile("services.yaml", { encoding: 'utf8' });
    // console.log(contents)
    const services = YAML.parse(contents)
    return services
}

export async function readServicesInput(services) {
    if (!services) {
        services = await getServices()
    }
    // console.log(services)
    await status(services)


    let input = await question('eg. all | dump | 1 pp_vs20olympx... Input: ')
    input = input.trim()
    if (input == 'all') {
        return services
    }

    if (input == 'dump') {
        return services.filter(s => fs.existsSync(path.join('../dump', s, 'combine.bson.gz')))
    }

    const arr = input.split(' ')

    const ans = []
    for (const i of arr) {
        if (i == '') {
            continue
        }

        if (i.startsWith('pp_')) {
            ans.push(i)
        } else {
            const idx = parseInt(i) -1
            const service = services[idx]
            ans.push(service)
        }
    }

    return ans
}

const mongourls = new Map([
        // ["127.0.0.1", "mongodb://myUserAdmin:doudou123456@127.0.0.1:27017/?authSource=admin"],
        ["doudou-test", "mongodb://myUserAdmin:doudoutestA1!@doudou-test:27017/?authSource=admin"],
        ["doudou-prod","mongodb://myUserAdmin:nc4IGBnmeHcJ3xZuIDePSY21y4Izq0brdRfaWmnhgMjmy6HZitZ6mwQCIa78cqH@doudou-prod:27017/?authSource=admin"],    
        ["dou-idr-prod", "mongodb://myUserAdmin:ohpeiGaP8shohphaihuoxeuv1Tei7eisaithaib6iiphu7iej2Yai5ee9ailoo8A@dou-idr-prod:27017/?authSource=admin"],
])
export function getmongourl(host) {
    if (host == "127.0.0.1" || host == "localhost") {
        return YAML.parse(fs.readFileSync("/data/game/bin/config/grpc_route.yaml").toString()).mongo
    }
    return mongourls.get(host)
}

export async function readHostsInput() {
    let i = 0
    const hosts = [
        "127.0.0.1",
        "doudou-test",
        "doudou-prod",
        "rp-hk-dev",
        "dou-idr-prod",
    ]
    for (const s of hosts) {
        i++

        const istr = chalk.blue(i.toString().padStart(3))

        const msg = `${istr}) ${s}`

        console.log(msg)
    }

    let input = await question('Input: ')
    input = input.trim()

    const idx = parseInt(input) -1

    const host = hosts[idx]
    if (host.endsWith("prod")) {
        let input = await question(`操作正式环境(${host})? [yes/no]: `)
        if (input != 'yes') {
            process.exit(0)            
        }
    }
    return host
}


function pad(timeEl, total = 2, str = '0') {
    return timeEl.toString().padStart(total, str)
}
export function formatData(timer) {
    const year = timer.getFullYear()
    const month = timer.getMonth() + 1 // 由于月份从0开始，因此需加1
    const day = timer.getDate()
    const hour = timer.getHours()
    const minute = timer.getMinutes()
    const second = timer.getSeconds()
    // return `${pad(year, 4)}-${pad(month)}-${pad(day)} ${pad(hour)}:${pad(minute)}:${pad(second)}`
    return `${pad(year, 4)}${pad(month)}${pad(day)}.${pad(hour)}${pad(minute)}${pad(second)}`
}


export async function runRemote(host, service) {
    const cwd = process.cwd()
    const pth = path.join("./", service)
    if (await fs.pathExists(pth)) {
        await $`go build  -gcflags="all=-N -l"  -o ../bin/${service} ./${pth}`
    }


    const bin = service
    // const host = "127.0.0.1"

    const remotepth = path.join(cwd, "../bin", bin+".new")
    await $`rsync -vcPz ../bin/${bin} root@${host}:${remotepth}`
 
    await $`ssh root@${host} "
    cd ${path.join(cwd, "../bin")}
    pkill -x ${procName(bin)}
    sleep 1
    mv ${bin} ${bin}.${formatData(new Date())}
    mv ${bin}.new ${bin}
    launch ./${bin}
    "`
}

export async function dumpGameSpinData(game) {
    if (fs.existsSync(path.join('../dump', game, 'combine.bson.gz'))) {
        return 
    }

    const mongouri = getmongourl("127.0.0.1")
    let spindata = 'rawSpinData'
    await $`mongodump --gzip -d ${game} -c ${spindata} -q '{"selected": true}' -o ../dump/ ${mongouri}`
    // await $`mongodump --gzip -d ${game} -c ${spindata} -o ../dump/ ${mongouri}`

    await $`mongodump --gzip -d ${game} -c combine -o ../dump/  ${mongouri}`
}
