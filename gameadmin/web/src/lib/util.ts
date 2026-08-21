import moment from 'moment';
import {App} from 'vue';
import {sprintf} from 'sprintf-js'
import {tip} from './tip';
import copy from 'copy-to-clipboard';
import crypto from "crypto-js"
import {JSEncrypt} from "jsencrypt";
import {useStore} from "@/pinia";

export function initGolbal(app: App) {
	let d = app.config.globalProperties

	// 用于element-plus table组件数据格式化
	// function(row, column, cellValue, index)
	d.goldSignedFormater = (row, col, c) => ut.fmtGoldSigned(c)
	d.dateFormater = (row, col, c) => ut.fmtDate(c)
	d.date1Formater = (row, col, c) => ut.fmtDate1(c)
	d.dateSecondFormater = (row, col, c) => ut.fmtDateSecond(c)
	d.goldFormater = (row, col, c) => ut.fmtGold(c)
	d.boolFormatter = (row, col, c) => ut.fmtBool(c)
	d.percentFormatter = (row, col, c) => ut.fmtPercent(c)

	d.goldTenThousand = (val) => { return val * 1e4 }

	d.imageUrl = (url) => ut.imageUrl(url)
	d.copyText = (text) => ut.copyTextToClipBoard(text)
	d.openIMG = (url) => ut.openIMG(url)
	d.percentage = (val, val2) => {
		if (val2 === 0) {
			return '0.00%'
		}
		return (val / val2 * 100).toFixed(2) + "%"
	}
}

