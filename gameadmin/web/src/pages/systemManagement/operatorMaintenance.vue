<template>
    <div class="adminList">
        <div class="searchView">
            <el-space wrap>
                <el-button type="primary" @click="clickAddDialog" circle>
                    <el-icon>
                        <Plus/>
                    </el-icon>
                </el-button>
                <div class="table-icon" @click="xiazai({OperatorId:0})">
                    <el-icon>
                        <Download/>
                    </el-icon>
                </div>
            </el-space>
        </div>
        <!-- 数据 -->
        <el-table :data="tableData" class="elTable" v-loading="loading"
                  :header-cell-style="{ background: '#F5F7FA', color: '#333333' }">
            <el-table-column prop="Id" label="Id"/>
            <el-table-column prop="Name" :label="$t('商户AppID')"/>
            <el-table-column prop="CurrencyKey" :label="$t('货币')"/>

            <el-table-column prop="Status" :label="$t('是否开启')">
                <template #default="scope">
                    <el-switch v-model="scope.row.Status" :active-value="0" :inactive-value="1"
                               :active-text="$t('是')" inline-prompt :inactive-text="$t('否')"
                               :before-change="() => changeStatus(scope.row)"/>
                </template>
            </el-table-column>
            <el-table-column :label="$t('操作')">
                <template #default="scope">
                    <el-space wrap style="margin: 0">
                        <div class="table-icon" @click="editOperator(scope.row)">
                            <el-icon>
                                <Setting/>
                            </el-icon>
                        </div>
                        <div class="table-icon" @click="xiazai({OperatorId:scope.row.Id})">
                            <el-icon>
                                <Download/>
                            </el-icon>
                        </div>
                    </el-space>
                </template>
            </el-table-column>

        </el-table>
        <!-- 分页 -->
        <Pagination :Count='Count' @paginationEmit='paginationEmit'></Pagination>
        <!-- 添加弹框 -->
        <el-dialog v-model="addDialog" @closed="resetForm(addFormRef)" :title="$t('管理员')"
                   :width="store.viewModel === 2 ? '100%' : '650px'">
            <el-form ref="addFormRef" :model="addForm" :rules="addRules" label-position="top" label-width="100px">
                <el-form-item :label="$t('AppID')" prop="Name">
                    <el-input v-model="addForm.Name" :placeholder="$t('请输入')"/>
                </el-form-item>
                <el-form-item :label="$t('商户回调地址')" prop="Address">
                    <el-input v-model.trim="addForm.Address" clearable :placeholder="$t('请输入')" size="default"/>
                </el-form-item>
                <el-form-item :label="$t('机器人识别码')" prop="Robot">
                    <el-input v-model.trim="addForm.Robot" clearable :placeholder="$t('请输入')" size="default"/>
                </el-form-item>
                <el-form-item :label="$t('chatid')" prop="ChatID">
                    <el-input v-model.trim="addForm.ChatID" clearable :placeholder="$t('请输入')" size="default"/>
                </el-form-item>
                <el-form-item :label="$t('白名单')" prop="WhitelistString">
                    <el-input
                        v-model="addForm.WhitelistString"
                        autosize
                        type="textarea"
                        :placeholder="$t('请用英文逗号进行分隔')"
                    />
                </el-form-item>
                <el-form-item :label="$t('货币')" prop="CurrencyKey">
                    <el-input
                        v-model="addForm.CurrencyKey"
                        autosize
                        :placeholder="$t('请用货币缩写进行搜索')"
                    />
                </el-form-item>
                <el-form-item :label="$t('钱包类型')" prop="WalletMode">
                    <el-select
                        v-model="addForm.WalletMode"
                        :placeholder="$t('请选择钱包类型')"
                        style="width: 240px"
                    >
                        <el-option
                            v-for="item in walletOptions"
                            :key="item.value"
                            :label="item.label"
                            :value="item.value"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item :label="$t('权限设置')" prop="MenuIds">
                    <el-tree ref="treeRef" :data="menuTreeList" show-checkbox accordion node-key="ID"
                             :props="{children: 'Children',label: 'Title',}"/>
                </el-form-item>

                <el-form-item :label="$t('游戏列表')">
                    <div style="margin-bottom: 1rem">
                        <el-switch v-model="gameListSwitch" inline-prompt @change="gameListSwitchChange"
                                   :inactive-text="$t('勾选方案')"
                                   :active-text="$t('搜索游戏')"/>
                        <div style="width: 100%">
                            <el-input v-model="gameListSearch" v-if="gameListSwitch" style="width: 240px" clearable :placeholder="$t('请搜索需要勾选的游戏')"
                                      @change="gameListSearchChange"/>
                            <el-select v-else
                                       v-model="gameListSelect"
                                       :placeholder="$t('请选择你需要的方案')" clearable
                                       style="width: 240px" @change="gameListSelectChange"
                            >
                                <el-option
                                    v-for="(item,index) in gameListSelectList"
                                    :key="item.ID"
                                    :label="item.Name"
                                    :value="item.ID+'++' + item.ExcluedGameIds"
                                />
                            </el-select>
                        </div>
                    </div>
                    <el-collapse v-model="gameListCollapse" style="width: 100%;">
                        <el-collapse-item :title="$t('已显示游戏列表')" name="1">
                            <el-checkbox-group v-model="gameListChecked">
                                <el-checkbox v-for="item in gameListNew" :label="item.ID"
                                             style="width: 150px;overflow: hidden;">{{item.Name}}</el-checkbox>
                            </el-checkbox-group>
                        </el-collapse-item>
                    </el-collapse>
                    <div style="width: 100%;margin-top: 1rem" >
                        <el-space fill>
                            <el-form-item :label="$t('重新保存至到我的方案中')">
                                <el-input v-model="scheme" style="width: 240px" :disabled="!gameListSwitch" clearable
                                          :placeholder="$t('您可以为新的方案命名')">
                                </el-input>
                                <el-button type="primary" :disabled="!scheme.trim()" @click="gameListCheckedSave"
                                           style="margin-bottom: 0;margin-left: 1rem">
                                    {{!gameListSwitch?$t('重新保存方案'):$t('保存方案')}}
                                </el-button>
                            </el-form-item>
                        </el-space>
                    </div>
                </el-form-item>

                <el-form-item :label="$t('备注')" prop="Remark">
                    <el-input v-model="addForm.Remark" :placeholder="$t('请输入')"/>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="AddAdminer">{{ $t('添加') }}</el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog v-model="editDialog" :title="$t('管理员')" :width="store.viewModel === 2 ? '85%' : '650px'">
            <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-position="top" label-width="100px">
                <el-form-item :label="$t('商户AppID')" prop="WalletMode">
                    <el-tag type="info">
                        {{ editForm.AppID }}
                    </el-tag>
                </el-form-item>
                <el-form-item :label="$t('商户AppSecret')" prop="WalletMode">
                    <el-tag type="info">
                        {{ editForm.AppSecret }}
                    </el-tag>
                </el-form-item>

                <el-form-item :label="$t('钱包类型')" prop="WalletMode">
                    <el-tag type="info">
                        {{ walletOptions.find(i=>i.value === editForm.WalletMode)?.label }}
                    </el-tag>
                </el-form-item>
                <el-form-item :label="$t('商户回调地址')" prop="Address">
                    <el-input v-model.trim="editForm.Address" clearable :placeholder="$t('请输入')" size="default"/>
                </el-form-item>
                <el-form-item :label="$t('机器人识别码')" prop="Robot">
                    <el-input v-model.trim="editForm.Robot" clearable :placeholder="$t('请输入')" size="default"/>
                </el-form-item>
                <el-form-item :label="$t('chatid')" prop="ChatID">
                    <el-input v-model.trim="editForm.ChatID" clearable :placeholder="$t('请输入')" size="default"/>
                </el-form-item>
                <el-form-item :label="$t('白名单')" prop="WhitelistString">
                    <el-input
                        v-model="editForm.WhitelistString"
                        autosize
                        type="textarea"
                        :placeholder="$t('请用英文逗号进行分隔')"
                    />
                </el-form-item>
                <el-form-item :label="$t('权限设置')" prop="MenuIds">
                    <el-tree ref="editTreeRef" :data="menuTreeList" show-checkbox accordion node-key="ID"
                             :props="{children: 'Children',label: 'Title',}"/>
                </el-form-item>
                <el-form-item :label="$t('游戏列表')">
                    <div style="margin-bottom: 1rem">
                        <el-switch v-model="gameListSwitch" inline-prompt @change="gameListSwitchChange"
                                   :inactive-text="$t('勾选方案')"
                                   :active-text="$t('搜索游戏')"/>
                        <div style="width: 100%">
                            <el-input v-model="gameListSearch" v-if="gameListSwitch" style="width: 240px" clearable :placeholder="$t('请搜索需要勾选的游戏')"
                                      @change="gameListSearchChange"/>
                            <el-select v-else
                                       v-model="gameListSelect"
                                       :placeholder="$t('请选择你需要的方案')" clearable
                                       style="width: 240px" @change="gameListSelectChange"
                            >
                                <el-option
                                    v-for="(item,index) in gameListSelectList"
                                    :key="item.ID"
                                    :label="item.Name"
                                    :value="item.ID+'++' + item.ExcluedGameIds"
                                />
                            </el-select>
                        </div>
                    </div>
                    <el-collapse v-model="gameListCollapse" style="width: 100%;">
                        <el-collapse-item :title="$t('已显示游戏列表')" name="1">
                            <el-checkbox-group v-model="gameListChecked">
                                <el-checkbox v-for="item in gameListNew" :label="item.ID"
                                             style="width: 150px;overflow: hidden;">{{item.Name}}</el-checkbox>
                            </el-checkbox-group>
                        </el-collapse-item>
                    </el-collapse>
                    <div style="width: 100%;margin-top: 1rem" >
                        <el-space fill>
                            <el-form-item :label="$t('重新保存至到我的方案中')">
                                <el-input v-model="scheme" style="width: 240px" :disabled="!gameListSwitch" clearable
                                          :placeholder="$t('您可以为新的方案命名')">
                                </el-input>
                                <el-button type="primary" :disabled="!scheme.trim()" @click="gameListCheckedSave"
                                           style="margin-bottom: 0;margin-left: 1rem">
                                    {{!gameListSwitch?$t('重新保存方案'):$t('保存方案')}}
                                </el-button>
                            </el-form-item>
                        </el-space>
                    </div>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="saveAdminer">{{ $t('添加') }}</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>

