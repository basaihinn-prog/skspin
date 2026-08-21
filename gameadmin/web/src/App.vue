<template>
	<el-config-provider :size="store.viewModel === 1 ? 'default' : 'small'" :locale="code[store.language]">
		<router-view />
	</el-config-provider>
</template>

<script setup lang='ts'>
import {ref, onMounted, watch, reactive} from 'vue'
import { ElConfigProvider } from 'element-plus'
import { Client } from '@/lib/client';
import { DevicePixelRatio } from '@/lib/DevicePixelRatio';
import router from '@/router'
import { useStore } from '@/pinia/index';
import { useI18n } from 'vue-i18n';

import en from "element-plus/es/locale/lang/en";
import da from "element-plus/es/locale/lang/da";
import de from "element-plus/es/locale/lang/de";
import es from "element-plus/es/locale/lang/es";
import fi from "element-plus/es/locale/lang/fi";
import fr from "element-plus/es/locale/lang/fr";
import id from "element-plus/es/locale/lang/id";
import it from "element-plus/es/locale/lang/it";
import ja from "element-plus/es/locale/lang/ja";
import ko from "element-plus/es/locale/lang/ko";
import nl from "element-plus/es/locale/lang/nl";
import no from "element-plus/es/locale/lang/nb-no";
import pl from "element-plus/es/locale/lang/pl";
import pt from "element-plus/es/locale/lang/pt";
import ro from "element-plus/es/locale/lang/ro";
import ru from "element-plus/es/locale/lang/ru";
import sv from "element-plus/es/locale/lang/sv";
import tr from "element-plus/es/locale/lang/tr";
import vi from "element-plus/es/locale/lang/vi";
import zhCn from "element-plus/es/locale/lang/zh-cn";
import my from "element-plus/es/locale/lang/my";
import th from "element-plus/es/locale/lang/th";

let { t } = useI18n()
const store = useStore()
let langList = reactive(null)
if (store.Token) {
	Client.setToken(store.Token)
} else {
	router.push({ path: '/login' })
}

const tableSetting = JSON.parse(localStorage.getItem('tableSetConfig'));

if (tableSetting) {
	store.setTableSetConfig(tableSetting)
}
const selectedTimeZone = parseInt(localStorage.getItem('timeZone'));

if (selectedTimeZone) {
	store.setSelectedTimeZone(selectedTimeZone)
}



let code = {
	"en": en,
	"da": da,
	"de": de,
	"es": es,
	"fi": fi,
	"fr": fr,
	"idr": id,
	"it": it,
	"ja": ja,
	"ko": ko,
	"nl": nl,
	"no": no,
	"pl": pl,
	"pt": pt,
	"ro": ro,
	"ru": ru,
	"sv": sv,
	"th": th,
	"tr": tr,
	"vi": vi,
	"zh": zhCn,
	"my": my,
}



onMounted(() => {
	// getLang()
    // document.body.offsetWidth = devicePixelRatio
    // new DevicePixelRatio().init()
    onWidth()
	window.onresize = () => {
		return (() => {
			onWidth()
		})()
	}

	if (store.language == "zh"){
		document.documentElement.style.setProperty('--fontsSize', '14px')
		document.documentElement.style.setProperty('--labelSize', '14px')
	}else{
		document.documentElement.style.setProperty('--fontsSize', '12px')
		document.documentElement.style.setProperty('--labelSize', '12px')
	}
})

const onWidth = () => {
	store.setViewModel(document.body.clientWidth > 1000 ? 1 : 2)
}

</script>

<style lang="scss">
.el-statistic {
    --el-statistic-content-font-size: 22px !important;
}
.switchContainer{
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: flex-end;

}
.tip{
	width: 100%;
	font-size: 14px;
	display: flex;
	align-items: center;
	padding-bottom: 10px;
	.red {
		color: #f56c6c;
	}

}
.el-form--default.el-form--label-top {
    .el-form-item {
        .el-form-item__label{
            font-size: 11px;
            color: rgba(18,31,62,.74);
            letter-spacing: .92px;
            margin: 0 0 3px;
            line-height: normal;
        }
    }
}
.el-date-editor{
    //border: 1px solid rgba(18,31,62,.15);
	font-size: 14px;
}
.el-input,.el-date-editor,.el-select-v2{
    //border: 1px solid rgba(18,31,62,.15);
}


