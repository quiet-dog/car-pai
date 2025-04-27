import { AddCar, CarModel, createCarApi, deleteCarApi, EditCar, editCarApi, getCarListApi, getSelectCarApi, Select } from '@/api/manage/car';
import { onMounted, reactive, ref } from 'vue';
import { usePagination } from "@/hooks/usePagination"
import { Action, ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus';
import { AddArea, AreaModel, getAreaListApi } from '@/api/manage/area';

export function useCarHook() {

    const { paginationData, changeCurrentPage, changePageSize } = usePagination()

    const searchFormData = reactive({
        name: "",
        remark: "",
        page: 1,
        pageSize: 10,
        phone: "",
        carNum: "",
        startTime: 0,
        endTime: 0,
        areaId: 0,
    })
    // const formRef = ref<FormInstance>()

    const tableData = ref<CarModel[]>([])
    const loading = ref(false)

    const dialogVisible = ref(false)
    const title = ref("")
    const formData = reactive<EditCar>({
        id: 0,
        carNum: "",
        areaIds: [],
        name: "",
        phone: "",
        startTime: 0,
        endTime: 0,
        remark: "",
        color: "0",
        carType: "0",
        listType: "0",
        cardNo: "",
        deviceIds: [],
    })
    const kind = ref("add")
    const mpFormRules: FormRules = reactive({
        carNum: [
            { required: true, message: "请输入车牌号", trigger: "blur" }
        ],
        // deviceIds: [
        //     {
        //         required: true, message: "请选择区域设备", trigger: "blur", type: "array", validator: (rule: any, value: any, callback: any) => {
        //             if (value.length === 0) {
        //                 callback(new Error("请选择区域设备"))
        //             } else {
        //                 callback()
        //             }
        //         }
        //     }
        // ],
        startTime: [
            {
                required: true, message: "请选择时间", trigger: "change", type: "number", validator: (rule: any, value: any, callback: any) => {
                    if (value === 0) {
                        callback(new Error("请选择有效期"))
                    } else {
                        callback()
                    }
                }
            }
        ],
    })
    const selectAreaOption = ref<AreaModel[]>([])
    const timePicker = ref([])
    const select = ref<Select[]>([])
    const deviceIds = ref<number[][]>([])

    const handleChangeSelect = (val: number[][]) => {
        if (!val) {
            deviceIds.value = []
            return
        }
        let ids: number[] = []
        val.forEach((item: number[]) => {
            if (item.length > 1) {
                ids.push(item[1])
            }
        })
        deviceIds.value = val
    }

    const handleSearch = () => {
        getTable()
    }

    const getTable = () => {
        loading.value = true
        searchFormData.page = paginationData.currentPage
        searchFormData.pageSize = paginationData.pageSize
        getCarListApi(searchFormData).
            then(res => {
                tableData.value = res.data.list
                paginationData.total = res.data.total
            }).
            catch(err => { }).
            finally(() => {
                loading.value = false
            })
    }

    const handleSizeChange = (val: number) => {
        paginationData.pageSize = val
        getTable()
    }

    const handleCurrentChange = (val: number) => {
        paginationData.currentPage = val
        getTable()
    }

    const handleRefresh = () => {
        getTable()
    }

    const handleOpen = (t: string) => {
        kind.value = t
        if (t === "Add") {
            title.value = "新增车牌号"
        } else {
            title.value = "编辑车牌号"
        }
        dialogVisible.value = true
    }

    const handleClose = (formRef: FormInstance) => {
        formRef.resetFields();
        timePicker.value = []
        deviceIds.value = []
    }

    const closeDialog = () => {
        dialogVisible.value = false;
    }

    const operateAction = (formRef: FormInstance) => {
        formRef.validate(async (valid) => {
            if (valid) {
                if (kind.value === "Add") {
                    const res = await createCarApi({
                        name: formData.name,
                        remark: formData.remark,
                        areaIds: formData.areaIds,
                        phone: formData.phone,
                        carNum: formData.carNum,
                        startTime: formData.startTime,
                        endTime: formData.endTime,
                        carType: formData.carType,
                        listType: formData.listType,
                        color: formData.color,
                        cardNo: formData.cardNo,
                        deviceIds: deviceIds.value.map(item => item[1])
                    })
                    if (res.code === 0) {
                        ElMessage({ type: "success", message: res.msg })
                        getTable()
                    }
                } else if (kind.value === "Edit") {
                    // deviceIds.value = formData.devices?.map(item => [item.areaId, item.id]) || []
                    // console.log("deviceIds", deviceIds.value)
                    const res = await editCarApi({
                        id: formData.id,
                        name: formData.name,
                        remark: formData.remark,
                        areaIds: formData.areaIds,
                        phone: formData.phone,
                        carNum: formData.carNum,
                        startTime: formData.startTime,
                        endTime: formData.endTime,
                        carType: formData.carType,
                        listType: formData.listType,
                        color: formData.color,
                        cardNo: formData.cardNo,
                        deviceIds: deviceIds.value.map(item => item[1])
                    })
                    if (res.code === 0) {
                        ElMessage({ type: "success", message: res.msg })
                        getTable()
                        // // 替换数据
                        // const index = tableData.value.indexOf(res.data)
                        // tableData.value.splice(index, 1, res.data)
                    }
                }
                closeDialog()
            }
        })
    }

    const handleRow = (row: CarModel, t: string) => {
        if (t === "Edit") {
            deviceIds.value = row.devices?.map(item => [item.areaId, item.id]) || []
            console.log("deviceIds", deviceIds.value)
            formData.id = row.id
            formData.name = row.name
            formData.remark = row.remark
            formData.areaIds = row.areas?.map(item => item.id)
            formData.phone = row.phone
            formData.carNum = row.carNum
            formData.startTime = row.startTime
            formData.endTime = row.endTime
            timePicker.value = [row.startTime!, row.endTime!] as any
            formData.carType = row.carType
            formData.listType = row.listType
            formData.color = row.color
            handleOpen("Edit")
        } else {
            ElMessageBox.alert("确定删除该车牌号吗?", "提示", {
                confirmButtonText: "确定",
                callback: async (action: Action) => {
                    if (action === "confirm") {
                        const res = await deleteCarApi({ id: row.id })
                        if (res.code === 0) {
                            ElMessage({ type: "success", message: res.msg })
                            getTable()
                        } else {
                            ElMessage({ type: "error", message: res.msg })
                        }
                        dialogVisible.value = false
                    }
                }
            })
        }
    }

    const getSelectOption = () => {
        getAreaListApi({
            name: "",
            remark: "",
            page: 1,
            pageSize: 1000,
        }).then(res => {
            selectAreaOption.value = res.data.list
        }).then(res => { })
    }

    const handleSelectArea = (val: number) => {
        // searchFormData.a
    }


    // const handleChangePicker = (val: any) => {
    //     console.log("valu", val)
    //     if (val == null) {
    //         formData.startTime = 0
    //         formData.endTime = 0
    //         return
    //     }
    //     formData.startTime = val[0]
    //     formData.endTime = val[1]
    // }

    onMounted(() => {
        getTable()
        getSelectOption()
        getSelectCarApi().then(res => {
            select.value = res.data
        })
    })

    return {
        searchFormData,
        handleSearch,
        loading,
        paginationData,
        changeCurrentPage,
        changePageSize,
        handleSizeChange,
        handleCurrentChange,
        handleRefresh,
        dialogVisible,
        title,
        formData,
        handleOpen,
        handleClose,
        closeDialog,
        operateAction,
        mpFormRules,
        tableData,
        handleRow,
        getSelectOption,
        selectAreaOption,
        timePicker,
        select,
        handleChangeSelect,
        deviceIds
    }
}