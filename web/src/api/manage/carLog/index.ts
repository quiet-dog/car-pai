import { request } from "@/utils/service"

export interface CarLogModel extends AddCarLog {
    id: number
}

export interface AddCarLog {
    carNum:    string      // 车牌号
	deviceId:  number     // 设备ID
	uri:       string       // 图片地址
	plateType: string // 车牌类型
}

export interface EditCarLog extends CId,AddCarLog { }

export interface SearchCarLog extends PageInfo, AddCarLog { }

export type CarLogListData = ListData<CarLogModel[]>

export function getCarLogListApi(data: SearchCarLog) {
    return request<ApiResponseData<CarLogListData>>({
        url: "/manage/carLog/getCarLogList",
        method: "post",
        data
    })
}

export function createCarLogApi(data: AddCarLog) {
    return request<ApiResponseData<CarLogModel>>({
        url: "/manage/carLog/createCarLog",
        method: "post",
        data
    })
}

export function editCarLogApi(data: EditCarLog) {
    return request<ApiResponseData<CarLogModel>>({
        url: "/manage/carLog/editCarLog",
        method: "post",
        data
    })
}

export function deleteCarLogApi(data: CId) {
    return request<ApiResponseData<null>>({
        url: "/manage/carLog/deleteCarLog",
        method: "post",
        data
    })
}