.el-form-item--default{
    margin-bottom: 0;
}
.login-form{
    .el-form-item--default{
        margin-bottom: 18px;
    }
}
.el-input-group__prepend .el-input{
    box-shadow: none;
    border: none;
}
.searchView {
    width: auto;
    display: flex;
    align-items: flex-end;
		background: white;

	padding: 15px;
	border-radius: 15px;

	//align-items: center;
	//flex-wrap: wrap;
	//align-content: space-between;
	//row-gap: 10px;
	//column-gap: 16px;
	//padding-bottom: 10px;
    .el-space{
        margin-bottom: 12px !important;
				margin-right: 15px;
        .el-space__item{
            padding: 0 !important;
            &:last-child{
                //float: right;
                //min-height: 50px;
            }
        }
        .el-button{
            margin: 0;
        }

    }



.el-input{
	width: 150px;
}
.el-button+.el-button {
		margin: 0;
	}

	>.searchItem {
		font-size: 14px;
		display: flex;
		align-items: center;

		.el-input {
			width: 200px;
		}

		p {
			padding-right: 4px;
			white-space: nowrap;
		}

		.el-date-editor--daterange {
			width: 260px;
		}

		.red {
			color: #f56c6c;
		}
	}

}

.searchViewShadow{
	box-shadow: 0 2px 15px #bababa;
	-moz-box-shadow: 0 2px 15px #bababa;
	-webkit-box-shadow: 0 2px 15px #bababa;
	padding: 10px 8px;
	border-radius: 10px;
	background: #ffffff;
}

.page_table_context{
	width: 100%;
	height: auto;
	border-radius: 10px;
	background: #ffffff;
	margin-top: 20px;
	padding: 20px 0;
}

.cell, .cell * {
	font-size: 14px;
}
.font_size16{
	font-size: $size16;
}

.page_table_context .customTable {
	font-size: 14px;
}

.elPagination {
	padding-top: 10px;
}

.el-descriptions__body {
	overflow-x: scroll;
}

.el-popover.el-popper {
	min-width: 300px;
    max-height: 800px;
    overflow-y: auto;
    .el-card__body{
        padding: .6rem;
    }
}

.el-popconfirm {
	.el-popconfirm__action {
		button {
			min-width: 60px;
			min-height: 30px;
		}

		.el-button+.el-button {
			margin-left: 20px;
		}
	}

}

::-webkit-scrollbar {
	width: 2px;
	height: 2px;
}

::-webkit-scrollbar-thumb {
	border-radius: 1em;
	background-color: rgba(50, 50, 50, .3);
}
// ::-webkit-scrollbar-track {
//   border-radius: 1em;
//   background-color: rgba(50, 50, 50, .1);
// }
//.el-dialog__header,.el-dialog__footer{
//    padding: 21px 25px;
//    margin: 0;
//}
//
//.el-dialog__header{
//    border-bottom: 1px solid rgba(33,33,33,.12);
//}
//.el-dialog__footer{
//    border-top: 1px solid rgba(33,33,33,.12);
//    .el-button{
//        margin: 0;
//    }
//}
//.slot-dialog{
//    .el-dialog__body{
//        height: 100%;
//        .el-divider--vertical{
//            border-left-color: rgb(40, 40, 51);
//            border-left-width: 3px;
//        }
//    }
//}
.el-dialog__body{
    .el-form{
        .el-form-item{
            margin-bottom: 18px;
        }
    }
}

.text-center{
    text-align: center;
    display: flex;
    justify-content: center;
    margin-bottom: .5rem;
}
.flex_vbox .text-center{
    flex-wrap: wrap;
}

.dropdown-item-style{
    display: flex;
    justify-content: space-between;
    background: transparent;
}
.sort-table-div{
    display: flex; align-items: center
}

.sort-table-column{
    display: inline-flex;
    flex-direction: column;
    align-items: center;
    height: 14px;
    width: 24px;
    vertical-align: middle;
    cursor: pointer;
    overflow: initial;
    position: relative;
}

.sort-table-column i:first-child{
    top: -6px;
    border: solid 6px transparent;
    position: absolute;
    left: 7px;
    border-bottom-color:#a8abb2
}
.sort-table-column i:last-child{
    border: solid 6px transparent;
    position: absolute;
    left: 7px;
    bottom: -6px;
    border-top-color: #a8abb2
}
.sort-border-color-top{
    border-bottom-color: #409eff !important;
}
.sort-border-color-bottom{
    border-top-color: #409eff !important;
}
.el-form-item__label{
	font-size: $labelSize;
	padding-left: 10px;
	line-height: normal!important;
	display: flex;
	align-items: center;
	justify-content: flex-start;

}

.game_popper, .game_popper .el-select-dropdown__list, .game_popper .el-select-dropdown{
	width: 300px!important;
}


.operator_proper, .operator_proper .el-select-dropdown__list, .operator_proper .el-select-dropdown{
	width: 150px!important;
	overflow-x: hidden;
}


.manufacturer_container, .manufacturer_container .el-select-dropdown__item{
	padding: 3px 0;
	margin-bottom: 10px;
	height: auto;

}

.manufacturer_options_all{
	width: 90%;
	display: block;

	margin: 0 auto;
}
.manufacturer_options{
	width: 90%;
	display: block;
	border-radius: 5px;
	background: #e8f1ff;

	margin: 0 auto;
	color: var(--el-color-primary);
	line-height: 25px;
	text-align: center;
}

.createSucee.el-message-box{
	--el-messagebox-width: 650px;
}

</style>
