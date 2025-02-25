import { request } from "@/utils/service"

export interface AreaModel extends AddArea {
    id: number
}

export interface AddArea {
    name: string
    remark?: string
    userIds?: number[]
}

export interface EditArea extends CId,AddArea { }

export interface SearchArea extends PageInfo, AddArea { }

export type AreaListData = ListData<AreaModel[]>

export function getAreaListApi(data: SearchArea) {
    return request<ApiResponseData<AreaListData>>({
        url: "/manage/area/getAreaList",
        method: "post",
        data
    })
}

export function createAreaApi(data: AddArea) {
    return request<ApiResponseData<AreaModel>>({
        url: "/manage/area/createArea",
        method: "post",
        data
    })
}

export function editAreaApi(data: EditArea) {
    return request<ApiResponseData<AreaModel>>({
        url: "/manage/area/editArea",
        method: "post",
        data
    })
}

export function deleteAreaApi(data: CId) {
    return request<ApiResponseData<null>>({
        url: "/manage/area/deleteArea",
        method: "post",
        data
    })
}
