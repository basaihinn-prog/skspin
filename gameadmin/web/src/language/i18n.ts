import {createPinia, storeToRefs} from 'pinia'
import {useStore} from "@/pinia/index";
import {createI18n} from 'vue-i18n'
import {reactive} from "vue";
import {Client} from "@/lib/client";
import {AdminGroup} from "@/api/adminpb/group";
import {tip} from "@/lib/tip";
import { ref } from 'vue'
const pinia = createPinia()
const store = useStore(pinia);


const {language,lang} = storeToRefs(store)
let langList = ref(null)
const piniaStore = useStore()
const { setLanguage,setLangRadio,setLang } = piniaStore
const getLangList = (lang, key) => {

    let datas = {}
    langList.value?.forEach(obj => {
        datas[obj['zh']] = obj[lang.toLowerCase()]
    })
    datas = {
        'key': key,
        ...datas
    }
    return datas
}
const initi18n = async () => {
    let [data, err] = await Client.Do(AdminGroup.GetLang, {
        PageIndex: 1,
        PageSize: 10000
    })
    if (err) return tip.e(err)
    data.List = !data.List && data.List.length === 0 ? [] : data.List

    setLang(data.List);
    langList.value = data.List
}
const msgList = async () => {
    if (!langList.value) {
        await initi18n()
    }
    // Iterate over the language list
    const messages = {};
    if (langList.value && langList.value.length) {
        for (const [key, value] of Object.entries(langList.value[0])) {
            if (key !== 'Permission') {
                messages[key?.toLowerCase()] = getLangList(key.toUpperCase(), key);
            }
        }
    }
    return messages;
};
const initI18nInstance = async () => {
    const messages = await msgList();
    const locale = store.AdminInfo.language

    return createI18n({
        locale: locale,
        fallbackLocale: 'zh',
        globalInjection: true,
        legacy: false,
        silentTranslationWarn: true,
        silentFallbackWarn: true,
        missingWarn: false,
        fallbackWarn: false,
        messages: messages
    });
};
const i18n = await initI18nInstance();
export default i18n