let currencyList = [
	{
		"code": "CN",    //国家编码
		"en": "China",    //国家英文名称
		"cn": "中国",      //国家中文名称
		"currency_code": "CNY",    //货币编码
		"currency_cn": "人民币元",   //货币中文名称
		"currency_en": "Chinese Yuan",    //货币英文名称
		"symbol": "CN¥",    //货币符号
		"symbol_native": "CN¥",    //货币原生符号
		'lang':'zh',
		'chineseLang':'中文',
	},
	{
		"code": "US",
		"en": "United States of America (USA)",
		"cn": "美国",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "EU",
		"en": "European Union",
		"cn": "欧盟",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "GB",
		"en": "Great Britain (United Kingdom; England)",
		"cn": "英国",
		"currency_code": "GBP",
		"currency_cn": "英镑",
		"currency_en": "British Pound Sterling",
		"symbol": "£",
		"symbol_native": "£"
	},
	{
		"code": "JP",
		"en": "Japan",
		"cn": "日本",
		"currency_code": "JPY",
		"currency_cn": "日元",
		"currency_en": "Japanese Yen",
		"symbol": "¥",
		"symbol_native": "￥"
	},
	{
		"code": "DE",
		"en": "Germany",
		"cn": "德国",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "FR",
		"en": "France",
		"cn": "法国",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "IT",
		"en": "Italy",
		"cn": "意大利",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "ES",
		"en": "Spain",
		"cn": "西班牙",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "AU",
		"en": "Australia",
		"cn": "澳大利亚",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "CA",
		"en": "Canada",
		"cn": "加拿大",
		"currency_code": "CAD",
		"currency_cn": "加元",
		"currency_en": "Canadian Dollar",
		"symbol": "CA$",
		"symbol_native": "$"
	},
	{
		"code": "HK",
		"en": "Hong Kong",
		"cn": "香港",
		"currency_code": "HKD",
		"currency_cn": "港元",
		"currency_en": "Hong Kong Dollar",
		"symbol": "HK$",
		"symbol_native": "$"
	},
	{
		"code": "TW",
		"en": "Taiwan",
		"cn": "台湾",
		"currency_code": "TWD",
		"currency_en": "New Taiwan Dollar",
		"symbol": "NT$",
		"symbol_native": "NT$"
	},
	{
		"code": "MO",
		"en": "Macao",
		"cn": "澳门",
		"currency_code": "MOP",
		"currency_cn": "澳门元",
		"currency_en": "Macanese Pataca",
		"symbol": "MOP$",
		"symbol_native": "MOP$"
	},
	{
		"code": "AR",
		"en": "Argentina",
		"cn": "阿根廷",
		"currency_code": "ARS",
		"currency_en": "Argentine Peso",
		"symbol": "AR$",
		"symbol_native": "$"
	},
	{
		"code": "AD",
		"en": "Andorra",
		"cn": "安道尔",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "AE",
		"en": "United Arab Emirates",
		"cn": "阿联酋",
		"currency_code": "AED",
		"currency_en": "United Arab Emirates Dirham",
		"symbol": "AED",
		"symbol_native": "د.إ."
	},
	{
		"code": "AF",
		"en": "Afghanistan",
		"cn": "阿富汗",
		"currency_code": "AFN",
		"currency_en": "Afghan Afghani",
		"symbol": "Af",
		"symbol_native": "؋"
	},
	{
		"code": "AG",
		"en": "Antigua & Barbuda",
		"cn": "安提瓜和巴布达",
		"currency_code": "XCD"
	},
	{
		"code": "AI",
		"en": "Anguilla",
		"cn": "安圭拉",
		"currency_code": "XCD"
	},
	{
		"code": "AL",
		"en": "Albania",
		"cn": "阿尔巴尼亚",
		"currency_code": "ALL",
		"currency_cn": "列克",
		"currency_en": "Albanian Lek",
		"symbol": "ALL",
		"symbol_native": "Lek"
	},
	{
		"code": "AM",
		"en": "Armenia",
		"cn": "亚美尼亚",
		"currency_code": "AMD",
		"currency_en": "Armenian Dram",
		"symbol": "AMD",
		"symbol_native": "դր."
	},
	{
		"code": "AO",
		"en": "Angola",
		"cn": "安哥拉",
		"currency_code": "AOA"
	},
	{
		"code": "AQ",
		"en": "Antarctica",
		"cn": "南极洲",
		"currency_code": ""
	},
	{
		"code": "AS",
		"en": "American Samoa",
		"cn": "美属萨摩亚",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "AT",
		"en": "Austria",
		"cn": "奥地利",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "AW",
		"en": "Aruba",
		"cn": "阿鲁巴",
		"currency_code": "AWG"
	},
	{
		"code": "AX",
		"en": "Aland Island",
		"cn": "奥兰群岛",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "AZ",
		"en": "Azerbaijan",
		"cn": "阿塞拜疆",
		"currency_code": "AZN",
		"currency_en": "Azerbaijani Manat",
		"symbol": "man.",
		"symbol_native": "ман."
	},
	{
		"code": "BA",
		"en": "Bosnia & Herzegovina",
		"cn": "波黑",
		"currency_code": "BAM",
		"currency_en": "Bosnia-Herzegovina Convertible Mark",
		"symbol": "KM",
		"symbol_native": "KM"
	},
	{
		"code": "BB",
		"en": "Barbados",
		"cn": "巴巴多斯",
		"currency_code": "BBD",
		"currency_cn": "巴巴多斯元"
	},
	{
		"code": "BD",
		"en": "Bangladesh",
		"cn": "孟加拉",
		"currency_code": "BDT",
		"currency_en": "Bangladeshi Taka",
		"symbol": "Tk",
		"symbol_native": "৳"
	},
	{
		"code": "BE",
		"en": "Belgium",
		"cn": "比利时",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "BF",
		"en": "Burkina",
		"cn": "布基纳法索",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "BG",
		"en": "Bulgaria",
		"cn": "保加利亚",
		"currency_code": "BGN",
		"currency_en": "Bulgarian Lev",
		"symbol": "BGN",
		"symbol_native": "лв."
	},
	{
		"code": "BH",
		"en": "Bahrain",
		"cn": "巴林",
		"currency_code": "BHD",
		"currency_cn": "巴林第纳尔",
		"currency_en": "Bahraini Dinar",
		"symbol": "BD",
		"symbol_native": "د.ب."
	},
	{
		"code": "BI",
		"en": "Burundi",
		"cn": "布隆迪",
		"currency_code": "BIF",
		"currency_cn": "布隆迪法郎",
		"currency_en": "Burundian Franc",
		"symbol": "FBu",
		"symbol_native": "FBu"
	},
	{
		"code": "BJ",
		"en": "Benin",
		"cn": "贝宁",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "BL",
		"en": "Saint Barthélemy",
		"cn": "圣巴泰勒米岛",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "BM",
		"en": "Bermuda",
		"cn": "百慕大",
		"currency_code": "BMD"
	},
	{
		"code": "BN",
		"en": "Brunei",
		"cn": "文莱",
		"currency_code": "BND",
		"currency_en": "Brunei Dollar",
		"symbol": "BN$",
		"symbol_native": "$"
	},
	{
		"code": "BO",
		"en": "Bolivia",
		"cn": "玻利维亚",
		"currency_code": "BOB",
		"currency_en": "Bolivian Boliviano",
		"symbol": "Bs",
		"symbol_native": "Bs"
	},
	{
		"code": "BQ",
		"en": "Caribbean Netherlands",
		"cn": "荷兰加勒比区",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "BR",
		"en": "Brazil",
		"cn": "巴西",
		"currency_code": "BRL",
		"currency_en": "Brazilian Real",
		"symbol": "R$",
		"symbol_native": "R$"
	},
	{
		"code": "BS",
		"en": "The Bahamas",
		"cn": "巴哈马",
		"currency_code": "BSD",
		"currency_cn": "巴哈马元"
	},
	{
		"code": "BT",
		"en": "Bhutan",
		"cn": "不丹",
		"currency_code": "BTN"
	},
	{
		"code": "BV",
		"en": "Bouvet Island",
		"cn": "布韦岛",
		"currency_code": "NOK",
		"currency_cn": "挪威克朗",
		"currency_en": "Norwegian Krone",
		"symbol": "Nkr",
		"symbol_native": "kr"
	},
	{
		"code": "BW",
		"en": "Botswana",
		"cn": "博茨瓦纳",
		"currency_code": "BWP",
		"currency_en": "Botswanan Pula",
		"symbol": "BWP",
		"symbol_native": "P"
	},
	{
		"code": "BY",
		"en": "Belarus",
		"cn": "白俄罗斯",
		"currency_code": "BYR"
	},
	{
		"code": "BZ",
		"en": "Belize",
		"cn": "伯利兹",
		"currency_code": "BZD",
		"currency_en": "Belize Dollar",
		"symbol": "BZ$",
		"symbol_native": "$"
	},
	{
		"code": "CC",
		"en": "Cocos (Keeling) Islands",
		"cn": "科科斯群岛",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "CD",
		"en": "Democratic Republic of the Congo",
		"cn": "刚果（金）",
		"currency_code": "CDF",
		"currency_en": "Congolese Franc",
		"symbol": "CDF",
		"symbol_native": "FrCD"
	},
	{
		"code": "CF",
		"en": "Central African Republic",
		"cn": "中非",
		"currency_code": "XAF",
		"currency_cn": "中非金融合作法郎",
		"currency_en": "CFA Franc BEAC",
		"symbol": "FCFA",
		"symbol_native": "FCFA"
	},
	{
		"code": "CG",
		"en": "Republic of the Congo",
		"cn": "刚果（布）",
		"currency_code": "XAF",
		"currency_cn": "中非金融合作法郎",
		"currency_en": "CFA Franc BEAC",
		"symbol": "FCFA",
		"symbol_native": "FCFA"
	},
	{
		"code": "CH",
		"en": "Switzerland",
		"cn": "瑞士",
		"currency_code": "CHF",
		"currency_cn": "瑞士法郎",
		"currency_en": "Swiss Franc",
		"symbol": "CHF",
		"symbol_native": "CHF"
	},
	{
		"code": "CI",
		"en": "Cote d'Ivoire",
		"cn": "科特迪瓦",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "CK",
		"en": "Cook Islands",
		"cn": "库克群岛",
		"currency_code": "NZD",
		"currency_cn": "新西兰元",
		"currency_en": "New Zealand Dollar",
		"symbol": "NZ$",
		"symbol_native": "$"
	},
	{
		"code": "CL",
		"en": "Chile",
		"cn": "智利",
		"currency_code": "CLP",
		"currency_cn": "智利比索",
		"currency_en": "Chilean Peso",
		"symbol": "CL$",
		"symbol_native": "$"
	},
	{
		"code": "CM",
		"en": "Cameroon",
		"cn": "喀麦隆",
		"currency_code": "XAF",
		"currency_cn": "中非金融合作法郎",
		"currency_en": "CFA Franc BEAC",
		"symbol": "FCFA",
		"symbol_native": "FCFA"
	},
	{
		"code": "CO",
		"en": "Colombia",
		"cn": "哥伦比亚",
		"currency_code": "COP",
		"currency_cn": "哥伦比亚比索",
		"currency_en": "Colombian Peso",
		"symbol": "CO$",
		"symbol_native": "$"
	},
	{
		"code": "CR",
		"en": "Costa Rica",
		"cn": "哥斯达黎加",
		"currency_code": "CRC",
		"currency_cn": "哥斯达黎加科朗",
		"currency_en": "Costa Rican Colón",
		"symbol": "₡",
		"symbol_native": "₡"
	},
	{
		"code": "CU",
		"en": "Cuba",
		"cn": "古巴",
		"currency_code": "CUP",
		"currency_cn": "古巴比索"
	},
	{
		"code": "CV",
		"en": "Cape Verde",
		"cn": "佛得角",
		"currency_code": "CVE",
		"currency_en": "Cape Verdean Escudo",
		"symbol": "CV$",
		"symbol_native": "CV$"
	},
	{
		"code": "CW",
		"en": "Curacao",
		"cn": "库拉索",
		"currency_code": "ANG"
	},
	{
		"code": "CX",
		"en": "Christmas Island",
		"cn": "圣诞岛",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "CY",
		"en": "Cyprus",
		"cn": "塞浦路斯",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "CZ",
		"en": "Czech Republic",
		"cn": "捷克",
		"currency_code": "CZK",
		"currency_en": "Czech Republic Koruna",
		"symbol": "Kč",
		"symbol_native": "Kč"
	},
	{
		"code": "DJ",
		"en": "Djibouti",
		"cn": "吉布提",
		"currency_code": "DJF",
		"currency_cn": "吉布提法郎",
		"currency_en": "Djiboutian Franc",
		"symbol": "Fdj",
		"symbol_native": "Fdj"
	},
	{
		"code": "DK",
		"en": "Denmark",
		"cn": "丹麦",
		"currency_code": "DKK",
		"currency_cn": "丹麦克朗",
		"currency_en": "Danish Krone",
		"symbol": "Dkr",
		"symbol_native": "kr",
		'lang':'da',
		'chineseLang':'丹麦文',
	},
	{
		"code": "DM",
		"en": "Dominica",
		"cn": "多米尼克",
		"currency_code": "XCD"
	},
	{
		"code": "DO",
		"en": "Dominican Republic",
		"cn": "多米尼加",
		"currency_code": "DOP",
		"currency_cn": "多米尼加比索",
		"currency_en": "Dominican Peso",
		"symbol": "RD$",
		"symbol_native": "RD$"
	},
	{
		"code": "DZ",
		"en": "Algeria",
		"cn": "阿尔及利亚",
		"currency_code": "DZD",
		"currency_cn": "阿尔及利亚第纳尔",
		"currency_en": "Algerian Dinar",
		"symbol": "DA",
		"symbol_native": "د.ج."
	},
	{
		"code": "EC",
		"en": "Ecuador",
		"cn": "厄瓜多尔",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "EE",
		"en": "Estonia",
		"cn": "爱沙尼亚",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "EG",
		"en": "Egypt",
		"cn": "埃及",
		"currency_code": "EGP",
		"currency_cn": "埃及镑",
		"currency_en": "Egyptian Pound",
		"symbol": "EGP",
		"symbol_native": "ج.م."
	},
	{
		"code": "EH",
		"en": "Western Sahara",
		"cn": "西撒哈拉",
		"currency_code": "MAD",
		"currency_cn": "摩洛哥迪拉姆",
		"currency_en": "Moroccan Dirham",
		"symbol": "MAD",
		"symbol_native": "د.م."
	},
	{
		"code": "ER",
		"en": "Eritrea",
		"cn": "厄立特里亚",
		"currency_code": "ERN",
		"currency_en": "Eritrean Nakfa",
		"symbol": "Nfk",
		"symbol_native": "Nfk"
	},
	{
		"code": "ET",
		"en": "Ethiopia",
		"cn": "埃塞俄比亚",
		"currency_code": "ETB",
		"currency_en": "Ethiopian Birr",
		"symbol": "Br",
		"symbol_native": "Br"
	},
	{
		"code": "FI",
		"en": "Finland",
		"cn": "芬兰",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "FJ",
		"en": "Fiji",
		"cn": "斐济群岛",
		"currency_code": "FJD",
		"currency_cn": "斐济元"
	},
	{
		"code": "FK",
		"en": "Falkland Islands",
		"cn": "马尔维纳斯群岛（福克兰）",
		"currency_code": "FKP"
	},
	{
		"code": "FM",
		"en": "Federated States of Micronesia",
		"cn": "密克罗尼西亚联邦",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "FO",
		"en": "Faroe Islands",
		"cn": "法罗群岛",
		"currency_code": "DKK",
		"currency_cn": "丹麦克朗",
		"currency_en": "Danish Krone",
		"symbol": "Dkr",
		"symbol_native": "kr"
	},
	{
		"code": "GA",
		"en": "Gabon",
		"cn": "加蓬",
		"currency_code": "XAF",
		"currency_cn": "中非金融合作法郎",
		"currency_en": "CFA Franc BEAC",
		"symbol": "FCFA",
		"symbol_native": "FCFA"
	},
	{
		"code": "GD",
		"en": "Grenada",
		"cn": "格林纳达",
		"currency_code": "XCD"
	},
	{
		"code": "GE",
		"en": "Georgia",
		"cn": "格鲁吉亚",
		"currency_code": "GEL",
		"currency_en": "Georgian Lari",
		"symbol": "GEL",
		"symbol_native": "GEL"
	},
	{
		"code": "GF",
		"en": "French Guiana",
		"cn": "法属圭亚那",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "GG",
		"en": "Guernsey",
		"cn": "根西岛",
		"currency_code": "GBP",
		"currency_cn": "英镑",
		"currency_en": "British Pound Sterling",
		"symbol": "£",
		"symbol_native": "£"
	},
	{
		"code": "GH",
		"en": "Ghana",
		"cn": "加纳",
		"currency_code": "GHS",
		"currency_en": "Ghanaian Cedi",
		"symbol": "GH₵",
		"symbol_native": "GH₵"
	},
	{
		"code": "GI",
		"en": "Gibraltar",
		"cn": "直布罗陀",
		"currency_code": "GIP"
	},
	{
		"code": "GL",
		"en": "Greenland",
		"cn": "格陵兰",
		"currency_code": "DKK",
		"currency_cn": "丹麦克朗",
		"currency_en": "Danish Krone",
		"symbol": "Dkr",
		"symbol_native": "kr"
	},
	{
		"code": "GM",
		"en": "Gambia",
		"cn": "冈比亚",
		"currency_code": "GMD",
		"currency_cn": "法拉西"
	},
	{
		"code": "GN",
		"en": "Guinea",
		"cn": "几内亚",
		"currency_code": "GNF",
		"currency_en": "Guinean Franc",
		"symbol": "FG",
		"symbol_native": "FG"
	},
	{
		"code": "GP",
		"en": "Guadeloupe",
		"cn": "瓜德罗普",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "GQ",
		"en": "Equatorial Guinea",
		"cn": "赤道几内亚",
		"currency_code": "XAF",
		"currency_cn": "中非金融合作法郎",
		"currency_en": "CFA Franc BEAC",
		"symbol": "FCFA",
		"symbol_native": "FCFA"
	},
	{
		"code": "GR",
		"en": "Greece",
		"cn": "希腊",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "GS",
		"en": "South Georgia and the South Sandwich Islands",
		"cn": "南乔治亚岛和南桑威奇群岛",
		"currency_code": "GBP",
		"currency_cn": "英镑",
		"currency_en": "British Pound Sterling",
		"symbol": "£",
		"symbol_native": "£"
	},
	{
		"code": "GT",
		"en": "Guatemala",
		"cn": "危地马拉",
		"currency_code": "GTQ",
		"currency_cn": "格查尔",
		"currency_en": "Guatemalan Quetzal",
		"symbol": "GTQ",
		"symbol_native": "Q"
	},
	{
		"code": "GU",
		"en": "Guam",
		"cn": "关岛",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "GW",
		"en": "Guinea-Bissau",
		"cn": "几内亚比绍",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "GY",
		"en": "Guyana",
		"cn": "圭亚那",
		"currency_code": "GYD",
		"currency_cn": "圭亚那元"
	},
	{
		"code": "HM",
		"en": "Heard Island and McDonald Islands",
		"cn": "赫德岛和麦克唐纳群岛",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "HN",
		"en": "Honduras",
		"cn": "洪都拉斯",
		"currency_code": "HNL",
		"currency_cn": "伦皮拉",
		"currency_en": "Honduran Lempira",
		"symbol": "HNL",
		"symbol_native": "L"
	},
	{
		"code": "HR",
		"en": "Croatia",
		"cn": "克罗地亚",
		"currency_code": "HRK",
		"currency_en": "Croatian Kuna",
		"symbol": "kn",
		"symbol_native": "kn"
	},
	{
		"code": "HT",
		"en": "Haiti",
		"cn": "海地",
		"currency_code": "HTG",
		"currency_cn": "古德"
	},
	{
		"code": "HU",
		"en": "Hungary",
		"cn": "匈牙利",
		"currency_code": "HUF",
		"currency_cn": "福林",
		"currency_en": "Hungarian Forint",
		"symbol": "Ft",
		"symbol_native": "Ft"
	},
	{
		"code": "ID",
		"en": "Indonesia",
		"cn": "印尼",
		"currency_code": "IDR",
		"currency_cn": "印尼盾",
		"currency_en": "Indonesian Rupiah",
		"symbol": "Rp",
		"symbol_native": "Rp"
	},
	{
		"code": "IE",
		"en": "Ireland",
		"cn": "爱尔兰",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "IL",
		"en": "Israel",
		"cn": "以色列",
		"currency_code": "ILS",
		"currency_en": "Israeli New Sheqel",
		"symbol": "₪",
		"symbol_native": "₪"
	},
	{
		"code": "IM",
		"en": "Isle of Man",
		"cn": "马恩岛",
		"currency_code": "GBP",
		"currency_cn": "英镑",
		"currency_en": "British Pound Sterling",
		"symbol": "£",
		"symbol_native": "£"
	},
	{
		"code": "IN",
		"en": "India",
		"cn": "印度",
		"currency_code": "INR",
		"currency_cn": "卢比",
		"currency_en": "Indian Rupee",
		"symbol": "Rs",
		"symbol_native": "টকা"
	},
	{
		"code": "IO",
		"en": "British Indian Ocean Territory",
		"cn": "英属印度洋领地",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "IQ",
		"en": "Iraq",
		"cn": "伊拉克",
		"currency_code": "IQD",
		"currency_cn": "伊拉克第纳尔",
		"currency_en": "Iraqi Dinar",
		"symbol": "IQD",
		"symbol_native": "د.ع."
	},
	{
		"code": "IR",
		"en": "Iran",
		"cn": "伊朗",
		"currency_code": "IRR",
		"currency_cn": "伊朗里亚尔",
		"currency_en": "Iranian Rial",
		"symbol": "IRR",
		"symbol_native": "﷼"
	},
	{
		"code": "IS",
		"en": "Iceland",
		"cn": "冰岛",
		"currency_code": "ISK",
		"currency_cn": "冰岛克朗",
		"currency_en": "Icelandic Króna",
		"symbol": "Ikr",
		"symbol_native": "kr"
	},
	{
		"code": "JE",
		"en": "Jersey",
		"cn": "泽西岛",
		"currency_code": "GBP",
		"currency_cn": "英镑",
		"currency_en": "British Pound Sterling",
		"symbol": "£",
		"symbol_native": "£"
	},
	{
		"code": "JM",
		"en": "Jamaica",
		"cn": "牙买加",
		"currency_code": "JMD",
		"currency_cn": "牙买加元",
		"currency_en": "Jamaican Dollar",
		"symbol": "J$",
		"symbol_native": "$"
	},
	{
		"code": "JO",
		"en": "Jordan",
		"cn": "约旦",
		"currency_code": "JOD",
		"currency_cn": "约旦第纳尔",
		"currency_en": "Jordanian Dinar",
		"symbol": "JD",
		"symbol_native": "د.أ."
	},
	{
		"code": "KE",
		"en": "Kenya",
		"cn": "肯尼亚",
		"currency_code": "KES",
		"currency_cn": "肯尼亚先令",
		"currency_en": "Kenyan Shilling",
		"symbol": "Ksh",
		"symbol_native": "Ksh"
	},
	{
		"code": "KG",
		"en": "Kyrgyzstan",
		"cn": "吉尔吉斯斯坦",
		"currency_code": "KGS"
	},
	{
		"code": "KH",
		"en": "Cambodia",
		"cn": "柬埔寨",
		"currency_code": "KHR",
		"currency_cn": "瑞尔",
		"currency_en": "Cambodian Riel",
		"symbol": "KHR",
		"symbol_native": "៛"
	},
	{
		"code": "KI",
		"en": "Kiribati",
		"cn": "基里巴斯",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "KM",
		"en": "The Comoros",
		"cn": "科摩罗",
		"currency_code": "KMF",
		"currency_cn": "科摩罗法郎",
		"currency_en": "Comorian Franc",
		"symbol": "CF",
		"symbol_native": "FC"
	},
	{
		"code": "KN",
		"en": "St. Kitts & Nevis",
		"cn": "圣基茨和尼维斯",
		"currency_code": "XCD"
	},
	{
		"code": "KP",
		"en": "North Korea",
		"cn": "朝鲜",
		"currency_code": "KPW",
		"currency_cn": "朝鲜元"
	},
	{
		"code": "KR",
		"en": "South Korea",
		"cn": "韩国",
		"currency_code": "KRW",
		"currency_cn": "韩元",
		"currency_en": "South Korean Won",
		"symbol": "₩",
		"symbol_native": "₩"
	},
	{
		"code": "KW",
		"en": "Kuwait",
		"cn": "科威特",
		"currency_code": "KWD",
		"currency_cn": "科威特第纳尔",
		"currency_en": "Kuwaiti Dinar",
		"symbol": "KD",
		"symbol_native": "د.ك."
	},
	{
		"code": "KY",
		"en": "Cayman Islands",
		"cn": "开曼群岛",
		"currency_code": "KYD"
	},
	{
		"code": "KZ",
		"en": "Kazakhstan",
		"cn": "哈萨克斯坦",
		"currency_code": "KZT",
		"currency_en": "Kazakhstani Tenge",
		"symbol": "KZT",
		"symbol_native": "тңг."
	},
	{
		"code": "LA",
		"en": "Laos",
		"cn": "老挝",
		"currency_code": "LAK",
		"currency_cn": "基普"
	},
	{
		"code": "LB",
		"en": "Lebanon",
		"cn": "黎巴嫩",
		"currency_code": "LBP",
		"currency_cn": "黎巴嫩镑",
		"currency_en": "Lebanese Pound",
		"symbol": "L.L.",
		"symbol_native": "ل.ل."
	},
	{
		"code": "LC",
		"en": "St. Lucia",
		"cn": "圣卢西亚",
		"currency_code": "XCD"
	},
	{
		"code": "LI",
		"en": "Liechtenstein",
		"cn": "列支敦士登",
		"currency_code": "CHF",
		"currency_cn": "瑞士法郎",
		"currency_en": "Swiss Franc",
		"symbol": "CHF",
		"symbol_native": "CHF"
	},
	{
		"code": "LK",
		"en": "Sri Lanka",
		"cn": "斯里兰卡",
		"currency_code": "LKR",
		"currency_cn": "斯里兰卡卢比",
		"currency_en": "Sri Lankan Rupee",
		"symbol": "SLRs",
		"symbol_native": "SL Re"
	},
	{
		"code": "LR",
		"en": "Liberia",
		"cn": "利比里亚",
		"currency_code": "LRD",
		"currency_cn": "利比里亚元"
	},
	{
		"code": "LS",
		"en": "Lesotho",
		"cn": "莱索托",
		"currency_code": "LSL"
	},
	{
		"code": "LT",
		"en": "Lithuania",
		"cn": "立陶宛",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "LU",
		"en": "Luxembourg",
		"cn": "卢森堡",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "LV",
		"en": "Latvia",
		"cn": "拉脱维亚",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "LY",
		"en": "Libya",
		"cn": "利比亚",
		"currency_code": "LYD",
		"currency_cn": "利比亚第纳尔",
		"currency_en": "Libyan Dinar",
		"symbol": "LD",
		"symbol_native": "د.ل."
	},
	{
		"code": "MA",
		"en": "Morocco",
		"cn": "摩洛哥",
		"currency_code": "MAD",
		"currency_cn": "摩洛哥迪拉姆",
		"currency_en": "Moroccan Dirham",
		"symbol": "MAD",
		"symbol_native": "د.م."
	},
	{
		"code": "MC",
		"en": "Monaco",
		"cn": "摩纳哥",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "MD",
		"en": "Moldova",
		"cn": "摩尔多瓦",
		"currency_code": "MDL",
		"currency_en": "Moldovan Leu",
		"symbol": "MDL",
		"symbol_native": "MDL"
	},
	{
		"code": "ME",
		"en": "Montenegro",
		"cn": "黑山",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "MF",
		"en": "Saint Martin (France)",
		"cn": "法属圣马丁",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "MG",
		"en": "Madagascar",
		"cn": "马达加斯加",
		"currency_code": "MGA",
		"currency_en": "Malagasy Ariary",
		"symbol": "MGA",
		"symbol_native": "MGA"
	},
	{
		"code": "MH",
		"en": "Marshall islands",
		"cn": "马绍尔群岛",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "MK",
		"en": "Republic of Macedonia (FYROM)",
		"cn": "马其顿",
		"currency_code": "MKD",
		"currency_en": "Macedonian Denar",
		"symbol": "MKD",
		"symbol_native": "MKD"
	},
	{
		"code": "ML",
		"en": "Mali",
		"cn": "马里",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "MM",
		"en": "Myanmar (Burma)",
		"cn": "缅甸",
		"currency_code": "MMK",
		"currency_en": "Myanma Kyat",
		"symbol": "MMK",
		"symbol_native": "K"
	},
	{
		"code": "MN",
		"en": "Mongolia",
		"cn": "蒙古国",
		"currency_code": "MNT",
		"currency_cn": "图格里克"
	},
	{
		"code": "MP",
		"en": "Northern Mariana Islands",
		"cn": "北马里亚纳群岛",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "MQ",
		"en": "Martinique",
		"cn": "马提尼克",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "MR",
		"en": "Mauritania",
		"cn": "毛里塔尼亚",
		"currency_code": "MRO",
		"currency_cn": "乌吉亚"
	},
	{
		"code": "MS",
		"en": "Montserrat",
		"cn": "蒙塞拉特岛",
		"currency_code": "XCD"
	},
	{
		"code": "MT",
		"en": "Malta",
		"cn": "马耳他",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "MU",
		"en": "Mauritius",
		"cn": "毛里求斯",
		"currency_code": "MUR",
		"currency_cn": "毛里求斯卢比",
		"currency_en": "Mauritian Rupee",
		"symbol": "MURs",
		"symbol_native": "MURs"
	},
	{
		"code": "MV",
		"en": "Maldives",
		"cn": "马尔代夫",
		"currency_code": "MVR",
		"currency_cn": "马尔代夫卢比"
	},
	{
		"code": "MW",
		"en": "Malawi",
		"cn": "马拉维",
		"currency_code": "MWK"
	},
	{
		"code": "MX",
		"en": "Mexico",
		"cn": "墨西哥",
		"currency_code": "MXN",
		"currency_en": "Mexican Peso",
		"symbol": "MX$",
		"symbol_native": "$"
	},
	{
		"code": "MY",
		"en": "Malaysia",
		"cn": "马来西亚",
		"currency_code": "MYR",
		"currency_cn": "林吉特",
		"currency_en": "Malaysian Ringgit",
		"symbol": "RM",
		"symbol_native": "RM"
	},
	{
		"code": "MZ",
		"en": "Mozambique",
		"cn": "莫桑比克",
		"currency_code": "MZN",
		"currency_en": "Mozambican Metical",
		"symbol": "MTn",
		"symbol_native": "MTn"
	},
	{
		"code": "NA",
		"en": "Namibia",
		"cn": "纳米比亚",
		"currency_code": "NAD",
		"currency_en": "Namibian Dollar",
		"symbol": "N$",
		"symbol_native": "N$"
	},
	{
		"code": "NC",
		"en": "New Caledonia",
		"cn": "新喀里多尼亚",
		"currency_code": "XPF"
	},
	{
		"code": "NE",
		"en": "Niger",
		"cn": "尼日尔",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "NF",
		"en": "Norfolk Island",
		"cn": "诺福克岛",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "NG",
		"en": "Nigeria",
		"cn": "尼日利亚",
		"currency_code": "NGN",
		"currency_cn": "奈拉",
		"currency_en": "Nigerian Naira",
		"symbol": "₦",
		"symbol_native": "₦"
	},
	{
		"code": "NI",
		"en": "Nicaragua",
		"cn": "尼加拉瓜",
		"currency_code": "NIO",
		"currency_en": "Nicaraguan Córdoba",
		"symbol": "C$",
		"symbol_native": "C$"
	},
	{
		"code": "NL",
		"en": "Netherlands",
		"cn": "荷兰",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "NO",
		"en": "Norway",
		"cn": "挪威",
		"currency_code": "NOK",
		"currency_cn": "挪威克朗",
		"currency_en": "Norwegian Krone",
		"symbol": "Nkr",
		"symbol_native": "kr"
	},
	{
		"code": "NP",
		"en": "Nepal",
		"cn": "尼泊尔",
		"currency_code": "NPR",
		"currency_cn": "尼泊尔卢比",
		"currency_en": "Nepalese Rupee",
		"symbol": "NPRs",
		"symbol_native": "नेरू"
	},
	{
		"code": "NR",
		"en": "Nauru",
		"cn": "瑙鲁",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "NU",
		"en": "Niue",
		"cn": "纽埃",
		"currency_code": "NZD",
		"currency_cn": "新西兰元",
		"currency_en": "New Zealand Dollar",
		"symbol": "NZ$",
		"symbol_native": "$"
	},
	{
		"code": "NZ",
		"en": "New Zealand",
		"cn": "新西兰",
		"currency_code": "NZD",
		"currency_cn": "新西兰元",
		"currency_en": "New Zealand Dollar",
		"symbol": "NZ$",
		"symbol_native": "$"
	},
	{
		"code": "OM",
		"en": "Oman",
		"cn": "阿曼",
		"currency_code": "OMR",
		"currency_cn": "阿曼里亚尔",
		"currency_en": "Omani Rial",
		"symbol": "OMR",
		"symbol_native": "ر.ع."
	},
	{
		"code": "PA",
		"en": "Panama",
		"cn": "巴拿马",
		"currency_code": "PAB",
		"currency_cn": "巴拿马巴波亚",
		"currency_en": "Panamanian Balboa",
		"symbol": "B/.",
		"symbol_native": "B/."
	},
	{
		"code": "PE",
		"en": "Peru",
		"cn": "秘鲁",
		"currency_code": "PEN",
		"currency_en": "Peruvian Nuevo Sol",
		"symbol": "S/.",
		"symbol_native": "S/."
	},
	{
		"code": "PF",
		"en": "French polynesia",
		"cn": "法属波利尼西亚",
		"currency_code": "XPF"
	},
	{
		"code": "PG",
		"en": "Papua New Guinea",
		"cn": "巴布亚新几内亚",
		"currency_code": "PGK"
	},
	{
		"code": "PH",
		"en": "The Philippines",
		"cn": "菲律宾",
		"currency_code": "PHP",
		"currency_cn": "菲律宾比索",
		"currency_en": "Philippine Peso",
		"symbol": "₱",
		"symbol_native": "₱"
	},
	{
		"code": "PK",
		"en": "Pakistan",
		"cn": "巴基斯坦",
		"currency_code": "PKR",
		"currency_en": "Pakistani Rupee",
		"symbol": "PKRs",
		"symbol_native": "₨"
	},
	{
		"code": "PL",
		"en": "Poland",
		"cn": "波兰",
		"currency_code": "PLN",
		"currency_cn": "兹罗提",
		"currency_en": "Polish Zloty",
		"symbol": "zł",
		"symbol_native": "zł"
	},
	{
		"code": "PM",
		"en": "Saint-Pierre and Miquelon",
		"cn": "圣皮埃尔和密克隆",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "PN",
		"en": "Pitcairn Islands",
		"cn": "皮特凯恩群岛",
		"currency_code": "NZD",
		"currency_cn": "新西兰元",
		"currency_en": "New Zealand Dollar",
		"symbol": "NZ$",
		"symbol_native": "$"
	},
	{
		"code": "PR",
		"en": "Puerto Rico",
		"cn": "波多黎各",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "PS",
		"en": "Palestinian territories",
		"cn": "巴勒斯坦",
		"currency_code": "ILS",
		"currency_en": "Israeli New Sheqel",
		"symbol": "₪",
		"symbol_native": "₪"
	},
	{
		"code": "PT",
		"en": "Portugal",
		"cn": "葡萄牙",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "PW",
		"en": "Palau",
		"cn": "帕劳",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "PY",
		"en": "Paraguay",
		"cn": "巴拉圭",
		"currency_code": "PYG",
		"currency_cn": "巴拉圭瓜拉尼",
		"currency_en": "Paraguayan Guarani",
		"symbol": "₲",
		"symbol_native": "₲"
	},
	{
		"code": "QA",
		"en": "Qatar",
		"cn": "卡塔尔",
		"currency_code": "QAR",
		"currency_cn": "卡塔尔里亚尔",
		"currency_en": "Qatari Rial",
		"symbol": "QR",
		"symbol_native": "ر.ق."
	},
	{
		"code": "RE",
		"en": "Réunion",
		"cn": "留尼汪",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "RO",
		"en": "Romania",
		"cn": "罗马尼亚",
		"currency_code": "RON",
		"currency_en": "Romanian Leu",
		"symbol": "RON",
		"symbol_native": "RON"
	},
	{
		"code": "RS",
		"en": "Serbia",
		"cn": "塞尔维亚",
		"currency_code": "RSD",
		"currency_en": "Serbian Dinar",
		"symbol": "din.",
		"symbol_native": "дин."
	},
	{
		"code": "RU",
		"en": "Russian Federation",
		"cn": "俄罗斯",
		"currency_code": "RUB",
		"currency_cn": "卢布",
		"currency_en": "Russian Ruble",
		"symbol": "RUB",
		"symbol_native": "₽."
	},
	{
		"code": "RW",
		"en": "Rwanda",
		"cn": "卢旺达",
		"currency_code": "RWF",
		"currency_cn": "卢旺达法郎",
		"currency_en": "Rwandan Franc",
		"symbol": "RWF",
		"symbol_native": "FR"
	},
	{
		"code": "SA",
		"en": "Saudi Arabia",
		"cn": "沙特阿拉伯",
		"currency_code": "SAR",
		"currency_cn": "亚尔",
		"currency_en": "Saudi Riyal",
		"symbol": "SR",
		"symbol_native": "ر.س."
	},
	{
		"code": "SB",
		"en": "Solomon Islands",
		"cn": "所罗门群岛",
		"currency_code": "SBD",
		"currency_cn": "所罗门元"
	},
	{
		"code": "SC",
		"en": "Seychelles",
		"cn": "塞舌尔",
		"currency_code": "SCR",
		"currency_cn": "塞舌尔卢比"
	},
	{
		"code": "SD",
		"en": "Sudan",
		"cn": "苏丹",
		"currency_code": "SDG",
		"currency_en": "Sudanese Pound",
		"symbol": "SDG",
		"symbol_native": "SDG"
	},
	{
		"code": "SE",
		"en": "Sweden",
		"cn": "瑞典",
		"currency_code": "SEK",
		"currency_cn": "瑞典克朗",
		"currency_en": "Swedish Krona",
		"symbol": "Skr",
		"symbol_native": "kr"
	},
	{
		"code": "SG",
		"en": "Singapore",
		"cn": "新加坡",
		"currency_code": "SGD",
		"currency_cn": "新加坡元",
		"currency_en": "Singapore Dollar",
		"symbol": "S$",
		"symbol_native": "$"
	},
	{
		"code": "SH",
		"en": "St. Helena & Dependencies",
		"cn": "圣赫勒拿",
		"currency_code": "SHP"
	},
	{
		"code": "SI",
		"en": "Slovenia",
		"cn": "斯洛文尼亚",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "SJ",
		"en": "Svalbard and Jan Mayen",
		"cn": "斯瓦尔巴群岛和扬马延岛",
		"currency_code": "NOK",
		"currency_cn": "挪威克朗",
		"currency_en": "Norwegian Krone",
		"symbol": "Nkr",
		"symbol_native": "kr"
	},
	{
		"code": "SK",
		"en": "Slovakia",
		"cn": "斯洛伐克",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "SL",
		"en": "Sierra Leone",
		"cn": "塞拉利昂",
		"currency_code": "SLL",
		"currency_cn": "利昂"
	},
	{
		"code": "SM",
		"en": "San Marino",
		"cn": "圣马力诺",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "SN",
		"en": "Senegal",
		"cn": "塞内加尔",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "SO",
		"en": "Somalia",
		"cn": "索马里",
		"currency_code": "SOS",
		"currency_cn": "索马里先令",
		"currency_en": "Somali Shilling",
		"symbol": "Ssh",
		"symbol_native": "Ssh"
	},
	{
		"code": "SR",
		"en": "Suriname",
		"cn": "苏里南",
		"currency_code": "SRD"
	},
	{
		"code": "SS",
		"en": "South Sudan",
		"cn": "南苏丹",
		"currency_code": "SSP"
	},
	{
		"code": "ST",
		"en": "Sao Tome & Principe",
		"cn": "圣多美和普林西比",
		"currency_code": "STD"
	},
	{
		"code": "SV",
		"en": "El Salvador",
		"cn": "萨尔瓦多",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "SX",
		"en": "Sint Maarten",
		"cn": "荷属圣马丁",
		"currency_code": "ANG"
	},
	{
		"code": "SY",
		"en": "Syria",
		"cn": "叙利亚",
		"currency_code": "SYP",
		"currency_cn": "叙利亚镑",
		"currency_en": "Syrian Pound",
		"symbol": "SY£",
		"symbol_native": "ل.س."
	},
	{
		"code": "SZ",
		"en": "Swaziland",
		"cn": "斯威士兰",
		"currency_code": "SZL"
	},
	{
		"code": "TC",
		"en": "Turks & Caicos Islands",
		"cn": "特克斯和凯科斯群岛",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "TD",
		"en": "Chad",
		"cn": "乍得",
		"currency_code": "XAF",
		"currency_cn": "中非金融合作法郎",
		"currency_en": "CFA Franc BEAC",
		"symbol": "FCFA",
		"symbol_native": "FCFA"
	},
	{
		"code": "TF",
		"en": "French Southern Territories",
		"cn": "法属南部领地",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "TG",
		"en": "Togo",
		"cn": "多哥",
		"currency_code": "XOF",
		"currency_cn": "非共体法郎",
		"currency_en": "CFA Franc BCEAO",
		"symbol": "CFA",
		"symbol_native": "CFA"
	},
	{
		"code": "TH",
		"en": "Thailand",
		"cn": "泰国",
		"currency_code": "THB",
		"currency_location": "th-TH",
		"currency_en": "Thai Baht",
		"symbol": "฿",
		"symbol_native": "฿"
	},
	{
		"code": "TJ",
		"en": "Tajikistan",
		"cn": "塔吉克斯坦",
		"currency_code": "TJS"
	},
	{
		"code": "TK",
		"en": "Tokelau",
		"cn": "托克劳",
		"currency_code": "NZD",
		"currency_cn": "新西兰元",
		"currency_en": "New Zealand Dollar",
		"symbol": "NZ$",
		"symbol_native": "$"
	},
	{
		"code": "TL",
		"en": "Timor-Leste (East Timor)",
		"cn": "东帝汶",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "TM",
		"en": "Turkmenistan",
		"cn": "土库曼斯坦",
		"currency_code": "TMT"
	},
	{
		"code": "TN",
		"en": "Tunisia",
		"cn": "突尼斯",
		"currency_code": "TND",
		"currency_cn": "突尼斯第纳尔",
		"currency_en": "Tunisian Dinar",
		"symbol": "DT",
		"symbol_native": "د.ت."
	},
	{
		"code": "TO",
		"en": "Tonga",
		"cn": "汤加",
		"currency_code": "TOP",
		"currency_en": "Tongan Paʻanga",
		"symbol": "T$",
		"symbol_native": "T$"
	},
	{
		"code": "TR",
		"en": "Turkey",
		"cn": "土耳其",
		"currency_code": "TRY",
		"currency_en": "Turkish Lira",
		"symbol": "TL",
		"symbol_native": "TL"
	},
	{
		"code": "TT",
		"en": "Trinidad & Tobago",
		"cn": "特立尼达和多巴哥",
		"currency_code": "TTD",
		"currency_cn": "特立尼达多巴哥元",
		"currency_en": "Trinidad and Tobago Dollar",
		"symbol": "TT$",
		"symbol_native": "$"
	},
	{
		"code": "TV",
		"en": "Tuvalu",
		"cn": "图瓦卢",
		"currency_code": "AUD",
		"currency_cn": "澳大利亚元",
		"currency_en": "Australian Dollar",
		"symbol": "AU$",
		"symbol_native": "$"
	},
	{
		"code": "TZ",
		"en": "Tanzania",
		"cn": "坦桑尼亚",
		"currency_code": "TZS",
		"currency_cn": "坦桑尼亚先令",
		"currency_en": "Tanzanian Shilling",
		"symbol": "TSh",
		"symbol_native": "TSh"
	},
	{
		"code": "UA",
		"en": "Ukraine",
		"cn": "乌克兰",
		"currency_code": "UAH",
		"currency_en": "Ukrainian Hryvnia",
		"symbol": "₴",
		"symbol_native": "₴"
	},
	{
		"code": "UG",
		"en": "Uganda",
		"cn": "乌干达",
		"currency_code": "UGX",
		"currency_en": "Ugandan Shilling",
		"symbol": "USh",
		"symbol_native": "USh"
	},
	{
		"code": "UM",
		"en": "United States Minor Outlying Islands",
		"cn": "美国本土外小岛屿",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "UY",
		"en": "Uruguay",
		"cn": "乌拉圭",
		"currency_code": "UYU",
		"currency_en": "Uruguayan Peso",
		"symbol": "$U",
		"symbol_native": "$"
	},
	{
		"code": "UZ",
		"en": "Uzbekistan",
		"cn": "乌兹别克斯坦",
		"currency_code": "UZS",
		"currency_en": "Uzbekistan Som",
		"symbol": "UZS",
		"symbol_native": "UZS"
	},
	{
		"code": "VA",
		"en": "Vatican City (The Holy See)",
		"cn": "梵蒂冈",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "VC",
		"en": "St. Vincent & the Grenadines",
		"cn": "圣文森特和格林纳丁斯",
		"currency_code": "XCD"
	},
	{
		"code": "VE",
		"en": "Venezuela",
		"cn": "委内瑞拉",
		"currency_code": "VEF",
		"currency_en": "Venezuelan Bolívar",
		"symbol": "Bs.F.",
		"symbol_native": "Bs.F."
	},
	{
		"code": "VG",
		"en": "British Virgin Islands",
		"cn": "英属维尔京群岛",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "VI",
		"en": "United States Virgin Islands",
		"cn": "美属维尔京群岛",
		"currency_code": "USD",
		"currency_cn": "美元",
		"currency_en": "US Dollar",
		"symbol": "$",
		"symbol_native": "$"
	},
	{
		"code": "VN",
		"en": "Vietnam",
		"cn": "越南",
		"currency_code": "VND",
		"currency_cn": "越南盾",
		"currency_en": "Vietnamese Dong",
		"symbol": "₫",
		"symbol_native": "₫"
	},
	{
		"code": "VU",
		"en": "Vanuatu",
		"cn": "瓦努阿图",
		"currency_code": "VUV"
	},
	{
		"code": "WF",
		"en": "Wallis and Futuna",
		"cn": "瓦利斯和富图纳",
		"currency_code": "XPF"
	},
	{
		"code": "WS",
		"en": "Samoa",
		"cn": "萨摩亚",
		"currency_code": "WST"
	},
	{
		"code": "YE",
		"en": "Yemen",
		"cn": "也门",
		"currency_code": "YER",
		"currency_cn": "也门里亚尔",
		"currency_en": "Yemeni Rial",
		"symbol": "YR",
		"symbol_native": "ر.ي."
	},
	{
		"code": "YT",
		"en": "Mayotte",
		"cn": "马约特",
		"currency_code": "EUR",
		"currency_cn": "欧元",
		"currency_en": "Euro",
		"symbol": "€",
		"symbol_native": "€"
	},
	{
		"code": "ZA",
		"en": "South Africa",
		"cn": "南非",
		"currency_code": "ZAR",
		"currency_cn": "兰特",
		"currency_en": "South African Rand",
		"symbol": "R",
		"symbol_native": "R"
	},
	{
		"code": "ZM",
		"en": "Zambia",
		"cn": "赞比亚",
		"currency_code": "ZMW"
	},
	{
		"code": "ZW",
		"en": "Zimbabwe",
		"cn": "津巴布韦",
		"currency_code": "ZWL",
		"currency_en": "Zimbabwean Dollar",
		"symbol": "ZWL$",
		"symbol_native": "ZWL$"
	}
]
function getCurrencyCode(locale) {
    if (!locale) {
        return undefined; // 处理未提供地区信息的情况
    }
    const matchingLocale = currencyList.find(list => list.code === locale.toUpperCase());
    return matchingLocale?.currency_code; // 可选链式操作，优雅地处理 undefined 的 locale
}
export default class ut {
    static getCurrencyList() {
        return currencyList
    }
    static LangList = [
        {abbreviation: 'en', lanName: '英文',lanNameNew:'English'},
        {abbreviation: 'da', lanName: '丹麦文',lanNameNew:'dansk'},
        {abbreviation: 'de', lanName: '德文',lanNameNew:'Deutsch'},
        {abbreviation: 'es', lanName: '西班牙文',lanNameNew:'Español'},
        {abbreviation: 'fi', lanName: '芬兰文',lanNameNew:'suomi, suomen kieli'},
        {abbreviation: 'fr', lanName: '法文',lanNameNew:'français, langue française'},
        {abbreviation: 'idr', lanName: '印尼文',lanNameNew:'Bahasa Indonesia'},
        // {abbreviation: 'id', lanName: '印尼文',lanNameNew:'Bahasa Indonesia'},
        {abbreviation: 'it', lanName: '意大利文',lanNameNew:'Italiano'},
        {abbreviation: 'ja', lanName: '日文',lanNameNew:'日本語／にほんご'},
        {abbreviation: 'ko', lanName: '韩文',lanNameNew:'한국어'},
        {abbreviation: 'nl', lanName: '荷兰文',lanNameNew:'Nederlands'},
        {abbreviation: 'no', lanName: '挪威文',lanNameNew:'norsk'},
        {abbreviation: 'pl', lanName: '波兰文',lanNameNew:'polski'},
        {abbreviation: 'pt', lanName: '葡萄牙文',lanNameNew:'português'},
        {abbreviation: 'ro', lanName: '罗马尼亚文',lanNameNew:'română'},
        {abbreviation: 'ru', lanName: '俄文',lanNameNew:'русский'},
        {abbreviation: 'sv', lanName: '瑞典文',lanNameNew:'svenska'},
        // {abbreviation: 'th', lanName: '泰文',lanNameNew:'ไทย'},
        {abbreviation: 'th', lanName: '泰文',lanNameNew:'ภาษาไทย'},
        {abbreviation: 'tr', lanName: '土耳其文',lanNameNew:'Türkçe'},
        {abbreviation: 'vi', lanName: '越南文',lanNameNew:'Tiếng Việt'},
        {abbreviation: 'zh', lanName: '中文',lanNameNew:'中文'},
        {abbreviation: 'my', lanName: '缅甸文',lanNameNew:'မြန်မာ'},
    ]
	// 路径/网址join
	static UrlJoin(a?, ...b) {
		if (b == null || b == undefined || b.length == 0) {
			return a
		}

		if (a.endsWith("/")) {
			return a + this.UrlJoin(...b)
		}

		return a + "/" + this.UrlJoin(...b)
	}