<script lang='ts' setup>
import {onMounted, ref, reactive, toRefs, nextTick} from 'vue'
import ReturnErrorHandle, {Client} from '@/lib/client';
import {InfoFilled} from '@element-plus/icons-vue'
import {tip} from '@/lib/tip';
import {useStore} from '@/pinia/index';
import type {FormInstance, FormRules} from 'element-plus'
import {useI18n} from 'vue-i18n';
import {AddOperatorReq, AdminInfo} from '@/api/adminpb/info';
import {AdminGroup} from '@/api/adminpb/group';
import {ElNotification, ElTree} from "element-plus";
import {Setting} from "@element-plus/icons-vue";
import ut from "@/lib/util";
import {excel} from "@/lib/excel";
import {saveAs} from 'file-saver';
import {AdminGameCenter} from "@/api/gamepb/admin";
import {GameStatus} from "@/api/gamepb/customer";
import axios from "axios";
import {storeToRefs} from "pinia";
import {computed} from "vue/dist/vue";

let {t} = useI18n()
const store = useStore()
let tableData = reactive([])
let param = reactive({
    PageIndex: 1,
    PageSize: 20
})
let Count = ref(0)
let loading = ref(false)
let addFormName = ref("")
let addFormAppId = ref(addFormName)
let menuTreeList = reactive([])
const treeRef = ref<InstanceType<typeof ElTree>>()
const editTreeRef = ref<InstanceType<typeof ElTree>>()
let checkedArr: number[] = []
let addDialog = ref(false)
let editDialog = ref(false)
const addFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()
let addForm = reactive({
    Name: addFormName,
    AppID: addFormAppId,
    MenuIds: [],
    Remark: '',
    Address: '',
    WhitelistString: '',
    CurrencyKey: '',
    Robot: '',
    ChatID: '',
    WalletMode: 2,
    WhiteIps: [],
    ExcluedGameId:"",
    ExcluedGameIds:[]
    // Status: true
})
let editForm = reactive({
    AppID: '',
    AppSecret: '',
    Name: '',
    OperatorId: '',
    Robot: '',
    ChatID: '',
    MenuIds: [],
    Remark: '',
    Address: '',
    WalletMode: '',
    WhitelistString: '',
    WhiteIps: [],
    ExcluedGameId:"000000000000000000000000",
    ExcluedGameIds:[]
    // Status: true
})

