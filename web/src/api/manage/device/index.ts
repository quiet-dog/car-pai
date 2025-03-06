import { request } from "@/utils/service"
import { AreaModel } from "../area"

export interface DeviceModel extends AddDevice {
    id: number
}

export interface AddDevice {
    host: string
    port: string
    hikUsername?: string
    hikPassword?: string
    dhUsername?: string
    dhPassword?: string
    remark?: string
    type: string
    areaId: number
    rtsp: string
    model:string
}

export interface EditDevice extends CId,AddDevice { 
    area?:AreaModel
}

export interface SearchDevice extends PageInfo {
    areaId: number
}

export type DeviceListData = ListData<DeviceModel[]>

export function getDeviceListApi(data: SearchDevice) {
    return request<ApiResponseData<DeviceListData>>({
        url: "/manage/device/getDeviceList",
        method: "post",
        data
    })
}

export function createDeviceApi(data: AddDevice) {
    return request<ApiResponseData<DeviceModel>>({
        url: "/manage/device/createDevice",
        method: "post",
        data
    })
}

export function editDeviceApi(data: EditDevice) {
    return request<ApiResponseData<DeviceModel>>({
        url: "/manage/device/editDevice",
        method: "post",
        data
    })
}

export function deleteDeviceApi(data: CId) {
    return request<ApiResponseData<null>>({
        url: "/manage/device/deleteDevice",
        method: "post",
        data
    })
}
