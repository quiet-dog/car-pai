import { DeviceModel, createDeviceApi, deleteDeviceApi, editDeviceApi, getDeviceListApi } from '@/api/manage/device';
import { onMounted, reactive, ref } from 'vue';
import { usePagination } from "@/hooks/usePagination"
import { Action, ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus';
import { AreaListData, AreaModel, getAreaListApi } from '@/api/manage/area';

export function useDeviceHook() {

    const { paginationData, changeCurrentPage, changePageSize } = usePagination()

    const searchFormData = reactive({
        name: "",
        remark: "",
        page: 1,
        pageSize: 10,
        areaId: 0
    })
    const selectAreaOption = ref<AreaModel[]>()
    const selectAreaValue = ref<number>()
    // const formRef = ref<FormInstance>()

    const tableData = ref<DeviceModel[]>([])
    const loading = ref(false)

    const dialogVisible = ref(false)
    const title = ref("")
    let formData = ref<DeviceModel>({
        id: 0,
        host: "",
        port: "",
        hikUsername: "",
        hikPassword: "",
        dhUsername: "",
        dhPassword: "",
        remark: "",
        type: "海康",
        areaId: 0,
        rtsp: "",
        model: "DS-TCG225"
    })
    const kind = ref("add")
    const mpFormRules: FormRules = reactive({
        host: [
            {
                required: true, message: "请输入正确的主机地址", trigger: "blur", validator: (rule: any, value: any, callback: any) => {
                    // 判断这个地址是不是合法的
                    const reg = /^(http|https):\/\/([\w.]+\/?)\S*/
                    if (reg.test(value)) {
                        callback()
                    } else {
                        callback(new Error("请输入正确的主机地址"))
                    }
                }
            }
        ],
        port: [
            { required: true, message: "请输入正确的端口号", trigger: "blur", pattern: /^[0-9]*$/ }
        ],
        hikUsername: [
            { required: true, message: "请输入海康用户名", trigger: "blur" }
        ],
        hikPassword: [
            { required: true, message: "请输入海康密码", trigger: "blur" }
        ],
        dhUsername: [
            { required: true, message: "请输入大华用户名", trigger: "blur" }
        ],
        dhPassword: [
            { required: true, message: "请输入大华密码", trigger: "blur" }
        ],
        areaId: [
            {
                required: true, message: "请选择区域", trigger: "blur", type: "number", validator: (rule: any, value: any, callback: any) => {
                    console.log("aaaav", value)
                    if (value === 0) {
                        callback(new Error("请选择区域"))
                    } else {
                        callback()
                    }
                }
            }
        ],
        rtsp: [
            {
                required: true, message: "请输入rtsp地址", trigger: "blur"
            }
        ],
        model: [
            {
                required: true, message: "请选择型号", trigger: "blur"
            }
        ]
    })

    const handleSearch = () => {
        getTable()
    }

    const getTable = () => {
        loading.value = true
        searchFormData.page = paginationData.currentPage
        searchFormData.pageSize = paginationData.pageSize
        getDeviceListApi(searchFormData).
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
            title.value = "新增区域"
        } else {
            title.value = "编辑区域"
        }
        dialogVisible.value = true
    }

    const handleClose = (formRef: FormInstance) => {
        formData.value = {
            id: 0,
            host: "",
            port: "",
            hikUsername: "",
            hikPassword: "",
            dhUsername: "",
            dhPassword: "",
            remark: "",
            type: "海康",
            areaId: 0,
            rtsp: "",
            model: "DS-TCG225"
        }
        formRef.resetFields();
    }

    const closeDialog = () => {
        dialogVisible.value = false;
    }

    const operateAction = (formRef: FormInstance) => {
        formRef.validate(async (valid) => {
            if (valid) {
                if (formData.value.type == "海康") {
                    formData.value.dhUsername = ""
                    formData.value.dhPassword = ""
                } else {
                    formData.value.hikUsername = ""
                    formData.value.hikPassword = ""
                }

                if (kind.value === "Add") {
                    const res = await createDeviceApi({
                        host: formData.value.host,
                        port: formData.value.port,
                        hikUsername: formData.value.hikUsername,
                        hikPassword: formData.value.hikPassword,
                        dhUsername: formData.value.dhUsername,
                        dhPassword: formData.value.dhPassword,
                        remark: formData.value.remark,
                        type: formData.value.type,
                        areaId: Number(formData.value.areaId),
                        rtsp: formData.value.rtsp,
                        model: formData.value.model
                    })
                    if (res.code === 0) {
                        ElMessage({ type: "success", message: res.msg })
                        getTable()
                    }
                } else if (kind.value === "Edit") {
                    const res = await editDeviceApi({
                        id: formData.value.id,
                        host: formData.value.host,
                        port: formData.value.port,
                        hikUsername: formData.value.hikUsername,
                        hikPassword: formData.value.hikPassword,
                        dhUsername: formData.value.dhUsername,
                        dhPassword: formData.value.dhPassword,
                        remark: formData.value.remark,
                        type: formData.value.type,
                        areaId: Number(formData.value.areaId),
                        rtsp: formData.value.rtsp,
                        model: formData.value.model
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

    const handleRow = (row: DeviceModel, t: string) => {
        if (t === "Edit") {
            formData.value.id = row.id
            formData.value.host = row.host
            formData.value.port = row.port
            formData.value.hikUsername = row.hikUsername
            formData.value.hikPassword = row.hikPassword
            formData.value.dhUsername = row.dhUsername
            formData.value.dhPassword = row.dhPassword
            formData.value.remark = row.remark
            formData.value.type = row.type
            formData.value.areaId = row.areaId
            formData.value.rtsp = row.rtsp
            formData.value.model = row.model
            handleOpen("Edit")
        } else {
            ElMessageBox.alert("确定删除该区域吗?", "提示", {
                confirmButtonText: "确定",
                callback: async (action: Action) => {
                    if (action === "confirm") {
                        const res = await deleteDeviceApi({ id: row.id })
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

    const handleSelectArea = (val: number) => {
        searchFormData.areaId = val
        getTable()
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

    onMounted(() => {
        getTable()
        getSelectOption()
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
        selectAreaOption,
        handleSelectArea,
        selectAreaValue
    }
}