const isIPv4validator = (rule: any, value: any, callback: any) => {
    let whitelistString = value.split(',').filter(item => item !== "")
    const isIPv4 = new RegExp(/^((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})$/);
    let whiteBools = whitelistString.filter(t => isIPv4.test(t));
    if (whiteBools.length) {
        callback()
    } else {
        callback(new Error('白名单不合法'))
    }
}
const addRules = reactive<FormRules>({
    Name: [{required: true, message: t('商户AppID不能为空'), trigger: 'blur',}],
    Address: [{required: true, message: t('商户回调地址不能为空'), trigger: 'blur'}],
    WhitelistString: [{required: true, trigger: 'blur', validator: isIPv4validator,}],
    MenuIds: [{required: true, message: t('请给商户管理员配置权限'), trigger: 'change'}],
})
const editRules = reactive<FormRules>({
    Name: [{required: true, message: t('商户AppID不能为空'), trigger: 'blur'}],
    Address: [{required: true, message: t('商户回调地址不能为空'), trigger: 'blur'}],
    WhitelistString: [{required: true, validator: isIPv4validator, trigger: 'blur'}],
    MenuIds: [{required: true, message: t('请给商户管理员配置权限'), trigger: 'change'}],
})
/*function internationalizeData(data) {
    // 假设您已经有了语言翻译的资源，例如一个翻译函数 translate(key, language)
    // 递归遍历数据并翻译标题
    for (const item of data) {
        item.Title = t(item.Title).replace(/^\s+|\s+$/g, "");
        if (item.Children) {
            internationalizeData(item.Children);
        }
    }
    return data;
}*/
onMounted(async () => {
    menuTreeList = store.AdminInfo.MenuList
    console.log(store.AdminInfo.MenuList, 'menuTreeList');
    getList()
    queryList()
    getGroupList()
    // menuTreeList = internationalizeData(JSON.parse(localStorage.getItem('game_store')).AdminInfo.MenuList)
    // await getList()
    // await queryList()
    // await getGroupList()

});

