import type {FilterInterface} from "@/interface/filterInterface";

export interface manufacturerInterface extends FilterInterface{
    id: number
    label: string
    localHref: string
}
