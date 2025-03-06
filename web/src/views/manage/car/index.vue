<script lang='ts' setup>
import { ref } from 'vue';
import { useCarHook } from './index.tsx';
import { FormInstance } from 'element-plus';
import { formatDateTime } from '@/utils/index.ts';

defineOptions({
    name: "Car"
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
    selectAreaOption,
    timePicker,
    handleChangePicker
    // formRules, formRef, closeDialog, operateAction,handleClose
} = useCarHook();
</script>

<template>
    <div class="app-container">
        <!-- 搜索框 -->
        <el-card shadow="never" class="search-wrapper">
            <el-form ref="searchFormRef" :inline="true" :model="searchFormData">
                <el-form-item prop="carNum" label="车牌号">
                    <el-input v-model="searchFormData.carNum" placeholder="车牌号" />
                </el-form-item>
                <el-form-item prop="areaId" label="区域">
                    <el-select placeholder="请选择对应区域" style="width: 240px" clearable v-model="searchFormData.areaId"
                        collapse-tags-tooltip>
                        <el-option v-for="item in selectAreaOption" :key="item.id" :value="item.id"
                            :label="item.name" />
                    </el-select>
                </el-form-item>
                <el-form-item prop="remark" label="备注">
                    <el-input v-model="searchFormData.remark" placeholder="备注" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="Search" @click="handleSearch">查询</el-button>
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
                    <el-table-column prop="carNum" label="车牌号" min-width="200" />
                    <el-table-column prop="name" label="车主姓名" min-width="200" />
                    <el-table-column prop="phone" label="车主电话" min-width="200" />
                    <el-table-column prop="startTime" label="开始时间" min-width="200">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.startTime) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="endTime" label="结束时间" min-width="200">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.endTime) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="备注" min-width="200" />
                    <el-table-column label="操作" width="200">
                        <template #default="scope">
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
                <el-form-item label="车牌号" prop="carNum">
                    <el-input v-model="formData.carNum" placeholder="请输入车牌号" />
                </el-form-item>
                <el-form-item label="区域" prop="areaIds">
                    <el-select placeholder="请选择对应区域" v-model="formData.areaIds" collapse-tags-tooltip multiple
                        collapse-tags>
                        <el-option v-for="item in selectAreaOption" :key="item.id" :value="item.id"
                            :label="item.name" />
                    </el-select>
                </el-form-item>
                <el-form-item label="黑白名单" prop="listType">
                    <el-select v-model="formData.listType" placeholder="Select" size="large">
                        <el-option label="白名单" value="0" />
                        <el-option label="黑名单" value="1" />
                    </el-select>
                </el-form-item>
                <el-form-item label="卡号" prop="cardNo">
                    <el-input v-model="formData.cardNo" placeholder="请输入卡号" />
                </el-form-item>
                <el-form-item label="有效期" prop="startTime">
                    <el-date-picker @change="handleChangePicker" v-model="timePicker" type="datetimerange"
                        start-placeholder="开始时间" end-placeholder="结束时间" format="YYYY-MM-DD HH:mm:ss"
                        date-format="YYYY/MM/DD ddd" time-format="A hh:mm:ss" value-format="x" />
                </el-form-item>
                <el-form-item label="车牌类型" prop="carType">
                    <el-select v-model="formData.carType" placeholder="Select" size="large">
                        <el-option label="标准民用用车与军车" value="0" />
                        <el-option label="02式民用车牌" value="1" />
                        <el-option label="武警车" value="2" />
                        <el-option label="警车" value="3" />
                        <el-option label="民用车双行尾牌" value="4" />
                        <el-option label="使馆车牌" value="5" />
                        <el-option label="农用车牌" value="6" />
                        <el-option label="摩托车牌" value="7" />
                        <el-option label="新能源车牌" value="8" />
                    </el-select>
                </el-form-item>
                <el-form-item label="颜色" prop="color">
                    <el-select v-model="formData.color" placeholder="Select" size="large">
                        <el-option label="蓝色" value="0" />
                        <el-option label="黄色" value="1" />
                        <el-option label="白色" value="2" />
                        <el-option label="黑色" value="3" />
                        <el-option label="绿色" value="4" />
                        <el-option label="其它" value="5" />
                    </el-select>
                </el-form-item>
                <el-form-item label="车主姓名" prop="name">
                    <el-input v-model="formData.name" placeholder="请输入车主姓名" />
                </el-form-item>
                <el-form-item label="车主电话" prop="phone">
                    <el-input v-model="formData.phone" placeholder="请输入车主电话" />
                </el-form-item>

                <el-form-item label="备注" prop="remark">
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

<style scoped></style>