const getList = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.GetOperatorList, param)
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    Count.value = data?.AllCount
    tableData = data.AllCount === 0 ? [] : data.List
}
let groupList = reactive([])
const getGroupList = async () => {
    let [data, err] = await Client.Do(AdminGroup.GroupList, {
        PageIndex: 1,
        PageSize: 10000
    })
    if (err) {
        return tip.e(err)
    }
    groupList = data.Count === 0 ? [] : data.List
}

const changeOpenGoole = (IsOpenGoole, Pid) => {
    return new Promise(async (resolve) => {
        let [data, err] = await Client.Do(AdminInfo.UpdateOpenGoole, {
            IsOpenGoole: !IsOpenGoole, Pid
        })
        if (err) {
            tip.e(err)
            return resolve(false)
        }
        return resolve(true)
    })
}

function generateTree(data, selectedIds) {
    const result = [];
    for (const item of data) {
        if (selectedIds.includes(item.ID)) {
            // 复制选中的节点，注意要深拷贝
            const newItem = {...item, Children: []};
            result.push(newItem);
            if (item.Children && item.Children.length > 0) {
                newItem.Children = generateTree(item.Children, selectedIds);
            }
        }
    }
    return result;
}

const editOperator = async (params) => {
    try {
        editDialog.value = true
        editForm.OperatorId = params.Id
        editForm.WhitelistString = params.WhiteIps.toString()
        Object.assign(editForm, params);
        await nextTick(() => {
            editTreeRef.value!.setCheckedKeys(editForm.MenuIds, false)
        })
        await getExcludeGameAll()
        if (gameListSelectList.value.length) {
            gameListSwitch.value = false
            if (!params.ExcluedGameIds) {
                gameListSelect.value = null
                return tip.e(t('该商户未关联方案，请确认'))
            }
            const matchedScheme = gameListSelectList.value.find(
                item => params.ExcluedGameId === item.ID
            );
            if (!matchedScheme) {
                scheme.value = ''
                gameListSelect.value = null
                return tip.e(t('该商户未关联方案，请确认'))
            }
            const findValue = gameListSelectList.value.find(item => item.ID === params.ExcluedGameId)
            gameListSelect.value = params.ExcluedGameId + '++' +findValue.ExcluedGameIds.toString()
            gameListChecked.value = findValue.ExcluedGameIds
            scheme.value =  matchedScheme.Name;
        }
    } catch (error) {
        console.error('An error occurred while editing the operator:', error);
        tip.e(t('编辑商户时发生错误，请重试'));
    }
}
const saveAdminer = async () => {
    editForm.WhiteIps = editForm.WhitelistString.split(',')
    editForm.MenuIds = editTreeRef.value!.getCheckedKeys(false)
    if (!gameListSwitch.value) {
        let gameListSelectValues = gameListSelect.value?.split('++')[1].split(',')
        let gameListSelectID = gameListSelect.value?.split('++')[0]
        editForm.ExcluedGameIds = gameListSelectValues
        editForm.ExcluedGameId = gameListSelectID
    } else {
        editForm.ExcluedGameIds = gameListChecked.value
        editForm.ExcluedGameId = null
    }
    await editFormRef.value.validate(async (valid, fields) => {
        if (valid) {
            let [data, err] = await Client.Do(AdminGroup.EditOperator, editForm)
            if (err) {
                return tip.e(err)
            }
            tip.s(t('修改成功'))
            editDialog.value = false
            await getList()
        }
    })
}
let currencyList = ref([])
const walletOptions = [
    {
        label: t('单一钱包'),
        value: 2
    },
    {
        label: t('转账钱包'),
        value: 1
    }
]
let gameList = ref(null)
let scheme = ref('')
let gameListSwitch = ref(true)