	static isNull(s: any) {
		if (s == "" || s == undefined || s == null) {
			return true
		}
		return false
	}

	// 数据格式化接口 fmt == format
	static fmtGoldSigned(v) {
		if (v > 0) {
			return sprintf("+%0.2f", v / 1e4)
		}
		return sprintf("%0.2f", v / 1e4)
	}
	static fmtGold(v) {
		if (!v || v === Infinity) {
			return '0.00'
		}
		return sprintf("%0.2f", v / 1e4)
	}


    static toNumberWithComma(str) {
        // 将字符串转换为数字
        let num = Number(ut.fmtGold(str))
        let formatter = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "THB",
        });
        let a = formatter.format(num)
        let b = (a.split('THB')[0] + a.split('THB')[1]).trim()
        return b;
    }
    static toNumberWithCommaNormal(str) {
        // 将字符串转换为数字
        let num = Number(str);

        // 当数据转换后为NaN的情况下 页面展示转变为0
        if (isNaN(num)){
            num = 0
        }
        let formatter = new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "THB",
        });
        let a = formatter.format(num)
        let b = (a.split('THB')[0] + a.split('THB')[1]).trim()
        return b;
    }

    static symbolCurrency(str, locale,divided=true) {
        let number = Number(ut.fmtGold(str)); // 假设 ut.fmtGold 是有效的函数
        if (!divided) number = Number(str)
        let formatter = new Intl.NumberFormat(locale || navigator.languages[1].toUpperCase(), {
            style: 'currency',
            currency: getCurrencyCode(locale), // 使用 getCurrencyCode 获取适合地区的货币代码
        });
        // 合并格式化逻辑，一次创建包含符号的格式化器
        return formatter.format(number);
    }


    static isBase64(str) {
        return /^data:image\/[a-z]+;base64,/i.test(str);
    }

	static fmtBool(v) { return v ? "是" : "否" }
	static fmtDate(v, format?) {
		if (typeof v === 'number') {
			return moment.utc(new Date(v * 1000)).format("YYYY-MM-DD HH:mm");
		} else {
			if (!v || v.startsWith('0001') || v.startsWith('1970')) {
				return '/'
			}

            if (!format){
                format = "YYYY-MM-DD HH:mm"
            }
			return moment.utc(v).format(format);
		}
	}
	static fmtDate1(v) {
		if (typeof v === 'number') {
			return moment(new Date(v * 1000)).format("YYYY-MM-DD");
		} else {
			if (!v || v.startsWith('0001') || v.startsWith('1970')) {
				return '/'
			}
			return moment(v).format("YYYY-MM-DD");
		}
	}
	static fmtUTCDate(timeStrap: number){

		if (timeStrap){

			let StartTimeData = new Date(timeStrap)


			let TimeFullYear = StartTimeData.getFullYear()
			let TimeMonth = StartTimeData.getMonth()
			let TimeDate = StartTimeData.getDate()
			let TimeHours = StartTimeData.getHours()
			let TimeMinutes = StartTimeData.getMinutes()
			let TimeMilliseconds = StartTimeData.getSeconds()

			timeStrap = new Date(Date.UTC(TimeFullYear,TimeMonth,TimeDate,TimeHours,TimeMinutes,TimeMilliseconds)).getTime() / 1000
		}

		return timeStrap
	}




	static fmtSelectedUTCDate(timeStrap:string | number, handle?: string){

		const store = useStore()


		if (timeStrap){
			const timeZone = store.SelectedTimeZone

			if (typeof timeStrap == "string"){

				let StartTimeData = new Date(timeStrap)


				let TimeFullYear = StartTimeData.getFullYear()
				let TimeMonth = StartTimeData.getMonth()
				let TimeDate = StartTimeData.getDate()
				let TimeHours = StartTimeData.getHours()
				let TimeMinutes = StartTimeData.getMinutes()
				let TimeMilliseconds = StartTimeData.getSeconds()


				timeStrap = new Date(TimeFullYear,TimeMonth,TimeDate,TimeHours,TimeMinutes,TimeMilliseconds).getTime()
			}

			if(handle == 'reduce'){

				timeStrap -= timeZone
			}else{

				timeStrap += timeZone
			}
			timeStrap = new Date(timeStrap).getTime()/ 1000

		}
		return timeStrap
	}

	static fmtSelectedUTCDateFormat(timeStrap, type?){

		let date = "/"
		if (timeStrap){
			let TimeStrap =(ut.fmtSelectedUTCDate(timeStrap, type) as number) * 1000





			if (timeStrap){
				date = moment.utc(TimeStrap).format("YYYY-MM-DD HH:mm:ss")
			}

		}

		return date;
	}


    static fmtDateSecond(v) {
        if (typeof v === 'number') {
            return moment.utc(new Date(v * 1000)).format("YYYY-MM-DD HH:mm:ss");
        } else {
            if (!v || v.startsWith('0001') || v.startsWith('1970')) {
                return '/'
            }
            return moment.utc(v).format("YYYY-MM-DD HH:mm:ss");
        }
    }
	static fmtPercent = (c) => { return (c * 100).toFixed(2) + "%"; }
	static fmtPercentNoNum = (c) => { return (c * 100).toFixed(2); }
	// static fmtPercentFour = (c,num) => { return (c * 100).toPrecision(num) + "%"}
	static fmtPercentFour = (c,num) => {
        /*if (typeof c !== 'number') {
            throw new TypeError('输入参数必须为数字类型');
        }*/
        const percent = c * 100;
        const precision = new Intl.NumberFormat('zh-CN', {
            minimumFractionDigits: num,
            maximumFractionDigits: num,
        }).format(percent);
        return precision + "%";
    }
    static fmtPercentFourNoNum = (c,num) => {
        const percent = c * 100;
        const precision = Math.max(num, 2);
        return Number(percent).toPrecision(precision);
    }
	static imageUrl = (url) => {
		if (url.startsWith('http')) {
			return url
		}
		if (url.startsWith('icon:')) {
			return this.UrlJoin('/up', 'ico', url.split(':')[1])
		}
		return this.UrlJoin('/up', url.split(':')[1])
	}

	static copyTextToClipBoard(v) {
		// copy(v)
		navigator.clipboard.writeText
		let input = document.createElement("input");
		input.value = v;
		document.body.appendChild(input);
		input.select();
		document.execCommand("Copy");
		document.body.removeChild(input);
		tip.s(t("复制"))
	}

	static openIMG(url) {
		let node = document.createElement('a');
		node.target = '_blank';
		node.href = url;
		node.click();
		node.remove();
	}

	static getQueryString(name: string) {
		let url: any = window.location.href;
		let theRequest = new Object();
		if (url.indexOf("?") != -1) {
			url = url.split("?");
			url = url[url.length - 1];
			let strs = url.split("&");
			for (let i = 0; i < strs.length; i++) {
				theRequest[strs[i].split("=")[0]] = unescape(strs[i].split("=")[1]);
			}
		}
		return theRequest[name];
	}
}



