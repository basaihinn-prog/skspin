let fs = require("fs")
let path = require("path")
let readline = require('readline'); 


 // 递归获取文件列表
 function getfiles(rootdir, ext) {
    let files = []
    let match

    match = function (dir) {
        let subpaths = fs.readdirSync(dir)
        for (let i = 0; i < subpaths.length; ++i) {
            let p = subpaths[i]
            if (p[0] === '.') {
                continue;
            }

            p = path.join(dir, p);
            let stat = fs.statSync(p);
            if (stat.isDirectory()) {
                match(p);
            }
            else if (stat.isFile()) {
                if (path.extname(p) == ext) {
                    files.push(p)
                }
            }
        }
    }

    match(rootdir)

    return files
}


let files = getfiles("./pb/_gen/", ".go")

for (let i in files) {
    let file = files[i]

    let content = fs.readFileSync(file, "utf-8")

    // remove json tag
    content = content.replace(/json:"\w+,\w+"/g, "")

    // @gotags: bson:"_id"
    content = content.replace(/`\s*\/\/\s*@gotags:([^\n]+)/g, "$1`")

    content = content.replace(/(\w+)\s+"game\/pb\/_gen\/pb\/duck\/mongodb"/g, '$1 "duck/mongodb"')

    fs.writeFileSync(file, content, 'utf-8')
}


let ts_files = getfiles("./pb/_gen/", ".ts")

for (let i in ts_files) {
    let filename = ts_files[i]

    const file = readline.createInterface({ 
        input: fs.createReadStream(filename), 
        output: process.stdout, 
        terminal: false
    }); 
    
    let lines = []
    let lastInterface = ""
    let lastInterfaceIndex = -1

    let opts = {
        prefix: "",
        ignore: false,
        package: "",
    }

    file.on('line', (line) => { 
        if (line.startsWith("import")) {
            if (line.includes("ObjectID") || line.includes("TimeStamp")) {
                //
            } else {
                lines.push(line)
            }
        } else {
            if (line.startsWith("export const protobufPackage")) {
                let matchComment = line.match(/"(\w+)"/)
                if (matchComment) {
                    opts.package = matchComment[1]
                }
                return
            }
            if (line.startsWith("}")) { 
                if (!opts.ignore) {
                    lines.push(line) 
                } else {
                    lines
                }
                opts.prefix = ""
                opts.ignore = false

                return
            }

            if (line.startsWith("/*")) {
                let matchComment = line.match(/@ts:\s*(\w+)\((\w*)\)/)
                if (matchComment) {
                    if (matchComment[1] == 'prefix') {
                        opts.prefix = matchComment[2]
                    } else if (matchComment[1] == 'ignore') {
                        opts.ignore = true
                    }
                }

                if (!opts.ignore) {
                    lines.push(line)
                }

                return
            }

            if (opts.ignore) {
                return
            }

            line = line.replace(/ObjectID/g, "string")
            line = line.replace(/TimeStamp/g, "string")

            // 这部分是把ts_proto生成的service 替换为我们的样式。
            let matchInterface = line.match(/export interface (\w+)/)
            if (matchInterface) {
                lastInterface = matchInterface[1]
                lastInterfaceIndex = lines.length
            }

            let matchFunc = line.match(/(\w+)\(\w+:\s*(\w+)\):\sPromise<(\w+)>;*/)
            if (matchFunc) {
                if (lastInterfaceIndex != -1) {
                    lines[lastInterfaceIndex] = `export class ${lastInterface} {`
                }

                let route = ""
                if (opts.prefix) {
                    route += opts.prefix + "/"
                }
                if (opts.package) {
                    route += opts.package + "."
                }
                route += lastInterface + "/" + matchFunc[1]

                lines.push(`  static async ${matchFunc[1]}(client, req : ${matchFunc[2]}) : Promise<[${matchFunc[3]},any]> {`)
                lines.push(`    return await client.send("${route}", req)`)
                lines.push("  }")
            } else {
                lines.push(line)
            }
        }
    });

    file.on('close', ()=>{ 
        fs.writeFileSync(filename, lines.join("\n"), 'utf-8')
    });
}



let pb_files = getfiles("./pb/", ".proto")

let a = "a".charCodeAt(0)
let z = "z".charCodeAt(0)
function isLower(str) {
    let code = str.charCodeAt(0)
    return code >= a && code <= z
}
function check(match, filename, lineidx) {
    if (match) {
        if (isLower(match[1])) {
            console.log({
                msg: "为避免不必要的麻烦，请勿使用小写开头",
                file: filename+":"+lineidx,
                line: match[0],
            })
        }
        return true
    }
    return false
}

let matches = [
    /message\s+(\w+)/,
    /repeated\s+\w+\s+(\w+)\s*=\s*\w+;/,
    /\w+\s+(\w+)\s*=\s*\w+;/,
    /rpc\s+(\w+)/,
]

for (let i in pb_files) {
    let filename = pb_files[i]

    const file = readline.createInterface({ 
        input: fs.createReadStream(filename), 
        output: process.stdout, 
        terminal: false
    }); 
    
    let lineidx = 0
    file.on('line', (line) => { 
        lineidx++
        for (let m of matches) {
            if (check(line.match(m), filename, lineidx)) {
                break
            }
        }
        if (line.match(/^\s+\/\/[^\s@]/)) {
            console.log("请在双斜杆后加一个空格", filename + ":"+lineidx, "\t\t", line)
        }
    })
}

