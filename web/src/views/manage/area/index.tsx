import { AreaModel, createAreaApi, deleteAreaApi, editAreaApi, getAreaListApi } from '@/api/manage/area';
import { onMounted, reactive, ref } from 'vue';
import { usePagination } from "@/hooks/usePagination"
import { Action, ElMessage, ElMessageBox, FormInstance, FormRules } from 'element-plus';

export function useAreaHook() {

    const { paginationData, changeCurrentPage, changePageSize } = usePagination()

    const searchFormData = reactive({
        name: "",
        remark: "",
        page: 1,
        pageSize: 10
    })
    // const formRef = ref<FormInstance>()

    const tableData = ref<AreaModel[]>([])
    const loading = ref(false)

    const dialogVisible = ref(false)
    const title = ref("")
    const formData = reactive<AreaModel>({
        id: 0,
        name: "",
        remark: ""
    })
    const kind = ref("add")
    const mpFormRules: FormRules = reactive({
        name: [
            { required: true, message: "请输入区域名称", trigger: "blur" }
        ],
        remark: [
            { required: false, message: "请输入区域备注", trigger: "blur" }
        ]
    })

    const handleSearch = () => {
        getTable()
    }

    const getTable = () => {
        loading.value = true
        searchFormData.page = paginationData.currentPage
        searchFormData.pageSize = paginationData.pageSize
        getAreaListApi(searchFormData).
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
        formRef.resetFields();
    }

    const closeDialog = () => {
        dialogVisible.value = false;
    }

    const operateAction = (formRef: FormInstance) => {
        formRef.validate(async (valid) => {
            if (valid) {
                if (kind.value === "Add") {
                    const res = await createAreaApi({
                        name: formData.name,
                        remark: formData.remark
                    })
                    if (res.code === 0) {
                        ElMessage({ type: "success", message: res.msg })
                        getTable()
                    }
                } else if (kind.value === "Edit") {
                    const res = await editAreaApi({
                        id: formData.id,
                        name: formData.name,
                        remark: formData.remark
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

    const handleRow = (row: AreaModel, t: string) => {
        if (t === "Edit") {
            formData.id = row.id
            formData.name = row.name
            formData.remark = row.remark
            handleOpen("Edit")
        } else {
            ElMessageBox.alert("确定删除该区域吗?", "提示", {
                confirmButtonText: "确定",
                callback: async (action: Action) => {
                    if (action === "confirm") {
                        const res = await deleteAreaApi({ id: row.id })
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

    onMounted(() => {
        getTable()
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
        handleRow
    }
}