let throttleTimer: number | null
export const Throttle = (fn: Function, delay: number) :Function => {


    return (...args: unknown[]) => {

        if (throttleTimer) {
            return;
        }
        throttleTimer = setTimeout(() => {
            fn.apply(this, args);
            throttleTimer = null;
        }, delay);
    }
}



let publicKeys = "-----BEGIN PUBLIC KEY-----\n" +
	"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtDdVQ3Dt8gTpM8slvceE\n" +
	"xzGomzssHmImy09oN1U7gbXcZS4tHUrewoFabz9qjIqu8omVk51D2Nbidg1Gh4uk\n" +
	"AAj10NR+edAlLP22AVzmw6qnONhCT8fRE6OUOokdi9Jy3vfHKpeQYLX/fKxyX+rs\n" +
	"0chVuqw3aKBjWe3WJZm6kGS7KG0wpcs+ToEZ93p0ji8qvOB98O+l9txrOgHGKGZj\n" +
	"YfFWrdaxWa03vM6TV1eZzIjH/jbulMdYOxxDcrEqgwAkQ3ZcV5xGGELZuXfWkyMe\n" +
	"2WoLNaCpAbA1ObMM1Gi5v4uyJLSXMkw/N7CEKbFhPte1YHX3lQKo0SRuVCl/Nf9P\n" +
	"zQIDAQAB\n" +
	"-----END PUBLIC KEY-----"


export const PasswordRSAEncryption = (password) => {

	if (!password){
		password = generatePassword(8)
	}

	const encrypt = new JSEncrypt();
	encrypt.setPublicKey(publicKeys);

	let encryptorvalue = encrypt.encrypt(password)// 对数据进行加密



	return encryptorvalue;



}




export function generatePassword(length) {
	let password = '';
	const charset = 'abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ123456789!@#$%&';

	for (let i = 0; i < length; i++) {
		const randomIndex = Math.floor(Math.random() * charset.length);
		password += charset[randomIndex];
	}

	return password;
}
