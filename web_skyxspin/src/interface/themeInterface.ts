import type {FilterInterface} from "@/interface/filterInterface";

export interface ThemeInterface extends FilterInterface{
    id: number
    label: string
    gameSource: string
}
