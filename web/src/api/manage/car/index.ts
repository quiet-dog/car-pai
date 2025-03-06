import { request } from "@/utils/service"
import { AreaModel } from '../area/index';


export interface AddCar {
    carNum: string
    areaIds : number[]
    name?: string
    phone?: string
    startTime?: number
    endTime?: number
    remark?: string
    color:string
    carType:string
    listType:string
    cardNo?:string
}

export interface CarModel extends Omit<AddCar, "areaIds"> {
    id: number;
    areas:AreaModel[]
}

export interface EditCar extends CId,AddCar { }

export interface SearchCar extends PageInfo, Omit<AddCar,"areaIds"|"carType"|"listType"|"color"> { 
    areaId:number
}

export type CarListData = ListData<CarModel[]>

export function getCarListApi(data: SearchCar) {
    return request<ApiResponseData<CarListData>>({
        url: "/manage/car/getCarList",
        method: "post",
        data
    })
}

export function createCarApi(data: AddCar) {
    return request<ApiResponseData<CarModel>>({
        url: "/manage/car/createCar",
        method: "post",
        data
    })
}

export function editCarApi(data: EditCar) {
    return request<ApiResponseData<CarModel>>({
        url: "/manage/car/editCar",
        method: "post",
        data
    })
}

export function deleteCarApi(data: CId) {
    return request<ApiResponseData<null>>({
        url: "/manage/car/deleteCar",
        method: "post",
        data
    })
}
