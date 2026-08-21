
import * as XLSX from 'xlsx'

interface ExcelOption {
    key: string
    name: string
    fmt?: (v: any, row) => any
}

interface ExcelSheetOptions {
    data: any[]
    options?: ExcelOption[]
    sheetName: string
}

/*
    excel.dump(data.List, `金币记录-${this.currentpage}`, [
        { key: "HandleTypeName", name: "操作类型" },
        { key: "Expend", name: "支出金币", fmt: ut.fmtGold },
        { key: "InCome", name: "收入金币", fmt: ut.fmtGold },
        { key: "Change", name: "净增减", fmt: ut.fmtGold },
        { key: "CurrentGold", name: "当前金币", fmt: ut.fmtGold },
        { key: "Time", name: "时间" },
        { key: "Comment", name: "日志" },
    ]);
*/

export class excel {
    static dump(data: any[], filename, options: ExcelOption[]) {
        let workbook = XLSX.utils.book_new()

        this.add_sheet(workbook, { data: data, options: options, sheetName: "Sheet1" })

        var wopts = { bookType: 'xlsx', bookSST: false, type: 'array' };
        XLSX.writeFile(workbook, filename + ".xlsx", <XLSX.WritingOptions>wopts);
    }


    static dump_multi_sheet(filename: string, dataArr: ExcelSheetOptions[]) {
        let workbook = XLSX.utils.book_new()

        dataArr.forEach(v => this.add_sheet(workbook, v))

        var wopts = { bookType: 'xlsx', bookSST: false, type: 'array' };
        XLSX.writeFile(workbook, filename + ".xlsx", <XLSX.WritingOptions>wopts);
    }
    static dump_multi_sheet_new(filename: string, dataArr: ExcelSheetOptions[]) {
        let workbook = XLSX.utils.book_new()

        dataArr.forEach(v => this.add_sheet_multi(workbook, v))

        var wopts = { bookType: 'xlsx', bookSST: false, type: 'array' };
        XLSX.writeFile(workbook, filename + ".xlsx", <XLSX.WritingOptions>wopts);
    }
    static add_sheet_multi(workbook, sheetOptions) {
        const { data, options, sheetName, topData } = sheetOptions;
        const sheetData = [];

        // 添加表头信息
        sheetData.push(['商户', topData.operator, '', '', '', '账期',topData.billDate]);
        sheetData.push(['账单创建日期', topData.billCreationDate,'', '', '', '结算货币',topData.settlementCurrency]);
        sheetData.push([]);

        // 添加列名
        sheetData.push(options.map(opt => opt.name));

        // 添加详细数据
        data.forEach(row => {
            sheetData.push(options.map(opt => row[opt.key] || ""));
        });

        // 创建工作表并添加到工作簿中
        const worksheet = XLSX.utils.aoa_to_sheet(sheetData);
        XLSX.utils.book_append_sheet(workbook, worksheet, sheetName);
    }
    static dump_multi_sheet_add_top(filename: string, dataArr: ExcelSheetOptions[],topDataList) {
        let workbook = XLSX.utils.book_new()

        dataArr.forEach(v => this.add_sheet(workbook, v))
        topDataList.forEach(v => this.add_sheet(workbook, v))

        var wopts = { bookType: 'xlsx', bookSST: false, type: 'array' };
        console.log(wopts,'wopts');
        console.log(workbook.Sheets,'wopts');

        XLSX.writeFile(workbook, filename + ".xlsx", <XLSX.WritingOptions>wopts);
    }
    private static add_sheet(workbook, sheetInfo: ExcelSheetOptions) {
        let keys: string[] = []
        let k2name = {}

        if (!sheetInfo.options) {
            sheetInfo.options = []
            let one = sheetInfo.data[0]
            for (let k in one) {
                sheetInfo.options.push({
                    key: k,
                    name: k,
                })
            }
        }

        sheetInfo.options.forEach(e => {
            keys.push(e.key)
            k2name[e.key] = e.name
        })

        let newdata: any[] = []
        for (let i = 0; i < sheetInfo.data.length; i++) {
            let one = sheetInfo.data[i]


            let newone = {}

            for (let j = 0; j < keys.length; j++) {
                let k = keys[j]
                let fmt = sheetInfo.options[j].fmt
                if (fmt) {
                    newone[k] = fmt(one[k], one)
                } else {
                    newone[k] = one[k]
                }
            }
            newdata.push(newone)
        }

        let sheet = XLSX.utils.json_to_sheet(newdata, { header: keys, skipHeader: false })
        console.log(newdata,'newdata',sheet,'sheet');
        XLSX.utils.book_append_sheet(workbook, sheet, sheetInfo.sheetName);

        const range = XLSX.utils.decode_range(<string>sheet['!ref'])

        for (let c = range.s.c; c <= range.e.c; c++) {
            const header = XLSX.utils.encode_col(c) + '1'
            sheet[header].v = k2name[sheet[header].v]
        }
    }




    static excelResolveData(tableHeaderMap, data) {

        let generatorData = []
        for (const i in data) {

            let dataItem = data[i]
            let generatorItem = {}

            for (const key in dataItem) {
                generatorItem[tableHeaderMap[key]] = dataItem[key]
            }
            generatorData.push(generatorItem)

        }

        return generatorData
    }
    static DataGeneratorExcel(tableHeader, data, fileName) {

        let tableHeaderMap = []
        for(let i in tableHeader){

            tableHeaderMap.push({
                key:tableHeader[i].value,
                name: tableHeader[i].label
            })

        }





        let generatorData = []
        for (const i in data) {

            let dataItem = data[i]
            let generatorItem = {}

            for (const key in dataItem) {

                generatorItem[key] = dataItem[key].toString()


            }
            generatorData.push(generatorItem)

        }

        this.dump(generatorData, fileName, tableHeaderMap)
    }
}
