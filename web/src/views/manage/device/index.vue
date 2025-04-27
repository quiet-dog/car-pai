<script lang='ts' setup>
import { ref } from 'vue';
import { useDeviceHook } from './index.tsx';
import { FormInstance } from 'element-plus';
import { formatDateTime } from '@/utils/index.ts';

defineOptions({
    name: "Device"
})

const formRef = ref<FormInstance>()
const {
    searchFormData, handleSearch, loading,
    handleCurrentChange, handleSizeChange,
    paginationData, handleRefresh,
    dialogVisible, title, formData,
    handleOpen,
    handleClose, closeDialog,
    operateAction, mpFormRules,
    tableData, handleRow,
    selectAreaOption, handleSelectArea,
    selectAreaValue,
    opearaRemote
    // formRules, formRef, closeDialog, operateAction,handleClose
} = useDeviceHook();
</script>

<template>
    <div class="app-container">
        <!-- 搜索框 -->
        <el-card shadow="never" class="search-wrapper">
            <el-form ref="searchFormRef" :inline="true" :model="searchFormData">
                <el-form-item prop="name" label="区域选择">
                    <el-select clearable v-model="selectAreaValue" placeholder="选择区域筛选设备" style="width: 240px"
                        @change="handleSelectArea">
                        <el-option v-for="item in selectAreaOption" placeholder="请选择区域筛选" :label="item.name"
                            :value="item.id" />
                    </el-select>
                </el-form-item>
            </el-form>
        </el-card>

        <!-- 表格 -->
        <el-card v-loading="loading" shadow="never">
            <div class="toolbar-wrapper">
                <el-button type="primary" plain icon="plus" @click="handleOpen('Add')">新增</el-button>
                <div>
                    <el-tooltip content="刷新" effect="light">
                        <el-button type="primary" icon="RefreshRight" circle plain @click="handleRefresh" />
                    </el-tooltip>
                </div>
            </div>
            <div class="table-wrapper">
                <el-table :data="tableData" border stripe>
                    <el-table-column prop="host" label="主机地址" min-width="200" />
                    <el-table-column prop="port" label="端口" min-width="200" />
                    <el-table-column prop="hikUsername" label="海康用户名" min-width="200" />
                    <el-table-column prop="hikPassword" label="海康密码" min-width="200" />
                    <el-table-column prop="dhUsername" label="大华用户名" min-width="200" />
                    <el-table-column prop="dhPassword" label="大华密码" min-width="200" />
                    <!-- <el-table-column prop="rtsp" label="流地址" min-width="200" /> -->
                    <el-table-column prop="model" label="型号" min-width="200" />
                    <el-table-column prop="area.name" label="所属区域" min-width="200" />
                    <el-table-column prop="remark" label="备注" min-width="200" />
                    <el-table-column prop="createdAt" label="创建时间">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.createdAt) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="200">
                        <template #default="scope">
                            <!-- <el-dropdown @command="opearaRemote">
                                <span class="el-dropdown-link">
                                    操作
                                    <el-icon class="el-icon--right">
                                        <arrow-down />
                                    </el-icon>
                                </span>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item :command="{
                                            lockStatus: 'open',
                                            id: scope.row.id
                                        }">开闸</el-dropdown-item>
                                        <el-dropdown-item :command="{
                                            lockStatus: 'close',
                                            id: scope.row.id
                                        }">关闸</el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown> -->
                            <el-button type="text" @click="handleRow(scope.row, 'Edit')">编辑</el-button>
                            <el-button type="danger" @click="handleRow(scope.row, '')" text>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div class="pager-wrapper">
                <el-pagination background :layout="paginationData.layout" :page-sizes="paginationData.pageSizes"
                    :total="paginationData.total" :page-size="paginationData.pageSize"
                    :currentPage="paginationData.currentPage" @size-change="handleSizeChange"
                    @current-change="handleCurrentChange" />
            </div>
        </el-card>

        <!-- 编辑添加 -->
        <el-dialog v-model="dialogVisible" :title="title" @closed="handleClose(formRef!)" width="30%">
            <el-form ref="formRef" :model="formData" :rules="mpFormRules" label-width="100px" label-position="left"
                style="width: 95%; margin-top: 15px">
                <el-form-item label="设备类型" prop="type">
                    <el-radio-group v-model="formData.type">
                        <el-radio value="海康">海康</el-radio>
                        <el-radio value="大华">大华</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="主机地址" prop="host">
                    <el-input v-model="formData.host" placeholder="请输入主机地址" />
                </el-form-item>
                <el-form-item label="端口" prop="port">
                    <el-input v-model="formData.port" placeholder="请输入主机端口" />
                </el-form-item>
                <el-form-item label="区域选择" prop="areaId">
                    <el-select v-model="formData.areaId" placeholder="选择区域" style="width: 240px">
                        <el-option v-for="item in selectAreaOption" :label="item.name" :value="item.id" />
                    </el-select>
                </el-form-item>
                <!-- <el-form-item label="流地址" prop="rtsp">
                    <el-input v-model="formData.rtsp" placeholder="请输入流地址" />
                </el-form-item> -->
                <el-form-item v-if="formData.type == '海康'" label="型号" prop="model">
                    <el-radio-group v-model="formData.model">
                        <el-radio value="DS-TCG225">DS-TCG225</el-radio>
                        <el-radio value="DS-TCG205-E">DS-TCG205-E</el-radio>
                        <el-radio value="DS-TCG2A5-E">DS-TCG2A5-E</el-radio>
                        <el-radio value="DS-TCG2A5-B">DS-TCG2A5-B</el-radio>
                        <el-radio value="DS-2CD9125-KS">DS-2CD9125-KS</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item v-if="formData.type == '海康'" label="海康用户名" prop="hikUsername">
                    <el-input v-model="formData.hikUsername" placeholder="请输入海康用户名" />
                </el-form-item>
                <el-form-item v-if="formData.type == '海康'" label="海康密码" prop="hikPassword">
                    <el-input v-model="formData.hikPassword" placeholder="请输入备注" />
                </el-form-item>
                <el-form-item v-if="formData.type == '大华'" label="大华用户名" prop="dhUsername">
                    <el-input v-model="formData.dhUsername" placeholder="请输入大华用户名" />
                </el-form-item>
                <el-form-item v-if="formData.type == '大华'" label="大华密码" prop="dhPassword">
                    <el-input v-model="formData.dhPassword" placeholder="请输入备注" />
                </el-form-item>
                <el-form-item prop="remark" label="备注">
                    <el-input v-model="formData.remark" placeholder="请输入备注" />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="closeDialog">取消</el-button>
                    <el-button type="primary" @click="operateAction(formRef!)">确认</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<style scoped>
.example-showcase .el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
</style>