let gameListNew = ref(null)
let gameListCollapse = ref(false)
let gameListSearch = ref(null)
let gameListSelect = ref(null)
let gameListSelectList = ref(null)
let gameListChecked = ref([])
const getExcludeGameAll = async () => {
    loading.value = true
    let [data, err] = await Client.Do(AdminInfo.GetexcludeGameAll, {})
    loading.value = false
    if (err) {
        return tip.e(err)
    }
    gameListSelectList.value = ReturnErrorHandle(data,'Name')

}

const StatusMap: JsMap<GameStatus, string> = {
    [GameStatus.Open]: t('正常'),
    [GameStatus.Maintenance]: t('维护'),
    [GameStatus.Hide]: t('隐藏'),
    [GameStatus.Close]: t('关闭'),
};
const queryList = async () => {
    try {
        loading.value = true;
        const [data, err] = await Client.Do(AdminGameCenter.GameList, {});
        if (err) {
            tip.e(err);
            return;
        }
        let list  = ReturnErrorHandle(data.List,'Name')
        if (list) {
            gameListNew.value = gameList.value = sortGamesByName(list);
        }

    } catch (error) {
        tip.e(t('获取游戏列表时发生错误'));
        console.error('Query game list error:', error);
    } finally {
        loading.value = false;
    }
}

const sortGamesByName = (list) => {
    return list.filter(item=> item.Status === 0 ).sort((a, b) => {
        const nameA = a.Name.toUpperCase(); // Convert names to uppercase for comparison
        const nameB = b.Name.toUpperCase();
        return nameA.localeCompare(nameB); // Use localeCompare for proper sorting
    });
};

const gameListSwitchChange = () => {
    gameListSearch.value=''
    console.log(gameListSelect.value);
    gameListChecked.value = []
    scheme.value = ''
}

