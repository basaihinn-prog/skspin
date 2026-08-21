import { h, resolveComponent } from 'vue'
export default class VNodes {
    data = []

    h(...args) {
        const vnode = h(...args)
        this.data.push(vnode)
        return this
    }

    r(name, ...args) {
        const vnode = r(name, ...args)
        this.data.push(vnode)
        return this
    }
}

export function modelProp(v) {
    return {
        "modelValue": v.value,
        "onUpdate:modelValue": (value) => v.value = value,
    }
}

export function modelPropNumber(v) {
    return {
        "modelValue": v.value.toString(),
        "onUpdate:modelValue": (value) => v.value = parseInt(value),
    }
}

export function r(name, ...args) {
    return h(resolveComponent(name), ...args)
}