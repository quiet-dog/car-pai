import { CarLogModel, getCarLogListApi } from '@/api/manage/carLog';
import { onMounted, reactive, ref } from 'vue';
import { usePagination } from "@/hooks/usePagination"
import { Action, ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus';
import { getUsersApi } from '@/api/authority/user';
import { AreaModel, getAreaListApi } from '@/api/manage/area';

export function useCarLogHook() {

    const { paginationData, changeCurrentPage, changePageSize } = usePagination()

    const searchFormData = reactive({
        carNum: "",
        deviceId: 0,
        uri: "",
        plateType: "",
        page: 1,
        pageSize: 10,
        areaId: 0
    })
    // const formRef = ref<FormInstance>()

    const tableData = ref<CarLogModel[]>([])
    const loading = ref(false)

    const dialogVisible = ref(false)
    const title = ref("")
    // const formData = reactive<CarLogModel>({
    //     id: 0,
    //     name: "",
    //     remark: "",
    //     userIds: []
    // })
    const kind = ref("add")
    const mpFormRules: FormRules = reactive({
        name: [
            { required: true, message: "请输入区域名称", trigger: "blur" }
        ],
        remark: [
            { required: false, message: "请输入区域备注", trigger: "blur" }
        ]
    })
    const selectUserOption = ref()
    const selectAreaValue = ref<number>(0)
    const selectAreaOption = ref<AreaModel[]>()

    const handleSearch = () => {
        getTable()
    }

    const getTable = () => {
        loading.value = true
        searchFormData.page = paginationData.currentPage
        searchFormData.pageSize = paginationData.pageSize
        searchFormData.areaId = selectAreaValue.value
        getCarLogListApi(searchFormData).
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

    const handleSelectArea = (val: number) => {
        searchFormData.areaId = val
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



    // const handleOpen = (t: string) => {
    //     kind.value = t
    //     if (t === "Add") {
    //         title.value = "新增区域"
    //     } else {
    //         title.value = "编辑区域"
    //     }
    //     dialogVisible.value = true
    //     getUserList()
    // }

    // const handleClose = (formRef: FormInstance) => {
    //     formRef.resetFields();
    // }

    // const closeDialog = () => {
    //     dialogVisible.value = false;
    // }

    // const operateAction = (formRef: FormInstance) => {
    //     formRef.validate(async (valid) => {
    //         if (valid) {
    //             if (kind.value === "Add") {
    //                 const res = await createCarLogApi({
    //                     name: formData.name,
    //                     remark: formData.remark,
    //                     userIds: formData.userIds
    //                 })
    //                 if (res.code === 0) {
    //                     ElMessage({ type: "success", message: res.msg })
    //                     getTable()
    //                 }
    //             } else if (kind.value === "Edit") {
    //                 const res = await editCarLogApi({
    //                     id: formData.id,
    //                     name: formData.name,
    //                     remark: formData.remark,
    //                     userIds: formData.userIds
    //                 })
    //                 if (res.code === 0) {
    //                     ElMessage({ type: "success", message: res.msg })
    //                     getTable()
    //                 }
    //             }
    //             closeDialog()
    //         }
    //     })
    // }

    // const handleRow = (row: CarLogModel, t: string) => {
    //     if (t === "Edit") {
    //         formData.id = row.id
    //         formData.name = row.name
    //         formData.remark = row.remark
    //         formData.userIds = row.userIds
    //         handleOpen("Edit")
    //     } else {
    //         ElMessageBox.alert("确定删除该区域吗?", "提示", {
    //             confirmButtonText: "确定",
    //             callback: async (action: Action) => {
    //                 if (action === "confirm") {
    //                     const res = await deleteCarLogApi({ id: row.id })
    //                     if (res.code === 0) {
    //                         ElMessage({ type: "success", message: res.msg })
    //                         getTable()
    //                     } else {
    //                         ElMessage({ type: "error", message: res.msg })
    //                     }
    //                     dialogVisible.value = false
    //                 }
    //             }
    //         })
    //     }
    // }

    // const getUserList = () => {
    //     getUsersApi({ page: 1, pageSize: 1000 }).then(res => {
    //         selectUserOption.value = res.data.list.map(item => {
    //             return {
    //                 label: item.username,
    //                 value: item.id
    //             }
    //         })
    //     })
    // }



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
        // formData,
        // handleOpen,
        // handleClose,
        // closeDialog,
        // operateAction,
        mpFormRules,
        tableData,
        // handleRow,
        selectUserOption,
        selectAreaValue,
        handleSelectArea,
        selectAreaOption
    }
}