const gameListSearchChange = () => {
    if (!gameListSearch.value) {
        gameListNew.value = gameList.value
        return
    }

    const searchTerm = gameListSearch.value.toLowerCase(); // Convert search term to lowercase for case-insensitive matching
    gameListNew.value = gameListNew.value.filter(item => {
        const itemName = item.Name.toLowerCase(); // Convert item name to lowercase for case-insensitive matching
        return itemName.includes(searchTerm) || itemName.startsWith(searchTerm); // Check for inclusion or prefix match
    });
}
const arraysEqual = (arr1, arr2) => {
    if (arr1.length !== arr2.length) return false;
    for (let i = 0; i < arr1.length; i++) {
        if (arr1[i] !== arr2[i]) return false;
    }
    return true;
};
const gameListSelectChange = () => {
    console.log(gameListSelect.value);
    let gameListSelectValues = gameListSelect.value.split('++')[1].split(',')
    let gameListSelectID = gameListSelect.value.split('++')[0]
    gameListChecked.value = gameListNew.value.filter(item => gameListSelectValues.includes(item.ID)).map(item => item.ID)
    const matchedScheme = gameListSelectList.value.find(
        item => arraysEqual(item.ExcluedGameIds, gameListSelectValues) && gameListSelectID === item.ID
    );
    scheme.value = matchedScheme ? matchedScheme.Name : null;
    if (!gameListSelectValues.length) {
        gameListNew.value = gameList.value
        gameListChecked.value = []
    }
}
const gameListCheckedSave = async () => {
    if (!scheme.value) {
        return tip.w(t('请为新的方案命名'))
    }
    if (!gameListSwitch.value) {
        let gameListSelectID = gameListSelect.value.split('++')[0]
        if (!gameListSelectID) {
            tip.e(t('未找到游戏名称，请刷新页面后重试'))
        }
        let list = {"ID":gameListSelectID,"Name":scheme.value,"Remark":"","ExcluedGameIds":gameListChecked.value}
        let [data, err] = await Client.Do(AdminInfo.GetExcludeGameUp, list)
        if (err) {
            return tip.e(err)
        }
        tip.s(t('修改方案成功'))
    } else {
        let list = {"ID":null,"Name":scheme.value,"Remark":"","ExcluedGameIds":gameListChecked.value}
        let [data, err] = await Client.Do(AdminInfo.GetExcludeGameAdd, list)
        if (err) {
            return tip.e(err)
        }
        tip.s(t('保存方案成功'))
        await getExcludeGameAll()
    }

}
const clickAddDialog = () => {
    addDialog.value = true
    getExcludeGameAll()
    gameListChecked.value = []
    currencyList.value = ut.getCurrencyList().map(list => {
        return {
            value: list.currency_code,
            label: (list.currency_code || '') + '-' + list.en
        }
    })
}
const resetForm = () => {
    if (!addFormRef) return
    addFormRef.value.resetFields()
}
const AddAdminer = async () => {
    if (!addFormRef) return
    addForm.WhiteIps = addForm.WhitelistString.split(',')
    addForm.MenuIds = treeRef.value!.getCheckedKeys(false)
    if (!gameListSwitch.value) {
        let gameListSelectValues = gameListSelect.value.split('++')[1].split(',')
        let gameListSelectID = gameListSelect.value.split('++')[0]
        addForm.ExcluedGameIds = gameListSelectValues
        addForm.ExcluedGameId = gameListSelectID
    } else {
        addForm.ExcluedGameIds = gameListChecked.value
        addForm.ExcluedGameId = null
    }
    addForm['ExcluedGameIds'] = gameListChecked.value
    console.log(addForm);
    await addFormRef.value.validate(async (valid, fields) => {
        if (valid) {
            let [data, err] = await Client.Do(AdminInfo.AddOperator, addForm)
            if (err) {
                return tip.e(err)
            }
            tip.s(t('添加成功'))
            addDialog.value = false
            await getList()
        }
    })
}

const DelAdminer = async (Pid) => {
    let [data, err] = await Client.Do(AdminInfo.DelAdminer, {Pid})
    if (err) {
        return tip.e(err)
    }
    tip.s(t('删除成功'))
    await getList()
}

const paginationEmit = (PageIndex, PageSize) => {
    param.PageIndex = PageIndex
    param.PageSize = PageSize
    getList()
}
let xiazai = async (list) => {
    let [data, err] = await Client.Do(AdminGroup.DownloadOperatorData, list)
    if (err) {
        return tip.e(err)
    }
    let datas = data.List.map(t => {
        return {
            ...t,
            Url: window.location.href
        }
    })

    const row = datas[0]

    const txt = JSON.stringify(datas, null, 2)
    let strData = new Blob([txt], {type: 'text/plain;charset=utf-8'});
    saveAs(strData, `商户信息-${row.AppId}.txt`)
    // let str = '王佳伟Vue字符串保存到txt文件下载到电脑案例'
    // let strData = new Blob([str], { type: 'text/plain;charset=utf-8' });
    // saveAs(strData, "测试文件下载.txt");

    /*
    excel.dump(datas, `商户信息-${row.AppId}`, [
        {key: "AppId", name: "AppID"},
        {key: "AppSecret", name: "AppSecret"},
        {key: "Address", name: "商户回调地址"},
        {key: "WhitelistString", name: "白名单"},
        {key: "Url", name: "后台管理的地址"},
        {key: "AdminPassword", name: "商户管理员用户密码"},
        {key: "AdminUserName", name: "商户管理员用户名"},
        {key: "GoogleCode", name: "商户管理员谷歌验证码"},
        {key: "Qrcode", name: "商户管理员绑定谷歌的二维码"},
        {key: "ApiUrl", name: "Api调用地址"},
    ])
    */
}
const changeStatus = async (list) => {
    if (!list.Status) {
        if (window.confirm('是否冻结商户列表')) {
            let [data, err] = await Client.Do(AdminGroup.UpdateOperatorStatus, {OperatorId: list.Id, Status: 1})
            if (err) {
                return tip.e(err)
            }
            await getList()
        } else {
            list.Status = 0
        }
    } else {
        list.Status = 0
        let [data, err] = await Client.Do(AdminGroup.UpdateOperatorStatus, {OperatorId: list.Id, Status: list.Status})
        if (err) {
            return tip.e(err)
        }
    }
    tip.i('该操作在重新登录后生效')
}

</script>
<style scoped lang='scss'></style>
